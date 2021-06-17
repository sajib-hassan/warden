package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"

	"github.com/sajib-hassan/warden/pkg/auth/encryptor"
)

const bootstrapUserAccount1 = `
INSERT INTO users (id, mobile, pin, name, active, roles)
VALUES (DEFAULT, '01670209726', $L$?$L$ , 'Sajib', true, '{customer}');
`
const bootstrapUserAccount2 = `
INSERT INTO users (id, mobile, pin, name, active, roles)
VALUES (DEFAULT, '017194342671', $L$?$L$ , 'Hassan', true, '{customer}')
`

func init() {
	up := []string{
		bootstrapUserAccount1,
		bootstrapUserAccount2,
	}

	down := []string{
		`TRUNCATE users CASCADE`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("add bootstrap accounts")
		pin, _ := encryptor.GenerateFromPassword("54321")
		fmt.Println(pin)
		for _, q := range up {
			_, err := db.Exec(q, pin)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("truncate users cascading")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
