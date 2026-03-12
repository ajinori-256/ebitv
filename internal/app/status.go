package app

import (
	"log"
	"time"
)

func (a *App) SetStatus(text string, updating bool, autoExpire bool) {
	a.Mu.Lock()
	a.StatusText = text
	a.Updating = updating
	if updating {
		a.StatusKind = "progress"
	} else {
		a.StatusKind = "info"
		a.CopyProgress = 0
		if autoExpire {
			a.Toasts = append(a.Toasts, ToastItem{
				Title:     "通知",
				Message:   text,
				Kind:      "info",
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(StatusDisplayPeriod),
			})
		}
	}
	a.StatusChanged = time.Now()
	a.Mu.Unlock()
}

func (a *App) SetStatusWithKind(text, kind string, updating bool, autoExpire bool) {
	log.Printf("Status Update: %s - %s (%s)", "USB同期",text, kind)
	
	a.Mu.Lock()
	a.StatusText = text
	a.StatusKind = kind
	a.Updating = updating
	if !updating {
		a.CopyProgress = 0
		if autoExpire {
			a.Toasts = append(a.Toasts, ToastItem{
				Title:     "USB同期",
				Message:   text,
				Kind:      kind,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(StatusDisplayPeriod),
			})
		}
	}
	a.StatusChanged = time.Now()
	a.Mu.Unlock()
}

func (a *App) SetStatusWithTitle(text, title, kind string, updating bool, autoExpire bool) {
	log.Printf("Status Update: %s - %s (%s)", title, text, kind)
	a.Mu.Lock()
	a.StatusText = text
	a.StatusKind = kind
	a.Updating = updating
	if !updating {
		a.CopyProgress = 0
		if autoExpire {
			a.Toasts = append(a.Toasts, ToastItem{
				Title:     title,
				Message:   text,
				Kind:      kind,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(StatusDisplayPeriod),
			})
		}
	}
	a.StatusChanged = time.Now()
	a.Mu.Unlock()
}

func (a *App) SetCopyProgress(text string, progress float64) {
	log.Printf("Copy Progress: %s (%.2f)", text, progress)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	a.Mu.Lock()
	a.StatusText = text
	a.StatusKind = "progress"
	a.Updating = true
	a.CopyProgress = progress
	a.StatusChanged = time.Now()
	a.Mu.Unlock()
}

func (a *App) SetIdleStatus(text string) {
	a.Mu.Lock()
	a.StatusText = text
	a.StatusKind = "idle"
	a.Updating = false
	a.CopyProgress = 0
	a.StatusChanged = time.Now()
	a.Mu.Unlock()
}
