package notifier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type InfobipSMS struct {
	from     string
	endpoint string
	apiKey   string
}

func NewInfobipSMS(from string) *InfobipSMS {
	return &InfobipSMS{from: from}
}

func (c *InfobipSMS) initialize() error {
	c.endpoint = viper.GetString("SMS_INFOBIP_GATEWAY_URL")
	c.apiKey = viper.GetString("SMS_INFOBIP_GATEWAY_API_KEY")

	return validation.ValidateStruct(c,
		validation.Field(&c.from, validation.Required),
		validation.Field(&c.endpoint, validation.Required),
		validation.Field(&c.apiKey, validation.Required),
	)
}

func (c *InfobipSMS) deliver(to string, message string) error {
	//message = url.QueryEscape(message)
	data := map[string]string{
		"to":      to,
		"message": message,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", c.apiKey)

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
