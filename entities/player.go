package entities

import "ebiten-rpg/animations"

type PLayerState uint8

const (
	Down PLayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Health     uint
	Animations map[PLayerState]*animations.Animation
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animations[Right]
	}
	if dx < 0 {
		return p.Animations[Left]
	}
	if dy > 0 {
		return p.Animations[Down]
	}
	if dy < 0 {
		return p.Animations[Up]
	}

	return nil
}
