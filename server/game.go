package main

import (
    "log"
    "unicode/utf8"
);

type Game struct {
    players []Player;
    tiles [][]Tile;
}

var MIN_NUMBER_OF_PLAYERS = 2
var MAX_NUMBER_OF_PLAYERS = 4
var DEFAULT_NUMBER_OF_LETTERS_IN_HAND = 7
var letterSet string;

func AddPlayer(playerName string, players *[]Player) {
    // Add a player to the list of players for the
    // upcoming game play

    if len(*players) >= MAX_NUMBER_OF_PLAYERS {
        log.Fatal("No more players can be added to the Game.")
    }

    player := Player{name: playerName}

    for i := 0; i < DEFAULT_NUMBER_OF_LETTERS_IN_HAND; i++ {
        player.lettersInHand += PopLetterFromSet()
    }

    *players = append(*players, player)

}

func PopLetterFromSet() string {
    // Pop the last letter (right end) from the
    // letter string. The letter will be returned and the
    // occurrence will be removed from the string.
    if len(letterSet) < 1 {
        log.Fatal("Cannot pop letter from set. Empty.")
    }
    var letterToReturn = string(letterSet[utf8.RuneCountInString(letterSet)-1])
    letterSet = letterSet[:len(letterSet)-1]
    return letterToReturn
}
