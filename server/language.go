//English
//100 letters

package main

import (
    "log"
    "math/rand"
    "time"
);

type Letter struct {
    occurrences, point_value int
}

var letters_amount = 100
var letter_distribution = map[rune]Letter {
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

func get_full_letter_set() string {
    /* Return a randomly shuffled full initial set of letters. */
    full_letter_set := ""
    letter_count := 0
    for letter, letter_properties := range letter_distribution {
        letter_count += letter_properties.occurrences
        for i := 0; i < letter_properties.occurrences; i++ {
            full_letter_set += string(letter)
        }
    }
    if letter_count != letters_amount {
        log.Fatal("Letter distribution error! Is %d, expected %d\n", letter_count, letters_amount)
    }
    // Shuffle string
    rand.Seed(time.Now().UTC().UnixNano())
    random_indices := rand.Perm(letter_count)
    var full_shuffled_letter_set []uint8 = make([]uint8, letter_count)
    for original_index, new_random_index := range random_indices {
        full_shuffled_letter_set[new_random_index] = uint8(full_letter_set[original_index])
    }
    return string(full_shuffled_letter_set)
}
