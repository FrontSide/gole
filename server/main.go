package main

import (
    "log"
    "unicode/utf8"
);

type Player struct {
    name string;
    points int;
    letters_in_hand string;
}

var MAX_NUMBER_OF_PLAYERS = 4
var DEFAULT_NUMBER_OF_LETTERS_IN_HAND = 7
var letter_set string;
var players []Player;

func add_player(player_name string) {
    // Add a player to the list of players for the
    // upcoming game play

    if len(players) >= MAX_NUMBER_OF_PLAYERS {
        log.Fatal("No more players can be added to the Game.")
    }

    player := Player{name: player_name}

    for i := 0; i < DEFAULT_NUMBER_OF_LETTERS_IN_HAND; i++ {
        player.letters_in_hand += pop_letter_from_set()
    }

    log.Println(player)

}

func pop_letter_from_set() string {
    // Pop the last letter (right end) from the
    // letter string. The letter will be returned and the
    // occurrence will be removed from the string.
    if len(letter_set) < 1 {
        log.Fatal("Cannot pop letter from set. Empty.")
    }
    var letter_to_return = string(letter_set[utf8.RuneCountInString(letter_set)-1])
    letter_set = letter_set[:len(letter_set)-1]
    return letter_to_return
}

func init() {
    // Get the full english letter set
    letter_set = get_full_letter_set()
    log.Println("Init letter set:", letter_set)
}

func main() {
    add_player("mrman")
    add_player("mrman2")
    log.Println("Init letter set:", letter_set)

    place_letter(2, 3, 'a')
}
