package main

import (
    "log"
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
    // - Return game struct that has uuid set as game id
    // - log fatal if no game in array has the given uuid
    for _, game := range games {
        if game.Id == strings.TrimSpace(uuid) {
            return &game, nil
        }
    }
    return &Game{}, errors.New("Game with uuid " + uuid + " could not be found!")
}

func StartNewGame(playerNames ...string) string {
    // Initiate a new game
    // Requires:
    // - A list of player names (2-4 players are legal)
    // Guarantees:
    // - Creates a new game object and adds the players
    // - Trow an error if the number of players is illegal
    // - Return the uuid of the game if successful

    if len(playerNames) < 2 || len(playerNames) > 4 {
        log.Fatalf("%d is not a legal amount of players. Needs to be 2-3.", len(playerNames))
    }

    game := &Game{}
    game.Id = GetNewUUID()

    // Letter set needs to be generated before Players are added
    // since letters need to be taken off the set.
    game.LetterSet = GetFullLetterSet()

    for _, playerName := range playerNames {
        AddPlayer(playerName, game)
    }

    game.Tiles = GetCleanTiles()

    //First player in slice will have first turn
    game.PlayerIdxWithTurn = 0

    games = append(games, *game)

    return game.Id
}

func main() {

    StartWebServer()

    // Run web service // accept API calls
    gameUuid := StartNewGame("MrMan", "MrsWoman")
    //log.Println(GetGameByUUID(StartNewGame("MrMan", "MrsWoman")))
    game, _ := GetGameByUUID(gameUuid)
    log.Println(game)
    FinishTurn(game)
    /*
    log.Println("Init letter set:", letterSet)

    log.Println(game)

    LockLetters()
    PlaceLetter(7, 7, 'a')
    PlaceLetter(7, 8, 'x')
    log.Println(string(GetLetterFromTile(7, 8)))
    log.Println(string(GetVerticalWordAtTile(7, 8)))
    log.Println("'" + string(GetHorizontalWordAtTile(7, 8)) + "'")
    log.Println(HasLetterInHand(&game.players[0], 'e'))*/
}
