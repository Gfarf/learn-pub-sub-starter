package main

import (
	"fmt"
	"os"

	"github.com/Gfarf/learn-pub-sub-starter/internal/gamelogic"
	"github.com/Gfarf/learn-pub-sub-starter/internal/pubsub"
	"github.com/Gfarf/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connection Successful")
	gamelogic.PrintServerHelp()
	newChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", pubsub.SimpleQueueDurable)
gameFor:
	for {
		commands := gamelogic.GetInput()
		switch commands[0] {
		case "pause":
			fmt.Println("Pausing the game")
			err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "resume":
			fmt.Println("Resuming the game")
			err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "quit":
			fmt.Println("Quiting the game")
			break gameFor
		default:
			fmt.Println("Command not understood")
		}
	}
}
