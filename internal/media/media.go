package media

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Slide struct {
	Path string
	Img  *ebiten.Image
	W    int
	H    int
}

func CollectImagePaths(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリを読み込めません: %s: %w", dir, err)
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			paths = append(paths, filepath.Join(dir, e.Name()))
		}
	}

	sort.Strings(paths)
	return paths, nil
}

// NewSlidesFromPaths はパス情報だけを持つ Slide スライスを返す（画像は未読み込み）。
func NewSlidesFromPaths(paths []string) []Slide {
	slides := make([]Slide, len(paths))
	for i, p := range paths {
		slides[i] = Slide{Path: p}
	}
	return slides
}

// LoadSlide は 1 枚の Slide の画像を読み込む。
func LoadSlide(s *Slide) error {
	img, _, err := ebitenutil.NewImageFromFile(s.Path)
	if err != nil {
		return fmt.Errorf("画像を読み込めません: %s: %w", s.Path, err)
	}
	bounds := img.Bounds()
	s.Img = img
	s.W = bounds.Dx()
	s.H = bounds.Dy()
	return nil
}

// LoadSlidesFromPaths はすべての画像を一度に読み込む（後方互換用）。
func LoadSlidesFromPaths(paths []string) ([]Slide, error) {
	slides := NewSlidesFromPaths(paths)
	for i := range slides {
		if err := LoadSlide(&slides[i]); err != nil {
			return nil, err
		}
	}
	return slides, nil
}

// SlideSource はスライドデッキを管理するインターフェース。
// Next を呼び出すたびに次のスライドへ進む。
type SlideSource interface {
	// Advance は次のスライドへ進む。現在のスライドを「前のスライド」として記憶する。
	Advance()
	// CurrentSlide は現在表示中のスライドを返す。
	CurrentSlide() Slide
	// PreviousSlide は Advance 直前のスライドを返す。未設定の場合は false。
	PreviousSlide() (Slide, bool)
}

// defaultSlideSource は SlideSource のデフォルト実装。
type defaultSlideSource struct {
	Current Slide
	Previous     Slide
	HasPrevious  bool
	Index		int
	Paths 		[]string
}

func NewDefaultSlideSource(paths []string) *defaultSlideSource {
	s := &defaultSlideSource{
		Current: Slide{},
		Previous:     Slide{},
		HasPrevious:  false,
		Index:        -1,
		Paths:        paths,
	}
	s.Advance() // 最初のスライドをセット
	return s
}

func (s *defaultSlideSource) Advance() {
	if s.Index >= 0 && s.Index < len(s.Paths) {
		if s.Previous.Img != nil {
			s.Previous.Img.Deallocate()
		}
		s.Previous = s.Current
		s.HasPrevious = true
	}
	s.Index++
	if s.Index >= len(s.Paths) {
		s.Index = 0
	}
	if s.Paths == nil || len(s.Paths) == 0 {
		s.Current = Slide{}
		return
	}
	s.Current = Slide{Path: s.Paths[s.Index]}
	LoadSlide(&s.Current)


}

func (s *defaultSlideSource) CurrentSlide() Slide {
	return s.Current
}

func (s *defaultSlideSource) PreviousSlide() (Slide, bool){
 return s.Previous, s.HasPrevious
}

