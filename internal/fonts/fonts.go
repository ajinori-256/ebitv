package fonts

import (
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

func BuildFaces(ttf []byte) (font.Face, font.Face) {
	defaultFace := basicfont.Face7x13
	if len(ttf) == 0 {
		return defaultFace, defaultFace
	}

	tt, err := opentype.Parse(ttf)
	if err != nil {
		log.Printf("font parse warning: %v", err)
		return defaultFace, defaultFace
	}

	face14, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: 14, DPI: 72, Hinting: font.HintingVertical})
	if err != nil {
		log.Printf("font face warning: %v", err)
		return defaultFace, defaultFace
	}
	face16, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: 16, DPI: 72, Hinting: font.HintingVertical})
	if err != nil {
		log.Printf("font face warning: %v", err)
		return defaultFace, defaultFace
	}

	return face14, face16
}
