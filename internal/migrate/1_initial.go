package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const accountTable = `
CREATE TABLE users (
id serial NOT NULL,
mobile text NOT NULL UNIQUE,
pin text NOT NULL,
name text NOT NULL,
active boolean NOT NULL DEFAULT TRUE,
roles text[] NOT NULL DEFAULT '{"customer"}',
last_login timestamp with time zone NOT NULL DEFAULT current_timestamp,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone DEFAULT current_timestamp,
PRIMARY KEY (id)
)`

const tokenTable = `
CREATE TABLE tokens (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
user_id int NOT NULL REFERENCES users(id),
token text NOT NULL UNIQUE,
expiry timestamp with time zone NOT NULL,
mobile boolean NOT NULL DEFAULT FALSE,
identifier text,
PRIMARY KEY (id)
)`

func init() {
	up := []string{
		accountTable,
		tokenTable,
	}

	down := []string{
		`DROP TABLE tokens`,
		`DROP TABLE users`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
