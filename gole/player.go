package main

import (
	"errors"
	"math/rand"
	"fmt"
	"time"
	"log"
)

type Player struct {
	Name          string
	Points        int
	LettersInHand []Letter
}

func (player *Player) GetLetterFromHandById(letterId string) (Letter, error) {
	// Retrieve the letter with the given letterId from the player's hand
	// Guarantees:
	// - Return the letter struct from the player's LettersInHand slice
	//   that has the given letterId. nil as second return parameter.
	// - Return an empty struct and an error as second parameter if
	//   no letter with the given Id was found.

	log.Println("Get letter from hand by id :: ", letterId)

	for _, letter := range player.LettersInHand {
		if letter.Id == letterId {
			return letter, nil
		}
	}

	return Letter{}, errors.New(fmt.Sprintf(
		"Letter with ID %s could not be found in player's hand.", letterId))

}

func (player *Player) SortHand(letterIds []string) error {
	// Sort the player's letter slice according to the arrangement
	// of the letterIds slice.
	// Requires:
	// - A slice of letterIds that consists of only ids of the letters
	//   the player had in the hand. The amount of letterIds in this
	//   slice must match the amount of letters in the players current
	//   lettersInHand slice.
	// Guarantees:
	// - Rearranges the LettersInHand slice of the player accoring to the
	//   arrangement of the letterIds slice passed to this function.
	// - Returns an error or nil
	//   depending on whether the action was successful

	log.Println("Sort hand.")

	if len(letterIds) != len(player.LettersInHand) {
		return errors.New(
			"The amount of the given letter IDs does not match " +
			"the amount of letters in the player's hand.")
	}

	var rearrangedLettersInHand = make([]Letter, len(player.LettersInHand))

	for idx, letterId := range letterIds {
		letterToAdd, err := player.GetLetterFromHandById(letterId)
		if err != nil {
			return err
		}
		rearrangedLettersInHand[idx] = letterToAdd
	}

	player.LettersInHand = rearrangedLettersInHand

	log.Println("Sorted Hand::")
	log.Println(player.LettersInHand)

	return nil

}

func (player *Player) ShuffleHand() {
	// Randomly rearrange the array of letters in the player's hand
	// Guarantees:
	// - Rearranges the letter structs stored in the player's
	//   LettersInHand slice.
	// - Return nothing

	log.Println("Shuffle hand.")

	var rearrangedLettersInHand = make([]Letter, len(player.LettersInHand))

	// Get a slice of random indexes with the length of the player's
	// LettersInHand slice
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndexes := rand.Perm(len(player.LettersInHand))

	for originalIdx, newIdx := range randomIndexes {
		rearrangedLettersInHand[newIdx] = player.LettersInHand[originalIdx]
	}

	player.LettersInHand = rearrangedLettersInHand

	log.Println("Shuffled Hand::")
	log.Println(player.LettersInHand)

}

func (player *Player) ReplaceWildcard(letterId string, letterCharacter rune) error {
	// Replace the wildcard character on a letter with the given
	// id with a normal letter.
	// Requires:
	// - The given letterId needs to be the id of a wildcard letter
	//   in the hand of the currently active player and must unsubstituted
	//   i.e. this function can only be called once for one letterId.
	// Guarantees:
	// - Will replace the character on a wildcard letter in the letter
	//   struct with the given id
	// - Will return an error if the letterId does not refer to a letter
	//   in the active players hand with a wildcard character on it.
	// - Will return an error if the given letter character is not
	//   a valid character in the alphabet.
	isWildcardLetter, err := player.IsRawWildcardLetter(letterId)

	if err != nil {
		return err
	}

	if ! isWildcardLetter {
		return errors.New("Cannot replace letter on non-wildcard letter.")
	}

	// Check if the given letterCharacter is a valid letter in the
	// alphabet by trying to turn it into a Letter struct.
	_, err = GetLetterStructFromRune(letterCharacter)
	if err != nil {
		return err
	}

	for idx, letterInHand := range player.LettersInHand {
		if letterInHand.Id == letterId {
			player.LettersInHand[idx].Character = letterCharacter
			return nil
		}
	}

	return errors.New("The Wildcard Replacement has failed for an unknown reason.")

}

func (player *Player) IsRawWildcardLetter(letterId string) (bool, error) {
	// Tell whether the letter with the given letterId is a wildcard tile,
	// that has not yet been substituted with a real letter.
	// Guarantees:
	// - Return true or false to indicate whether the letter with given
	//   is is an unsubstituted wildcard letter
	// - Return an error as second return parameter it the given ID does
	//   not refer to a letter in the players hand.
	for _, letterInHand := range player.LettersInHand {
		if letterInHand.Id == letterId {
			return letterInHand.Character == WILDCARD_CHARACTER, nil
		}
	}
	return false, errors.New("Player has no letter with ID" + string(letterId) + " in hand.")
}

func (player *Player) PopLetterFromHand(letterId string) (Letter, error) {
	// Remove the letter with given letterId from the Player's hand
	// and return it
	// Guarantees:
	// - Remove the letter struct with the given Id
	//   from the LettersInHand slice
	// - Return the full Letter struct that has been removed from the hand
	// - Return an error if a letter with the given Id does not exist
	//   in the Player's LettersInHand slice. The player's hand is unmodified.
	for idx, letterInHand := range player.LettersInHand {
		if letterInHand.Id == letterId {
			player.LettersInHand = append(player.LettersInHand[:idx], player.LettersInHand[idx+1:]...)
			return letterInHand, nil
		}
	}

	return Letter{}, errors.New("Player has no letter with ID" + string(letterId) + " in hand.")

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
