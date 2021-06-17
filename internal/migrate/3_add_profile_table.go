package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const profileTable = `
CREATE TABLE profiles (
id serial NOT NULL,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
user_id int NOT NULL REFERENCES users(id),
date_of_birth timestamp with time zone NOT NULL,
nid text NOT NULL,
PRIMARY KEY (id)
)`

const bootstrapAccountProfiles = `
INSERT INTO profiles(user_id, date_of_birth, nid) 
VALUES(1, '1982-10-06', '1234567890');
INSERT INTO profiles(user_id, date_of_birth, nid) 
VALUES(2, '1982-10-07', '0987654321');
`

func init() {
	up := []string{
		profileTable,
		bootstrapAccountProfiles,
	}

	down := []string{
		`DROP TABLE profiles`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("create profile table")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("drop profile table")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
