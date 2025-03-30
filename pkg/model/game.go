package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Game represents the core Othello game logic
type Game struct {
	Board     *Board
	History   []Move
	GameOver  bool
	PassCount int // Track consecutive passes
	Winner    Piece
}

// Move represents a player's move
type Move struct {
	Position Position
	Piece    Piece
}

// NewGame creates a new Othello game
func NewGame() *Game {
	return &Game{
		Board:     NewBoard(),
		History:   []Move{},
		GameOver:  false,
		PassCount: 0,
		Winner:    Empty,
	}
}

// MakeMove attempts to place a piece at the given position
// Returns an error if the move is invalid
func (g *Game) MakeMove(row, col int) error {
	if g.GameOver {
		return errors.New("game is already over")
	}

	// Try to make the move
	if !g.Board.MakeMove(row, col) {
		return errors.New("invalid move")
	}

	// Record the move in history
	g.History = append(g.History, Move{
		Position: Position{Row: row, Col: col},
		Piece:    g.Board.CurrentPlayer,
	})

	// Reset pass count since a valid move was made
	g.PassCount = 0

	// Check game state after the move
	g.updateGameState()

	return nil
}

// Pass skips the current player's turn when they have no valid moves
func (g *Game) Pass() error {
	if g.GameOver {
		return errors.New("game is already over")
	}

	if g.Board.HasValidMove() {
		return errors.New("cannot pass when valid moves are available")
	}

	// Record the pass in history
	g.History = append(g.History, Move{
		Position: Position{Row: -1, Col: -1}, // -1,-1 indicates a pass
		Piece:    g.Board.CurrentPlayer,
	})

	// Increment pass count
	g.PassCount++

	// Switch turns
	g.Board.CurrentPlayer = g.Board.getOpponent()

	// Check game state after the pass
	g.updateGameState()

	return nil
}

// updateGameState checks if the game is over
func (g *Game) updateGameState() {
	if g.Board.IsGameOver() || g.PassCount >= 2 {
		g.GameOver = true
		g.Winner = g.Board.GetWinner()
	}
}

// GetValidMoves returns all valid moves for the current player
func (g *Game) GetValidMoves() []Position {
	return g.Board.GetValidMoves()
}

// HasValidMoves checks if the current player has any valid moves
func (g *Game) HasValidMove() bool {
	return g.Board.HasValidMove()
}

// GetScore returns the current score (black count, white count)
func (g *Game) GetScore() (int, int) {
	return g.Board.BlackCnt, g.Board.WhiteCnt
}

// GetGameStatus returns a string describing the current game state
func (g *Game) GetGameStatus() string {
	if g.GameOver {
		switch g.Winner {
		case Black:
			return "Game Over - Black Wins!"
		case White:
			return "Game Over - White Wins!"
		default:
			return "Game Over - It's a tie!"
		}
	}

	if g.Board.CurrentPlayer == Black {
		return "Black's turn"
	}
	return "White's turn"
}

// Reset restarts the game with a new board
func (g *Game) Reset() {
	g.Board = NewBoard()
	g.History = []Move{}
	g.GameOver = false
	g.PassCount = 0
	g.Winner = Empty
}

// FormatMove converts a position to human-readable form (e.g., "E4")
func FormatMove(row, col int) string {
	if row < 0 || col < 0 {
		return "Pass"
	}
	return fmt.Sprintf("%c%d", 'A'+col, row+1)
}

// ParseMove converts a string like "E4" to board coordinates
func ParseMove(move string) (int, int, error) {
	if len(move) < 2 {
		return -1, -1, errors.New("invalid move format")
	}

	// Handle pass
	if move == "pass" || move == "Pass" || move == "PASS" {
		return -1, -1, nil
	}

	col := int(move[0])
	if col >= 'a' && col <= 'h' {
		col -= 'a'
	} else if col >= 'A' && col <= 'H' {
		col -= 'A'
	} else {
		return -1, -1, errors.New("invalid column")
	}

	row := -1
	fmt.Sscanf(move[1:], "%d", &row)
	row-- // Convert from 1-based to 0-based

	if row < 0 || row >= 8 || col < 0 || col >= 8 {
		return -1, -1, errors.New("position out of bounds")
	}

	return row, col, nil
}

// GetPieceName returns a string representation of the piece
func GetPieceName(p Piece) string {
	switch p {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Empty"
	}
}

// MakeComputerMove makes a move for the computer player
func (g *Game) MakeComputerMove() bool {
	moves := g.GetValidMoves()
	if len(moves) == 0 {
		return false
	}

	// Simple AI: Choose a random valid move
	rand.Seed(time.Now().UnixNano())
	move := moves[rand.Intn(len(moves))]
	err := g.MakeMove(move.Row, move.Col)
	return err == nil
}

// GetCurrentPlayer returns the current player
func (g *Game) GetCurrentPlayer() Piece {
	return g.Board.CurrentPlayer
}
