package main

import (
    "log"
);

var games []Game

func init() {}

func GetGameByUUID(uuid string) *Game {
    // Return the game with the given ID if existent in games array.
    // Requires:
    // - a lower letter standard unix uuid as created for the games
    // Guarantees:
    // - Return game struct that has uuid set as game id
    // - log fatal if no game in array has the given uuid
    for _, game := range games {
        if game.id == uuid {
            return &game
        }
    }
    log.Fatalf("Game with uuid %s could not be found!", uuid)
    return &Game{}
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

    var game Game = Game{}
    game.id = GetNewUUID()

    // Letter set needs to be generated before Players are added
    // since letters need to be taken off the set.
    game.letterSet = GetFullLetterSet()

    for _, playerName := range playerNames {
        AddPlayer(playerName, &game)
    }

    game.tiles = GetCleanTiles()

    //First player in slice will have first turn
    game.playerIdxWithTurn = 0

    games = append(games, game)

    return game.id
}

func main() {
    // Run web service // accept API calls
    gameUuid := StartNewGame("MrMan", "MrsWoman")
    //log.Println(GetGameByUUID(StartNewGame("MrMan", "MrsWoman")))
    log.Println(GetGameByUUID(gameUuid))
    FinishTurn(GetGameByUUID(gameUuid))
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
