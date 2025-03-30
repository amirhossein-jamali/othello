package gui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/amirhossein-jamali/othello/pkg/ai"
	"github.com/amirhossein-jamali/othello/pkg/model"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Helper function for converting fixed point to int and calculating width/height
func fixedToIntWidth(r fixed.Rectangle26_6) int {
	return (r.Max.X - r.Min.X).Ceil()
}

func fixedToIntHeight(r fixed.Rectangle26_6) int {
	return (r.Max.Y - r.Min.Y).Ceil()
}

// GameState represents the current UI state
type GameState int

const (
	StateMainMenu GameState = iota
	StateGameMode
	StateColorSelect
	StateInGame
	StateGameOver
)

// GameMode represents the game mode
type GameMode int

const (
	ModeHumanVsHuman GameMode = iota + 1
	ModeHumanVsEasyAI
	ModeHumanVsMediumAI
	ModeHumanVsHardAI
)

// Game represents the main Ebiten game structure
type Game struct {
	gameState GameState
	gameMode  GameMode

	// Core game logic
	othelloGame *model.Game
	aiPlayer    *ai.Player

	// Resources
	resources *Resources

	// Display state
	selectedCellX  int
	selectedCellY  int
	validMoves     []model.Position
	lastMoveX      int
	lastMoveY      int
	computerAction bool
	lastActionTime time.Time

	// Animation
	animating      bool
	animationStart time.Time

	board       *model.Board
	gameOver    bool
	message     string
	colorChosen bool
}

// NewGame creates a new GUI game
func NewGame() *Game {
	return &Game{
		board:       model.NewBoard(),
		gameOver:    false,
		message:     "Choose your color: Press B for Black, W for White",
		colorChosen: false,
		resources:   NewResources(),
		gameState:   StateMainMenu,
	}
}

// Update handles game logic updates each frame
func (g *Game) Update() error {
	// Always check for main menu return key (escape) in any state except main menu
	if g.gameState != StateMainMenu && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameState = StateMainMenu
		return nil
	}

	switch g.gameState {
	case StateMainMenu:
		g.updateMainMenu()
	case StateGameMode:
		g.updateGameMode()
	case StateColorSelect:
		g.updateColorSelect()
	case StateInGame:
		g.updateGame()
	case StateGameOver:
		g.updateGameOver()
	}
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(BackgroundColor)

	switch g.gameState {
	case StateMainMenu:
		g.drawMainMenu(screen)
	case StateGameMode:
		g.drawGameMode(screen)
	case StateColorSelect:
		g.drawColorSelect(screen)
	case StateInGame:
		g.drawGame(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	}
}

// Layout returns the game's logical screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// updateMainMenu handles main menu interactions
func (g *Game) updateMainMenu() {
	// Process mouse clicks only when released to prevent accidental selections
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// Handle start button click
		x, y := ebiten.CursorPosition()
		startButtonRect := image.Rect(ScreenWidth/2-100, ScreenHeight/2-25, ScreenWidth/2+100, ScreenHeight/2+25)

		if startButtonRect.Min.X <= x && x <= startButtonRect.Max.X &&
			startButtonRect.Min.Y <= y && y <= startButtonRect.Max.Y {
			g.gameState = StateGameMode
		}
	}
}

// drawMainMenu renders the main menu
func (g *Game) drawMainMenu(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(BackgroundColor)

	// Title
	titleText := "Othello / Reversi"
	bounds, _ := font.BoundString(g.resources.GetLargeFont(), titleText)
	x := (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y := ScreenHeight / 3
	text.Draw(screen, titleText, g.resources.GetLargeFont(), x, y, TextColor)

	// Get mouse position for hover effect
	mouseX, mouseY := ebiten.CursorPosition()
	startButtonRect := image.Rect(ScreenWidth/2-100, ScreenHeight/2-25, ScreenWidth/2+100, ScreenHeight/2+25)
	buttonHovered := startButtonRect.Min.X <= mouseX && mouseX <= startButtonRect.Max.X &&
		startButtonRect.Min.Y <= mouseY && mouseY <= startButtonRect.Max.Y

	// Start button with hover effect
	buttonColor := ButtonColor
	if buttonHovered {
		buttonColor = HoverColor
	}
	drawRect(screen, startButtonRect, buttonColor)

	buttonText := "Start Game"
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), buttonText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight/2 + fixedToIntHeight(bounds)/2
	text.Draw(screen, buttonText, g.resources.GetNormalFont(), x, y, TextColor)
}

