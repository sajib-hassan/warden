package notifier

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type RobiSMS struct {
	from     string
	endpoint string
	username string
	password string
}

func NewRobiSMS(from string) *RobiSMS {
	return &RobiSMS{from: from}
}

func (c *RobiSMS) initialize() error {
	c.endpoint = viper.GetString("SMS_ROBI_GATEWAY_URL")
	c.username = viper.GetString("SMS_ROBI_GATEWAY_USERNAME")
	c.password = viper.GetString("SMS_ROBI_GATEWAY_PASSWORD")

	return validation.ValidateStruct(c,
		validation.Field(&c.from, validation.Required),
		validation.Field(&c.endpoint, validation.Required),
		validation.Field(&c.username, validation.Required),
		validation.Field(&c.password, validation.Required),
	)
}

func (c *RobiSMS) deliver(to string, message string) error {
	//message = url.QueryEscape(message)

	//Because Robi doesn't support + in the number
	to = strings.TrimLeft(to, "+")
	data := url.Values{
		"Username": {c.username},
		"Password": {c.password},
		"From":     {c.from},
		"To":       {to},
		"Message":  {message},
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	//req.SetBasicAuth("<API Username>", "<API Password>")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	//fmt.Println(string(body))
	return nil
}
