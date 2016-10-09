//English
//100 letters

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"gole/golelibs"
)

type LetterAttributes struct {
	Occurrences, PointValue int
}

type Letter struct {
	Id string
	Character  rune
	Attributes LetterAttributes
}

const WILDCARD_CHARACTER rune = '*'

var lettersAmount = 118
var letterDistribution = map[rune]LetterAttributes{
	WILDCARD_CHARACTER: {20, 0},
	'a':                {9, 1},
	'b':                {2, 3},
	'c':                {2, 3},
	'd':                {4, 2},
	'e':                {12, 1},
	'f':                {2, 4},
	'g':                {3, 2},
	'h':                {2, 4},
	'j':                {1, 8},
	'k':                {1, 5},
	'i':                {9, 1},
	'l':                {4, 1},
	'm':                {2, 3},
	'n':                {6, 1},
	'o':                {8, 1},
	'p':                {2, 3},
	'q':                {1, 10},
	'r':                {6, 1},
	's':                {4, 1},
	't':                {6, 1},
	'u':                {4, 1},
	'v':                {2, 4},
	'x':                {1, 8},
	'w':                {2, 4},
	'y':                {2, 4},
	'z':                {1, 10},
}

func GetFullLetterSet() ([]Letter, error) {
	/* Return a randomly shuffled full initial set of letters. */
	var fullLetterSet []Letter
	letterCount := 0
	for letter, letterProperties := range letterDistribution {
		letterCount += letterProperties.Occurrences
		for i := 0; i < letterProperties.Occurrences; i++ {
			letterStruct, err := GetLetterStructFromRune(letter)
			if err != nil {
				return []Letter{Letter{}}, err
			}
			fullLetterSet = append(fullLetterSet, letterStruct)
		}
	}
	if letterCount != lettersAmount {
		return []Letter{Letter{}}, errors.New(fmt.Sprintf("Letter distribution error! Is %d, expected %d\n", letterCount, lettersAmount))
	}
	// Shuffle string
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndexes := rand.Perm(letterCount)
	var fullShuffledLetterSet []Letter = make([]Letter, letterCount)
	for originalIndex, newRandomIndex := range randomIndexes {
		fullShuffledLetterSet[newRandomIndex] = fullLetterSet[originalIndex]
	}
	return fullShuffledLetterSet, nil
}

func GetLetterStructFromRune(letter rune) (Letter, error) {
	// Return the full letter struct for a letter passed as a rune
	// Requires:
	// - Given rune needs to be an existing key in the letterDistribution map
	// Guarantees:
	// - Return the Letter Struct representation of the given character (rune)
	// - Return en error if the given letter rune is not a valid
	//   character in the alphabet.
	var letterStruct Letter
	if _, ok := letterDistribution[letter]; !ok {
		return Letter{}, errors.New("Letter could not be found.")
	}
	letterStruct.Attributes = letterDistribution[letter]
	letterStruct.Character = letter
	letterStruct.Id = golelibs.GetNewUUID()
	return letterStruct, nil
}
