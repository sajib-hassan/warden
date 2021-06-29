package notifier

import "fmt"

type CLISMS struct{}

func NewCLISMS() *CLISMS {
	return &CLISMS{}
}

func (C *CLISMS) initialize() error {
	return nil
}

func (C *CLISMS) deliver(to string, message string) error {
	fmt.Println(to, message)
	return nil
}
