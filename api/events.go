package nertivia

import "fmt"

type Event interface {
	Get() Event
}

type UserEvent struct {
	User *User
	CommonServersArr []string
	CommonFriendsArr []string
	IsBlocked bool
}

type User struct {
	ID string `json:"_id"`
	Avatar string
	Admin int
	Badges []int
	Username string
	UniqueID string
	Tag string
	Created int
	About map[string]string `json:"about_me"`
}

func (u UserEvent) Get() Event {
	return u
}

func (u *UserEvent) String() string {
	return fmt.Sprint(u.User.Username,"#",u.User.Tag)
}
