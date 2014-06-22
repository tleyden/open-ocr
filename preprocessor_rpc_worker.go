package ocrworker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type PreprocessorRpcWorker struct {
	rabbitConfig    RabbitConfig
	conn            *amqp.Connection
	channel         *amqp.Channel
	tag             string
	Done            chan error
	bindingKey      string
	preprocessorMap map[string]Preprocessor
}

const preprocessor_tag = "preprocessor" // TODO: should be unique for each worker instance (eg, uuid)

func NewPreprocessorRpcWorker(rc RabbitConfig, preprocessor string) (*PreprocessorRpcWorker, error) {

	preprocessorMap := make(map[string]Preprocessor)
	preprocessorMap[PREPROCESSOR_STROKE_WIDTH_TRANSFORM] = StrokeWidthTransformer{}
	preprocessorMap[PREPROCESSOR_IDENTITY] = IdentityPreprocessor{}

	_, ok := preprocessorMap[preprocessor]
	if !ok {
		return nil, fmt.Errorf("No preprocessor found for: %q", preprocessor)
	}

	preprocessorRpcWorker := &PreprocessorRpcWorker{
		rabbitConfig:    rc,
		conn:            nil,
		channel:         nil,
		tag:             preprocessor_tag,
		Done:            make(chan error),
		bindingKey:      preprocessor,
		preprocessorMap: preprocessorMap,
	}
	return preprocessorRpcWorker, nil
}

func (w PreprocessorRpcWorker) Run() error {

	var err error

	logg.LogTo("PREPROCESSOR_WORKER", "Run() called...")

	logg.LogTo("PREPROCESSOR_WORKER", "dialing %q", w.rabbitConfig.AmqpURI)
	w.conn, err = amqp.Dial(w.rabbitConfig.AmqpURI)
	if err != nil {
		return err
	}

	go func() {
		fmt.Printf("closing: %s", <-w.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	logg.LogTo("PREPROCESSOR_WORKER", "got Connection, getting Channel")
	w.channel, err = w.conn.Channel()
	if err != nil {
		return err
	}

	if err = w.channel.ExchangeDeclare(
		w.rabbitConfig.Exchange,     // name of the exchange
		w.rabbitConfig.ExchangeType, // type
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	// just call the queue the same name as the binding key, since
	// there is no reason to have a different name.
	queueName := w.bindingKey

	queue, err := w.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	if err = w.channel.QueueBind(
		queue.Name,              // name of the queue
		w.bindingKey,            // bindingKey
		w.rabbitConfig.Exchange, // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "Queue bound to Exchange, starting Consume (consumer tag %q, binding key: %v)", preprocessor_tag, w.bindingKey)
	deliveries, err := w.channel.Consume(
		queue.Name,       // name
		preprocessor_tag, // consumerTag,
		true,             // noAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	go w.handle(deliveries, w.Done)

	return nil
}

func (w *PreprocessorRpcWorker) Shutdown() error {
	// will close() the deliveries channel
	if err := w.channel.Cancel(w.tag, true); err != nil {
		return fmt.Errorf("Worker cancel failed: %s", err)
	}

	if err := w.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logg.LogTo("PREPROCESSOR_WORKER", "Shutdown OK")

	// wait for handle() to exit
	return <-w.Done
}

func (w *PreprocessorRpcWorker) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		logg.LogTo(
			"PREPROCESSOR_WORKER",
			"got %d byte delivery: [%v]. Routing key: %s Reply to: %v",
			len(d.Body),
			d.DeliveryTag,
			d.RoutingKey,
			d.ReplyTo,
		)

		err := w.handleDelivery(d)
		if err != nil {
			msg := "Error handling delivery in preprocessor.  Error: %v"
			logg.LogError(fmt.Errorf(msg, err))
		}

	}
	logg.LogTo("PREPROCESSOR_WORKER", "handle: deliveries channel closed")
	done <- fmt.Errorf("handle: deliveries channel closed")
}

func (w *PreprocessorRpcWorker) preprocessImage(ocrRequest *OcrRequest) error {

	descriptor := w.bindingKey // eg, "stroke-width-transform"
	preprocessor := w.preprocessorMap[descriptor]
	logg.LogTo("PREPROCESSOR_WORKER", "Preproces %v via %v", ocrRequest, descriptor)

	err := preprocessor.preprocess(ocrRequest)
	if err != nil {
		msg := "Error doing %s on: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, descriptor, ocrRequest, err)
		logg.LogError(fmt.Errorf(errMsg))
		return err
	}
	return nil

}

func (w *PreprocessorRpcWorker) strokeWidthTransform(ocrRequest *OcrRequest) error {

	// write bytes to a temp file

	tmpFileNameInput, err := createTempFileName()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameInput)

	tmpFileNameOutput, err := createTempFileName()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameOutput)

	err = saveBytesToFileName(ocrRequest.ImgBytes, tmpFileNameInput)
	if err != nil {
		return err
	}

	// run DecodeText binary on it (if not in path, print warning and do nothing)
	darkOnLightSetting := "1" // todo: this should be passed as a param.
	out, err := exec.Command(
		"DetectText",
		tmpFileNameInput,
		tmpFileNameOutput,
		darkOnLightSetting,
	).CombinedOutput()
	if err != nil {
		logg.LogFatal("Error running command: %s.  out: %s", err, out)
	}
	logg.LogTo("PREPROCESSOR_WORKER", "output: %v", string(out))

	// read bytes from output file into ocrRequest.ImgBytes
	resultBytes, err := ioutil.ReadFile(tmpFileNameOutput)
	if err != nil {
		return err
	}

	ocrRequest.ImgBytes = resultBytes

	return nil

}

func (w *PreprocessorRpcWorker) handleDelivery(d amqp.Delivery) error {

	ocrRequest := OcrRequest{}
	err := json.Unmarshal(d.Body, &ocrRequest)
	if err != nil {
		msg := "Error unmarshaling json: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, string(d.Body), err)
		logg.LogError(fmt.Errorf(errMsg))
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "ocrRequest before: %v", ocrRequest)
	routingKey := ocrRequest.nextPreprocessor(w.rabbitConfig.RoutingKey)
	logg.LogTo("PREPROCESSOR_WORKER", "publishing with routing key %q", routingKey)
	logg.LogTo("PREPROCESSOR_WORKER", "ocrRequest after: %v", ocrRequest)

	err = w.preprocessImage(&ocrRequest)
	if err != nil {
		msg := "Error preprocessing image: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, ocrRequest, err)
		logg.LogError(fmt.Errorf(errMsg))
		return err
	}

	ocrRequestJson, err := json.Marshal(ocrRequest)
	if err != nil {
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "sendRpcResponse to: %v", routingKey)

	if err := w.channel.Publish(
		w.rabbitConfig.Exchange, // publish to an exchange
		routingKey,              // routing to 0 or more queues
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(ocrRequestJson),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			ReplyTo:         d.ReplyTo,
			CorrelationId:   d.CorrelationId,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return err
	}
	logg.LogTo("PREPROCESSOR_WORKER", "handleDelivery succeeded")

	return nil
}
