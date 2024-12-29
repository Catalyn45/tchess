package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("0 - local game")
	fmt.Println("1 - LAN game")

	fmt.Print("insert: ")
	var input int
	fmt.Scanf("%d", &input)

	var game *Game
	isLan := input == 1
	if isLan {
		fmt.Println("0 - host game")
		fmt.Println("1 - connect to another game")

		fmt.Print("insert: ")
		fmt.Scanf("%d", &input)
		fmt.Scanf("%d", &input)
		isHost := input == 0

		var playerStartsFirst bool
		var connection net.Conn
		if isHost {
			// host server

			playerStartsFirst = true
			listener, err := net.Listen("tcp", ":1234")
			if err != nil {
				fmt.Println("Error starting server:", err)
				return
			}

			defer listener.Close()

			connection, err = listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				return
			}

			defer connection.Close()

		} else {
			// connect to server
			playerStartsFirst = false

			var address string

			fmt.Print("insert address: ")
			fmt.Scanf("%s", &address)
			fmt.Scanf("%s", &address)

			var err error
			connection, err = net.Dial("tcp", address + ":1234")
			if err != nil {
				fmt.Println("Error connecting to server:", err)
				return
			}

			defer connection.Close()
		}

		game = newGame(newHumanPlayer(), newLanPlayer(connection), playerStartsFirst)
	} else {
		game = newGame(newHumanPlayer(), newHumanPlayer(), true)
	}

	game.Start()
}