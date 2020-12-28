package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"gopkg.in/irc.v3"
)

const grantType = `authorization_code`
const tokenURL = `https://id.twitch.tv/oauth2/token`
const userTokenURL = `https://id.twitch.tv/oauth2/authorize`

var hc = http.DefaultClient

// requests pending from the authorisation call
// var requests = make(map[string]clientToken) // @todo log incoming requests

// mu to protect the requests map from read/writes
// var mu = &sync.RWMutex{}

// New twitch service
func New(ctx context.Context, clientID string, clientSecret string, refreshUserToken string) {
	log.Print("new twitch plugin started")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// todo handle web page for redirect getting this token
	_, _ = authoriseBot(ctx, clientID)

	token, err := getClientToken(ctx)
	if err != nil {
		log.Print(err.Error())
	}

	token, err = refreshToken(ctx, clientID, clientSecret, refreshUserToken)
	if err != nil {
		panic(err.Error())
	}

	// IRC chat
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Token: %s", token.AccessToken)

	config := irc.ClientConfig{
		Nick: strings.ToLower(`TestTommy647Bot`),
		Pass: "oauth:" + token.AccessToken,
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			log.Printf("Message: %s - %s", m.Command, m.String())

			switch m.Command {
			case irc.RPL_WELCOME:
				c.Writef("JOIN #%s", "tommy647uk")
			case "PING":
				msg := &irc.Message{
					Command: "PONG",
					Params: []string{
						m.Params[0],
						m.Trailing(),
					},
				}
				log.Printf("PING: %s", msg.String())
				c.WriteMessage(msg)
			case "PRIVMSG":
				// if its not from a joined channel, ignore
				if !c.FromChannel(m) {
					return
				}
				usr, ok := m.Tags.GetTag("display-name")
				if !ok && m.Prefix.User != "" {
					usr = m.Prefix.User
				}
				c.WriteMessage(&irc.Message{
					Command: "PRIVMSG",
					Params: []string{
						m.Params[0],
						"Thanks for the message @",
						usr,
						":",
						m.Trailing(),
					},
				})
			}
		}),
	}

	client := irc.NewClient(conn, config)
	// request additional irc.v3 data on messages
	client.CapRequest(`twitch.tv/tags`, false)
	client.CapRequest(`twitch.tv/membership`, false)
	client.CapRequest(`twitch.tv/commands`, false)
	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:

				msg := &irc.Message{
					Tags:    nil,
					Prefix:  nil,
					Command: "PRIVMSG",
					Params: []string{
						"#tommy647uk",
						":Test Message",
						"Hello",
					},
				}
				log.Printf("sending a test message: %s", msg.String())
				if err := client.WriteMessage(msg); err != nil {
					log.Println("ERROR: ", err.Error())
				}
			}
		}
	}()
	if err := client.RunContext(ctx); err != nil {
		panic(err.Error())
	}
}

// Handle http requests for the authorisation callback
func Handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world!"))
}

func getToken(ctx context.Context) (*clientToken, error) {

	return &clientToken{}, nil
}

// authoriseBot using twitch oauth - opens browser for the client to log in to twitch, and choose to
// authorise the application to perform actions as them
func authoriseBot(ctx context.Context, clientID string, scopes ...string) (string, error) {
	var requestID = uuid.New()
	// if we are given no scopes, add the defaults
	if len(scopes) == 0 {
		scopes = append(scopes, scopeChatRead, scopeChatEdit)
	}
	var token string
	req, err := http.NewRequest(http.MethodGet, userTokenURL, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add(`client_id`, clientID)
	q.Add(`redirect_uri`, fmt.Sprintf(`http://localhost:8080/twitch?id=%s`, requestID.String())) // @todo update this so we can control the host
	q.Add(`response_type`, `code`)
	q.Add(`scope`, strings.Join(scopes, " "))

	req.URL.RawQuery = q.Encode()

	log.Println(`https://id.twitch.tv` + req.URL.RequestURI())

	if err := browser.OpenURL(fmt.Sprintf("%s%s", `https://id.twitch.tv`, req.URL.RequestURI())); err != nil {
		return "", err
	}
	// @todo: need to wait for the token via HTTP
	return token, nil
}

func refreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*clientToken, error) {
	// get a token
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		panic(err.Error())
	}
	q := req.URL.Query()
	q.Add(`client_id`, clientID)
	q.Add(`client_secret`, clientSecret)
	q.Add(`refresh_token`, refreshToken)
	q.Add(`grant_type`, `refresh_token`)

	req.URL.RawQuery = q.Encode()

	log.Println(req.URL.RequestURI())

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("status: %d", resp.StatusCode)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response: %s", d)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed")
	}
	var token clientToken

	buf := bytes.NewBuffer(d)

	if err = json.NewDecoder(buf).Decode(&token); err != nil {
		return nil, err
	}
	log.Printf("Token: %#v", token)
	log.Printf("Token valid for: %s", token.ExpiresIn)
	return &token, nil
	return nil, nil
}

func getClientToken(ctx context.Context) (*clientToken, error) {
	// get a token
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		panic(err.Error())
	}
	q := req.URL.Query()
	q.Add(`client_id`, clientID)
	q.Add(`client_secret`, clientSecret)
	q.Add(`code`, botUserToken)
	q.Add(`grant_type`, grantType)
	q.Add(`redirect_uri`, `http://localhost`)

	req.URL.RawQuery = q.Encode()

	log.Println(req.URL.RequestURI())

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer resp.Body.Close()
	log.Printf("status: %d", resp.StatusCode)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response: %s", d)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed")
	}
	var token clientToken

	buf := bytes.NewBuffer(d)

	if err = json.NewDecoder(buf).Decode(&token); err != nil {
		return nil, err
	}
	log.Printf("Token: %#v", token)
	log.Printf("Token valid for: %s", token.ExpiresIn)
	return &token, nil
}
