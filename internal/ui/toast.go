package ui

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const toastFadeOutDuration = 450 * time.Millisecond

func DrawToast(screen *ebiten.Image, winW, winH int, title, message, kind string, updating bool, progress float64, changedAt, expiresAt time.Time, stackIndex int, titleFont, bodyFont font.Face) {
	alpha := toastAlpha(expiresAt)
	if alpha <= 0 {
		return
	}

	w := 520.0
	h := 104.0
	if updating {
		h = 124
	}
	x := float64(winW) - w - 22
	y := float64(winH) - h - 22 - float64(stackIndex)*116

	if !(changedAt.IsZero() || updating) {
		age := time.Since(changedAt)
		if age < 220*time.Millisecond {
			p := float64(age) / float64(220*time.Millisecond)
			ease := 1 - (1-p)*(1-p)
			x += (1 - ease) * 26
		}
	}

	shadow := alphaColor(color.RGBA{0, 0, 0, 110}, alpha)
	vector.DrawFilledRect(screen, float32(x+4), float32(y+4), float32(w), float32(h), shadow, false)

	bg := alphaColor(color.RGBA{22, 24, 30, 232}, alpha)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bg, false)

	accent := color.RGBA{91, 155, 255, 255}
	switch kind {
	case "success":
		accent = color.RGBA{47, 201, 126, 255}
	case "error":
		accent = color.RGBA{255, 97, 97, 255}
	case "progress":
		accent = color.RGBA{255, 179, 64, 255}
	}
	accent = alphaColor(accent, alpha)
	vector.DrawFilledRect(screen, float32(x), float32(y), 6, float32(h), accent, false)

	titleColor := alphaColor(color.RGBA{245, 247, 250, 255}, alpha)
	bodyColor := alphaColor(color.RGBA{178, 186, 200, 255}, alpha)
	text.Draw(screen, title, titleFont, int(x+20), int(y+34), titleColor)
	text.Draw(screen, message, bodyFont, int(x+20), int(y+62), bodyColor)

	if updating {
		barX := x + 20
		barY := y + h - 24
		barW := w - 52
		vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(barW), 8, alphaColor(color.RGBA{58, 64, 77, 255}, alpha), false)

		p := progress
		if p <= 0 {
			pulse := 0.2 + 0.2*math.Sin(float64(time.Now().UnixNano())/1.4e8)
			p = pulse
		}
		if p > 1 {
			p = 1
		}
		vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(barW*p), 8, accent, false)

	}
}

func toastAlpha(expiresAt time.Time) float64 {
	if expiresAt.IsZero() {
		return 1
	}
	remaining := time.Until(expiresAt)
	if remaining <= 0 {
		return 0
	}
	if remaining >= toastFadeOutDuration {
		return 1
	}
	return float64(remaining) / float64(toastFadeOutDuration)
}

func alphaColor(c color.RGBA, a float64) color.RGBA {
	if a >= 1 {
		return c
	}
	if a <= 0 {
		return color.RGBA{R: c.R, G: c.G, B: c.B, A: 0}
	}
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: uint8(float64(c.A) * a)}
}
