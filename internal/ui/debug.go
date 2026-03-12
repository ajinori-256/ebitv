package ui

import (
	"image/color"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type DebugPanelData struct {
	AnimStyle          int
	Updating           bool
	Progress           float64
	IsWindowFocused    bool
	WindowWidth        int
	WindowHeight       int
	StatusKind         string
	StatusText         string
	CurrentSlidePath   string
	HasPreviousSlide   bool
	ToastCount         int
	ChangeInterval     time.Duration
	TransitionDuration time.Duration
}

func DrawDebugPanel(screen *ebiten.Image, data DebugPanelData, titleFont, bodyFont font.Face) {
	panelX := 20.0
	panelY := 20.0
	panelW := 520.0
	panelH := 238.0
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelW), float32(panelH), color.RGBA{18, 21, 28, 220}, false)
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), 4, float32(panelH), color.RGBA{91, 155, 255, 255}, false)

	text.Draw(screen, "Debug Overlay", titleFont, int(panelX+16), int(panelY+30), color.RGBA{245, 247, 250, 255})
	text.Draw(screen, "Dキー: ON-OFF切替", bodyFont, int(panelX+16), int(panelY+54), color.RGBA{190, 198, 210, 255})

	animName := []string{"Fade", "Slide", "Zoom"}[data.AnimStyle%3]
	state := "OFF"
	if data.Updating {
		state = "ON"
	}

	currentSlide := "(none)"
	if data.CurrentSlidePath != "" {
		currentSlide = filepath.Base(data.CurrentSlidePath)
	}

	text.Draw(screen, "Animation: "+animName+" / Progress: "+state+" "+formatPercent(data.Progress), bodyFont, int(panelX+16), int(panelY+84), color.RGBA{225, 232, 242, 255})
	text.Draw(screen, "Window: "+strconv.Itoa(data.WindowWidth)+"x"+strconv.Itoa(data.WindowHeight)+" / Focus: "+strconv.FormatBool(data.IsWindowFocused), bodyFont, int(panelX+16), int(panelY+108), color.RGBA{225, 232, 242, 255})
	text.Draw(screen, "Slide: "+currentSlide+" / HasPrev: "+strconv.FormatBool(data.HasPreviousSlide), bodyFont, int(panelX+16), int(panelY+132), color.RGBA{225, 232, 242, 255})
	text.Draw(screen, "Status: ["+data.StatusKind+"] "+trimText(data.StatusText, 42), bodyFont, int(panelX+16), int(panelY+156), color.RGBA{225, 232, 242, 255})
	text.Draw(screen, "Toasts: "+strconv.Itoa(data.ToastCount), bodyFont, int(panelX+16), int(panelY+180), color.RGBA{225, 232, 242, 255})
	text.Draw(screen, "Interval: "+data.ChangeInterval.String()+" / Transition: "+data.TransitionDuration.String(), bodyFont, int(panelX+16), int(panelY+204), color.RGBA{225, 232, 242, 255})
}

func formatPercent(v float64) string {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	pct := int(math.Round(v * 100))
	return strconv.Itoa(pct) + "%"
}

func trimText(s string, max int) string {
	if max <= 3 || len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
