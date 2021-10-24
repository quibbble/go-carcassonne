# Go-carcassonne

Go-carcassonne is a Golang implementation of the game logic for the board game [Carcassonne](https://boardgamegeek.com/boardgame/822/carcassonne). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

Please check out [Quibbble.com](https://quibbble.com/annex) if you wish to view and play a live version of this game which utilizes this project along with a custom UI.

## Game API

To play a game create a new Carcassonne instance:
```go
carcassonne, err := NewCarcassonne(bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"} // must contain at least 2 and at most 5 teams
})
```

To rotate the play tile (the tile about to be place by the current team) do the following action:
```go
carcassonne.Do(bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "RotateRight", // can also be "RotateLeft"
})
```

To place the play tile on the board do the following action:
```go
carcassonne.Do(bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceTile",
    Details: PlaceTileActionDetails{
        X: 0,
        Y: 1,
    },
})
```

To place a token on the last placed tile do the following action:
```go
carcassonne.Do(bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceToken",
    Details: PlaceTokenActionDetails{
        Pass: false, // true if you wish to pass placing a token
        X: 0,
        Y: 1,
        Type: "Knight", // can be "Farmer", "Knight", "Thief", or "Monk"
        Side: "Top", // if "Knight" or "Thief" can be "Top", "Right", "Bottom", "Left"; if "Farmer" can be "TopA", "TopB", "RightA", ...; if "Monk" then ""
    },
})
```