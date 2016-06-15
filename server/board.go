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

func TileHasTripleWordEffect(verticalIdx int, horizontalIdx int) bool {
    return verticalIdx % 7 == 0 && horizontalIdx % 7 == 0 && !(verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2)
}

func init() {
    /* Init board. Create tiles. */
    tiles = make([][]Tile, VERTICAL_TILES_AMOUNT)
    for verticalIdx := 0; verticalIdx < VERTICAL_TILES_AMOUNT; verticalIdx++ {
        tiles[verticalIdx] = make([]Tile, HORIZONTAL_TILES_AMOUNT)
        for horizontalIdx := 0; horizontalIdx < HORIZONTAL_TILES_AMOUNT; horizontalIdx++ {
            var tile Tile = Tile{}
            if TileHasTripleWordEffect(verticalIdx, horizontalIdx) {
                tile.effect = TRIPLE_WORD_TILE_EFFECT
            } else {
                tile.effect = NO_TILE_EFFECT
            }
            tiles[verticalIdx][horizontalIdx] = tile
        }
    }
}

func GetLetterFromTile(verticalTileIdx int, horizontalTileIdx int) rune {
    // Return letter at given tile.
    //
    // Guarantees:
    // - Returns 0 if tile is empty
    // - Returns 0 if tile idx out of range

    if verticalTileIdx < 0 || horizontalTileIdx < 0 || verticalTileIdx > VERTICAL_TILES_AMOUNT-1 || horizontalTileIdx > HORIZONTAL_TILES_AMOUNT-1 {
        return 0
    }

    return tiles[verticalTileIdx][horizontalTileIdx].letter
}

func GetHorizontalWordAtTile(verticalTileIdx int, horizontalTileIdx int) string {
    // Get the horizontal word (read from left to right)
    // that the letter on the given tile is a part of (if any).

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerLeftTileOfWord int = 0
    var outerRightTileOfWord int = HORIZONTAL_TILES_AMOUNT-1

    // Go to left outer tile of horizontal word at this tile
    for horizontalLoopIdx := horizontalTileIdx-1; horizontalLoopIdx >= 0; horizontalLoopIdx-- {
        if GetLetterFromTile(verticalTileIdx, horizontalLoopIdx) == 0 {
            outerLeftTileOfWord = horizontalLoopIdx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for horizontalLoopIdx := outerLeftTileOfWord+1; horizontalLoopIdx <= HORIZONTAL_TILES_AMOUNT; horizontalLoopIdx++ {
        if GetLetterFromTile(verticalTileIdx, horizontalLoopIdx) == 0 {
            outerRightTileOfWord = horizontalLoopIdx - 1
        }
    }

    // Read word from outer left to outer right tile
    var word string;
    for loopIdx := outerLeftTileOfWord; loopIdx < outerRightTileOfWord; loopIdx++ {
        word += string(GetLetterFromTile(verticalTileIdx, loopIdx))
    }

    return word

}

func GetVerticalWordAtTile(verticalTileIdx int, horizontalTileIdx int) string {
    // Get the vertical word (read from top to bottom)
    // that the letter on the given tile is a part of (if any).

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerTopTileOfWord int = 0
    var outerBottomTileOfWord int = VERTICAL_TILES_AMOUNT-1

    // Go to top outer tile of vertical word at this tile
    for verticalLoopIdx := verticalTileIdx-1; verticalLoopIdx >= 0; verticalLoopIdx-- {
        if GetLetterFromTile(verticalLoopIdx, horizontalTileIdx) == 0 {
            outerTopTileOfWord = verticalLoopIdx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for verticalLoopIdx := outerTopTileOfWord+1; verticalLoopIdx <= VERTICAL_TILES_AMOUNT; verticalLoopIdx++ {
        if GetLetterFromTile(verticalLoopIdx, horizontalTileIdx) == 0 {
            outerBottomTileOfWord = verticalLoopIdx - 1
        }
    }

    // Read word from outer left to outer right tile
    var word string;
    for loopIdx := outerTopTileOfWord; loopIdx < outerBottomTileOfWord; loopIdx++ {
        word += string(GetLetterFromTile(loopIdx, horizontalTileIdx))
    }

    return word

}

func IsLegalPlacement(verticalTileIdx int, horizontalTileIdx int, letter rune) (bool, string) {
    // check whether the given letter can be placed on the given tile
    // without actually placing it there.
    //
    // returns true if placement is possible, otherwise false
    // returns a string as second parameter describing why the placement is
    // illegal, otherwise emptystring.

    if ! IsLegalLetter(letter) {
        return false, "Character illegal."
    }

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx) != 0 {
        return false, "Tile occupied"
    }

    // Except for when the first letter is placed,
    // a new letter must always be adjacing at least one more.
    if !(verticalTileIdx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontalTileIdx == (HORIZONTAL_TILES_AMOUNT-1)/2) &&
       GetLetterFromTile(verticalTileIdx + 1 , horizontalTileIdx) +
       GetLetterFromTile(verticalTileIdx, horizontalTileIdx + 1) +
       GetLetterFromTile(verticalTileIdx - 1 , horizontalTileIdx) +
       GetLetterFromTile(verticalTileIdx, horizontalTileIdx - 1) == 0 {
       return false, "Not connected to word"
   }

    return true, ""

}

func PlaceLetter(verticalTileIdx int, horizontalTileIdx int, letter rune) {
    // add a letter to the board.
    // throw an error if placement of the leter is not legal

    if isLegal, reason := IsLegalPlacement(verticalTileIdx, horizontalTileIdx, letter); !isLegal {
        log.Fatal("Cannot place letter. ", reason)
    }

    tiles[verticalTileIdx][horizontalTileIdx].letter = letter
}

func RemoveLetter(verticalTileIdx int, horizontalTileIdx int) {
    // Remove one single letter from the board that has
    // not been locked yet
}

func LockLetters() {
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

func CheckAllWords() bool {
    // Check whether all words on the board are
    // in the dictionary
    return false
}
