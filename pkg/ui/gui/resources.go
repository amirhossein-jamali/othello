package gui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// Resources manages the game's visual resources
type Resources struct {
	fonts      *FontResources
	boardImage *ebiten.Image
}

// FontResources contains font faces for different sizes
type FontResources struct {
	smallFont  font.Face
	normalFont font.Face
	largeFont  font.Face
}

// NewResources creates and initializes resources
func NewResources() *Resources {
	res := &Resources{
		fonts: &FontResources{
			smallFont:  basicfont.Face7x13,
			normalFont: basicfont.Face7x13,
			largeFont:  basicfont.Face7x13,
		},
		boardImage: ebiten.NewImage(BoardSize, BoardSize),
	}

	// Initialize the board image
	res.initBoardImage()

	return res
}

// initBoardImage creates the base board image
func (r *Resources) initBoardImage() {
	// Fill with board color
	r.boardImage.Fill(BoardColor)

	// Draw grid lines with increased width
	for i := 0; i <= 8; i++ {
		// Calculate position as float for more precise drawing
		pos := float64(i * CellSize)

		// Draw horizontal grid lines
		for j := 0; j < GridLineWidth; j++ {
			offset := float64(j)
			ebitenutil_DrawLine(r.boardImage, 0, pos+offset, BoardSize, pos+offset, GridColor)
		}

		// Draw vertical grid lines
		for j := 0; j < GridLineWidth; j++ {
			offset := float64(j)
			ebitenutil_DrawLine(r.boardImage, pos+offset, 0, pos+offset, BoardSize, GridColor)
		}
	}

	// Add cell borders for better visualization
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			x := col * CellSize
			y := row * CellSize

			// Draw darker border inside each cell
			drawCellBorder(r.boardImage, x, y, CellSize, CellSize, 1, CellBorderColor)
		}
	}

	// Traditional marker dots removed from the board
}

// drawCellBorder draws a border inside a cell
func drawCellBorder(dst *ebiten.Image, x, y, width, height, borderWidth int, clr color.Color) {
	// Top border
	drawRectangle(dst, image.Rect(x, y, x+width, y+borderWidth), clr)
	// Left border
	drawRectangle(dst, image.Rect(x, y, x+borderWidth, y+height), clr)
	// Right border
	drawRectangle(dst, image.Rect(x+width-borderWidth, y, x+width, y+height), clr)
	// Bottom border
	drawRectangle(dst, image.Rect(x, y+height-borderWidth, x+width, y+height), clr)
}

// drawRectangle draws a filled rectangle
func drawRectangle(dst *ebiten.Image, rect image.Rectangle, clr color.Color) {
	rectImg := ebiten.NewImage(rect.Dx(), rect.Dy())
	rectImg.Fill(clr)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(rectImg, op)
}

// GetBoardImage returns the board image
func (r *Resources) GetBoardImage() *ebiten.Image {
	return r.boardImage
}

// GetSmallFont returns the small font
func (r *Resources) GetSmallFont() font.Face {
	return r.fonts.smallFont
}

// GetNormalFont returns the normal font
func (r *Resources) GetNormalFont() font.Face {
	return r.fonts.normalFont
}

// GetLargeFont returns the large font
func (r *Resources) GetLargeFont() font.Face {
	return r.fonts.largeFont
}

// ebitenutil_DrawLine draws a line (replacement for ebitenutil.DrawLine)
func ebitenutil_DrawLine(dst *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	length := math.Hypot(x2-x1, y2-y1)
	if length == 0 {
		return
	}

	angle := math.Atan2(y2-y1, x2-x1)

	width := 1.0 // Line width
	height := float64(math.Max(1, length))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x1, y1)
	op.ColorM.ScaleWithColor(clr)

	lineImg := ebiten.NewImage(1, 1)
	lineImg.Fill(color.White)
	dst.DrawImage(lineImg, op)
}

// drawCircle fills a circle at the specified position
func drawCircle(dst *ebiten.Image, centerX, centerY, radius int, clr color.Color) {
	diameter := radius * 2
	circleImg := ebiten.NewImage(diameter, diameter)

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - radius)
			dy := float64(y - radius)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= float64(radius) {
				circleImg.Set(x, y, clr)
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(centerX-radius), float64(centerY-radius))
	dst.DrawImage(circleImg, op)
}
