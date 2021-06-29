package notifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type SSLSMS struct {
	endpoint string
	from     string
	username string
	password string
}

func NewSSLSMS(from string) *SSLSMS {
	return &SSLSMS{from: from}
}

func (c *SSLSMS) initialize() error {
	c.endpoint = viper.GetString("SMS_SSL_GATEWAY_URL")
	c.username = viper.GetString("SMS_SSL_GATEWAY_USERNAME")
	c.password = viper.GetString("SMS_SSL_GATEWAY_PASSWORD")

	return validation.ValidateStruct(c,
		validation.Field(&c.from, validation.Required),
		validation.Field(&c.endpoint, validation.Required),
		validation.Field(&c.username, validation.Required),
		validation.Field(&c.password, validation.Required),
	)
}

func (c *SSLSMS) deliver(to string, message string) error {
	message = url.QueryEscape(message)
	cSmsId := "dpl" + strconv.FormatInt(time.Now().UnixNano(), 10)
	data := fmt.Sprintf(
		"user=%s&pass=%s&sms[0][0]=%s&sms[0][1]=%s&sms[0][2]=%s&sid=%s",
		c.username, c.password, to, message, cSmsId, c.from,
	)

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBufferString(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
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
