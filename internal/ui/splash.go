package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type SplashData struct {
	Logo         *ebiten.Image
	WindowWidth  int
	WindowHeight int
	Alpha        float64
}

func DrawSplash(screen *ebiten.Image, data SplashData, titleFont, bodyFont font.Face) {
	_ = titleFont
	_ = bodyFont
	alpha := clamp01(data.Alpha)
	if alpha <= 0 || data.Logo == nil {
		return
	}

	ww := float32(data.WindowWidth)
	hh := float32(data.WindowHeight)
	if ww <= 0 || hh <= 0 {
		return
	}

	bg := alphaColor(color.RGBA{R: 8, G: 11, B: 18, A: 245}, alpha)
	vector.DrawFilledRect(screen, 0, 0, ww, hh, bg, false)

	bounds := data.Logo.Bounds()
	logoW := float64(bounds.Dx())
	logoH := float64(bounds.Dy())
	if logoW <= 0 || logoH <= 0 {
		return
	}

	maxW := float64(data.WindowWidth) * 0.55
	maxH := float64(data.WindowHeight) * 0.28
	scale := min(maxW/logoW, maxH/logoH)
	if scale > 1 {
		scale = 1
	}

	drawW := logoW * scale
	drawH := logoH * scale
	offX := (float64(data.WindowWidth) - drawW) / 2
	offY := (float64(data.WindowHeight) - drawH) / 2

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offX, offY)
	op.ColorScale.Scale(1, 1, 1, float32(alpha))
	screen.DrawImage(data.Logo, op)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
