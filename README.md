# Go-carcassonne

Go-carcassonne is a [Go](https://golang.org) implementation of the board game [Carcassonne](https://en.wikipedia.org/wiki/Carcassonne_(board_game)). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

Check out [carcassonne.quibbble.com](https://carcassonne.quibbble.com) if you wish to view and play a live version of this game which utilizes this project along with a separate custom UI.

[![Quibbble Carcassonne](https://i.imgur.com/sLVp1x0.png)](https://quibbble.com/annex)

## Usage

To play a game create a new Carcassonne instance:
```go
builder := Builder{}
game, err := builder.Create(&bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"}, // must contain at least 2 and at most 5 teams
    MoreOptions: CarcassonneMoreOptions{
        Seed: 123, // seed used to generate deterministic randomness
    }
})
```

To rotate the play tile (the tile about to be placed by the current team) do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "RotateTileRight", // can also be "RotateTileLeft"
})
```

To place the play tile on the board do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceTile",
    MoreDetails: PlaceTileActionDetails{
        X: 0,
        Y: 1,
    },
})
```

To place a token on the last placed tile do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceToken",
    MoreDetails: PlaceTokenActionDetails{
        Pass: false, // true if you wish to pass placing a token
        X: 0,
        Y: 1,
        Type: "Knight", // can be "Farmer", "Knight", "Thief", or "Monk"
        Side: "Top", // if "Knight" or "Thief" can be "Top", "Right", "Bottom", "Left"; if "Farmer" can be "TopA", "TopB", "RightA", ...; if "Monk" then ""
    },
})
```

To get the current state of the game call the following:
```go
snapshot, err := game.GetSnapshot("TeamA")
```
