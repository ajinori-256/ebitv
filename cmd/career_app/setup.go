package main

import (
	"flag"
	"os"
	"path/filepath"

	appPkg "career_ad/internal/app"
	"log"
)

type Config struct {
	DebugUI    bool
	MockUSBDir string
}

func parseArgs() Config {
	debugUI := flag.Bool("debug_ui", false, "デバッグオーバーレイを表示")
	mockUSBDir := flag.String("mock_usb_dir", "", "USBの代わりに同期元として扱うディレクトリ（data/ を含む）")
	flag.Parse()
	return Config{
		DebugUI:    *debugUI,
		MockUSBDir: *mockUSBDir,
	}
}

func setupApp() (*appPkg.App, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	adRoot := filepath.Join(home, "ad")
	dataDir := filepath.Join(adRoot, "data")
	configPath := filepath.Join(adRoot, "config.ini")

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf("データディレクトリを作成できません: %v", err)
	}

	return appPkg.NewApp(dataDir, configPath, defaultConfigINI, smallFont, largeFont, splashLogo)
}
