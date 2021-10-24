package go_carcassonne

import (
	bg "github.com/quibbble/go-boardgame"
)

const key = "carcassonne"

type Builder struct{}

func (b *Builder) Create(options bg.BoardGameOptions) (bg.BoardGame, error) {
	return NewCarcassonne(options)
}

func (b *Builder) Key() string {
	return key
}
