package main

import (
    "log"
    "fmt"
    "errors"
    "strings"
);

var games []Game

func init() {}

func GetGameByUUID(uuid string) (*Game, error) {
    // Return the game with the given ID if existent in games array.
    // Requires:
    // - a lower letter standard unix uuid as created for the games
    // Guarantees:
    // - Return reference to game struct that has uuid set as game id
    // - log fatal if no game in array has the given uuid
    for idx, _ := range games {
        if games[idx].Id == strings.TrimSpace(uuid) {
            return &games[idx], nil
        }
    }
    return &Game{}, errors.New("Game with uuid " + uuid + " could not be found!")
}

func StartNewGame(playerNames ...string) (string, error) {
    // Initiate a new game
    // Requires:
    // - A list of player names (2-4 players are legal)
    // Guarantees:
    // - Creates a new game object and adds the players
    // - Trow an error if the number of players is illegal
    // - Return the uuid of the game if successful

    if len(playerNames) < 2 || len(playerNames) > 4 {
        return "", errors.New(fmt.Sprintf("%d is not a legal amount of players. Needs to be 2-3.", len(playerNames)))
    }

    game := &Game{}
    game.Id = GetNewUUID()

    // Letter set needs to be generated before Players are added
    // since letters need to be taken off the set.
    var err error
    game.LetterSet, err = GetFullLetterSet()
    if err != nil {
        return "", err
    }

    for _, playerName := range playerNames {
        log.Printf("Add player %s to Game %s\n", playerName, game.Id)
        AddPlayer(playerName, game)
    }

    game.Tiles = GetCleanTiles()

    // Update placement legality of whole board
    game.UpdatePlacementLegalityOfAllTiles()

    //First player in slice will have first turn
    game.PlayerIdxWithTurn = 0

    games = append(games, *game)

    return game.Id, nil
}

func main() {
    StartWebServer()
}
