package usb

import (
	"os"
	"os/user"
	"path/filepath"
	"sort"
)

func FindCandidates() []string {
	roots := scanRoots()
	var matches []string

	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path1 := filepath.Join(root, entry.Name())
			if IsPayloadDir(path1) {
				matches = append(matches, path1)
			}

			subEntries, err := os.ReadDir(path1)
			if err != nil {
				continue
			}
			for _, sub := range subEntries {
				if !sub.IsDir() {
					continue
				}
				path2 := filepath.Join(path1, sub.Name())
				if IsPayloadDir(path2) {
					matches = append(matches, path2)
				}
			}
		}
	}

	if len(matches) == 0 {
		return nil
	}
	sort.Strings(matches)

	uniq := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, m := range matches {
		if seen[m] {
			continue
		}
		seen[m] = true
		uniq = append(uniq, m)
	}
	return uniq
}

func scanRoots() []string {
	currentUser, err := user.Current()
	if err != nil || currentUser.Username == "" {
		return []string{"/media", "/run/media", "/mnt"}
	}
	return []string{
		"/media",
		filepath.Join("/media", currentUser.Username),
		"/run/media",
		filepath.Join("/run/media", currentUser.Username),
		"/mnt",
	}
}

func IsPayloadDir(path string) bool {
	// config.ini は任意（なければ埋め込みデフォルトを使用）
	// data/ ディレクトリの存在のみ必須条件とする
	dataPath := filepath.Join(path, "data")
	info, err := os.Stat(dataPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