// updateGameMode handles game mode selection
func (g *Game) updateGameMode() {
	// Process mouse clicks only when released to prevent accidental selections
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Calculate button positions and check clicks
		buttonHeight := 50
		buttonSpacing := 20
		buttonY := ScreenHeight / 3

		for mode := ModeHumanVsHuman; mode <= ModeHumanVsHardAI; mode++ {
			buttonRect := image.Rect(ScreenWidth/2-150, buttonY, ScreenWidth/2+150, buttonY+buttonHeight)

			if buttonRect.Min.X <= x && x <= buttonRect.Max.X &&
				buttonRect.Min.Y <= y && y <= buttonRect.Max.Y {
				g.startGame(mode)
				return
			}

			buttonY += buttonHeight + buttonSpacing
		}
	}
}

// drawGameMode renders the game mode selection screen
func (g *Game) drawGameMode(screen *ebiten.Image) {
	// Title
	titleText := "Select Game Mode"
	bounds, _ := font.BoundString(g.resources.GetLargeFont(), titleText)
	x := (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y := ScreenHeight / 5
	text.Draw(screen, titleText, g.resources.GetLargeFont(), x, y, TextColor)

	// Mode buttons
	buttonHeight := 50
	buttonSpacing := 20
	buttonY := ScreenHeight / 3

	modes := []string{
		"Human vs Human",
		"Human vs Computer (Easy)",
		"Human vs Computer (Medium)",
		"Human vs Computer (Hard)",
	}

	// Get cursor position for hover effect
	cx, cy := ebiten.CursorPosition()

	for _, modeText := range modes {
		buttonRect := image.Rect(ScreenWidth/2-150, buttonY, ScreenWidth/2+150, buttonY+buttonHeight)

		// Check if mouse is over this button for hover effect
		hover := buttonRect.Min.X <= cx && cx <= buttonRect.Max.X &&
			buttonRect.Min.Y <= cy && cy <= buttonRect.Max.Y

		// Draw button with appropriate color
		buttonColor := ButtonColor
		if hover {
			buttonColor = HoverColor
		}
		drawRect(screen, buttonRect, buttonColor)

		// Draw text
		bounds, _ = font.BoundString(g.resources.GetNormalFont(), modeText)
		x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
		y = buttonY + buttonHeight/2 + fixedToIntHeight(bounds)/3
		text.Draw(screen, modeText, g.resources.GetNormalFont(), x, y, TextColor)

		buttonY += buttonHeight + buttonSpacing
	}
}

// startGame initializes a new game with the selected mode
func (g *Game) startGame(mode GameMode) {
	g.gameMode = mode

	// If playing against computer, show color selection screen
	if mode != ModeHumanVsHuman {
		g.gameState = StateColorSelect
		return
	}

	// For human vs human, just start the game directly
	g.initializeGame(model.Black) // Player 1 is always black in human vs human
}

// initializeGame sets up the game with the selected color for the human player
func (g *Game) initializeGame(humanColor model.Piece) {
	g.othelloGame = model.NewGame()
	g.gameState = StateInGame
	g.validMoves = g.othelloGame.GetValidMoves()
	g.selectedCellX = -1
	g.selectedCellY = -1
	g.lastMoveX = -1
	g.lastMoveY = -1
	g.computerAction = false
	g.animating = false

	// Create AI if playing against computer
	if g.gameMode != ModeHumanVsHuman {
		var difficulty string
		switch g.gameMode {
		case ModeHumanVsEasyAI:
			difficulty = ai.Easy
		case ModeHumanVsMediumAI:
			difficulty = ai.Medium
		case ModeHumanVsHardAI:
			difficulty = ai.Hard
		}

		aiColor := model.White
		if humanColor == model.White {
			aiColor = model.Black
		}

		g.aiPlayer = ai.NewPlayer(difficulty, aiColor)

		// If AI is black, let it make the first move
		if aiColor == model.Black {
			g.makeComputerMove()
		}
	}
}

// updateColorSelect handles color selection screen interactions
func (g *Game) updateColorSelect() {
	// Process mouse clicks only when released to prevent accidental selections
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		pieceY := ScreenHeight * 3 / 4

		// Check if black piece was clicked
		blackX := ScreenWidth / 3
		if math.Sqrt(math.Pow(float64(x-blackX), 2)+math.Pow(float64(y-pieceY), 2)) <= 40 {
			g.initializeGame(model.Black)
			return
		}

		// Check if white piece was clicked
		whiteX := ScreenWidth * 2 / 3
		if math.Sqrt(math.Pow(float64(x-whiteX), 2)+math.Pow(float64(y-pieceY), 2)) <= 40 {
			g.initializeGame(model.White)
			return
		}
	}
}

