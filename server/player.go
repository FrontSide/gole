package main

import (
    "errors"
)

type Player struct {
    Name string;
    Points int;
    LettersInHand []Letter;
}

func (player *Player) RemoveLetterFromHand(letter rune) error {
    // Remove the first occurrenct of the given letter from the players hand
    // Throw an error if the letter does not exist in the set
    var idxOfLetterToRemove = -1
    for idx, letterInHand := range player.LettersInHand {
        if letterInHand.Character == letter {
            idxOfLetterToRemove = idx
        }
    }

    if idxOfLetterToRemove == -1 {
        return errors.New("Player has no letter " + string(letter) + " in hand.")
    }

    player.LettersInHand = append(player.LettersInHand[:idxOfLetterToRemove], player.LettersInHand[idxOfLetterToRemove+1:]...)
    return nil

}
