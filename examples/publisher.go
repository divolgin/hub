package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/provider/nats"
	"github.com/mitchellh/go-homedir"
	"github.com/pborman/uuid"
)

var envFilePath string

func init() {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Couldn't ascertain home directory: %s", err.Error())
	}
	flag.StringVar(&envFilePath, "env_path", homeDir+"/.env", "path to the .env file")
	flag.Parse()
	loadEnvironment(envFilePath)
}

func loadEnvironment(path string) {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("Couldn't load the env file (%s) properly: %s", path, err.Error())
	}
}

func main() {
	// create connection
	config := nats.Config{
		User:           os.Getenv("NATS_USER"),
		Password:       os.Getenv("NATS_PASSWORD"),
		Host:           os.Getenv("NATS_HOST"),
		Port:           os.Getenv("NATS_PORT"),
		ServiceName:    os.Getenv("SERVICE_NAME"),
		DefaultTimeout: time.Second * 5,
	}
	conn, err := nats.NewConnection(config.ConnectionUrl(), config)
	if err != nil {
		log.Fatalf("Connecting to NATs resulted in an error: %s", err.Error())
	}

	// publish
	msg := &hub.Message{
		Id:         uuid.New(),
		Topic:      "test",
		Reply:      "test.RES",
		IsResponse: false,
		Payload: hub.Payload{
			Error: "",
			Data:  []byte{},
		},
	}
	if err := conn.Publish(msg); err != nil {
		log.Fatalf("Publishing msg to nats produced error: %s", err.Error())
	}
	time.Sleep(time.Second * 5)
}
