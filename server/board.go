package main

import (
    "log"
)

var VERTICAL_TILES_AMOUNT int = 15
var HORIZONTAL_TILES_AMOUNT int = 15

type Tile struct {
    locked bool
    letter rune
    effect SpecialTileEffect
}

type SpecialTileEffect int
const (
    DOUBLE_LETTER_TILE_EFFECT SpecialTileEffect = iota
    TRIPLE_LETTER_TILE_EFFECT
    DOUBLE_WORD_TILE_EFFECT
    TRIPLE_WORD_TILE_EFFECT
    NO_TILE_EFFECT
)

var tiles [][]Tile;

func tile_has_triple_word_effect(vertical_idx int, horizontal_idx int) bool {
    return vertical_idx % 7 == 0 && horizontal_idx % 7 == 0 && !(vertical_idx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontal_idx == (HORIZONTAL_TILES_AMOUNT-1)/2)
}

func init() {
    /* Init board. Create tiles. */
    tiles = make([][]Tile, VERTICAL_TILES_AMOUNT)
    for vertical_idx := 0; vertical_idx < VERTICAL_TILES_AMOUNT; vertical_idx++ {
        tiles[vertical_idx] = make([]Tile, HORIZONTAL_TILES_AMOUNT)
        for horizontal_idx := 0; horizontal_idx < HORIZONTAL_TILES_AMOUNT; horizontal_idx++ {
            var tile Tile = Tile{}
            if tile_has_triple_word_effect(vertical_idx, horizontal_idx) {
                tile.effect = TRIPLE_WORD_TILE_EFFECT
            } else {
                tile.effect = NO_TILE_EFFECT
            }
            tiles[vertical_idx][horizontal_idx] = tile
        }
    }
}

func get_letter_from_tile(vertical_tile_idx int, horizontal_tile_idx int) rune {
    // Return letter at given tile.
    // Return 0 if no letter on tile.
    return tiles[vertical_tile_idx][horizontal_tile_idx].letter
}

func get_horizontal_word_at_tile(vertical_tile_idx int, horizontal_tile_idx int) string {
    // Get the horizontal word (read from left to right)
    // that the letter on the given tile is a part of (if any).

    if get_letter_from_tile(vertical_tile_idx, horizontal_tile_idx) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outer_left_tile_of_word int = 0
    var outer_right_tile_of_word int = HORIZONTAL_TILES_AMOUNT-1

    // Go to left outer tile of horizontal word at this tile
    for horizontal_loop_idx := horizontal_tile_idx-1; horizontal_loop_idx >= 0; horizontal_loop_idx-- {
        if get_letter_from_tile(vertical_tile_idx, horizontal_loop_idx) == 0 {
            outer_left_tile_of_word = horizontal_loop_idx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for horizontal_loop_idx := outer_left_tile_of_word+1; horizontal_loop_idx <= HORIZONTAL_TILES_AMOUNT; horizontal_loop_idx++ {
        if get_letter_from_tile(vertical_tile_idx, horizontal_loop_idx) == 0 {
            outer_right_tile_of_word = horizontal_loop_idx - 1
        }
    }

    // Read word from outer left to outer right tile
    var word string;
    for loop_idx := outer_left_tile_of_word; loop_idx < outer_right_tile_of_word; loop_idx++ {
        word += string(get_letter_from_tile(vertical_tile_idx, loop_idx))
    }

    return word

}

func get_vertical_word_at_tile(vertical_tile_idx int, horizontal_tile_idx int) string {
    // Get the vertical word (read from top to bottom)
    // that the letter on the given tile is a part of (if any).

    if get_letter_from_tile(vertical_tile_idx, horizontal_tile_idx) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outer_top_tile_of_word int = 0
    var outer_bottom_tile_of_word int = VERTICAL_TILES_AMOUNT-1

    // Go to top outer tile of vertical word at this tile
    for vertical_loop_idx := vertical_tile_idx-1; vertical_loop_idx >= 0; vertical_loop_idx-- {
        if get_letter_from_tile(vertical_loop_idx, horizontal_tile_idx) == 0 {
            outer_top_tile_of_word = vertical_loop_idx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for vertical_loop_idx := outer_top_tile_of_word+1; vertical_loop_idx <= VERTICAL_TILES_AMOUNT; vertical_loop_idx++ {
        if get_letter_from_tile(vertical_loop_idx, horizontal_tile_idx) == 0 {
            outer_bottom_tile_of_word = vertical_loop_idx - 1
        }
    }

    // Read word from outer left to outer right tile
    var word string;
    for loop_idx := outer_top_tile_of_word; loop_idx < outer_bottom_tile_of_word; loop_idx++ {
        word += string(get_letter_from_tile(loop_idx, horizontal_tile_idx))
    }

    return word

}

func place_letter(vertical_tile_idx int, horizontal_tile_idx int, letter rune) {
    // add a letter to the board.
    // throw an error if the letter cannot be placed there
    // (tile occupied)
    if get_letter_from_tile(vertical_tile_idx, horizontal_tile_idx) != 0 {
        log.Fatal("Cannot place letter. Tile occupied")
    }

    if ! is_legal_letter(letter) {
        log.Fatal("Cannot place letter. Character illegal.")
    }

    tiles[vertical_tile_idx][horizontal_tile_idx].letter = letter
}

func remove_letter(vertical_tile_idx int, horizontal_tile_idx int) {
    // Remove one single letter from the board that has
    // not been locked yet
}

func lock_letters() {
    // lock all letters on the board so they cannot be
    // removed by the player anymore
    for _, column := range tiles {
        for _, tile := range column {
            if tile.letter != 0 {
                tile.locked = true
            }
        }
    }
}

func check_all_words() bool {
    // Check whether all words on the board are
    // in the dictionary
    return false
}
