package main

import (
    "log"
);

var games []Game

func init() {
    // Get the full english letter set
    letterSet = GetFullLetterSet()
    log.Println("Init letter set:", letterSet)
}

func StartNewGame(...playerNames string) string {
    // Initiate a new game
    // Requires:
    // - A list of player names (2-4 players are legal)
    // Guarantees:
    // - Creates a new game object and adds the players
    // - Trow an error if the number of players is illegal
    // - Return the uuid of the game if successful

    if playerNames.len() < 2 || playerNames.len() > 4 {
        log.Fatal("%d is not a legal amount of players. Needs to be 2-3.", playerNames.len)
    }

    var game Game = Game{}
    game.id = GetNewUUID()

    for playerName := range playerNames {
        AddPlayer(playerName, &game.players)
    }

    return game.id
}

func main() {
    // Run web service // accept API calls
    log.Println(game)

    log.Println("Init letter set:", letterSet)

    log.Println(game)

    LockLetters()
    PlaceLetter(7, 7, 'a')
    PlaceLetter(7, 8, 'x')
    log.Println(string(GetLetterFromTile(7, 8)))
    log.Println(string(GetVerticalWordAtTile(7, 8)))
    log.Println("'" + string(GetHorizontalWordAtTile(7, 8)) + "'")
    log.Println(HasLetterInHand(&game.players[0], 'e'))
}
