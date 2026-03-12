package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/logo.png
var logoPNG []byte

var splashLogo *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(logoPNG))
	if err != nil {
		log.Printf("logo load warning: %v", err)
		return
	}
	splashLogo = ebiten.NewImageFromImage(img)
}
