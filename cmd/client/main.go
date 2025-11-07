package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Gfarf/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Gfarf/learn-pub-sub-starter/internal/pubsub"
	"github.com/Gfarf/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connection Successful")
	gamelogic.PrintClientHelp()
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, username), routing.PauseKey, 1)
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Closing server from signal...")
}
