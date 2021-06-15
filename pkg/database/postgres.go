// Package database implements postgres connection and queries.
package database

import (
	"context"
	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
)

// DBConn returns a postgres connection pool.
func DBConn() (*pg.DB, error) {
	viper.SetDefault("db_network", "tcp")
	viper.SetDefault("db_addr", "localhost:5432")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_database", "customer_db")

	db := pg.Connect(&pg.Options{
		Network:  viper.GetString("db_network"),
		Addr:     viper.GetString("db_addr"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_password"),
		Database: viper.GetString("db_database"),
	})

	if err := checkConn(db); err != nil {
		return nil, err
	}

	if viper.GetBool("db_debug") {
		db.AddQueryHook(pgdebug.NewDebugHook())
	}

	return db, nil
}

func checkConn(db *pg.DB) error {
	return db.Ping(context.Background())
}
