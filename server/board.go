VERTICAL_TILES_AMOUNT = 15
HORIZONTAL_TILES_AMOUNT = 15

type Tile struct {
    locked bool
    letter rune
}

var tiles =

func place_letter(vertical_tile, horizontal_tile, letter) {
    // add a letter to the board.
    // throw an error if the letter cannot be placed there
    // (tile occupied)
}

func remove_letter(vertical_tile, horizontal_tile) {
    // Remove one single letter from the board that has
    // not been locked yet
}

func lock_letters() {
    // lock all letters on the board so they cannot be
    // removed by the player anymore
}

func check_all_words() bool {
    // Check whether all words on the board are
    // in the dictionary
}
