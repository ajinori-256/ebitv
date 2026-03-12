package fileops

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDirContents(srcDir, dstDir string, onProgress func(done, total int, current string)) error {
	total, err := CountFiles(srcDir)
	if err != nil {
		return err
	}
	done := 0
	if onProgress != nil {
		onProgress(done, total, "")
	}
	return copyDirContentsRec(srcDir, dstDir, total, &done, onProgress)
}

// RemoveUnmatchedFiles 削除する: dstDir 内にあるが srcDir に存在しないファイル/ディレクトリ
func RemoveUnmatchedFiles(dstDir, srcDir string) error {
	srcSet := make(map[string]bool)
	err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(srcDir, path)
		srcSet[relPath] = true
		return nil
	})
	if err != nil {
		return err
	}

	return removeUnmatchedRec(dstDir, srcDir, srcSet, ".")
}

func removeUnmatchedRec(dstDir, srcDir string, srcSet map[string]bool, relPath string) error {
	entries, err := os.ReadDir(filepath.Join(dstDir, relPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		relSubPath := filepath.Join(relPath, entry.Name())
		if relSubPath == "." {
			continue
		}
		if !srcSet[relSubPath] {
			fullPath := filepath.Join(dstDir, relSubPath)
			if err := os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else if entry.IsDir() {
			if err := removeUnmatchedRec(dstDir, srcDir, srcSet, relSubPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyDirContentsRec(srcDir, dstDir string, total int, done *int, onProgress func(done, total int, current string)) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0o755); err != nil {
				return err
			}
			if err := copyDirContentsRec(srcPath, dstPath, total, done, onProgress); err != nil {
				return err
			}
			continue
		}
		if err := CopyFile(srcPath, dstPath); err != nil {
			return err
		}
		(*done)++
		if onProgress != nil {
			onProgress(*done, total, srcPath)
		}
	}
	return nil
}

func CountFiles(root string) (int, error) {
	total := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			total++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}

func CopyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
