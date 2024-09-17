package spritesheet

import "image"

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	Tilesize      int
}

func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w, h, t,
	}
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.Tilesize
	y := (index / s.WidthInTiles) * s.Tilesize

	return image.Rect(x, y, x+s.Tilesize, y+s.Tilesize)
}
