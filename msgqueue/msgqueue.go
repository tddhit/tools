package msgqueue

type Producer interface {
	Publish(topic string, body []byte) error
}

type Consumer interface {
	Messages(topic string) <-chan *Message
}

type Message struct {
	Body []byte
}
