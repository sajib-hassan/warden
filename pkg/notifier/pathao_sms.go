package notifier

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type PathaoSMS struct {
	from     string
	endpoint string
	username string
	password string
}

func NewPathaoSMS(from string) *PathaoSMS {
	return &PathaoSMS{from: from}
}

func (c *PathaoSMS) initialize() error {
	c.endpoint = viper.GetString("SMS_PATHAO_GATEWAY_URL")
	c.username = viper.GetString("SMS_PATHAO_GATEWAY_USERNAME")
	c.password = viper.GetString("SMS_PATHAO_GATEWAY_PASSWORD")

	return validation.ValidateStruct(c,
		validation.Field(&c.from, validation.Required),
		validation.Field(&c.endpoint, validation.Required),
		validation.Field(&c.username, validation.Required),
		validation.Field(&c.password, validation.Required),
	)
}

func (c *PathaoSMS) deliver(to string, message string) error {
	message = url.QueryEscape(message)
	data := url.Values{
		"user":    {c.username},
		"pass":    {c.password},
		"from":    {c.from},
		"to":      {to},
		"message": {message},
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