// drawColorSelect renders the color selection screen
func (g *Game) drawColorSelect(screen *ebiten.Image) {
	// Fill with background color
	screen.Fill(BackgroundColor)

	// Title
	titleText := "Choose Your Color"
	bounds, _ := font.BoundString(g.resources.GetLargeFont(), titleText)
	x := (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y := ScreenHeight / 4
	text.Draw(screen, titleText, g.resources.GetLargeFont(), x, y, TextColor)

	// Instructions
	instructionText := "Click on a piece to select your color"
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), instructionText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight / 2
	text.Draw(screen, instructionText, g.resources.GetNormalFont(), x, y, TextColor)

	// Draw black and white pieces
	pieceY := ScreenHeight * 3 / 4
	pieceRadius := 40

	// Black piece with label
	blackLabel := "Black"
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), blackLabel)
	blackLabelX := ScreenWidth/3 - fixedToIntWidth(bounds)/2
	text.Draw(screen, blackLabel, g.resources.GetNormalFont(), blackLabelX, pieceY-pieceRadius-10, TextColor)

	// Draw black piece with highlight effect when hovered
	mouseX, mouseY := ebiten.CursorPosition()
	blackX := ScreenWidth / 3
	blackHovered := math.Sqrt(math.Pow(float64(mouseX-blackX), 2)+math.Pow(float64(mouseY-pieceY), 2)) <= float64(pieceRadius)

	// Draw highlight circle if hovered
	if blackHovered {
		// Draw a slightly larger highlight circle behind the piece
		drawCircle(screen, blackX, pieceY, pieceRadius+5, HighlightColor)
	}

	// Always draw the piece
	drawCircle(screen, blackX, pieceY, pieceRadius, BlackPieceColor)

	// White piece with label
	whiteLabel := "White"
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), whiteLabel)
	whiteLabelX := ScreenWidth*2/3 - fixedToIntWidth(bounds)/2
	text.Draw(screen, whiteLabel, g.resources.GetNormalFont(), whiteLabelX, pieceY-pieceRadius-10, TextColor)

	// Draw white piece with highlight effect when hovered
	whiteX := ScreenWidth * 2 / 3
	whiteHovered := math.Sqrt(math.Pow(float64(mouseX-whiteX), 2)+math.Pow(float64(mouseY-pieceY), 2)) <= float64(pieceRadius)

	// Draw highlight circle if hovered
	if whiteHovered {
		// Draw a slightly larger highlight circle behind the piece
		drawCircle(screen, whiteX, pieceY, pieceRadius+5, HighlightColor)
	}

	// Always draw the piece
	drawCircle(screen, whiteX, pieceY, pieceRadius, WhitePieceColor)

	// Draw a hint to press ESC to return to main menu
	escText := "Press ESC to return to main menu"
	bounds, _ = font.BoundString(g.resources.GetSmallFont(), escText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight - 20
	text.Draw(screen, escText, g.resources.GetSmallFont(), x, y, TextColor)
}

