package ai

import (
	"math"
	"math/rand"
	"time"

	"github.com/amirhossein-jamali/othello/pkg/model"
)

// Difficulty levels for AI
const (
	Easy   = "easy"
	Medium = "medium"
	Hard   = "hard"
)

// Player represents an AI player
type Player struct {
	Difficulty string
	Piece      model.Piece
}

// NewPlayer creates a new AI player with the specified difficulty
func NewPlayer(difficulty string, piece model.Piece) *Player {
	return &Player{
		Difficulty: difficulty,
		Piece:      piece,
	}
}

// GetMove returns the AI's chosen move
func (p *Player) GetMove(board *model.Board) (int, int, error) {
	switch p.Difficulty {
	case Easy:
		return p.getRandomMove(board)
	case Medium:
		return p.getMediumMove(board)
	case Hard:
		return p.getHardMove(board)
	default:
		return p.getRandomMove(board)
	}
}

// getRandomMove returns a random valid move
func (p *Player) getRandomMove(board *model.Board) (int, int, error) {
	moves := board.GetValidMoves()
	if len(moves) == 0 {
		return -1, -1, nil
	}

	rand.Seed(time.Now().UnixNano())
	move := moves[rand.Intn(len(moves))]
	return move.Row, move.Col, nil
}

// getMediumMove uses a simple heuristic to choose a move
func (p *Player) getMediumMove(board *model.Board) (int, int, error) {
	moves := board.GetValidMoves()
	if len(moves) == 0 {
		return -1, -1, nil
	}

	// Evaluate each move and choose the best one
	bestScore := math.MinInt32
	var bestMove model.Position

	for _, move := range moves {
		boardCopy := board.Clone()
		boardCopy.MakeMove(move.Row, move.Col)
		score := p.evaluatePosition(boardCopy)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove.Row, bestMove.Col, nil
}

// getHardMove uses minimax algorithm with alpha-beta pruning
func (p *Player) getHardMove(board *model.Board) (int, int, error) {
	moves := board.GetValidMoves()
	if len(moves) == 0 {
		return -1, -1, nil
	}

	bestScore := math.MinInt32
	var bestMove model.Position

	for _, move := range moves {
		boardCopy := board.Clone()
		boardCopy.MakeMove(move.Row, move.Col)
		score := p.minimax(boardCopy, 4, math.MinInt32, math.MaxInt32, false)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove.Row, bestMove.Col, nil
}

// minimax implements the minimax algorithm with alpha-beta pruning
func (p *Player) minimax(board *model.Board, depth int, alpha, beta int, maximizing bool) int {
	if depth == 0 || board.IsGameOver() {
		return p.evaluatePosition(board)
	}

	moves := board.GetValidMoves()
	if len(moves) == 0 {
		return p.evaluatePosition(board)
	}

	if maximizing {
		maxScore := math.MinInt32
		for _, move := range moves {
			boardCopy := board.Clone()
			boardCopy.MakeMove(move.Row, move.Col)
			score := p.minimax(boardCopy, depth-1, alpha, beta, false)
			maxScore = max(maxScore, score)
			alpha = max(alpha, score)
			if beta <= alpha {
				break
			}
		}
		return maxScore
	} else {
		minScore := math.MaxInt32
		for _, move := range moves {
			boardCopy := board.Clone()
			boardCopy.MakeMove(move.Row, move.Col)
			score := p.minimax(boardCopy, depth-1, alpha, beta, true)
			minScore = min(minScore, score)
			beta = min(beta, score)
			if beta <= alpha {
				break
			}
		}
		return minScore
	}
}

// evaluatePosition returns a score for the current board position
func (p *Player) evaluatePosition(board *model.Board) int {
	// Simple evaluation: count pieces with weights
	var score int
	weights := [8][8]int{
		{100, -20, 10, 5, 5, 10, -20, 100},
		{-20, -50, -2, -2, -2, -2, -50, -20},
		{10, -2, -1, -1, -1, -1, -2, 10},
		{5, -2, -1, -1, -1, -1, -2, 5},
		{5, -2, -1, -1, -1, -1, -2, 5},
		{10, -2, -1, -1, -1, -1, -2, 10},
		{-20, -50, -2, -2, -2, -2, -50, -20},
		{100, -20, 10, 5, 5, 10, -20, 100},
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := board.GetPiece(i, j)
			if piece == p.Piece {
				score += weights[i][j]
			} else if piece != model.Empty {
				score -= weights[i][j]
			}
		}
	}

	return score
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
