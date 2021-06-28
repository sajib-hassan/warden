package dbconn

import (
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect returns a postgres connection pool.
func Connect() error {
	//viper.SetDefault("mongo_uri", "mongodb://test:test@localhost:27017")
	//viper.SetDefault("mongo_db_name", "test_db")

	//logrus.Info(viper.GetString("mongo_uri"))
	//logrus.Info(viper.GetString("mongo_db_name"))
	//logrus.Info(viper.AllKeys())

	return mgm.SetDefaultConfig(
		nil,
		viper.GetString("mongo_db_name"),
		options.Client().ApplyURI(viper.GetString("mongo_uri")),
	)
}