// updateGame handles in-game interactions
func (g *Game) updateGame() {
	if g.othelloGame.GameOver {
		g.gameState = StateGameOver
		return
	}

	// Don't process input during animation
	if g.animating {
		// Check if animation is complete
		elapsed := time.Since(g.animationStart).Seconds()
		if elapsed >= AnimationDuration {
			g.animating = false
		} else {
			return
		}
	}

	// Handle computer's turn
	if g.isComputerTurn() {
		// Add a delay before computer's move for better UX
		if !g.computerAction {
			g.computerAction = true
			g.lastActionTime = time.Now()
			return
		}

		// Wait a bit before making the move
		if time.Since(g.lastActionTime).Milliseconds() < 800 {
			return
		}

		g.makeComputerMove()
		g.computerAction = false
		return
	}

	// Handle human player's input
	g.handlePlayerInput()
}

// isComputerTurn checks if it's the computer's turn
func (g *Game) isComputerTurn() bool {
	return g.gameMode != ModeHumanVsHuman &&
		g.othelloGame.Board.CurrentPlayer == g.aiPlayer.Piece
}

// makeComputerMove processes the AI's move
func (g *Game) makeComputerMove() {
	// If AI has no valid moves, pass
	if !g.othelloGame.HasValidMove() {
		g.othelloGame.Pass()
		g.validMoves = g.othelloGame.GetValidMoves()
		return
	}

	// Get AI's move
	row, col, err := g.aiPlayer.GetMove(g.othelloGame.Board)
	if err != nil {
		g.othelloGame.Pass()
		g.validMoves = g.othelloGame.GetValidMoves()
		return
	}

	// Make the move
	err = g.othelloGame.MakeMove(row, col)
	if err == nil {
		g.lastMoveX = col
		g.lastMoveY = row
		g.animating = true
		g.animationStart = time.Now()
	}

	// Update valid moves for next player
	g.validMoves = g.othelloGame.GetValidMoves()
}

// handlePlayerInput processes user input during the game
func (g *Game) handlePlayerInput() {
	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Update selected cell based on mouse position
	g.updateSelectedCell(mouseX, mouseY)

	// Handle board click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.handleBoardClick()
	}
}

// updateSelectedCell updates the currently selected cell based on mouse position
func (g *Game) updateSelectedCell(mouseX, mouseY int) {
	// Check if mouse is over the board
	if mouseX >= BoardMarginX && mouseX < BoardMarginX+BoardSize &&
		mouseY >= BoardMarginY && mouseY < BoardMarginY+BoardSize {
		// Calculate cell coordinates
		col := (mouseX - BoardMarginX) / CellSize
		row := (mouseY - BoardMarginY) / CellSize

		g.selectedCellX = col
		g.selectedCellY = row
	} else {
		g.selectedCellX = -1
		g.selectedCellY = -1
	}
}

// handleBoardClick processes a click on the board
func (g *Game) handleBoardClick() {
	// Ensure a cell is selected
	if g.selectedCellX < 0 || g.selectedCellY < 0 {
		return
	}

	// Check if the move is valid
	isValid := false
	for _, move := range g.validMoves {
		if move.Row == g.selectedCellY && move.Col == g.selectedCellX {
			isValid = true
			break
		}
	}

	if !isValid {
		return
	}

	// Make the move
	err := g.othelloGame.MakeMove(g.selectedCellY, g.selectedCellX)
	if err == nil {
		g.lastMoveX = g.selectedCellX
		g.lastMoveY = g.selectedCellY
		g.animating = true
		g.animationStart = time.Now()
		g.validMoves = g.othelloGame.GetValidMoves()

		// If the next player has no valid moves, pass automatically
		if !g.othelloGame.HasValidMove() && !g.othelloGame.GameOver {
			g.othelloGame.Pass()
			g.validMoves = g.othelloGame.GetValidMoves()
		}
	}
}

