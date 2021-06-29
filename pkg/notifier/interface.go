package notifier

type SMSChannel interface {
	initialize() error
	deliver(to string, message string) error
}
