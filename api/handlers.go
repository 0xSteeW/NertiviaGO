package nertivia

type Handler interface {
	Get() Handler
}
type MessageCreate struct {

}

func (mc *MessageCreate) Get() Handler {
	return mc
}
