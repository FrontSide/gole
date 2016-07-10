package main

import (
    "log"
    "math"
    "errors"
    "encoding/json"
)

var VERTICAL_TILES_AMOUNT int = 15
var HORIZONTAL_TILES_AMOUNT int = 15

type Tile struct {
    IsLocked bool
    Letter Letter
    Effect SpecialTileEffect
    PlacementIsLegal bool
}

func TileSliceToString(tiles []Tile) string {

    tilesJson, err := json.Marshal(tiles)
    if err != nil {
        log.Fatal("Could not convert tiles slice to Json")
    }
    return string(tilesJson)

}

type SpecialTileEffect int
const (
    DOUBLE_LETTER_TILE_EFFECT SpecialTileEffect = iota
    TRIPLE_LETTER_TILE_EFFECT
    DOUBLE_WORD_TILE_EFFECT
    TRIPLE_WORD_TILE_EFFECT
    NO_TILE_EFFECT
    CENTER_TILE_EFFECT // not really an effect but a special tile
)

func TileHasTripleWordEffect(verticalIdx int, horizontalIdx int) bool {
    // Return a bool that indicates whether a tile at given
    // index has a 'triple word' effect, according to the
    // original scrabble board.
    return verticalIdx % 7 == 0 && horizontalIdx % 7 == 0 && !(verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2)
}

func TileHasDoubleWordEffect(verticalIdx int, horizontalIdx int) bool {
    // Return a bool that indicates whether a tile at given
    // index has a 'double word' effect, according to the
    // original scrabble board.


    if verticalIdx == horizontalIdx || verticalIdx == HORIZONTAL_TILES_AMOUNT-1-horizontalIdx {
        return !(TileHasTripleWordEffect(verticalIdx, horizontalIdx) ||
                 TileHasTripleLetterEffect(verticalIdx, horizontalIdx) ||
                 TileHasDoubleLetterEffect(verticalIdx, horizontalIdx) ||
                 TileIsCenterTile(verticalIdx, horizontalIdx))
    }
    return false

}

