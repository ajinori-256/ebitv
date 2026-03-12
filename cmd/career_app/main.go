package main

import (
	"context"
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	app := parseArgs()
	log.Printf("career_app start")
	log.Printf("アプリケーションを起動します: debug_ui=%v mock_usb_dir=%q", app.DebugUI, app.MockUSBDir)

	g, err := setupApp()
	if err != nil {
		log.Fatal(err)
	}
	g.DebugUI = app.DebugUI
	if app.MockUSBDir != "" {
		g.MockUSBDirs = []string{app.MockUSBDir}
	}

	ebiten.SetWindowTitle("career_app")
	ebiten.SetWindowSize(g.WinW, g.WinH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go g.WatchUSBDevices(ctx)

	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
	cancel()
}
