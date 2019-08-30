package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // Support remotes in Viper
)

var (
	v *viper.Viper
)

// API returns current API configuration
func API() (ret struct {
	Port int
}) {
	ret.Port = v.GetInt("api.port")
	return
}

// DB returns current database configuration
func DB() (ret struct {
	Host     string
	Port     int
	User     string
	Password string
	Schema   string
}) {
	ret.Host = v.GetString("db.host")
	ret.Port = v.GetInt("db.port")
	ret.User = v.GetString("db.user")
	ret.Password = v.GetString("db.password")
	ret.Schema = v.GetString("db.schema")
	return
}

// Cache returns current cache configuration
func Cache() (ret struct {
	Host     string
	Port     int
	Database int
}) {
	ret.Host = v.GetString("redis.host")
	ret.Port = v.GetInt("redis.port")
	ret.Database = v.GetInt("redis.database")
	return
}

// LoadConfiguration binds to the Consul service and keeps monitoring for changes
func LoadConfiguration() {
	setupLogging()

	v = viper.New()

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	v.AddRemoteProvider("consul", "example.com:8500", "yanbapp-api-"+env)
	v.SetConfigType("json")

	err := v.ReadRemoteConfig()
	if err != nil {
		log.Fatalln("Error reading configuration: " + err.Error())
	}

	log.Printf("Loaded confs: %+v\n", v.AllSettings())
}