func TileHasTripleLetterEffect(verticalIdx int, horizontalIdx int) bool {
    // Return a bool that indicates whether a tile at given
    // index has a 'triple letter' effect, according to the
    // original scrabble board.

    if verticalIdx == 1 || verticalIdx == VERTICAL_TILES_AMOUNT - 2 {
        return math.Abs(float64(horizontalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 2
    }

    if horizontalIdx == 1 || horizontalIdx == HORIZONTAL_TILES_AMOUNT-2 {
        return math.Abs(float64(verticalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 2
    }



    return math.Abs(float64(verticalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 2 &&
           math.Abs(float64(horizontalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 2

}

func TileHasDoubleLetterEffect(verticalIdx int, horizontalIdx int) bool {
    // Return a bool that indicates whether a tile at given
    // index has a 'double letter' effect, according to the
    // original scrabble board.
    if verticalIdx == 0 || verticalIdx == VERTICAL_TILES_AMOUNT-1 {
        return (horizontalIdx + 1) % 4 == 0 && !TileHasTripleLetterEffect(verticalIdx, horizontalIdx)
    }

    if horizontalIdx == 0 || horizontalIdx == HORIZONTAL_TILES_AMOUNT-1 {
        return (verticalIdx + 1) % 4 == 0 && !TileHasTripleLetterEffect(verticalIdx, horizontalIdx)
    }

    //If the the column index is adjacent to the middle column
    if math.Abs(float64(verticalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 1 {
        return math.Abs(float64(horizontalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 1 ||
               math.Abs(float64(horizontalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 5
    }

    //If the the row index is adjacent to the middle row
    if math.Abs(float64(horizontalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 1 {
        return math.Abs(float64(verticalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 1 ||
               math.Abs(float64(verticalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 5
    }

    //If the the colum index is the middle
    if verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 {
        return math.Abs(float64(horizontalIdx - (VERTICAL_TILES_AMOUNT-1)/2)) == 3
    }

    //If the the row index is the middle
    if horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2 {
        return math.Abs(float64(verticalIdx - (HORIZONTAL_TILES_AMOUNT-1)/2)) == 3
    }

    return false

}

func TileIsCenterTile(verticalIdx int, horizontalIdx int) bool {
    // Return a bool that indicates whether a tile at given
    // index has it the oard's center tile
    return (verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2)
}

func GetCleanTiles() [][]Tile {
    // Create a initial 2-dimensional array of board tiles.
    // Requires:
    // - The constants required to create the tile array need to be availvable.
    //   (VERTICAL_TILES_AMOUNT, HORIZONTAL_TILES_AMOUNT, etc.)
    // Guarantees:
    // - Return a 2-dimesnisonal array with elements of type 'Tile'
    // - The tiles are empty, unlocked and are initiated with their tile effects

    var tiles = make([][]Tile, VERTICAL_TILES_AMOUNT)
    for verticalIdx := 0; verticalIdx < VERTICAL_TILES_AMOUNT; verticalIdx++ {
        tiles[verticalIdx] = make([]Tile, HORIZONTAL_TILES_AMOUNT)
        for horizontalIdx := 0; horizontalIdx < HORIZONTAL_TILES_AMOUNT; horizontalIdx++ {
            var tile Tile = Tile{}
            if TileHasTripleWordEffect(verticalIdx, horizontalIdx) {
                tile.Effect = TRIPLE_WORD_TILE_EFFECT
            } else if TileHasDoubleWordEffect(verticalIdx, horizontalIdx) {
                tile.Effect = DOUBLE_WORD_TILE_EFFECT
            } else if TileHasTripleLetterEffect(verticalIdx, horizontalIdx) {
                tile.Effect = TRIPLE_LETTER_TILE_EFFECT
            } else if TileHasDoubleLetterEffect(verticalIdx, horizontalIdx) {
                tile.Effect = DOUBLE_LETTER_TILE_EFFECT
            } else if TileIsCenterTile(verticalIdx, horizontalIdx) {
                tile.Effect = CENTER_TILE_EFFECT
            } else {
                tile.Effect = NO_TILE_EFFECT
            }
            tiles[verticalIdx][horizontalIdx] = tile
        }
    }
    return tiles
}

func GetLetterFromTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) (Letter, error) {
    // Return letter at given tile.
    //
    // Guarantees:
    // - Returns empty letter struct and error if tile is empty
    // - Returns empty letter struct and error if tile idx out of range

    if verticalTileIdx < 0 || horizontalTileIdx < 0 || verticalTileIdx > VERTICAL_TILES_AMOUNT-1 || horizontalTileIdx > HORIZONTAL_TILES_AMOUNT-1 {
        return Letter{}, errors.New("Index out of bounds")
    }

    letter := tiles[verticalTileIdx][horizontalTileIdx].Letter
    if letter == (Letter{}) {
        return letter, errors.New("No letter on tile")
    }

    return letter, nil
}

func GetHorizontalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) []Tile {
    // Get the horizontal word (read from left to right)
    // that the letter on the given tile is a part of (if any).

    var err error
    _, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
    if err != nil {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerLeftTileOfWord int = 0
    var outerRightTileOfWord int = HORIZONTAL_TILES_AMOUNT-1

    // Go to left outer tile of horizontal word at this tile
    for horizontalLoopIdx := horizontalTileIdx-1; horizontalLoopIdx >= 0; horizontalLoopIdx-- {
        _, err = GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles)
        if err != nil {
            outerLeftTileOfWord = horizontalLoopIdx + 1
            break
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for horizontalLoopIdx := outerLeftTileOfWord+1; horizontalLoopIdx <= HORIZONTAL_TILES_AMOUNT; horizontalLoopIdx++ {
        _, err := GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles)
        if err != nil {
            outerRightTileOfWord = horizontalLoopIdx
            break;
        }
    }

    var fullRow []Tile = tiles[verticalTileIdx]

    return fullRow[outerLeftTileOfWord:outerRightTileOfWord]

}

func GetVerticalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) []Tile {
    // Get the vertical word (read from top to bottom)
    // that the letter on the given tile is a part of (if any).

    log.Printf("Get vertical word at tile %d,%d", verticalTileIdx, horizontalTileIdx)

    var err error
    _, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
    if err != nil {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerTopTileOfWord int = 0
    var outerBottomTileOfWord int = VERTICAL_TILES_AMOUNT-1

    // Go to top outer tile of vertical word at this tile
    for verticalLoopIdx := verticalTileIdx; verticalLoopIdx >= 0; verticalLoopIdx-- {
        _, err = GetLetterFromTile(verticalLoopIdx, horizontalTileIdx, tiles)
        if err != nil {
            outerTopTileOfWord = verticalLoopIdx + 1
            break
        }
    }

    log.Printf("The outer top tile is: %d", outerTopTileOfWord)

    // Go to bottom outer tile of horizontal word at given tile
    for verticalLoopIdx := outerTopTileOfWord; verticalLoopIdx < VERTICAL_TILES_AMOUNT; verticalLoopIdx++ {
        _, err := GetLetterFromTile(verticalLoopIdx, horizontalTileIdx, tiles)
        if err != nil {
            outerBottomTileOfWord = verticalLoopIdx - 1
            break
        }
    }

    log.Printf("The outer bottom tile is: %d", outerBottomTileOfWord)


    var verticalWordTiles []Tile
    for _, horizontalTiles := range tiles[outerTopTileOfWord:outerBottomTileOfWord] {
        verticalWordTiles = append(verticalWordTiles, horizontalTiles[horizontalTileIdx])
        log.Printf("Append letter %c", horizontalTiles[horizontalTileIdx].Letter.Character)
    }
    return verticalWordTiles

}

func IsLegalPlacement(verticalTileIdx int, horizontalTileIdx int, letter rune, tiles [][]Tile) (bool, string) {
    // check whether the given letter can be placed on the given tile
    // without actually placing it there.
    //
    // returns true if placement is possible, otherwise false
    // returns a string as second parameter describing why the placement is
    // illegal, otherwise emptystring.

    if ! IsLegalLetter(letter) {
        return false, "Character illegal."
    }

    var err error
    _, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
    if err == nil {
        return false, "Tile occupied"
    }

    var letterOnConnectedTiles bool
    // Inspect all the tiles connected to the one at the given coordninates
    // for placed letters.
    _, err = GetLetterFromTile(verticalTileIdx + 1 , horizontalTileIdx, tiles)
    if err == nil {
        letterOnConnectedTiles = true
    }

    _, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx + 1, tiles)
    if err == nil {
        letterOnConnectedTiles = true
    }

    _, err = GetLetterFromTile(verticalTileIdx - 1 , horizontalTileIdx, tiles)
    if err == nil {
        letterOnConnectedTiles = true
    }

    _, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx - 1, tiles)
    if err == nil {
        letterOnConnectedTiles = true
    }

    // Except for when the first letter is placed,
    // a new letter must always be adjacing at least one more.
    if tiles[verticalTileIdx][horizontalTileIdx].Effect != CENTER_TILE_EFFECT &&
       !letterOnConnectedTiles {
       return false, "Not connected to word"
   }

    return true, ""

}

func LockLetters(tiles [][]Tile) {
    // lock all letters on the board so they cannot be
    // removed by the player anymore
    for _, column := range tiles {
        for _, tile := range column {
            if tile.Letter.Character != 0 {
                tile.IsLocked = true
            }
        }
    }
}


func (game *Game) UpdatePlacementLegalityOfAllTiles() {
    // Check all tiles for placement legality
    // The PlacementIsLegal bool on every tile is updated
    // Placement is legal on a tile if it is connected to another word
    // or if it's the middle tile AND
    // if it doesn't already have a latter on it

    for verticalIdx, tileRow := range game.Tiles {
        for horizontalIdx, _ := range tileRow {
            // Since range copies the values of the data structure
            // we are iterating over (not a reference anymore),
            // we an not just simply make changes on the new structure.
            // Instead we need to actually change the value on the original
            // struct that we want to manipulate
            game.Tiles[verticalIdx][horizontalIdx].PlacementIsLegal, _ =
                IsLegalPlacement(verticalIdx, horizontalIdx, 'a',  game.Tiles)
        }
    }

}
