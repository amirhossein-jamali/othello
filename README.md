# Othello Game

A Go implementation of the classic Othello/Reversi board game with both GUI and console interfaces.

## Features

- Graphical user interface with animations and visual feedback
- Console-based text interface for terminal play
- AI opponent with configurable difficulty
- Standard Othello/Reversi rules implementation
- Cross-platform compatibility

## Installation

### Prerequisites

- Go 1.18 or higher
- Dependencies (automatically installed via Go modules):
  - [Ebitengine](https://github.com/hajimehoshi/ebiten) for GUI
  - Other dependencies as listed in `go.mod`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/amirhossein-jamali/othello.git
cd othello

# Build the game
go build -o othello ./cmd/main.go
```

## Usage

Run the game with GUI (default):

```bash
./othello
```

Run the game in console mode:

```bash
./othello -console
```

## Game Rules

Othello (also known as Reversi) is a strategy board game played on an 8×8 grid:

1. The game begins with four discs placed in the center: two black and two white, arranged diagonally
2. Black moves first
3. Players take turns placing discs on the board with their assigned color facing up
4. A valid move must:
   - Be placed adjacent to an opponent's disc
   - Flip at least one opponent's disc
5. Discs are flipped when they are surrounded at the ends of a line (horizontal, vertical, or diagonal) by the discs of the opponent
6. If a player cannot make a valid move, their turn is skipped
7. The game ends when neither player can make a valid move
8. The player with the most discs on the board wins

## Project Structure

```
othello/
├── cmd/
│   └── main.go         # Application entry point
├── pkg/
│   ├── ai/
│   │   └── player.go   # AI opponent implementation
│   ├── model/
│   │   ├── board.go    # Game board model and logic
│   │   └── game.go     # Game state and rules
│   └── ui/
│       ├── console/    # Terminal-based interface
│       └── gui/        # Graphical interface using Ebitengine
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # This file
```

## License

[MIT License](LICENSE)

## Acknowledgments

- [Ebitengine](https://github.com/hajimehoshi/ebiten) for the game engine
- Original Othello game invented by Goro Hasegawa 