// drawGame renders the game board and pieces
func (g *Game) drawGame(screen *ebiten.Image) {
	// Draw the board
	boardImage := g.resources.GetBoardImage()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(BoardMarginX), float64(BoardMarginY))
	screen.DrawImage(boardImage, op)

	// Draw pieces
	g.drawPieces(screen)

	// Draw valid moves
	g.drawValidMoves(screen)

	// Draw selected cell highlight
	g.drawSelectedCell(screen)

	// Draw last move indicator
	g.drawLastMove(screen)

	// Draw history panel
	g.drawHistoryPanel(screen)

	// Draw status bar
	g.drawStatusBar(screen)
}

// drawPieces renders all pieces on the board
func (g *Game) drawPieces(screen *ebiten.Image) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := g.othelloGame.Board.GetPiece(row, col)
			if piece == model.Empty {
				continue
			}

			// Calculate piece position
			centerX := BoardMarginX + col*CellSize + CellSize/2
			centerY := BoardMarginY + row*CellSize + CellSize/2
			radius := (CellSize / 2) - 4

			// Choose color based on piece type
			var pieceColor color.Color
			if piece == model.Black {
				pieceColor = BlackPieceColor
			} else {
				pieceColor = WhitePieceColor
			}

			drawCircle(screen, centerX, centerY, radius, pieceColor)
		}
	}
}

// drawValidMoves highlights valid moves
func (g *Game) drawValidMoves(screen *ebiten.Image) {
	for _, move := range g.validMoves {
		centerX := BoardMarginX + move.Col*CellSize + CellSize/2
		centerY := BoardMarginY + move.Row*CellSize + CellSize/2
		radius := CellSize / 4

		drawCircle(screen, centerX, centerY, radius, ValidMoveColor)
	}
}

// drawSelectedCell highlights the cell under the cursor
func (g *Game) drawSelectedCell(screen *ebiten.Image) {
	if g.selectedCellX >= 0 && g.selectedCellY >= 0 {
		// Check if this cell is a valid move
		isValidMove := false
		for _, move := range g.validMoves {
			if move.Row == g.selectedCellY && move.Col == g.selectedCellX {
				isValidMove = true
				break
			}
		}

		// Only highlight if it's a valid move
		if isValidMove {
			x := BoardMarginX + g.selectedCellX*CellSize
			y := BoardMarginY + g.selectedCellY*CellSize

			rect := image.Rect(x, y, x+CellSize, y+CellSize)
			drawRect(screen, rect, HighlightColor)
		}
	}
}

// drawLastMove highlights the last move played
func (g *Game) drawLastMove(screen *ebiten.Image) {
	if g.lastMoveX >= 0 && g.lastMoveY >= 0 {
		x := BoardMarginX + g.lastMoveX*CellSize
		y := BoardMarginY + g.lastMoveY*CellSize

		// Draw small markers at the corners of the cell
		markerSize := 5
		markerColor := color.RGBA{255, 255, 0, 255} // Yellow

		// Top left
		drawRect(screen, image.Rect(x, y, x+markerSize, y+markerSize), markerColor)
		// Top right
		drawRect(screen, image.Rect(x+CellSize-markerSize, y, x+CellSize, y+markerSize), markerColor)
		// Bottom left
		drawRect(screen, image.Rect(x, y+CellSize-markerSize, x+markerSize, y+CellSize), markerColor)
		// Bottom right
		drawRect(screen, image.Rect(x+CellSize-markerSize, y+CellSize-markerSize, x+CellSize, y+CellSize), markerColor)
	}
}

