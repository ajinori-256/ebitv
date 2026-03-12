package app

import (
	"image"
	"path/filepath"
	"sync"
	"time"

	"career_ad/internal/config"
	"career_ad/internal/media"
	"career_ad/internal/ui"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ChangeInterval      = 5 * time.Second
	TransitionDuration  = 700 * time.Millisecond
	StatusDisplayPeriod = 3 * time.Second
	DefaultWidth        = 720
	DefaultHeight       = 1280
	SplashDuration      = 2200 * time.Millisecond
	SplashFadeDuration  = 550 * time.Millisecond
)

type Slide = media.Slide

type ToastItem struct {
	Title     string
	Message   string
	Kind      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type App struct {
	Mu sync.RWMutex

	Source      media.SlideSource
	SplashLogo  *ebiten.Image
	StartedAt   time.Time
	LastChanged time.Time
	AnimStyle   int
	Transition  time.Time

	WinW int
	WinH int

	DataDir    string
	ConfigPath string
	MockUSBDirs []string

	Updating      bool
	StatusText    string
	StatusKind    string
	CopyProgress  float64
	StatusChanged time.Time
	Toasts        []ToastItem

	SmallFont font.Face
	LargeFont font.Face

	Cfg               config.Config
	DefaultConfigData []byte
	DebugUI           bool
}

func NewApp(dataDir, configPath string, defaultCfgData []byte, smallFont, largeFont font.Face, splashLogo *ebiten.Image) (*App, error) {
	
	cfg, _ := config.LoadFile(configPath, defaultCfgData)
	paths, err := media.CollectImagePaths(dataDir)
	if err != nil {
		return nil, err
	}

	return &App{
		Source:        media.NewDefaultSlideSource(paths),
		SplashLogo:    splashLogo,
		StartedAt:     time.Now(),
		LastChanged:   time.Now(),
		WinW:          DefaultWidth,
		WinH:          DefaultHeight,
		DataDir:       dataDir,
		ConfigPath:    configPath,
		StatusText:    "USBデバイス待機中",
		StatusKind:    "idle",
		StatusChanged: time.Now(),
		Toasts:            make([]ToastItem, 0, 8),
		SmallFont:         smallFont,
		LargeFont:         largeFont,
		Cfg:               cfg,
		DefaultConfigData: defaultCfgData,
	}, nil
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		a.Mu.Lock()
		a.DebugUI = !a.DebugUI
		a.Mu.Unlock()
	}

	now := time.Now()

	a.Mu.Lock()
	if now.Sub(a.LastChanged) >= a.Cfg.ChangeInterval && ebiten.IsFocused() {
		a.Source.Advance()
		a.LastChanged = now
		a.Transition = now
		a.AnimStyle = (a.AnimStyle + 1) % 3
	}
	if len(a.Toasts) > 0 {
		filtered := a.Toasts[:0]
		for _, t := range a.Toasts {
			if t.ExpiresAt.IsZero() || now.Before(t.ExpiresAt) {
				filtered = append(filtered, t)
			}
		}
		a.Toasts = filtered
	}
	a.Mu.Unlock()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(image.Black)

	a.Mu.RLock()
	animStyle := a.AnimStyle
	transition := a.Transition
	status := a.StatusText
	statusKind := a.StatusKind
	updating := a.Updating
	copyProgress := a.CopyProgress
	statusChanged := a.StatusChanged
	toasts := append([]ToastItem(nil), a.Toasts...)
	winW := a.WinW
	winH := a.WinH
	startedAt := a.StartedAt
	a.Mu.RUnlock()

	currentSlide := a.Source.CurrentSlide()
	previousSlide, hasPrev := a.Source.PreviousSlide()

	title := "ステータス"
	if currentSlide.Path != "" {
		title = filepath.Base(currentSlide.Path)
	}

	stack := make([]ToastItem, 0, 6)
	if updating {
		stack = append(stack, ToastItem{
			Title:     title,
			Message:   status,
			Kind:      statusKind,
			CreatedAt: statusChanged,
		})
	}
	for i := len(toasts) - 1; i >= 0; i-- {
		stack = append(stack, toasts[i])
		if len(stack) >= 4 {
			break
		}
	}

	if a.Source.CurrentSlide().Img == nil {
		a.DrawToast(screen, winW, winH, "画像なし", "USB接続を待機中", "info", false, 0, statusChanged, time.Time{}, len(stack))
		for i, t := range stack {
			a.DrawToast(screen, winW, winH, t.Title, t.Message, t.Kind, updating && i == 0, copyProgress, t.CreatedAt, t.ExpiresAt, i)
		}
		return
	}

	now := time.Now()
	if hasPrev && !transition.IsZero() {
		progress := now.Sub(transition)
		if progress < a.Cfg.TransitionDuration {
			p := float64(progress) / float64(a.Cfg.TransitionDuration)
			switch animStyle {
			case 0:
				a.DrawSlide(screen, previousSlide, winW, winH, 0, 1.0, 1.0-p)
				a.DrawSlide(screen, currentSlide, winW, winH, 0, 1.0, p)
			case 1:
				a.DrawSlide(screen, previousSlide, winW, winH, -float64(winW)*p, 1.0, 1.0)
				a.DrawSlide(screen, currentSlide, winW, winH, float64(winW)*(1.0-p), 1.0, 1.0)
			default:
				a.DrawSlide(screen, previousSlide, winW, winH, 0, 1.0+0.08*p, 1.0-p)
				a.DrawSlide(screen, currentSlide, winW, winH, 0, 0.92+0.08*p, p)
			}
		} else {
			a.DrawSlide(screen, currentSlide, winW, winH, 0, 1.0, 1.0)
		}
	} else {
		a.DrawSlide(screen, currentSlide, winW, winH, 0, 1.0, 1.0)
	}

	for i, t := range stack {
		a.DrawToast(screen, winW, winH, t.Title, t.Message, t.Kind, updating && i == 0, copyProgress, t.CreatedAt, t.ExpiresAt, i)
	}

	if a.DebugUI {
		a.DrawDebugPanel(screen, ui.DebugPanelData{
			AnimStyle:          animStyle,
			Updating:           updating,
			Progress:           copyProgress,
			IsWindowFocused:    ebiten.IsFocused(),
			WindowWidth:        winW,
			WindowHeight:       winH,
			StatusKind:         statusKind,
			StatusText:         status,
			CurrentSlidePath:   currentSlide.Path,
			HasPreviousSlide:   hasPrev,
			ToastCount:         len(toasts),
			ChangeInterval:     a.Cfg.ChangeInterval,
			TransitionDuration: a.Cfg.TransitionDuration,
		})
	}

	if alpha := splashAlpha(startedAt); alpha > 0 {
		a.DrawSplash(screen, ui.SplashData{
			Logo:         a.SplashLogo,
			WindowWidth:  winW,
			WindowHeight: winH,
			Alpha:        alpha,
		})
	}
}

