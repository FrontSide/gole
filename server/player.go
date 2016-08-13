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

func (player *Player) AddLetterToHand(letter Letter) error {
    // Add a letter to the hand of the player
    // Requires:
    // - Full letter struct of letter to be added
    // Guarantees:
    // - Throw an error if the maximum of letters in hand
    //   would exceed by adding this letter

    if len(player.LettersInHand) >= MAX_NUMBER_OF_LETTERS_IN_HAND {
        return errors.New("Cannot add letter to player hand. Maximum reached.")
    }

    player.LettersInHand = append(player.LettersInHand, letter)
    return nil
}
