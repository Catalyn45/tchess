package main

import (
	"fmt"
	"strconv"
)

type Board struct {
	table [8][8] *Piece
	status *[]string
}

func newBoard(status *[]string) *Board {
	table := [8][8]*Piece{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			table[i][j] = newEmptyPiece()
		}
	}

	board := &Board {
		table: table,
		status: status,
	}

	board.setupPieces()

	return board
}

func (self *Board) setupPieces() {
	for i := 0; i < 8; i++ {
		self.table[1][i] = newPawn(false)
		self.table[6][i] = newPawn(true)
	}

	self.table[0][0] = newRook(false)
	self.table[0][7] = newRook(false)

	self.table[7][0] = newRook(true)
	self.table[7][7] = newRook(true)


	self.table[0][1] = newKnight(false)
	self.table[0][6] = newKnight(false)

	self.table[7][1] = newKnight(true)
	self.table[7][6] = newKnight(true)


	self.table[0][2] = newBishop(false)
	self.table[0][5] = newBishop(false)

	self.table[7][2] = newBishop(true)
	self.table[7][5] = newBishop(true)


	self.table[0][3] = newQueen(false)
	self.table[0][4] = newKing(false)

	self.table[7][3] = newQueen(true)
	self.table[7][4] = newKing(true)
}

func (self *Board) Reset() {
	for _, rows := range self.table {
		for _, piece := range rows {
			piece.isTarget = false
			piece.isSelected = false
		}
	}
}

func (self *Board) ShowStatusLine(statusLine *int) {
	if *statusLine >= len(*self.status) {
		return
	}

	fmt.Printf("                   %s", (*self.status)[*statusLine])

	*statusLine++
}

func (self *Board) Draw() {
    fmt.Println("\033[H\033[2J")

	fmt.Print("   ");

	statusLine := 0

	for j := 0; j < 8; j++ {
		fmt.Printf("%c ", 'a' + j);
	}
	self.ShowStatusLine(&statusLine)
	fmt.Println()

	fmt.Print("  ─");

	for j := 0; j < 8; j++ {
		if self.table[0][j].IsSelectedOrTarget() {
			fmt.Print("═─");
		} else {
			fmt.Print("──");
		}
	}

	self.ShowStatusLine(&statusLine)
	fmt.Println();

	for i := 0; i < 8; i++ {
		fmt.Print(strconv.Itoa(i + 1), " ")

		for j := 0; j < 8; j++ {
			piece := self.table[i][j]

			separator := "│"
			if piece.IsSelectedOrTarget() || (j > 0 && self.table[i][j - 1].IsSelectedOrTarget()) {
				separator = "‖"
			}

			fmt.Printf("%s%s", separator, piece.symbol)
		}


		if self.table[i][7].IsSelectedOrTarget() {
			fmt.Print("‖");
		} else {
			fmt.Print("│");
		}

		fmt.Print(" ", strconv.Itoa(i + 1))
		self.ShowStatusLine(&statusLine)
		fmt.Println()

		fmt.Print("  ─")
		for j := 0; j < 8; j++ {
			if self.table[i][j].IsSelectedOrTarget() || (i < 7 && self.table[i+1][j].IsSelectedOrTarget()) {
				fmt.Print("═─")
			} else {
				fmt.Print("──")
			}
		}

		self.ShowStatusLine(&statusLine)
		fmt.Println()
	}

	fmt.Print("   ");
	for j := 0; j < 8; j++ {
		fmt.Printf("%c ", 'a' + j);
	}

	self.ShowStatusLine(&statusLine)
	fmt.Println()

	self.ShowStatusLine(&statusLine)
	fmt.Println()
}

func (self *Board) selectPiece(line int, column int, isPlayer bool) bool {
	piece := self.table[line][column]

	if piece.isPlayer != isPlayer {
		return false
	}

	if piece.kind == PIECE_PAWN {
		if !self.SetPawnTargets(piece, line, column) {
			return false
		}
	} else if piece.kind == PIECE_KNIGHT {
		if !self.SetKnightTargets(piece, line, column) {
			return false
		}
	} else if piece.kind == PIECE_ROOK {
		if !self.SetRookTargets(piece, line, column) {
			return false
		}
	} else if piece.kind == PIECE_BISHOP {
		if !self.SetBishopTargets(piece, line, column) {
			return false
		}
	} else if piece.kind == PIECE_QUEEN {
		rookResult := self.SetRookTargets(piece, line, column)
		bishopResult := self.SetBishopTargets(piece, line, column)

		if !rookResult && !bishopResult {
			return false
		}
	} else if piece.kind == PIECE_KING {
		if !self.SetKingTargets(piece, line, column) {
			return false
		}
	}

	piece.Select()

	return true
}

func (self *Board) movePiece(sourceLine int, sourceColumn, destinationLine, destinationColumn int, isPlayer bool) error {
	destinationPlace := self.table[destinationLine][destinationColumn]
	if !destinationPlace.isTarget {
		return fmt.Errorf("The target position contains a piece of yours")
	}

	sourcePiece := self.table[sourceLine][sourceColumn]

	self.table[destinationLine][destinationColumn] = sourcePiece

	self.table[sourceLine][sourceColumn] = newEmptyPiece()

	if self.IsCheck(isPlayer) {
		// rollback
		self.table[sourceLine][sourceColumn] = sourcePiece
		self.table[destinationLine][destinationColumn] = destinationPlace

		return fmt.Errorf("The move leaves your king in check")
	}

	// Pawn transform in queen
	if (isPlayer && destinationLine == 0) || (!isPlayer && destinationLine == 7) {
		self.table[destinationLine][destinationColumn] = newQueen(isPlayer)
	}

	return nil
}

