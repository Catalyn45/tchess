package main

import (
	"fmt"
)

type Game struct {
	board *Board
	isPlayerTurn bool
	isEnded bool
	status *[]string
	player Player
	enemy Player
	currentPlayer Player
}

func newGame(player Player, enemy Player, playerStartsFirst bool) *Game {
	status := []string{
		"  Uppercase -> you",
		"  Lowercase -> enemy",
		"  ",
		"  P/p -> Pawn",
		"R/r -> Rook",
		"  N/n -> kNight",
		"B/n -> Bishop",
		"  Q/q -> Queen",
		"K/k -> King",
		"  ",
		"───",
		"  │ │ -> empty spot",
		"───",
		"  ",
		"─═─",
		"  ‖ ‖ -> selected place / possible target",
		"─═─",
		"",
		"  turn:",
		"                     status:",
	}

	return &Game{
		board: newBoard(&status),
		isPlayerTurn: playerStartsFirst,
		isEnded: false,
		status: &status,
		player: player,
		enemy: enemy,
		currentPlayer: nil,
	}
}

func (self *Game) Start() {
	self.setTurn(self.isPlayerTurn)

	for {
		self.board.Draw()
		if self.isEnded {
			break
		}

		line, column, err := self.currentPlayer.GetSelectionCoordinates()
		if err != nil {
			self.setStatus(err.Error())
			continue
		}

		self.board.Reset()

		if !self.coordinatesValid(line, column) {
			self.setStatus("Coordinates are out of range")
			self.board.Reset()
			continue
		}

		if !self.board.selectPiece(line, column, self.isPlayerTurn) {
			self.setStatus("The piece can't be moved in any position")
			self.board.Reset()
			continue
		}

		self.setStatus("Select the target place")
		self.board.Draw()

		destLine, destColumn, err := self.currentPlayer.GetTargetCoordinates()
		if err != nil {
			self.setStatus(err.Error())
			self.board.Reset()
			continue
		}

		if !self.coordinatesValid(destLine, destColumn) {
			self.setStatus("Coordinates are out of range")
			self.board.Reset()
			continue
		}

		err = self.board.movePiece(line, column, destLine, destColumn, self.isPlayerTurn)
		if err != nil {
			self.setStatus(err.Error())
			self.board.Reset()
			continue
		}

		self.board.Reset()

		if self.isPlayerTurn {
			err := self.enemy.NotifySelect(line, column)
			if err != nil {
				self.setStatus(err.Error())
				self.board.Reset()
				continue
			}

			err = self.enemy.NotifyTarget(destLine, destColumn)
			if err != nil {
				self.setStatus(err.Error())
				self.board.Reset()
				continue
			}
		}

		self.setTurn(!self.isPlayerTurn)

	}
}

func (self *Game) setTurn(isPlayer bool) {
	var message string
	var statusMessage string

	if isPlayer {
		self.currentPlayer = self.player
		message = "  turn: you"
		statusMessage = "select a piece"
	} else {
		self.currentPlayer = self.enemy
		message = "  turn: enemy"
		statusMessage = "wait for enemy"
	}

	(*self.status)[18] = message
	self.setStatus(statusMessage)

	self.isPlayerTurn = isPlayer
}

func (self *Game) setStatus(message string) {
	statusTemplate := "                     status: %s"

	(*self.status)[19] = fmt.Sprintf(statusTemplate, message)
}

func (self *Game) coordinatesValid(line int, column int) bool {
	if line < 0 {
		return false
	}

	if line > 7 {
		return false
	}

	if column < 0 {
		return false
	}

	if column > 7 {
		return false
	}

	return true
}
