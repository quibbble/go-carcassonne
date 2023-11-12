package go_carcassonne

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

// state holds all necessary game objects and high level game logic
type state struct {
	turn            string
	teams           []string
	winners         []string
	playTiles       map[string]*tile // teams to the tiles to place onto the board at the start of any given turn
	lastPlacedTiles map[string]*tile // the tiles that were last placed by each team
	board           *board
	boardTokens     []*token       // a list of tokens currently on the board
	tokens          map[string]int // number of tokens each team can play
	scores          map[string]int // points of each team
	deck            *deck
}

func newState(teams []string, random *rand.Rand) *state {
	tokens := make(map[string]int)
	scores := make(map[string]int)
	playTiles := make(map[string]*tile)
	lastPlacedTiles := make(map[string]*tile)
	for _, team := range teams {
		tokens[team] = 7
		scores[team] = 0
	}
	deck := newDeck(random)
	for _, team := range teams {
		tile, _ := deck.Draw()
		playTiles[team] = tile
		lastPlacedTiles[team] = nil
	}
	return &state{
		turn:            teams[0],
		teams:           teams,
		winners:         make([]string, 0),
		playTiles:       playTiles,
		lastPlacedTiles: lastPlacedTiles,
		board:           newBoard(),
		boardTokens:     make([]*token, 0),
		tokens:          tokens,
		scores:          scores,
		deck:            deck,
	}
}

func (s *state) RotateTileRight(team string) error {
	if s.playTiles[team] == nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot rotate tile"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	s.playTiles[team].RotateRight()
	return nil
}

func (s *state) RotateTileLeft(team string) error {
	if s.playTiles[team] == nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot rotate tile"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	s.playTiles[team].RotateLeft()
	return nil
}

