package comm

import (
	"errors"
	"log"

	"github.com/marcsantiago/gojay"
	"golang.org/x/net/websocket"
)

// A basic message for our WebSocket app

type Message struct {
	Message  string
	UserName string
}

func (m *Message) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "message":
		return dec.AddString(&m.Message)
	case "userName":
		return dec.AddString(&m.UserName)
	}
	return nil
}

func (m *Message) NKeys() int {
	return 2
}

func (m *Message) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey("message", m.Message)
	enc.AddStringKey("userName", m.UserName)
}

func (m *Message) IsNil() bool {
	return m == nil
}

// Here are defined our communication types
type Sender chan gojay.MarshalerJSONObject

func (s Sender) MarshalStream(enc *gojay.StreamEncoder) {
	select {
	case <-enc.Done():
		return
	case m := <-s:
		enc.AddObject(m)
	}
}

type Receiver chan *Message

func (s Receiver) UnmarshalStream(dec *gojay.StreamDecoder) error {
	m := &Message{}
	if err := dec.AddObject(m); err != nil {
		return err
	}
	s <- m
	return nil
}

type SenderReceiver struct {
	Send    Sender
	Receive Receiver
	Dec     *gojay.StreamDecoder
	Enc     *gojay.StreamEncoder
	Conn    *websocket.Conn
}

func (c *SenderReceiver) SetReceiver() {
	c.Receive = Receiver(make(chan *Message))
	c.Dec = gojay.Stream.BorrowDecoder(c.Conn)
	go func() { _ = c.Dec.DecodeStream(c.Receive) }()
}

func (c *SenderReceiver) SetSender(nCons int) {
	c.Send = Sender(make(chan gojay.MarshalerJSONObject))
	c.Enc = gojay.Stream.BorrowEncoder(c.Conn).NConsumer(nCons).LineDelimited()
	go c.Enc.EncodeStream(c.Send)
}

func (c *SenderReceiver) SendMessage(m gojay.MarshalerJSONObject) error {
	select {
	case <-c.Enc.Done():
		return errors.New("sender closed")
	case c.Send <- m:
		log.Print("message sent by client: ", m)
		return nil
	}
}

func (c *SenderReceiver) OnMessage(f func(*Message)) error {
	for {
		select {
		case <-c.Dec.Done():
			return errors.New("receiver closed")
		case m := <-c.Receive:
			f(m)
		}
	}
}

func (c *SenderReceiver) Init(sender int) *SenderReceiver {
	c.SetSender(sender)
	c.SetReceiver()
	return c
}
