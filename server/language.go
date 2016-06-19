//English
//100 letters

package main

import (
    "log"
    "math/rand"
    "time"
);

type Letter struct {
    occurrences, pointValue int
}

var lettersAmount = 100
var letterDistribution = map[rune]Letter {
    ' ': {2, 0},
    'a': {9, 1},
    'b': {2, 3},
    'c': {2, 3},
    'd': {4, 2},
    'e': {12, 1},
    'f': {2, 4},
    'g': {3, 2},
    'h': {2, 4},
    'j': {1, 8},
    'k': {1, 5},
    'i': {9, 1},
    'l': {4, 1},
    'm': {2, 3},
    'n': {6, 1},
    'o': {8, 1},
    'p': {2, 3},
    'q': {1, 10},
    'r': {6, 1},
    's': {4, 1},
    't': {6, 1},
    'u': {4, 1},
    'v': {2, 4},
    'x': {1, 8},
    'w': {2, 4},
    'y': {2, 4},
    'z': {1, 10},
}

func GetFullLetterSet() string {
    /* Return a randomly shuffled full initial set of letters. */
    fullLetterSet := ""
    letterCount := 0
    for letter, letterProperties := range letterDistribution {
        letterCount += letterProperties.occurrences
        for i := 0; i < letterProperties.occurrences; i++ {
            fullLetterSet += string(letter)
        }
    }
    if letterCount != lettersAmount {
        log.Fatal("Letter distribution error! Is %d, expected %d\n", letterCount, lettersAmount)
    }
    // Shuffle string
    rand.Seed(time.Now().UTC().UnixNano())
    randomIndices := rand.Perm(letterCount)
    var fullShuffledLetterSet []uint8 = make([]uint8, letterCount)
    for originalIndex, newRandomIndex := range randomIndices {
        fullShuffledLetterSet[newRandomIndex] = uint8(fullLetterSet[originalIndex])
    }
    return string(fullShuffledLetterSet)
}

func IsLegalLetter(letter rune) bool {
    // Check whether a given letter is an existing character in the letter set
    for legalLetter, _ := range letterDistribution {
        if rune(legalLetter) == letter {
            return true
        }
    }
    return false
}

func GetLetterFromRune(letter rune) Letter {
    // Return the full letter struct for a letter
    // passed as a rune
    return letterDistribution[string(letter)]
}