// drawHistoryPanel renders the game history in a dedicated panel
func (g *Game) drawHistoryPanel(screen *ebiten.Image) {
	// Draw panel background
	panelRect := image.Rect(HistoryPanelX, HistoryPanelY, HistoryPanelX+HistoryPanelW, HistoryPanelY+HistoryPanelH)
	drawRect(screen, panelRect, PanelBackColor)

	// Draw panel border
	borderWidth := 2
	// Top border
	drawRect(screen, image.Rect(panelRect.Min.X, panelRect.Min.Y, panelRect.Max.X, panelRect.Min.Y+borderWidth), PanelBorderColor)
	// Left border
	drawRect(screen, image.Rect(panelRect.Min.X, panelRect.Min.Y, panelRect.Min.X+borderWidth, panelRect.Max.Y), PanelBorderColor)
	// Right border
	drawRect(screen, image.Rect(panelRect.Max.X-borderWidth, panelRect.Min.Y, panelRect.Max.X, panelRect.Max.Y), PanelBorderColor)
	// Bottom border
	drawRect(screen, image.Rect(panelRect.Min.X, panelRect.Max.Y-borderWidth, panelRect.Max.X, panelRect.Max.Y), PanelBorderColor)

	// Draw panel title
	titleText := "Game History"
	bounds, _ := font.BoundString(g.resources.GetNormalFont(), titleText)
	titleX := HistoryPanelX + (HistoryPanelW-fixedToIntWidth(bounds))/2
	titleY := HistoryPanelY + HistoryTitleH/2 + fixedToIntHeight(bounds)/3
	text.Draw(screen, titleText, g.resources.GetNormalFont(), titleX, titleY, TextColor)

	// Draw horizontal separator under title
	drawRect(screen, image.Rect(HistoryPanelX, HistoryPanelY+HistoryTitleH, HistoryPanelX+HistoryPanelW, HistoryPanelY+HistoryTitleH+1), PanelBorderColor)

	// List moves history
	moveCount := len(g.othelloGame.History)
	if moveCount > 0 {
		// Show column headers
		headerY := HistoryPanelY + HistoryTitleH + HistoryItemH
		text.Draw(screen, "Move", g.resources.GetSmallFont(), HistoryPanelX+15, headerY, TextColor)
		text.Draw(screen, "Black", g.resources.GetSmallFont(), HistoryPanelX+80, headerY, BlackMoveColor)
		text.Draw(screen, "White", g.resources.GetSmallFont(), HistoryPanelX+180, headerY, WhiteMoveColor)

		// Draw separator under headers
		drawRect(screen, image.Rect(HistoryPanelX, headerY+5, HistoryPanelX+HistoryPanelW, headerY+6), PanelBorderColor)

		// Calculate number of moves to display (2 moves per row - one black, one white)
		maxMoves := 16 // Maximum number of moves to display
		startMove := 0
		if moveCount > maxMoves {
			startMove = moveCount - maxMoves
		}

		// Display moves in pairs (black and white)
		moveNum := startMove/2 + 1
		rowY := headerY + HistoryItemH

		for i := startMove; i < moveCount; i += 2 {
			// Move number
			moveNumStr := fmt.Sprintf("%d.", moveNum)
			text.Draw(screen, moveNumStr, g.resources.GetSmallFont(), HistoryPanelX+15, rowY, TextColor)

			// Black's move
			if i < moveCount {
				move := g.othelloGame.History[i]
				moveStr := "Pass"
				if move.Position.Row >= 0 {
					moveStr = model.FormatMove(move.Position.Row, move.Position.Col)
				}
				text.Draw(screen, moveStr, g.resources.GetSmallFont(), HistoryPanelX+80, rowY, BlackMoveColor)
			}

			// White's move
			if i+1 < moveCount {
				move := g.othelloGame.History[i+1]
				moveStr := "Pass"
				if move.Position.Row >= 0 {
					moveStr = model.FormatMove(move.Position.Row, move.Position.Col)
				}
				text.Draw(screen, moveStr, g.resources.GetSmallFont(), HistoryPanelX+180, rowY, WhiteMoveColor)
			}

			moveNum++
			rowY += HistoryItemH
		}
	} else {
		// No moves yet
		noMovesText := "No moves yet"
		bounds, _ := font.BoundString(g.resources.GetSmallFont(), noMovesText)
		x := HistoryPanelX + (HistoryPanelW-fixedToIntWidth(bounds))/2
		y := HistoryPanelY + HistoryTitleH + 50
		text.Draw(screen, noMovesText, g.resources.GetSmallFont(), x, y, TextColor)
	}
}

