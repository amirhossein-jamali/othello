package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/amirhossein-jamali/othello/pkg/ui/console"
	"github.com/amirhossein-jamali/othello/pkg/ui/gui"
)

func main() {
	// Parse command line flags
	useConsole := flag.Bool("console", false, "Run in console mode")
	flag.Parse()

	if *useConsole {
		fmt.Println("Starting Othello in console mode...")
		game := console.NewConsoleGame()
		game.Run()
	} else {
		fmt.Println("Starting Othello in GUI mode...")
		gui.RunGame()
	}

	os.Exit(0)
}

func showHelp() {
	fmt.Println("Othello / Reversi Game")
	fmt.Println("----------------------")
	fmt.Println("A classic board game where players compete to control the board with their pieces.")
	fmt.Println("\nUsage:")
	fmt.Println("  othello [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -mode=gui     Run in graphical mode (default)")
	fmt.Println("  -mode=console Run in text-based console mode")
	fmt.Println("  -help         Show this help information")
}
