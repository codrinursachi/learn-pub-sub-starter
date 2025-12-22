package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ successfully.")
	channel, err := connection.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
		return
	}

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		fmt.Println("Failed to publish message:", err)
		return
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Sending pause message")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				fmt.Println("Failed to publish pause message:", err)
			} else {
				fmt.Println("Game paused.")
			}
		case "resume":
			fmt.Println("Sending resume message")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				fmt.Println("Failed to publish resume message:", err)
			} else {
				fmt.Println("Game resumed.")
			}
		case "quit":
			fmt.Println("Shutting down Peril server...")
			return
		default:
			fmt.Println("Unknown command:", input[0])
		}
	}

}
