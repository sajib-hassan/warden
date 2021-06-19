package seeder

import (
	"fmt"
	"log"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/internal/db/repos"
	"github.com/sajib-hassan/warden/pkg/auth/encryptor"
)

// UsersSeed seeds user data
func (s Seed) UsersSeed() {
	pin1, _ := encryptor.GenerateFromPassword("54321")
	pin2, _ := encryptor.GenerateFromPassword("65432")
	user1 := &usingpin.User{
		Mobile: "01670209726",
		Pin:    pin1,
		Name:   "Sajib",
		Active: true,
		Roles:  []string{"customer"},
	}

	user2 := &usingpin.User{
		Mobile: "01794342671",
		Pin:    pin2,
		Name:   "Hassan",
		Active: false,
		Roles:  []string{"customer"},
	}

	ur := repos.NewUserStore()
	err := ur.Create(user1)
	if err != nil {
		log.Fatal(err)
	}

	err = ur.Create(user2)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted 2 documents into users collection!\n")
}
