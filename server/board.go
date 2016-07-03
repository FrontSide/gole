package main

import (
    "log"
    "encoding/json"
    "math"
)

var VERTICAL_TILES_AMOUNT int = 15
var HORIZONTAL_TILES_AMOUNT int = 15

type Tile struct {
    IsLocked bool
    Letter rune
    Effect SpecialTileEffect
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
        log.Printf("v::%d,h::%d", verticalIdx, horizontalIdx)
        log.Println(math.Abs(float64(verticalIdx - (VERTICAL_TILES_AMOUNT-1)/2)))
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

func GetLetterFromTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) rune {
    // Return letter at given tile.
    //
    // Guarantees:
    // - Returns 0 if tile is empty
    // - Returns 0 if tile idx out of range

    if verticalTileIdx < 0 || horizontalTileIdx < 0 || verticalTileIdx > VERTICAL_TILES_AMOUNT-1 || horizontalTileIdx > HORIZONTAL_TILES_AMOUNT-1 {
        return 0
    }

    return tiles[verticalTileIdx][horizontalTileIdx].Letter
}

func GetHorizontalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) []Tile {
    // Get the horizontal word (read from left to right)
    // that the letter on the given tile is a part of (if any).

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerLeftTileOfWord int = 0
    var outerRightTileOfWord int = HORIZONTAL_TILES_AMOUNT-1

    // Go to left outer tile of horizontal word at this tile
    for horizontalLoopIdx := horizontalTileIdx-1; horizontalLoopIdx >= 0; horizontalLoopIdx-- {
        if GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles) == 0 {
            outerLeftTileOfWord = horizontalLoopIdx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for horizontalLoopIdx := outerLeftTileOfWord+1; horizontalLoopIdx <= HORIZONTAL_TILES_AMOUNT; horizontalLoopIdx++ {
        if GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles) == 0 {
            outerRightTileOfWord = horizontalLoopIdx - 1
        }
    }

    var fullRow []Tile = tiles[verticalTileIdx]

    return fullRow[outerLeftTileOfWord:outerRightTileOfWord]

}

func GetVerticalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) []Tile {
    // Get the vertical word (read from top to bottom)
    // that the letter on the given tile is a part of (if any).

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles) == 0 {
        log.Fatal("Cannot retrieve horizontal word. Initial tile is empty.")
    }

    var outerTopTileOfWord int = 0
    var outerBottomTileOfWord int = VERTICAL_TILES_AMOUNT-1

    // Go to top outer tile of vertical word at this tile
    for verticalLoopIdx := verticalTileIdx-1; verticalLoopIdx >= 0; verticalLoopIdx-- {
        if GetLetterFromTile(verticalLoopIdx, horizontalTileIdx, tiles) == 0 {
            outerTopTileOfWord = verticalLoopIdx + 1
            break;
        }
    }

    // Go to right outer tile of horizontal word at given tile
    for verticalLoopIdx := outerTopTileOfWord+1; verticalLoopIdx <= VERTICAL_TILES_AMOUNT; verticalLoopIdx++ {
        if GetLetterFromTile(verticalLoopIdx, horizontalTileIdx, tiles) == 0 {
            outerBottomTileOfWord = verticalLoopIdx - 1
        }
    }

    return tiles[outerTopTileOfWord:outerBottomTileOfWord][horizontalTileIdx]

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

    if GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles) != 0 {
        return false, "Tile occupied"
    }

    // Except for when the first letter is placed,
    // a new letter must always be adjacing at least one more.
    if tiles[verticalTileIdx][horizontalTileIdx].Effect != CENTER_TILE_EFFECT &&
       GetLetterFromTile(verticalTileIdx + 1 , horizontalTileIdx, tiles) +
       GetLetterFromTile(verticalTileIdx, horizontalTileIdx + 1, tiles) +
       GetLetterFromTile(verticalTileIdx - 1 , horizontalTileIdx, tiles) +
       GetLetterFromTile(verticalTileIdx, horizontalTileIdx - 1, tiles) == 0 {
       return false, "Not connected to word"
   }

    return true, ""

}

func LockLetters(tiles [][]Tile) {
    // lock all letters on the board so they cannot be
    // removed by the player anymore
    for _, column := range tiles {
        for _, tile := range column {
            if tile.Letter != 0 {
                tile.IsLocked = true
            }
        }
    }
}


func GetLegalPlacementMapAsJson(game *Game) string {
    // Return a two-dimensional slice containing information,
    // whether a letter can be placed on the tile of the according
    // slice indices, as JSON.

    var legalPlacementMap = make([][]bool, VERTICAL_TILES_AMOUNT)
    for verticalIdx, tileRow := range game.Tiles {
        var legalPlacementMapRow = make([]bool, HORIZONTAL_TILES_AMOUNT)
        for horizontalIdx, _ := range tileRow {
            legalPlacementMapRow[horizontalIdx], _ = IsLegalPlacement(verticalIdx, horizontalIdx, 'a',  game.Tiles)
        }
        legalPlacementMap[verticalIdx] = legalPlacementMapRow
    }

    legalPlacementMapJson, err := json.Marshal(legalPlacementMap)
    if err != nil {
        log.Fatal(err)
    }

    return string(legalPlacementMapJson)

}

func GetBoardAsJson(game *Game) string {
    // Converts a given two-dimensional
    // slice of tiles into a json representation
    // of the board with all the information the tile
    // struct has as well.
    boardJson, err := json.Marshal(game.Tiles)
    if err != nil {
        log.Fatal(err)
    }

    return string(boardJson)
}
