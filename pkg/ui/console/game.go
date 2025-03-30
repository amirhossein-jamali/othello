package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amirhossein-jamali/othello/pkg/ai"
	"github.com/amirhossein-jamali/othello/pkg/model"
)

// ConsoleGame represents the console-based game interface
type ConsoleGame struct {
	game        *model.Game
	reader      *bufio.Reader
	aiPlayer    *ai.Player
	playerColor model.Piece
	gameMode    string
}

// NewConsoleGame creates a new console-based game
func NewConsoleGame() *ConsoleGame {
	return &ConsoleGame{
		game:        model.NewGame(),
		reader:      bufio.NewReader(os.Stdin),
		playerColor: model.Empty, // Will be set during initialization
	}
}

// Run starts the console game loop
func (c *ConsoleGame) Run() {
	fmt.Println("Welcome to Othello!")

	// First select game mode
	c.selectGameMode()

	// Then select player color if playing against AI
	if c.gameMode != "human" {
		c.selectPlayerColor()
	}

	fmt.Println("\nGame started! Enter moves in the format 'A1', 'B2', etc.")
	fmt.Println("Type 'quit' to exit the game.")

	for !c.game.GameOver {
		c.displayBoard()
		c.displayStatus()

		// If it's AI's turn and we're playing against AI
		if c.gameMode != "human" && c.game.Board.CurrentPlayer != c.playerColor {
			fmt.Println("AI is thinking...")
			row, col, _ := c.aiPlayer.GetMove(c.game.Board)

			if row < 0 || col < 0 {
				fmt.Println("AI passes their turn.")
				c.game.Pass()
				continue
			}

			move := model.FormatMove(row, col)
			fmt.Printf("AI places at %s\n", move)
			c.game.MakeMove(row, col)
			continue
		}

		if !c.game.HasValidMove() {
			fmt.Println("No valid moves available. Press Enter to pass...")
			c.reader.ReadString('\n')
			err := c.game.Pass()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}

		move, err := c.getPlayerMove()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if move == "quit" {
			break
		}

		row, col, err := model.ParseMove(move)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		err = c.game.MakeMove(row, col)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	c.displayBoard()
	c.displayGameOver()
}

// selectGameMode lets the player choose the game mode
func (c *ConsoleGame) selectGameMode() {
	for {
		fmt.Println("\nSelect game mode:")
		fmt.Println("1. Human vs Human")
		fmt.Println("2. Human vs Easy AI")
		fmt.Println("3. Human vs Medium AI")
		fmt.Println("4. Human vs Hard AI")
		fmt.Print("Enter choice (1-4): ")

		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			c.gameMode = "human"
			return
		case "2":
			c.gameMode = ai.Easy
			return
		case "3":
			c.gameMode = ai.Medium
			return
		case "4":
			c.gameMode = ai.Hard
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// selectPlayerColor lets the player choose their color
func (c *ConsoleGame) selectPlayerColor() {
	for {
		fmt.Println("\nSelect your color:")
		fmt.Println("B. Black (moves first)")
		fmt.Println("W. White (moves second)")
		fmt.Print("Enter choice (B/W): ")

		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		switch input {
		case "B":
			c.playerColor = model.Black
			c.aiPlayer = ai.NewPlayer(c.gameMode, model.White)
			return
		case "W":
			c.playerColor = model.White
			c.aiPlayer = ai.NewPlayer(c.gameMode, model.Black)
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// displayBoard shows the current state of the board
func (c *ConsoleGame) displayBoard() {
	fmt.Println("\n  A B C D E F G H")
	fmt.Println("  ---------------")
	for i := 0; i < 8; i++ {
		fmt.Printf("%d|", i+1)
		for j := 0; j < 8; j++ {
			piece := c.game.Board.GetPiece(i, j)
			switch piece {
			case model.Black:
				fmt.Print("B ")
			case model.White:
				fmt.Print("W ")
			default:
				if c.isValidMove(i, j) {
					fmt.Print("* ")
				} else {
					fmt.Print(". ")
				}
			}
		}
		fmt.Printf("|%d\n", i+1)
	}
	fmt.Println("  ---------------")
	fmt.Println("  A B C D E F G H")
}

// displayStatus shows the current game status
func (c *ConsoleGame) displayStatus() {
	blackCount, whiteCount := c.game.GetScore()
	fmt.Printf("\nScore: Black: %d, White: %d\n", blackCount, whiteCount)

	currentPlayer := "Black"
	if c.game.Board.CurrentPlayer == model.White {
		currentPlayer = "White"
	}
	fmt.Printf("%s's turn\n", currentPlayer)

	if !c.game.HasValidMove() {
		fmt.Println("No valid moves available!")
	}
}

// displayGameOver shows the final game result
func (c *ConsoleGame) displayGameOver() {
	fmt.Println("\nGame Over!")
	blackCount, whiteCount := c.game.GetScore()
	fmt.Printf("Final Score - Black: %d, White: %d\n", blackCount, whiteCount)

	switch c.game.Winner {
	case model.Black:
		fmt.Println("Black wins!")
	case model.White:
		fmt.Println("White wins!")
	default:
		fmt.Println("It's a tie!")
	}
}

// getPlayerMove reads and validates player input
func (c *ConsoleGame) getPlayerMove() (string, error) {
	fmt.Print("Enter your move: ")
	move, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	move = strings.TrimSpace(move)
	move = strings.ToUpper(move)

	if move == "QUIT" {
		return "quit", nil
	}

	return move, nil
}

// isValidMove checks if a move is valid
func (c *ConsoleGame) isValidMove(row, col int) bool {
	moves := c.game.GetValidMoves()
	for _, move := range moves {
		if move.Row == row && move.Col == col {
			return true
		}
	}
	return false
}
