package main

import (
    "strings"
);

type Player struct {
    name string;
    points int;
    lettersInHand string;
}

func HasLetterInHand(player *Player, letter rune) bool {
    // Check if the given player has the given letter in their hand
    return strings.Contains(player.lettersInHand, string(letter))
}
