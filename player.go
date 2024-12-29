package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Player interface {
	GetSelectionCoordinates() (line int, column int, err error)
	GetTargetCoordinates() (line int, column int, err error)
	NotifySelect(sourceLine int, sourceColumn int) error
	NotifyTarget(destinationLine int, destinationColumn int) error
}

type HumanPlayer struct {
}

func newHumanPlayer() *HumanPlayer {
	return &HumanPlayer{}
}

func (self *HumanPlayer) getCoordinates(message string) (int, int, error) {
	fmt.Print(message)

	var input string
	fmt.Scan(&input)

	if len(input) != 2 {
		return 0, 0, fmt.Errorf("invalid position, you need to specify line and column (ex. 7a)")
	}

	line, err := strconv.Atoi(string(input[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid position, you need to specify line and column (ex. 7a)")
	}

	if input[1] < 'a' || input[1] > 'z' {
		return 0, 0, fmt.Errorf("invalid position, you need to specify line and column (ex. 7a)")
	}

	column := int(input[1] - 'a')

	return line - 1, column, nil
}

func (self *HumanPlayer) GetSelectionCoordinates() (line int, column int, err error) {
	return self.getCoordinates("Enter position (or -1 to forfeit): ")
}

func (self *HumanPlayer) GetTargetCoordinates() (line int, column int, err error) {
	return self.getCoordinates("Enter position (or 0 to cancel): ")
}

func (self *HumanPlayer) NotifySelect(sourceLine int, sourceColumn int) error { return nil }

func (self *HumanPlayer) NotifyTarget(destinationLine int, destinationColumn int) error { return nil }

type LanPlayer struct {
	connection net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func newLanPlayer(connection net.Conn) *LanPlayer {
	return &LanPlayer{
		connection: connection,
		reader: bufio.NewReader(connection),
		writer: bufio.NewWriter(connection),
	}
}

func (self *LanPlayer) GetSelectionCoordinates() (line int, column int, err error) {
	runeLine, _, err := self.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	runeColumn, _, err := self.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}
	
	return 7 - int(runeLine), 7 - int(runeColumn), nil
}

func (self *LanPlayer) GetTargetCoordinates() (line int, column int, err error) {
	runeLine, _, err := self.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	runeColumn, _, err := self.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	return 7 - int(runeLine), 7 - int(runeColumn), nil
}

func (self *LanPlayer) NotifySelect(sourceLine int, sourceColumn int) error {
	_, err := self.writer.WriteRune(rune(sourceLine))
	if err != nil {
		return err
	}

	_, err = self.writer.WriteRune(rune(sourceColumn))
	if err != nil {
		return err
	}

	self.writer.Flush()

	return nil
}

func (self *LanPlayer) NotifyTarget(destinationLine int, destinationColumn int) error {
	_, err := self.writer.WriteRune(rune(destinationLine))
	if err != nil {
		return err
	}

	_, err = self.writer.WriteRune(rune(destinationColumn))
	if err != nil {
		return err
	}

	self.writer.Flush()

	return nil
}