func (self *Board) SetPawnTargets(piece *Piece, line int, column int) bool {
	if piece.isPlayer && line == 0 {
		return false
	}

	if !piece.isPlayer && line == 7 {
		return false
	}

	direction := 1
	if piece.isPlayer {
		direction = -1
	}

	if self.table[line+direction][column].kind == PIECE_EMPTY {
		self.table[line+direction][column].SetAsTarget()

		if (piece.isPlayer && line == 6) || (!piece.isPlayer && line == 1) {
			if self.table[line+2*direction][column].kind == PIECE_EMPTY {
				self.table[line+2*direction][column].SetAsTarget()
			}
		}
	}

	if column < 7 {
		nextPiece := self.table[line+direction][column+1]
		if nextPiece.isPlayer != piece.isPlayer && nextPiece.kind != PIECE_EMPTY {
			self.table[line+direction][column+1].SetAsTarget()
		}
	}

	if column > 0 {
		nextPiece := self.table[line+direction][column-1]
		if nextPiece.isPlayer != piece.isPlayer && nextPiece.kind != PIECE_EMPTY {
			self.table[line+direction][column-1].SetAsTarget()
		}
	}

	return true
}

func (self *Board) SetKnightTargets(piece *Piece, line int, column int) bool {
	positions := [8][2]int{
		{ 1, -2},
		{ 1,  2},
		{-1, -2},
		{-1,  2},
		{ 2, -1},
		{ 2,  1},
		{-2, -1},
		{-2,  1},
	}

	atLeastOnce := false
	for _, position := range positions {
		lineDelta := position[0]
		columnDelta := position[1]

		newLine := line + lineDelta
		newColumn := column + columnDelta

		if newLine < 0 || newLine > 7 {
			continue
		}

		if newColumn < 0 || newColumn > 7 {
			continue
		}

		destinationPiece := self.table[newLine][newColumn]

		if destinationPiece.kind != PIECE_EMPTY && destinationPiece.isPlayer == piece.isPlayer {
			continue
		}

		destinationPiece.SetAsTarget()
		atLeastOnce = true
	}

	return atLeastOnce
}

func (self *Board) SetRookTargets(piece *Piece, line int, column int) bool {
	atLeastOnce := false
	for delta := -1; delta <= 1; delta += 2 {
		for newLine := line + delta; newLine >= 0 && newLine <= 7; newLine += delta {
			destinationPiece := self.table[newLine][column]

			if destinationPiece.kind != PIECE_EMPTY && destinationPiece.isPlayer == piece.isPlayer {
				break
			}

			destinationPiece.SetAsTarget()
			atLeastOnce = true

			if destinationPiece.kind != PIECE_EMPTY {
				break
			}
		}

		for newColumn := column + delta; newColumn >= 0 && newColumn <= 7; newColumn += delta {
			destinationPiece := self.table[line][newColumn]

			if destinationPiece.kind != PIECE_EMPTY && destinationPiece.isPlayer == piece.isPlayer {
				break
			}

			destinationPiece.SetAsTarget()
			atLeastOnce = true

			if destinationPiece.kind != PIECE_EMPTY {
				break
			}
		}
	}
	
	return atLeastOnce
}

func (self *Board) SetBishopTargets(piece *Piece, line int, column int) bool {
	atLeastOnce := false
	for lineDelta := -1; lineDelta <= 1; lineDelta += 2 {
		for columnDelta := -1; columnDelta <= 1; columnDelta += 2 {
			newLine := line + lineDelta
			newColumn := column + columnDelta

			for newLine >= 0 && newLine <= 7 && newColumn >= 0 && newColumn <= 7 {
				destinationPiece := self.table[newLine][newColumn]

				if destinationPiece.kind != PIECE_EMPTY && destinationPiece.isPlayer == piece.isPlayer {
					break
				}

				destinationPiece.SetAsTarget()
				atLeastOnce = true

				if destinationPiece.kind != PIECE_EMPTY {
					break
				}

				newLine += lineDelta
				newColumn += columnDelta
			}
		}
	}

	return atLeastOnce
}

func (self *Board) SetKingTargets(piece *Piece, line int, column int) bool {
	atLeastOnce := false
	for lineDelta := -1; lineDelta <= 1; lineDelta++ {
		for columnDelta := -1; columnDelta <= 1; columnDelta++ {
			if lineDelta == 0 && columnDelta == 0 {
				continue
			}

			newLine := line + lineDelta
			newColumn := column + columnDelta

			if newLine < 0 || newLine > 7 {
				continue
			}

			if newColumn < 0 || newColumn > 7 {
				continue
			}

			destinationPiece := self.table[newLine][newColumn]

			if destinationPiece.kind != PIECE_EMPTY && destinationPiece.isPlayer == piece.isPlayer {
				continue
			}

			destinationPiece.SetAsTarget()
			atLeastOnce = true
		}
	}

	return atLeastOnce
}

func (self *Board) IsCheck(isPlayer bool) bool {
	var king *Piece
	for line := 0; line < 8; line++ {
		for column := 0; column < 8; column++ {
			piece := self.table[line][column]

			if piece.kind == PIECE_KING && piece.isPlayer == isPlayer {
				king = piece
			} else if piece.kind != PIECE_EMPTY && piece.isPlayer != isPlayer {
				self.selectPiece(line, column, !isPlayer)
			}
		}
	}

	isCheck := king.isTarget

	self.Reset()

	return isCheck
}
