package nertiviago

type State struct {
	CustomStatuses    [][]string     `json:"customStatusArr"`
	DMs               []dm           `json:"dms"`
	MemberStatuses    [][]string     `json:"memberStatusArr"`
	LoginResponse     string         `json:"message"`
	MutedChannels     []string `json:"mutedChannels"`
	Notifications     []notification
	ProgramActivities []interface{}  `json:"programActivityArr"`
	ServerMembers     []serverMember `json:"serverMembers"`
	ServerRoles       []serverRole   `json:"serverRoles"`
	Settings          settings
	ThisClient        StateUser `json:"user"`
}

type settings struct {
	GDriveLinked bool     `json:"GDriveLinked"`
	CustomEmojis []emoji `json:"customEmojis"`
}

type emoji struct {
	V int `json:"__v"`
	HiddenID string `json:"_id"`
	ID string `json:"emojiID"`
	Gif bool `json:"gif"`
	Name string `json:"name"`
	User string `json:"user"`
}

type dm struct {
	ChannelID    string   `json:"channelID"`
	LastMessaged int      `json:"lastMessaged"`
	Recipients   []Member `json:"recipients"`
}

type notification struct {
	ChannelID     string `json:"channelID"`
	Count         int
	LastMessageID string `json:"lastMessageID"`
	Mentioned     bool
	Recipient     string
	Sender        *User
	Type          string
}

type StateUser struct {
	HiddenID        string `json:"_id"`
	Admin           int
	Avatar          string `json:"avatar,omitempty"`
	Bot             bool
	Friends         []*User
	Servers         []interface{}
	Status          int
	SurveyCompleted bool `json:"survey_completed"`
	Tag             string
	ID              string `json:"uniqueID"`
	Username        string
}

type serverMember struct {
	Member   *Member
	Roles    []string
	ServerID string `json:"server_id"`
	Type     string
}

type serverRole struct {
	Bot         string `json:"bot,omitempty"`
	Default     bool
	Deletable   bool
	ID          string
	Name        string
	Order       int
	Permissions int
	ServerID    string `json:"server_id"`
}

type Member struct {
	Avatar   string `json:"avatar,omitempty"`
	Bot      bool
	Tag      string
	ID       string `json:"uniqueID"`
	Username string
}

type programActivity struct {
}
