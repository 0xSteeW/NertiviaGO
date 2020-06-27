package nertivia

import (
	"fmt"
)

type Event interface {
	Get() Event
}

type MessageSend struct {
	Message string `json:"message"`
	TempID  string `json:"tempID"`
}

type UserEvent struct {
	User             *User
	CommonServersArr []string
	CommonFriendsArr []string
	IsBlocked        bool
}

type User struct {
	ID       string `json:"_id"`
	Avatar   string
	Admin    int
	Badges   []int
	Username string
	UniqueID string
	Tag      string
	Created  int
	About    map[string]string `json:"about_me"`
}

type ChannelEvent struct {
	Status    bool
	ChannelID string
	Messages  []map[string]interface{}
}

func (c ChannelEvent) Get() Event {
	return c
}

func (u UserEvent) Get() Event {
	return u
}

func (u *User) String() string {
	return fmt.Sprint(u.Username, "#", u.Tag)
}

//Server event
type ServerEvent struct {
	Name string
	Avatar string
	DefaultChannel string `json:"default_channel_id"`
	ID string `json:"server_id"`
	Created int
	Banner string
}

func (s *ServerEvent) Get() Event {
	return s
}

const (
	OnMessageCreate = "receiveMessage"
	OnButtonClick = "messageButtonClicked"
)
