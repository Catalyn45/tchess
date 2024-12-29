package main

const (
	PIECE_EMPTY = 0
	PIECE_PAWN = 1
	PIECE_ROOK = 2
	PIECE_KNIGHT = 3
	PIECE_BISHOP = 4
	PIECE_QUEEN = 5
	PIECE_KING = 6
)

type Piece struct {
	kind int
	symbol string
	isTarget bool
	isSelected bool
	isPlayer bool
}

func (self *Piece) Select() {
	self.isSelected = true
}

func (self *Piece) UnSelect() {
	self.isSelected = false
}

func (self *Piece) SetAsTarget() {
	self.isTarget = true
}

func (self *Piece) UnsetAsTarget() {
	self.isTarget = false
}

func (self *Piece) IsSelectedOrTarget() bool {
	return self.isTarget || self.isSelected
}

func newEmptyPiece() *Piece {
	return &Piece{
		kind: PIECE_EMPTY,
		symbol: " ",
		isTarget: false,
		isSelected: false,
		isPlayer: false,
	}
}

func newPawn(isPlayer bool) *Piece {
	symbol := "p"
	if isPlayer {
		symbol = "P"
	}

	return &Piece{
		kind: PIECE_PAWN,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

func newRook(isPlayer bool) *Piece {
	symbol := "r"
	if isPlayer {
		symbol = "R"
	}

	return &Piece{
		kind: PIECE_ROOK,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

func newKnight(isPlayer bool) *Piece {
	symbol := "n"
	if isPlayer {
		symbol = "N"
	}

	return &Piece{
		kind: PIECE_KNIGHT,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

func newBishop(isPlayer bool) *Piece {
	symbol := "b"
	if isPlayer {
		symbol = "B"
	}

	return &Piece{
		kind: PIECE_BISHOP,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

func newQueen(isPlayer bool) *Piece {
	symbol := "q"
	if isPlayer {
		symbol = "Q"
	}

	return &Piece{
		kind: PIECE_QUEEN,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

func newKing(isPlayer bool) *Piece {
	symbol := "k"
	if isPlayer {
		symbol = "K"
	}

	return &Piece{
		kind: PIECE_KING,
		symbol: symbol,
		isTarget: false,
		isSelected: false,
		isPlayer: isPlayer,
	}
}

