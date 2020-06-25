package nertivia

import (
	"encoding/json"
	"fmt"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"io/ioutil"
	"log"
	"nertivia/globals"
	"net/http"
)

type Session struct {
	Token string
	Client *gosocketio.Client
	State *State
}


type sidResponse struct {
	SID string
	Upgrades interface{}
}
// New creates a new session struct with provided token
func New(token string) *Session {
	sess := new(Session)
	sess.Token = token
	return sess
}

func getSID() (string,error) {
	end := globals.ReadConstants()
	res, err := http.Get(end.EndpointURL)
	if err != nil {
		return "",err
	}
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	sidResp := new(sidResponse)
	format := b[4:len(b)-4]
	err = json.Unmarshal(format,sidResp)
	if err != nil {
		fmt.Println(err)
	}
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
	url := gosocketio.GetUrl(end.EndpointURL, 443, true)
	fmt.Println(url)
	client, err := gosocketio.Dial(url,transport.GetDefaultWebsocketTransport())
	if err != nil {
		return err
	}
	authData := getAuth(s.Token)
	err = client.Emit("POST",authData)
	if err != nil {
		return err
	}
	s.Client = client
	return nil
}

func (s *Session) GetUser(userID string) *UserEvent {
	user := new(UserEvent)
	end := globals.ReadConstants()
	err := s.Request(user,end.EndpointUser,"/",userID)
	if err != nil {
		log.Fatal(err)
		return &UserEvent{}
	}
	return user
}

func formatParams(strings []string) string {
	var final string
	for _, str := range strings {
		final = final + str
	}
	return final
}
func (s *Session) Request(gettable Event,endpoint string, params ...string) error {
	url := fmt.Sprint(endpoint,formatParams(params))
	fmt.Println(url)
	request, err := http.NewRequest("GET",url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization",s.Token)
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
	err = json.Unmarshal(body, gettable)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Session) Close() {
	s.Client.Close()
}