func (s *state) PlaceTile(team string, tile *tile, x, y int) error {
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot play on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if s.playTiles[team] == nil || !tile.equals(s.playTiles[team]) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot place tile %+v", team, tile),
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	if err := s.board.Place(tile, x, y); err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidAction,
		}
	}
	s.lastPlacedTiles[team] = tile
	s.playTiles[team] = nil

	// if there are no tokens to place or cannot place token anywhere skip place token action here
	if s.tokens[s.turn] == 0 || len(s.targets()) <= 1 {
		if err := s.PlaceToken(s.turn, true, 0, 0, "", ""); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) PlaceToken(team string, pass bool, x, y int, typ, side string) error {
	if len(s.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s game already completed", key),
			Status: bgerr.StatusGameOver,
		}
	}
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("currently %s's turn", s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if s.playTiles[team] != nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot place token"),
			Status: bgerr.StatusInvalidAction,
		}
	}
	// try placing token
	if !pass {
		if s.lastPlacedTiles[team].X != x || s.lastPlacedTiles[team].Y != y {
			return &bgerr.Error{
				Err:    fmt.Errorf("cannot place token on tile at %d,%d", s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y),
				Status: bgerr.StatusInvalidAction,
			}
		}
		if !contains(TokenTypes, typ) {
			return &bgerr.Error{
				Err:    fmt.Errorf("invalid token type %s", typ),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if (typ == Thief || typ == Knight) && !contains(Sides, side) {
			return &bgerr.Error{
				Err:    fmt.Errorf("invalid side %s with token %s", side, typ),
				Status: bgerr.StatusInvalidActionDetails,
			}
		} else if typ == Farmer && !contains(FarmSides, side) {
			return &bgerr.Error{
				Err:    fmt.Errorf("invalid farm side %s with token %s", side, typ),
				Status: bgerr.StatusInvalidActionDetails,
			}
		} else if typ == Monk && s.lastPlacedTiles[team].Center != Cloister {
			return &bgerr.Error{
				Err:    fmt.Errorf("cannot place %s on tile that does not contain %s", Monk, Cloister),
				Status: bgerr.StatusInvalidAction,
			}
		}
		if s.tokens[team] <= 0 {
			return &bgerr.Error{
				Err:    fmt.Errorf("not enough tokens to place for team %s", team),
				Status: bgerr.StatusInvalidAction,
			}
		}
		// check to ensure token does not connect to pre-existing tokens in given structure
		switch typ {
		case Thief:
			road, err := s.board.generateRoad(x, y, side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			tokens := tokensInStructure(s.boardTokens, road)
			if len(tokens) > 0 {
				return &bgerr.Error{
					Err:    fmt.Errorf("cannot place token on road that is already claimed"),
					Status: bgerr.StatusInvalidAction,
				}
			}
		case Knight:
			city, err := s.board.generateCity(x, y, side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			tokens := tokensInStructure(s.boardTokens, city)
			if len(tokens) > 0 {
				return &bgerr.Error{
					Err:    fmt.Errorf("cannot place token on city that is already claimed"),
					Status: bgerr.StatusInvalidAction,
				}
			}
		case Farmer:
			farm, err := s.board.generateFarm(x, y, side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			tokens := tokensInStructure(s.boardTokens, farm)
			if len(tokens) > 0 {
				return &bgerr.Error{
					Err:    fmt.Errorf("cannot place token on farm that is already claimed"),
					Status: bgerr.StatusInvalidAction,
				}
			}
		}
		// add the token
		s.tokens[team]--
		token := newToken(x, y, team, typ, side)
		s.boardTokens = append(s.boardTokens, token)
	}
	// score completed cities
	citySides := make([]string, 0)
	for _, side := range Sides {
		if s.lastPlacedTiles[team].Sides[side] == City {
			citySides = append(citySides, side)
		}
	}
	if len(citySides) > 0 {
		if s.lastPlacedTiles[team].ConnectedCitySides {
			citySides = citySides[:1]
		}
		for _, citySide := range citySides {
			city, err := s.board.generateCity(s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y, citySide)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			if city.complete {
				// add to completed list in board
				s.board.completeCities = append(s.board.completeCities, city)

				// check if token inside city
				inside := tokensInStructure(s.boardTokens, city)
				if len(inside) > 0 {
					// score and add points
					points, err := scoreCity(city)
					if err != nil {
						return &bgerr.Error{
							Err:    err,
							Status: bgerr.StatusInvalidAction,
						}
					}
					winners := pointsWinners(inside)
					for _, winner := range winners {
						s.scores[winner] += points
					}
					// remove inside from board and add back to tokens pile
					for _, token := range inside {
						s.tokens[token.Team]++
					}
					s.boardTokens = removeTokens(s.boardTokens, inside...)
					// set color of completed
					for _, n := range city.nodes {
						for _, side := range n.sides {
							n.tile.Teams[side] = winners
						}
					}
				}
			}
		}
	}
	// score completed roads
	roadSides := make([]string, 0)
	for _, side := range Sides {
		if s.lastPlacedTiles[team].Sides[side] == Road {
			roadSides = append(roadSides, side)
		}
	}
	if len(roadSides) > 0 {
		if len(roadSides) <= 2 {
			roadSides = roadSides[:1]
		}
		for _, roadSide := range roadSides {
			road, err := s.board.generateRoad(s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y, roadSide)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			if road.complete {
				inside := tokensInStructure(s.boardTokens, road)
				if len(inside) > 0 {
					// score and add points
					points, err := scoreRoad(road)
					if err != nil {
						return &bgerr.Error{
							Err:    err,
							Status: bgerr.StatusInvalidAction,
						}
					}
					winners := pointsWinners(inside)
					for _, winner := range winners {
						s.scores[winner] += points
					}
					// remove inside from board and add back to tokens pile
					for _, token := range inside {
						s.tokens[token.Team]++
					}
					s.boardTokens = removeTokens(s.boardTokens, inside...)
					// add to completed list in board
					s.board.completeRoads = append(s.board.completeRoads, road)
					// set color of completed
					for _, n := range road.nodes {
						for _, side := range n.sides {
							n.tile.Teams[side] = winners
						}
					}
				}
			}
		}
	}
	// score completed cloister
	cloisterLocationsToCheck := [][]int{
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y - 1},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y - 1},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y - 1}}
	for _, location := range cloisterLocationsToCheck {
		tile := s.board.tile(location[0], location[1])
		if tile != nil && tile.Center == Cloister {
			count, err := s.board.tilesSurroundingCloister(location[0], location[1])
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			if count == 8 {
				for _, token := range s.boardTokens {
					if token.Type == Monk && token.X == location[0] && token.Y == location[1] {
						// add to score
						s.scores[token.Team] += count + 1
						// remove inside from board and add back to tokens pile
						s.tokens[token.Team]++
						s.boardTokens = removeTokens(s.boardTokens, token)
						// set color of completed
						tile.CenterTeam = token.Team
						break
					}
				}
			}
		}
	}

	// draw tile for player
	if !s.deck.Empty() {
		tile, _ := s.deck.Draw()
		s.playTiles[s.turn] = tile
	}

	tilesInHands := 0
	for _, tile := range s.playTiles {
		if tile != nil {
			tilesInHands++
		}
	}

	if tilesInHands > 0 {
		// next turn
		for idx, team := range s.teams {
			if team == s.turn {
				s.turn = s.teams[(idx+1)%len(s.teams)]
				break
			}
		}

		// edge case where play tile isn't playable so re-draw
		if !s.board.playable(s.playTiles[s.turn]) {
			if !s.deck.Empty() {
				tried := []*tile{s.playTiles[s.turn]}
				retryLimit := s.deck.Size()
				for i := 0; i < retryLimit; i++ {
					tile, _ := s.deck.Draw()
					if s.board.playable(tile) {
						s.playTiles[s.turn] = tile
						s.deck.Add(tried...)
						tried = nil
						break
					} else {
						tried = append(tried, tile)
					}
				}
				if tried != nil {
					// edge case where no tile in the deck is playable so end the game instead
					if err := s.score(); err != nil {
						return err
					}
				}
			} else {
				// edge case where tiles still remain but cannot be played so end the game instead
				if err := s.score(); err != nil {
					return err
				}
			}
		}
	} else {
		// all tiles have been played so score
		if err := s.score(); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) score() error {
	// score incomplete roads, cities, and cloister and score farms
	for len(s.boardTokens) > 0 {
		token := s.boardTokens[0]
		switch token.Type {
		case Knight:
			city, err := s.board.generateCity(token.X, token.Y, token.Side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			// score and add points
			points, err := scoreCity(city)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			inside := tokensInStructure(s.boardTokens, city)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of incomplete
			for _, n := range city.nodes {
				for _, side := range n.sides {
					n.tile.Teams[side] = winners
				}
			}
		case Thief:
			road, err := s.board.generateRoad(token.X, token.Y, token.Side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			// score and add points
			points, err := scoreRoad(road)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			inside := tokensInStructure(s.boardTokens, road)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of incomplete
			for _, n := range road.nodes {
				for _, side := range n.sides {
					n.tile.Teams[side] = winners
				}
			}
		case Monk:
			tile := s.board.tile(token.X, token.Y)
			if tile != nil && tile.Center == Cloister {
				count, err := s.board.tilesSurroundingCloister(token.X, token.Y)
				if err != nil {
					return &bgerr.Error{
						Err:    err,
						Status: bgerr.StatusInvalidAction,
					}
				}
				s.scores[token.Team] += count + 1
				// remove inside from board and add back to tokens pile
				s.tokens[token.Team]++
				s.boardTokens = removeTokens(s.boardTokens, token)
				// set color of incomplete
				tile.CenterTeam = token.Team
			}
		case Farmer:
			farm, err := s.board.generateFarm(token.X, token.Y, token.Side)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			// score and add points
			points, err := scoreFarm(farm, s.board.completeCities)
			if err != nil {
				return &bgerr.Error{
					Err:    err,
					Status: bgerr.StatusInvalidAction,
				}
			}
			inside := tokensInStructure(s.boardTokens, farm)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of farmland
			for _, n := range farm.nodes {
				// get number of city sides
				numCities := 0
				for _, section := range n.tile.Sides {
					if section == City {
						numCities++
					}
				}
				// edge case where two adjacent disconnected city sections leads to uncolored farmland between them
				if !n.tile.ConnectedCitySides && numCities == 2 {
					for _, farmSide := range FarmSides {
						n.tile.FarmTeams[farmSide] = winners
					}
				} else {
					// otherwise, do normal coloring
					for _, farmSide := range n.sides {
						n.tile.FarmTeams[farmSide] = winners
					}
				}
			}
		}
	}
	// winner is team with the highest score
	max := 0
	winners := make([]string, 0)
	for p, score := range s.scores {
		if score > max {
			max = score
			winners = []string{p}
		} else if score == max {
			winners = append(winners, p)
		}
	}
	s.winners = winners
	return nil
}

func (s *state) SetWinners(winners []string) error {
	for _, winner := range winners {
		if !contains(s.teams, winner) {
			return &bgerr.Error{
				Err:    fmt.Errorf("winner not in teams"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
	}
	s.winners = winners
	return nil
}

func (s *state) targets() []*bg.BoardGameAction {
	targets := make([]*bg.BoardGameAction, 0)
	if s.playTiles[s.turn] != nil {
		// add rotating tile as valid targets
		targets = append(targets, &bg.BoardGameAction{
			Team:       s.turn,
			ActionType: ActionRotateTileLeft,
		}, &bg.BoardGameAction{
			Team:       s.turn,
			ActionType: ActionRotateTileRight,
		})
		// find all valid places to play tile
		emptySpaces := s.board.getEmptySpaces()
		for _, emptySpace := range emptySpaces {
			valid := true
			for _, side := range Sides {
				if emptySpace.adjacent[side] != nil && emptySpace.adjacent[side].Sides[AcrossSide[side]] != s.playTiles[s.turn].Sides[side] {
					valid = false
				}
			}
			if valid {
				targets = append(targets, &bg.BoardGameAction{
					Team:       s.turn,
					ActionType: ActionPlaceTile,
					MoreDetails: PlaceTileActionDetails{
						X: emptySpace.X,
						Y: emptySpace.Y,
					},
				})
			}
		}
	} else {
		// find all valid places to play token
		targets = append(targets, &bg.BoardGameAction{
			Team:       s.turn,
			ActionType: ActionPlaceToken,
			MoreDetails: PlaceTokenActionDetails{
				Pass: true,
			},
		})
		if s.lastPlacedTiles[s.turn].Center == Cloister {
			targets = append(targets, &bg.BoardGameAction{
				Team:       s.turn,
				ActionType: ActionPlaceToken,
				MoreDetails: PlaceTokenActionDetails{
					X:    s.lastPlacedTiles[s.turn].X,
					Y:    s.lastPlacedTiles[s.turn].Y,
					Type: Monk,
				},
			})
		}
		for _, side := range Sides {
			switch s.lastPlacedTiles[s.turn].Sides[side] {
			case Road:
				// check if road is already claimed
				road, _ := s.board.generateRoad(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, side)
				if len(tokensInStructure(s.boardTokens, road)) == 0 {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Thief,
							Side: side,
						},
					})
				}
				// check if farmland A is claimed
				farm, _ := s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchA))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchA),
						},
					})
				}
				// check if farmland B is claimed
				farm, _ = s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchB))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchB),
						},
					})
				}
			case City:
				// check if city is already claimed
				city, _ := s.board.generateCity(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, side)
				if len(tokensInStructure(s.boardTokens, city)) == 0 {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Knight,
							Side: side,
						},
					})
				}
			case Farm:
				// check if farmland is claimed
				farm, _ := s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchA))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchA),
						},
					}, &bg.BoardGameAction{
						Team:       s.turn,
						ActionType: ActionPlaceToken,
						MoreDetails: PlaceTokenActionDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchB),
						},
					})
				}
			}
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	if s.playTiles[s.turn] == nil {
		message = fmt.Sprintf("%s must place a token", s.turn)
	}
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}

