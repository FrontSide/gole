package main

import (
    "strings"
);

type Player struct {
    Name string;
    Points int;
    LettersInHand string;
}

func HasLetterInHand(player *Player, letter rune) bool {
    // Check if the given player has the given letter in their hand
    return strings.Contains(player.LettersInHand, string(letter))
}
