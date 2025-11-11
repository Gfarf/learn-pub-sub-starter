package main

import (
	"fmt"
	"log"
	"os"

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
	gs := gamelogic.NewGameState(username)
	newChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(gs, newChannel),
	)
	if err != nil {
		log.Fatalf("could not subscribe to moves: %v", err)
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		string(routing.WarRecognitionsPrefix),
		string(routing.WarRecognitionsPrefix)+".*",
		pubsub.SimpleQueueDurable,
		handlerWar(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war: %v", err)
	}
	for {
		commands := gamelogic.GetInput()
		switch commands[0] {
		case "spawn":
			gs.CommandSpawn(commands)
		case "move":
			mv, err := gs.CommandMove(commands)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(newChannel, string(routing.ExchangePerilTopic), routing.ArmyMovesPrefix+"."+username, mv)
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				fmt.Printf("Move correctly published: %v\n", mv)
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("invalid command, try again")
		}
	}
	/* wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Closing server from signal...")*/
}
