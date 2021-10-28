package go_carcassonne

import (
	bg "github.com/quibbble/go-boardgame"
	"time"
)

const key = "carcassonne"

type Builder struct{}

func (b *Builder) Create(options bg.BoardGameOptions, seed ...int64) (bg.BoardGame, error) {
	if len(seed) > 0 {
		return NewCarcassonne(options, seed[0])
	}
	return NewCarcassonne(options, time.Now().Unix())
}

func (b *Builder) Key() string {
	return key
}
