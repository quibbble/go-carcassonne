package go_carcassonne

import (
	"fmt"
	"math/rand"
)

type deck struct {
	tiles  []*tile
	random *rand.Rand
}

func newDeck(random *rand.Rand) *deck {
	d := make([]*tile, 0)
	for _, tileAmount := range tiles {
		for i := 0; i < tileAmount.amount; i++ {
			d = append(d, tileAmount.tile.copy())
		}
	}
	result := &deck{
		tiles:  d,
		random: random,
	}
	result.Shuffle()
	return result
}

func (d *deck) Shuffle() {
	for i := 0; i < len(d.tiles); i++ {
		r := d.random.Intn(len(d.tiles))
		if i != r {
			d.tiles[r], d.tiles[i] = d.tiles[i], d.tiles[r]
		}
	}
}

func (d *deck) Empty() bool {
	return len(d.tiles) == 0
}

func (d *deck) Add(tiles ...*tile) {
	d.tiles = append(d.tiles, tiles...)
	d.Shuffle()
}

func (d *deck) Draw() (*tile, error) {
	size := len(d.tiles)
	if size <= 0 {
		return nil, fmt.Errorf("cannot draw from empty deck")
	}
	tile := d.tiles[size-1]
	d.tiles = d.tiles[:size-1]
	return tile, nil
}
