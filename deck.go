package go_carcassonne

import (
	"fmt"
	"math/rand"
)

type deck struct {
	tiles []*tile
}

func newDeck() *deck {
	tiles := make([]*tile, 0)
	for tile, num := range Tiles {
		for i := 0; i < num; i++ {
			tiles = append(tiles, tile.copy())
		}
	}
	deck := &deck{
		tiles: tiles,
	}
	deck.Shuffle()
	return deck
}

func (d *deck) Shuffle() {
	for i := 0; i < len(d.tiles); i++ {
		r := rand.Intn(len(d.tiles))
		if i != r {
			d.tiles[r], d.tiles[i] = d.tiles[i], d.tiles[r]
		}
	}
}

func (d *deck) Empty() bool {
	if len(d.tiles) <= 0 {
		return true
	}
	return false
}

func (d *deck) Add(tiles ...*tile) {
	for _, tile := range tiles {
		d.tiles = append(d.tiles, tile)
	}
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
