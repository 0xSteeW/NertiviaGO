package nertivia

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
	"io/ioutil"
	"log"
	"nertivia/globals"
	"net/http"
)

type Session struct {
	Token   string
	Client  *gosocketio.Client
	State   *State
	Channel *gosocketio.Channel
}

type sidResponse struct {
	SID      string
	Upgrades interface{}
}

// New creates a new session struct with provided token
func New(token string) *Session {
	sess := new(Session)
	sess.Token = token
	return sess
}

func getSID() (string, error) {
	end := globals.ReadConstants()
	res, err := http.Get("https://" + end.EndpointURL + "/socket.io/?EIO=3&transport=polling")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	sidResp := new(sidResponse)
	format := b[4 : len(b)-4]
	err = json.Unmarshal(format, sidResp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sidResp.SID)
	return sidResp.SID, nil
}

type auth struct {
	Authentication map[string]string
}

func getAuth(token string) *auth {
	a := new(auth)
	a.Authentication = make(map[string]string)
	a.Authentication["token"] = token
	return a
}

func (s *Session) Ping() bool {
	return s.Client.IsAlive()
}

// Open creates a websocket with socket.io using the provided token.
func (s *Session) Open() error {
	end := globals.ReadConstants()
	client, err := gosocketio.Dial(
		gosocketio.AddrWebsocket(end.EndpointURL, 443, true),
		transport.DefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = client.On(gosocketio.OnConnection, func(channel *gosocketio.Channel) { fmt.Println("Connected to Nertivia websocket.") })
	if err != nil {
		log.Fatal(err)
		return err
	}
	auth := make(map[string]string)
	auth["token"] = s.Token
	client.Emit("authentication", auth)
	err = client.On("success", func(channel *gosocketio.Channel, data interface{}) { fmt.Println("Authorized") })
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.Channel = client.Channel
	s.Client = client
	return nil
}

// Handlers

func (s *Session) OnMessage(handler func(*Session, *MessageCreate)) error {
	fmt.Println("Adding handler")
	messageCreate := NewMessageHandler()
	err := s.Client.On(OnMessageCreate, func(c *gosocketio.Channel, data interface{}) {
		body, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, messageCreate)
		if err != nil {
			log.Fatal(err)
			return
		}
		handler(s, messageCreate)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) GetUser(userID string) (*UserEvent, error) {
	user := new(UserEvent)
	end := globals.ReadConstants()
	err := s.Request(user, end.EndpointUser, "/", userID)
	if err != nil {
		return &UserEvent{}, err
	}
	return user, nil
}

func (s *Session) GetChannel(channelID string) (*ChannelEvent, error) {
	channel := new(ChannelEvent)
	end := globals.ReadConstants()
	err := s.Request(channel, end.EndpointChannel, "/", channelID)
	if err != nil {
		return &ChannelEvent{}, err
	}
	return channel, nil
}

func (s *Session) ChannelMessageSend(channelID string, message string) error {
	channel, err := s.GetChannel(channelID)
	end := globals.ReadConstants()
	if err != nil {
		return errors.New("could not find channel " + channelID)
	}
	data := MessageSend{Message: message, TempID: "0"}
	dataByte, _ := json.Marshal(data)
	scode, err := s.Send(dataByte, end.EndpointChannel, "/", channel.ChannelID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(scode)
	return nil
}

func formatParams(strings []string) string {
	var final string
	for _, str := range strings {
		final = final + str
	}
	return final
}

func (s *Session) Send(data []byte, endpoint string, params ...string) (int, error) {
	url := fmt.Sprint(endpoint, formatParams(params))
	fmt.Println(string(data))
	bodyPost := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, bodyPost)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", s.Token)
	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	return response.StatusCode, err
}

func (s *Session) Request(event Event, endpoint string, params ...string) error {
	url := fmt.Sprint(endpoint, formatParams(params))
	fmt.Println(url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", s.Token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, event)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Session) Close() {
	s.Client.Close()
}
