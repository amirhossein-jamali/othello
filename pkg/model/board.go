package model

// Direction represents a direction for searching on the board
type Direction struct {
	DRow, DCol int
}

// All eight possible directions to check
var Directions = []Direction{
	{-1, -1}, {-1, 0}, {-1, 1}, // Up-left, Up, Up-right
	{0, -1}, {0, 1}, // Left, Right
	{1, -1}, {1, 0}, {1, 1}, // Down-left, Down, Down-right
}

// Piece represents a game piece (disc) on the board
type Piece int8

const (
	Empty Piece = iota
	Black
	White
)

// Position represents a board position with row and column
type Position struct {
	Row, Col int
}

// Board represents the Othello game board
type Board struct {
	Cells         [][]Piece
	CurrentPlayer Piece
	Size          int
	BlackCnt      int
	WhiteCnt      int
}

// NewBoard creates a new Othello board with the initial setup
func NewBoard() *Board {
	size := 8
	cells := make([][]Piece, size)
	for i := range cells {
		cells[i] = make([]Piece, size)
	}

	// Initialize center pieces
	cells[3][3] = White
	cells[3][4] = Black
	cells[4][3] = Black
	cells[4][4] = White

	return &Board{
		Cells:         cells,
		CurrentPlayer: Black,
		Size:          size,
		BlackCnt:      2,
		WhiteCnt:      2,
	}
}

// IsValidPosition checks if a position is on the board
func (b *Board) IsValidPosition(row, col int) bool {
	return row >= 0 && row < b.Size && col >= 0 && col < b.Size
}

// GetPiece returns the piece at the given position
func (b *Board) GetPiece(row, col int) Piece {
	if !b.IsValidPosition(row, col) {
		return Empty
	}
	return b.Cells[row][col]
}

// IsValidMove checks if placing a piece at the given position is valid
func (b *Board) IsValidMove(row, col int) bool {
	if !b.IsValidPosition(row, col) || b.Cells[row][col] != Empty {
		return false
	}

	opponent := b.getOpponent()

	// Check in all eight directions
	for _, dir := range Directions {
		r, c := row+dir.DRow, col+dir.DCol

		// Must have at least one opponent piece adjacent
		if !b.IsValidPosition(r, c) || b.Cells[r][c] != opponent {
			continue
		}

		r += dir.DRow
		c += dir.DCol

		// Keep going in this direction
		for b.IsValidPosition(r, c) {
			// If we find our own piece, this is a valid move
			if b.Cells[r][c] == b.CurrentPlayer {
				return true
			}

			// If we find an empty cell, this direction is invalid
			if b.Cells[r][c] == Empty {
				break
			}

			r += dir.DRow
			c += dir.DCol
		}
	}

	return false
}

// GetValidMoves returns all valid moves for the current player
func (b *Board) GetValidMoves() []Position {
	var moves []Position

	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.IsValidMove(i, j) {
				moves = append(moves, Position{Row: i, Col: j})
			}
		}
	}

	return moves
}

// MakeMove applies a move to the board and updates the current player
func (b *Board) MakeMove(row, col int) bool {
	if !b.IsValidMove(row, col) {
		return false
	}

	b.Cells[row][col] = b.CurrentPlayer
	b.flipPieces(row, col)
	b.CurrentPlayer = b.getOpponent()
	return true
}

// getOpponent returns the opposite player
func (b *Board) getOpponent() Piece {
	if b.CurrentPlayer == Black {
		return White
	}
	return Black
}

// flipPieces flips the appropriate pieces on the board
func (b *Board) flipPieces(row, col int) {
	opponent := b.getOpponent()
	flipped := 0

	// Check all eight directions and flip pieces
	for _, dir := range Directions {
		r, c := row+dir.DRow, col+dir.DCol
		var toFlip []Position

		if !b.IsValidPosition(r, c) || b.Cells[r][c] != opponent {
			continue
		}

		toFlip = append(toFlip, Position{Row: r, Col: c})
		r += dir.DRow
		c += dir.DCol

		foundOwn := false
		for b.IsValidPosition(r, c) {
			if b.Cells[r][c] == Empty {
				break
			}

			if b.Cells[r][c] == b.CurrentPlayer {
				foundOwn = true
				break
			}

			toFlip = append(toFlip, Position{Row: r, Col: c})
			r += dir.DRow
			c += dir.DCol
		}

		if foundOwn {
			for _, pos := range toFlip {
				b.Cells[pos.Row][pos.Col] = b.CurrentPlayer
				flipped++
			}
		}
	}

	// Update piece counts
	if b.CurrentPlayer == Black {
		b.BlackCnt += flipped + 1
		b.WhiteCnt -= flipped
	} else {
		b.WhiteCnt += flipped + 1
		b.BlackCnt -= flipped
	}
}

// HasValidMove checks if the current player has any valid moves
func (b *Board) HasValidMove() bool {
	return len(b.GetValidMoves()) > 0
}

// IsGameOver checks if the game is over (no valid moves for either player)
func (b *Board) IsGameOver() bool {
	// Save current player
	current := b.CurrentPlayer

	// Check if current player has moves
	hasMove := b.HasValidMove()
	if hasMove {
		return false
	}

	// Switch to opponent and check their moves
	b.CurrentPlayer = b.getOpponent()
	opponentHasMove := b.HasValidMove()

	// Restore current player
	b.CurrentPlayer = current

	return !opponentHasMove
}

// GetWinner returns the winner (or Empty if tie)
func (b *Board) GetWinner() Piece {
	if b.BlackCnt > b.WhiteCnt {
		return Black
	}
	if b.WhiteCnt > b.BlackCnt {
		return White
	}
	return Empty // Tie
}

// Clone creates a deep copy of the board
func (b *Board) Clone() *Board {
	newBoard := &Board{
		CurrentPlayer: b.CurrentPlayer,
		Size:          b.Size,
		BlackCnt:      b.BlackCnt,
		WhiteCnt:      b.WhiteCnt,
	}

	newBoard.Cells = make([][]Piece, b.Size)
	for i := 0; i < b.Size; i++ {
		newBoard.Cells[i] = make([]Piece, b.Size)
		for j := 0; j < b.Size; j++ {
			newBoard.Cells[i][j] = b.Cells[i][j]
		}
	}

	return newBoard
}
