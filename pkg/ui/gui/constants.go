package gui

import "image/color"

// Screen dimensions
const (
	ScreenWidth   = 950 // افزایش عرض برای تاریخچه بازی
	ScreenHeight  = 650 // افزایش ارتفاع برای فضای بیشتر
	BoardSize     = 480 // Board is square
	CellSize      = BoardSize / 8
	BoardMarginX  = 60                            // حاشیه سمت چپ تخته بازی
	BoardMarginY  = 80                            // حاشیه بالای تخته بازی
	GridLineWidth = 2                             // Increased line width for better visibility
	HistoryPanelX = BoardMarginX + BoardSize + 30 // مکان پنل تاریخچه
	HistoryPanelY = BoardMarginY
	HistoryPanelW = 320       // عرض پنل تاریخچه
	HistoryPanelH = BoardSize // ارتفاع پنل تاریخچه
	HistoryTitleH = 30        // ارتفاع عنوان پنل تاریخچه
	HistoryItemH  = 25        // ارتفاع هر آیتم در تاریخچه
)

// Animation constants
const (
	AnimationDuration = 0.3 // seconds
)

// Colors
var (
	BackgroundColor  = color.RGBA{40, 40, 40, 255}
	BoardColor       = color.RGBA{34, 139, 34, 255} // Forest Green
	GridColor        = color.RGBA{0, 70, 0, 255}    // Darker Green for better contrast
	BlackPieceColor  = color.RGBA{0, 0, 0, 255}
	WhitePieceColor  = color.RGBA{240, 240, 240, 255}
	ValidMoveColor   = color.RGBA{50, 255, 50, 230} // Much brighter green with higher opacity
	HighlightColor   = color.RGBA{220, 220, 0, 150} // Brighter yellow highlight
	ButtonColor      = color.RGBA{0, 100, 0, 255}
	HoverColor       = color.RGBA{0, 120, 0, 255}
	TextColor        = color.White
	CellBorderColor  = color.RGBA{0, 50, 0, 255}      // Color for cell borders
	PanelBackColor   = color.RGBA{30, 70, 30, 230}    // Background color for panels
	PanelBorderColor = color.RGBA{60, 180, 60, 255}   // Border color for panels
	BlackMoveColor   = color.RGBA{180, 180, 180, 255} // تاریخچه حرکت های سیاه
	WhiteMoveColor   = color.RGBA{255, 255, 255, 255} // تاریخچه حرکت های سفید
	MarkerDotColor   = color.RGBA{20, 20, 20, 200}    // رنگ نقاط راهنما روی تخته
)