// get the tokens that fall in the structure
func tokensInStructure(tokens []*token, structure *structure) []*token {
	tokensInside := make([]*token, 0)
	for _, token := range tokens {
		for _, n := range structure.nodes {
			// check if token type matches section type and token on node
			if StructureTypeToTokenType[structure.typ] == token.Type &&
				n.tile.X == token.X && n.tile.Y == token.Y && contains(n.sides, token.Side) {
				tokensInside = append(tokensInside, token)
			}
		}
	}
	return tokensInside
}

// create a new list that has removed toRemove from original
func removeTokens(original []*token, toRemove ...*token) []*token {
	newTokens := make([]*token, 0)
	for _, token := range original {
		found := false
		for _, rm := range toRemove {
			if token.X == rm.X && token.Y == rm.Y && token.Side == rm.Side &&
				token.Type == rm.Type && token.Team == rm.Team {
				found = true
			}
		}
		if !found {
			newTokens = append(newTokens, token)
		}
	}
	return newTokens
}

// given a list of tokens, get the teams(s) with the most tokens
func pointsWinners(tokens []*token) []string {
	max := 0
	winners := make([]string, 0)
	tally := make(map[string]int)
	for _, token := range tokens {
		tally[token.Team]++
	}
	for team, count := range tally {
		if count > max {
			winners = []string{team}
			max = count
		} else if count == max {
			winners = append(winners, team)
		}
	}
	sort.Strings(winners)
	return winners
}