func (a *App) DrawSlide(screen *ebiten.Image, s Slide, winW, winH int, shiftX, scaleBoost, alpha float64) {
	if s.Img == nil {
		return
	}
	imgW, imgH := float64(s.W), float64(s.H)
	if imgW == 0 || imgH == 0 || alpha <= 0 {
		return
	}

	ww, hh := float64(winW), float64(winH)
	scale := min(ww/imgW, hh/imgH) * scaleBoost
	drawW := imgW * scale
	drawH := imgH * scale
	offX := (ww-drawW)/2 + shiftX
	offY := (hh - drawH) / 2

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offX, offY)
	op.ColorScale.Scale(1, 1, 1, float32(alpha))
	screen.DrawImage(s.Img, op)
}

func (a *App) DrawToast(screen *ebiten.Image, winW, winH int, title, message, kind string, updating bool, progress float64, changedAt, expiresAt time.Time, stackIndex int) {
	ui.DrawToast(screen, winW, winH, title, message, kind, updating, progress, changedAt, expiresAt, stackIndex, a.LargeFont, a.SmallFont)
}

func (a *App) DrawDebugPanel(screen *ebiten.Image, data ui.DebugPanelData) {
	ui.DrawDebugPanel(screen, data, a.LargeFont, a.SmallFont)
}

func (a *App) DrawSplash(screen *ebiten.Image, data ui.SplashData) {
	ui.DrawSplash(screen, data, a.LargeFont, a.SmallFont)
}

func splashAlpha(startedAt time.Time) float64 {
	if startedAt.IsZero() {
		return 0
	}
	elapsed := time.Since(startedAt)
	if elapsed <= 0 {
		return 1
	}
	if elapsed >= SplashDuration {
		return 0
	}
	fadeStart := SplashDuration - SplashFadeDuration
	if elapsed <= fadeStart {
		return 1
	}
	remaining := SplashDuration - elapsed
	return float64(remaining) / float64(SplashFadeDuration)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	a.Mu.Lock()
	a.WinW = outsideWidth
	a.WinH = outsideHeight
	a.Mu.Unlock()
	return outsideWidth, outsideHeight
}
