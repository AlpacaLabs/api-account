package configuration

import (
	"encoding/json"

	configuration "github.com/AlpacaLabs/go-config"
	"github.com/rs/xid"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	flagForGrpcPort = "grpc_port"
	flagForHTTPPort = "http_port"
)

type Config struct {
	// AppName is a low cardinality identifier for this service.
	AppName string

	// AppID is a unique identifier for the instance (pod) running this app.
	AppID string

	KafkaConfig configuration.KafkaConfig
	SQLConfig   configuration.SQLConfig

	// GrpcPort controls what port our gRPC server runs on.
	GrpcPort int

	// HTTPPort controls what port our HTTP server runs on.
	HTTPPort int
}

func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Could not marshal config to string: %v", err)
	}
	return string(b)
}

func LoadConfig() Config {
	c := Config{
		AppName:  "api-account",
		AppID:    xid.New().String(),
		GrpcPort: 8081,
		HTTPPort: 8083,
	}

	c.KafkaConfig = configuration.LoadKafkaConfig()
	c.SQLConfig = configuration.LoadSQLConfig()

	flag.Int(flagForGrpcPort, c.GrpcPort, "gRPC port")
	flag.Int(flagForHTTPPort, c.HTTPPort, "gRPC HTTP port")

	flag.Parse()

	viper.BindPFlag(flagForGrpcPort, flag.Lookup(flagForGrpcPort))
	viper.BindPFlag(flagForHTTPPort, flag.Lookup(flagForHTTPPort))

	viper.AutomaticEnv()

	c.GrpcPort = viper.GetInt(flagForGrpcPort)
	c.HTTPPort = viper.GetInt(flagForHTTPPort)

	return c
}
