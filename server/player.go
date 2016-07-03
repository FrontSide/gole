package main

type Player struct {
    Name string;
    Points int;
    LettersInHand []Letter;
}

func HasLetterInHand(player *Player, letter rune) bool {
    // Check if the given player has the given letter in their hand
    for _, letterInHand := range player.LettersInHand {
        if letterInHand.Character == letter {
            return true
        }
    }
    return false
}
