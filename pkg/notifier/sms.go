package notifier

import "github.com/spf13/viper"

const (
	SMS_USING_PATHAO  = "pathao"
	SMS_USING_SSL     = "ssl"
	SMS_USING_ROBI    = "robi"
	SMS_USING_INFOBIP = "infobip"
	SMS_USING_CLI     = "cli"
)

type SMSChannel interface {
	initialize() error
	deliver(to string, message string) error
}

type SMSClient struct {
	fromName string
	channel  SMSChannel
}

func NewSMSClient() (*SMSClient, error) {
	fromName := viper.GetString("SMS_FROM_MASKING_NAME")
	channelName := viper.GetString("SMS_DEFAULT_CHANNEL_NAME")
	c := &SMSClient{
		fromName: fromName,
	}
	err := c.SetChannel(channelName)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *SMSClient) SetChannel(channelName string) error {
	if channelName == "" {
		channelName = viper.GetString("SMS_DEFAULT_CHANNEL_NAME")
	}

	switch channelName {
	case SMS_USING_SSL:
		c.channel = NewSSLSMS(c.fromName)
	case SMS_USING_ROBI:
		c.channel = NewRobiSMS(c.fromName)
	case SMS_USING_PATHAO:
		c.channel = NewPathaoSMS(c.fromName)
	case SMS_USING_INFOBIP:
		c.channel = NewInfobipSMS(c.fromName)
	default:
		c.channel = NewCLISMS()
	}

	return c.channel.initialize()
}

func (c *SMSClient) Send(to, message string) error {
	to = formatBDMobile(to)
	return c.channel.deliver(to, message)
}

func formatBDMobile(m string) string {
	return m[len(m)-11:]
}
