package main

import (
	"embed"
	"log"

	"career_ad/internal/fonts"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

//go:embed assets/fonts/MPLUS1p-Regular.ttf
var fontFS embed.FS

var (
	smallFont font.Face = basicfont.Face7x13
	largeFont font.Face = basicfont.Face7x13
)

func init() {
	b, err := fontFS.ReadFile("assets/fonts/MPLUS1p-Regular.ttf")
	if err != nil {
		log.Printf("font load warning: %v", err)
		return
	}

	smallFont, largeFont = fonts.BuildFaces(b)
}
