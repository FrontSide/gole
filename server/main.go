package main

import (
    "log"
);

func init() {
    // Get the full english letter set
    letterSet = GetFullLetterSet()
    log.Println("Init letter set:", letterSet)
}

func main() {

    var game Game = Game{}

    log.Println(game)

    AddPlayer("mrman", &game.players)
    AddPlayer("mrman2", &game.players)
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
