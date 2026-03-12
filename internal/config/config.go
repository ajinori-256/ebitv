package config

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ChangeInterval     time.Duration
	TransitionDuration time.Duration
}

var Default = Config{
	ChangeInterval:     5 * time.Second,
	TransitionDuration: 700 * time.Millisecond,
}

func Parse(data []byte) (Config, error) {
	cfg := Default
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "[") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "change_interval_sec":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.ChangeInterval = time.Duration(n) * time.Second
			}
		case "transition_duration_ms":
			if n, err := strconv.Atoi(val); err == nil && n >= 0 {
				cfg.TransitionDuration = time.Duration(n) * time.Millisecond
			}
		}
	}
	return cfg, scanner.Err()
}

func LoadFile(path string, defaultData []byte) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Parse(defaultData)
	}
	return Parse(data)
}
