package mailtrap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/donohutcheon/gowebserver/provider/mail"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/joho/godotenv"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password""`
	Host string `json:"domain""`
	SMTPPorts []int `json:"smtp_ports"`
}

type MailTrap struct {
	mail.Client
	state  *state.ServerState
	config Config
}

const mailTrapURL = "https://mailtrap.io/api/v1/inboxes.json?api_token="

func New(state *state.ServerState) *MailTrap {
	m := new(MailTrap)
	m.state = state

	err := godotenv.Load()
	if err != nil {
		// TODO: Use proper logger
		fmt.Printf("Could not load environment files. %s", err.Error())
	}

	var apiToken = os.Getenv("MAILTRAP_API_TOKEN")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, mailTrapURL+apiToken, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var configs []Config
	err = json.Unmarshal(bytes, &configs)
	if err != nil {
		panic(err)
	}

	m.config = configs[0]
	config := &m.config

	if len(config.Username) == 0 {
		panic("mailtrap username is empty")
	}
	if len(config.Password) == 0 {
		panic("mailtrap username is empty")
	}
	if len(config.Host) == 0 {
		panic("mailtrap host is empty")
	}
	if len(config.SMTPPorts) == 0 {
		panic("mailtrap ports are empty")
	}
	fmt.Printf("%v", m)

	return m
}

func (m *MailTrap) SendMail(to []string, from, subject, message string) error {
	config := &m.config

	toList := strings.Join(to, ",")
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, toList, subject, message)
	auth := smtp.CRAMMD5Auth(config.Username, config.Password)
	addr := config.Host + ":" + strconv.Itoa(config.SMTPPorts[3])

	err := smtp.SendMail(addr, auth, from, to, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
