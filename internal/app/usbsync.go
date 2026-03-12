package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"career_ad/internal/config"
	"career_ad/internal/fileops"
	"career_ad/internal/media"
	"career_ad/internal/usb"
)

const usbScanInterval = 2 * time.Second

func (a *App) WatchUSBDevices(ctx context.Context) {
	ticker := time.NewTicker(usbScanInterval)
	defer ticker.Stop()

	seen := map[string]bool{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			candidates := a.collectSyncCandidates()
			if len(candidates) == 0 {
				a.SetIdleStatus("USBデバイス待機中")
			}
			active := make(map[string]bool, len(candidates))
			for _, src := range candidates {
				active[src] = true
				if seen[src] {
					continue
				}
				if err := a.syncFromUSB(src); err != nil {
					a.SetStatusWithKind("データ更新失敗: "+err.Error(), "error", false, true)
				} else {
					a.SetStatusWithKind("データ更新完了", "success", false, true)
				}
				seen[src] = true
			}
			for src := range seen {
				if !active[src] {
					delete(seen, src)
				}
			}
		}
	}
}

func (a *App) collectSyncCandidates() []string {
	candidates := append([]string(nil), usb.FindCandidates()...)

	a.Mu.RLock()
	mockDirs := append([]string(nil), a.MockUSBDirs...)
	a.Mu.RUnlock()

	for _, dir := range mockDirs {
		if dir == "" {
			continue
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if usb.IsPayloadDir(abs) {
			candidates = append(candidates, abs)
		}
	}

	if len(candidates) == 0 {
		return nil
	}
	sort.Strings(candidates)
	uniq := candidates[:0]
	seen := map[string]bool{}
	for _, c := range candidates {
		if seen[c] {
			continue
		}
		seen[c] = true
		uniq = append(uniq, c)
	}
	return uniq
}

func (a *App) syncFromUSB(sourceRoot string) error {
	a.SetCopyProgress("USBデータの確認中", 0.05)

	srcConfig := filepath.Join(sourceRoot, "config.ini")
	srcData := filepath.Join(sourceRoot, "data")

	if err := os.MkdirAll(filepath.Dir(a.ConfigPath), 0o755); err != nil {
		return fmt.Errorf("保存先ディレクトリ作成失敗: %w", err)
	}

	a.SetCopyProgress("config.iniをコピー中", 0.1)
	if _, err := os.Stat(srcConfig); err == nil {
		// USBにconfig.iniがあればコピー
		if err := fileops.CopyFile(srcConfig, a.ConfigPath); err != nil {
			return fmt.Errorf("config.iniコピー失敗: %w", err)
		}
	} else {
		// USBになければ埋め込みデフォルトをそのまま保存
		if err := os.WriteFile(a.ConfigPath, a.DefaultConfigData, 0o644); err != nil {
			return fmt.Errorf("デフォルトconfig保存失敗: %w", err)
		}
	}

	// config を再読み込み
	if cfg, err := config.LoadFile(a.ConfigPath, a.DefaultConfigData); err == nil {
		a.Mu.Lock()
		a.Cfg = cfg
		a.Mu.Unlock()
	}

	a.SetCopyProgress("保存先ディレクトリを準備中", 0.15)
	if err := os.MkdirAll(a.DataDir, 0o755); err != nil {
		return fmt.Errorf("保存先data作成失敗: %w", err)
	}

	totalFiles, err := fileops.CountFiles(srcData)
	if err != nil {
		return fmt.Errorf("dataファイル数取得失敗: %w", err)
	}
	if totalFiles == 0 {
		a.SetCopyProgress("コピー対象ファイルがありません", 0.85)
	}

	if err := fileops.CopyDirContents(srcData, a.DataDir, func(done, total int, _ string) {
		if total <= 0 {
			return
		}
		phase := 0.15 + 0.55*(float64(done)/float64(total))
		a.SetCopyProgress(fmt.Sprintf("dataをコピー中 %d/%d", done, total), phase)
	}); err != nil {
		return fmt.Errorf("dataコピー失敗: %w", err)
	}

	newFileList, err := media.CollectImagePaths(a.DataDir)
	if err != nil {
		return fmt.Errorf("画像パスの収集に失敗: %w", err)
	}
	if len(newFileList) == 0 {
		a.SetCopyProgress("コピー完了: 画像ファイルが見つかりません", 0.9)
	} else {
		a.SetCopyProgress("コピー完了: 画像ファイルを読み込み中", 0.9)
	}

	a.Mu.Lock()
	a.LastChanged = time.Now()
	a.Transition = time.Time{}
	newSlides := media.NewDefaultSlideSource(newFileList)
	a.Source = newSlides
	a.LastChanged = time.Now()
	a.Transition = time.Time{}
	a.Mu.Unlock()

	a.SetCopyProgress("不要ファイルを削除中", 0.75)
	if err := fileops.RemoveUnmatchedFiles(a.DataDir, srcData); err != nil {
		return fmt.Errorf("不要ファイル削除失敗: %w", err)
	}

	a.SetCopyProgress("画像を読み込み中", 0.85)


	return nil
}
