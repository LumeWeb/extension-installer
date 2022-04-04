package main

import (
	"bufio"
	"fmt"
	"github.com/lumeweb/extension-installer/src/platform"
	_ "github.com/mediabuyerbot/go-crx3"
	"os"
)

func main() {
	platform.SetConsoleTitle("Lume Web Installer")
	fmt.Println("Welcome to the Lume Web Extension Installer! A pretty UX is coming soon...")
	fmt.Println("")
	fmt.Println("Press enter twice to install, or type q and press enter to abort.")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if input.Text() == "q" {
		os.Exit(0)
	}

	platform.StartInstall()
}

/*func handleInstall() {
	appWindow.SetContent(
		container.NewWithoutLayout(
			container.New(
				layout.NewVBoxLayout()., container.New(
					layout.NewHBoxLayout(),
					canvas.NewText("Status:", color.Black),
					canvas.NewText("abc", color.Black),
				),
				widget.NewProgressBarInfinite(),
			),
		),
	)
	//go platform.StartInstall(&appWindow)
}
*/
