package nertiviago

type Handler interface {
	GetHandler() string
}
type MessageCreate struct {
	Message *Message
}

type Message struct {
	ChannelID string
	Created   int64
	Creator   *User
	Mentions  []*User
	Content   string `json:"message"`
	Quotes    []string
	ID        string `json:"message_id"`
}

func NewMessageHandler() *MessageCreate {
	return &MessageCreate{}
}
func (mc *MessageCreate) GetHandler() string {
	return "test"
}

type HandlerFunction func(Handler)