// drawStatusBar renders the score and current player
func (g *Game) drawStatusBar(screen *ebiten.Image) {
	blackCount, whiteCount := g.othelloGame.GetScore()

	// Draw scores
	scoreText := fmt.Sprintf("Black: %d   White: %d", blackCount, whiteCount)
	bounds, _ := font.BoundString(g.resources.GetNormalFont(), scoreText)
	x := (BoardMarginX + BoardSize/2) - fixedToIntWidth(bounds)/2
	y := 40
	text.Draw(screen, scoreText, g.resources.GetNormalFont(), x, y, TextColor)

	// Draw current player indicator
	statusText := g.othelloGame.GetGameStatus()
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), statusText)
	x = (BoardMarginX + BoardSize/2) - fixedToIntWidth(bounds)/2
	y = ScreenHeight - 30
	text.Draw(screen, statusText, g.resources.GetNormalFont(), x, y, TextColor)
}

// updateGameOver handles game over screen interactions
func (g *Game) updateGameOver() {
	// Process mouse clicks only when released to prevent accidental selections
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// Check for menu button click
		x, y := ebiten.CursorPosition()
		menuButtonRect := image.Rect(ScreenWidth/2-100, ScreenHeight-100, ScreenWidth/2+100, ScreenHeight-50)

		if menuButtonRect.Min.X <= x && x <= menuButtonRect.Max.X &&
			menuButtonRect.Min.Y <= y && y <= menuButtonRect.Max.Y {
			g.gameState = StateMainMenu
		}
	}
}

// drawGameOver renders the game over screen
func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Draw the final board state
	g.drawGame(screen)

	// Draw a semi-transparent overlay
	overlayColor := color.RGBA{0, 0, 0, 200}
	drawRect(screen, image.Rect(0, 0, ScreenWidth, ScreenHeight), overlayColor)

	// Draw game over text
	gameOverText := "Game Over"
	bounds, _ := font.BoundString(g.resources.GetLargeFont(), gameOverText)
	x := (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y := ScreenHeight/3 - fixedToIntHeight(bounds)/2
	text.Draw(screen, gameOverText, g.resources.GetLargeFont(), x, y, TextColor)

	// Draw the result
	blackCount, whiteCount := g.othelloGame.GetScore()
	var resultText string
	if blackCount > whiteCount {
		resultText = "Black Wins!"
	} else if whiteCount > blackCount {
		resultText = "White Wins!"
	} else {
		resultText = "It's a Tie!"
	}

	bounds, _ = font.BoundString(g.resources.GetLargeFont(), resultText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight/3 + 50
	text.Draw(screen, resultText, g.resources.GetLargeFont(), x, y, TextColor)

	// Draw the score
	scoreText := fmt.Sprintf("Final Score: Black %d - White %d", blackCount, whiteCount)
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), scoreText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight/3 + 100
	text.Draw(screen, scoreText, g.resources.GetNormalFont(), x, y, TextColor)

	// Get mouse position for hover effect
	mouseX, mouseY := ebiten.CursorPosition()
	menuButtonRect := image.Rect(ScreenWidth/2-100, ScreenHeight-100, ScreenWidth/2+100, ScreenHeight-50)
	buttonHovered := menuButtonRect.Min.X <= mouseX && mouseX <= menuButtonRect.Max.X &&
		menuButtonRect.Min.Y <= mouseY && mouseY <= menuButtonRect.Max.Y

	// Draw menu button with hover effect
	buttonColor := ButtonColor
	if buttonHovered {
		buttonColor = HoverColor
	}
	drawRect(screen, menuButtonRect, buttonColor)

	menuText := "Main Menu"
	bounds, _ = font.BoundString(g.resources.GetNormalFont(), menuText)
	x = (ScreenWidth - fixedToIntWidth(bounds)) / 2
	y = ScreenHeight - 75 + fixedToIntHeight(bounds)/3
	text.Draw(screen, menuText, g.resources.GetNormalFont(), x, y, TextColor)
}

// drawRect draws a filled rectangle
func drawRect(dst *ebiten.Image, rect image.Rectangle, clr color.Color) {
	rectImg := ebiten.NewImage(rect.Dx(), rect.Dy())
	rectImg.Fill(clr)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(rectImg, op)
}

// RunGame starts the GUI game
func RunGame() {
	game := NewGame()

	// Configure the window
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Othello / Reversi")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
