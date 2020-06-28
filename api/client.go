package nertiviago

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
	"io/ioutil"
	"log"
	"nertiviago/globals"
	"net/http"
	"time"
)

type Session struct {
	Token   string
	Client  *gosocketio.Client
	State   State
	Timeout time.Duration
}

type sidResponse struct {
	SID      string
	Upgrades interface{}
}

// New creates a new session struct with provided token
func New(token string, timeout ...int) *Session {
	sess := new(Session)
	sess.Token = token
	if timeout != nil {
		sess.Timeout = time.Duration(timeout[0])
	} else {
		sess.Timeout = 10
	}
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
	auth := make(map[string]string)
	auth["token"] = s.Token
	client.Emit("authentication", auth)
	logged := make(chan bool, 1)
	err = client.On("success", func(channel *gosocketio.Channel, data interface{}) {
		state, _ := json.Marshal(data)
		s.Client = client
		updateState := new(State)
		err := json.Unmarshal(state, updateState)
		s.State = *updateState
		if err != nil {
			fmt.Println(err)
			return
		}
		logged <- true
	})
	select {
	case <-logged:
		return nil
	case <-time.After(s.Timeout * time.Second):
		return errors.New(fmt.Sprint("could not receive an authorization in ", s.Timeout, " seconds"))
	}
}

// Handlers

func (s *Session) OnMessage(handler func(*Session, *MessageCreate)) error {
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

func (s *Session) GetServer(serverID string) (*ServerEvent, error) {
	server := new(ServerEvent)
	end := globals.ReadConstants()
	err := s.Request(server, end.EndpointServer, "/", serverID)
	return server, err
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

func createButtonPayload(buttons []string) *buttonPayload {
	bp := new(buttonPayload)
	for _, button := range buttons {
		bp.add(button, button)
	}
	return bp
}

func (s *Session) ChannelMessageSendWithButtons(channelID string, message string, buttons ...string) error {
	end := globals.ReadConstants()
	dataRaw := createButtonPayload(buttons)
	dataRaw.Message = message
	dataRaw.TempID = "0"
	data, err := json.Marshal(dataRaw)
	if err != nil {
		return err
	}
	_, err = s.Send(data, end.EndpointChannel, "/", channelID)
	return err
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
	b, _ := ioutil.ReadAll(response.Body)
	return response.StatusCode, errors.New(string(b))
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
