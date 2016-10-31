package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"
    "fmt"
)

var VERTICAL_TILES_AMOUNT int = 15
var HORIZONTAL_TILES_AMOUNT int = 15

type Tile struct {
	IsLocked         bool
	Letter           Letter
	Effect           SpecialTileEffect
	PlacementIsLegal bool
}

func TileSliceToString(tiles []Tile) string {
	// Return the json representation of a slice
	// of tiles as string.

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

func GetIdStringForTileSlice(tiles []Tile) string {
    // Get a string of concatenated IDs from all
    // letters in a slice of tiles.
    // This function can be used to easily compare the
    // quality of two slices of tiles without the need
    // to compare the whole tile slice.
    //
    // Requires:
    // - A slice of tiles with Letters in them
    // Guarantees:
    // - Get the IDs of all Letters from all tiles
    //   in the given slice and concatenates them to one string
    // - Return the string concatenation of Letter IDs

    var idConcatenation string;

    for _, tile := range tiles {
        idConcatenation += tile.Letter.Id
    }

    return idConcatenation

}

func TileHasTripleWordEffect(verticalIdx int, horizontalIdx int) bool {
	// Return a bool that indicates whether a tile at given
	// index has a 'triple word' effect, according to the
	// original scrabble board.
	return verticalIdx%7 == 0 &&
		horizontalIdx%7 == 0 &&
		!(verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 &&
			horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2)
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

	if verticalIdx == 1 || verticalIdx == VERTICAL_TILES_AMOUNT-2 {
		return math.Abs(float64(horizontalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 2
	}

	if horizontalIdx == 1 || horizontalIdx == HORIZONTAL_TILES_AMOUNT-2 {
		return math.Abs(float64(verticalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 2
	}

	return math.Abs(float64(verticalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 2 &&
		math.Abs(float64(horizontalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 2

}

func TileHasDoubleLetterEffect(verticalIdx int, horizontalIdx int) bool {
	// Return a bool that indicates whether a tile at given
	// index has a 'double letter' effect, according to the
	// original scrabble board.
	if verticalIdx == 0 || verticalIdx == VERTICAL_TILES_AMOUNT-1 {
		return (horizontalIdx+1)%4 == 0 && !TileHasTripleLetterEffect(verticalIdx, horizontalIdx)
	}

	if horizontalIdx == 0 || horizontalIdx == HORIZONTAL_TILES_AMOUNT-1 {
		return (verticalIdx+1)%4 == 0 && !TileHasTripleLetterEffect(verticalIdx, horizontalIdx)
	}

	//If the the column index is adjacent to the middle column
	if math.Abs(float64(verticalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 1 {
		return math.Abs(float64(horizontalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 1 ||
			math.Abs(float64(horizontalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 5
	}

	//If the the row index is adjacent to the middle row
	if math.Abs(float64(horizontalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 1 {
		return math.Abs(float64(verticalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 1 ||
			math.Abs(float64(verticalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 5
	}

	//If the the colum index is the middle
	if verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 {
		return math.Abs(float64(horizontalIdx-(VERTICAL_TILES_AMOUNT-1)/2)) == 3
	}

	//If the the row index is the middle
	if horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2 {
		return math.Abs(float64(verticalIdx-(HORIZONTAL_TILES_AMOUNT-1)/2)) == 3
	}

	return false

}

func TileIsCenterTile(verticalIdx int, horizontalIdx int) bool {
	// Return a bool that indicates whether a tile at given
	// index has it the oard's center tile
	return (verticalIdx == (VERTICAL_TILES_AMOUNT-1)/2 && horizontalIdx == (HORIZONTAL_TILES_AMOUNT-1)/2)
}

func AreValidBoardCoordinates(verticalIdx int, horizontalIdx int) bool {
	// Tell whether the given indexes are within the bounds of the board
	return (verticalIdx < VERTICAL_TILES_AMOUNT) &&
		(horizontalIdx < HORIZONTAL_TILES_AMOUNT) &&
		(verticalIdx > -1) && (horizontalIdx > -1)
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

func GetHorizontalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) (bool, []Tile, int) {
	// Get the horizontal word (read from left to right)
	// that the letter on the given tile is a part of (if any).
	// Requires:
	// - valid vertical and horizontal index of tile on board
	// - two-dimensional array of tiles from board
	// Guarantees:
	// - return boolean (1st value) that describes whether there is
	//   a series (>0) of horizontal letters around the given tile.
	// - return a tile-array (2nd value) with the letters around the tile
	//   in same order as appearing on the board from left to right
	//   (including the letter on the given tile) if existing.
	// - return the vertical index (3rd value) of the top outer tile of the
	//   vertical word on the board if a vertical word has been found
	// - return false (1st value), an empty tile array (2nd value)
	//   and -1 (3rd value) if there are no tiles horizontally adjacent
	//   to the given tile or if the given tile is already empty.

	log.Printf("Get horizontal word at tile %d,%d", verticalTileIdx, horizontalTileIdx)

	// Make sure the given/initial tile is not empty
	var err error
	_, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
	if err != nil {
		log.Printf("Cannot retrieve horizontal word. Initial tile is empty. v:%d,h:%d", verticalTileIdx, horizontalTileIdx)
		return false, []Tile{}, -1
	}

	//Make sure there is at least one horizontally adjacent non-empty tile around the given tile
	_, noTileLeftErr := GetLetterFromTile(verticalTileIdx, horizontalTileIdx-1, tiles)
	_, noTileRightErr := GetLetterFromTile(verticalTileIdx, horizontalTileIdx+1, tiles)
	if noTileLeftErr != nil && noTileRightErr != nil {
		log.Printf("Cannot retrieve horizontal word. Adjacent tiles empty. v:%d,h:%d", verticalTileIdx, horizontalTileIdx)
		return false, []Tile{}, -1
	}

	var outerLeftTileOfWord int = 0
	var outerRightTileOfWord int = HORIZONTAL_TILES_AMOUNT - 1

	// Go to left outer tile of horizontal word at this tile
	for horizontalLoopIdx := horizontalTileIdx - 1; horizontalLoopIdx >= 0; horizontalLoopIdx-- {
		_, err = GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles)
		if err != nil {
			outerLeftTileOfWord = horizontalLoopIdx + 1
			break
		}
	}

	log.Printf("The outer left tile is: %d", outerLeftTileOfWord)

	// Go to right outer tile of horizontal word at given tile
	for horizontalLoopIdx := outerLeftTileOfWord + 1; horizontalLoopIdx <= HORIZONTAL_TILES_AMOUNT; horizontalLoopIdx++ {
		_, err := GetLetterFromTile(verticalTileIdx, horizontalLoopIdx, tiles)
		if err != nil {
			outerRightTileOfWord = horizontalLoopIdx
			break
		}
	}

	log.Printf("The outer right tile is: %d", outerRightTileOfWord)

	var fullRow []Tile = tiles[verticalTileIdx]

	return true, fullRow[outerLeftTileOfWord:outerRightTileOfWord], outerLeftTileOfWord

}

func GetVerticalWordAtTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) (bool, []Tile, int) {
	// Get the vertical word (read from top to bottom)
	// that the letter on the given tile is a part of (if any).
	// Requires:
	// - valid vertical and horizontal index of tile on board
	// - two-dimensional array of tiles from board
	// Guarantees:
	// - return boolean (1st value) that describes whether there is a
	//   series (>0) of vertical letters around the given tile.
	// - return a tile-array (2nd value) with the letters around the tile
	//   in same order as appearing on the board from top to bottom
	//   (including the letter on the given tile) if existing.
	// - return the vertical index (3rd value) of the top outer tile of the
	//   vertical word on the board if a vertical word has been found
	// - return false (1st value) and an empty tile array (2nd value)
	//   and -1 (3rd value) if there are no tiles vertically adjacent
	//   to the given tile or if the given tile is already empty.

	log.Printf("Get vertical word at tile %d,%d", verticalTileIdx, horizontalTileIdx)

	// Make sure the initial tile is not empty
	var err error
	_, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
	if err != nil {
		log.Printf("Cannot retrieve vertical word. Initial tile is empty. v:%d,h:%d", verticalTileIdx, horizontalTileIdx)
		return false, []Tile{}, -1
	}

	//Make sure there are vertically adjacent non-empty tiles around the given tile
	_, noTileAboveErr := GetLetterFromTile(verticalTileIdx-1, horizontalTileIdx, tiles)
	_, noTileBelowErr := GetLetterFromTile(verticalTileIdx+1, horizontalTileIdx, tiles)
	if noTileAboveErr != nil && noTileBelowErr != nil {
		log.Printf("Cannot retrieve vertical word. Adjacent tiles empty. v:%d,h:%d", verticalTileIdx, horizontalTileIdx)
		return false, []Tile{}, -1
	}

	var outerTopTileOfWord int = 0
	var outerBottomTileOfWord int = VERTICAL_TILES_AMOUNT - 1

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
	for verticalLoopIdx := outerTopTileOfWord; verticalLoopIdx <= VERTICAL_TILES_AMOUNT; verticalLoopIdx++ {
		_, err := GetLetterFromTile(verticalLoopIdx, horizontalTileIdx, tiles)
		if err != nil {
			outerBottomTileOfWord = verticalLoopIdx
			break
		}
	}

	log.Printf("The outer bottom tile is: %d", outerBottomTileOfWord)

	var verticalWordTiles []Tile
	for _, horizontalTiles := range tiles[outerTopTileOfWord:outerBottomTileOfWord] {
		verticalWordTiles = append(verticalWordTiles, horizontalTiles[horizontalTileIdx])
		log.Printf("Append letter %c", horizontalTiles[horizontalTileIdx].Letter.Character)
	}
	return true, verticalWordTiles, outerTopTileOfWord

}

type WordOnBoard struct {
    firstLetterXIdx int
    firstLetterYIdx int
    lastLetterXIdx  int
    lastLetterYIdx  int
    wordTiles       []Tile
}

func (game *Game) GetNewWordsFromBoard() ([]WordOnBoard, error) {
	// Get all unlocked words from the game board

	// Requires:
	// - A reference to a game instance as base variable
	//
	// Guarantees:
	// - Scans across all tiles to find unlocked ones and retrieve
	//   the horizontal and vertical word for each unlocked tile
	//   each unique word will only be accounted for once
	// - Return each unique found word in an array of WordOnBoard structs.
    //   A word is defined as a slice of letters with a length of at least 2.
	// - Return en empty WordOnBoard slice and an error
	//   if no word has been found or if an unlocked tile is not connected
	//   to the center tile through other tiles.

	var newWords []WordOnBoard

	// Each found words will be stored as a string in this variable
	// consisting of a concatenation of its letter ids.
	// This way the uniqueness of a word can be easily determined
	// (since two new words might be the same
	//  but consist of different letter tiles)
	// without the need to execute expensive slice comparisons.
	var newWordsIds []string

	// Get all unlocked tiles
	for verticalIdx, column := range game.Tiles {

        // Each row (i.e. each vertialIdx) can only contain
        // one new horizontal word,
        // So we can save some computing power by stopping to look for
        // horizonal words in the row once one has been found.
        var horizontalWordFoundInRow bool

		for horizontalIdx, tile := range column {
			if tile.Letter != (Letter{}) && !tile.IsLocked {

				if !IsConnectedToCenterTile(verticalIdx, horizontalIdx, game.Tiles, nil) {
					return nil, errors.New(fmt.Sprintf("Tile v:%d,h:%d is isolated from the center tile.", verticalIdx, horizontalIdx))
				}

                var ignoreHorizontalWord, ignoreVerticalWord bool
                horizontalWordFirstLetterXIdx := horizontalIdx
                verticalWordFirstLetterYIdx := verticalIdx
                var hasHorizontalWord, hasVerticalWord bool
                var horizontalWordTiles, verticalWordTiles []Tile

                if ! horizontalWordFoundInRow {
			        hasHorizontalWord, horizontalWordTiles, horizontalWordFirstLetterXIdx = GetHorizontalWordAtTile(verticalIdx, horizontalIdx, game.Tiles)
                } else {
                    ignoreHorizontalWord = true
                }

                hasVerticalWord, verticalWordTiles, verticalWordFirstLetterYIdx = GetVerticalWordAtTile(verticalIdx, horizontalIdx, game.Tiles)

                // Make sure that each word that has been found
                // is only registered i.e. accounted for once
                for _, existingIdConcatenation := range newWordsIds {
                    if existingIdConcatenation == GetIdStringForTileSlice(horizontalWordTiles) {
                        ignoreHorizontalWord = true
                    }
                    if existingIdConcatenation == GetIdStringForTileSlice(verticalWordTiles) {
                        ignoreVerticalWord = true
                    }
                }

                if hasHorizontalWord && ! ignoreHorizontalWord {
                    word := WordOnBoard{
                        firstLetterXIdx: horizontalWordFirstLetterXIdx,
                        firstLetterYIdx: verticalIdx,
                        lastLetterXIdx: horizontalWordFirstLetterXIdx + len(horizontalWordTiles) - 1,
                        lastLetterYIdx: verticalIdx,
                        wordTiles: horizontalWordTiles,
                    }
                    newWords = append(newWords, word)
                    newWordsIds = append(newWordsIds, GetIdStringForTileSlice(horizontalWordTiles))
                }

                if hasVerticalWord && ! ignoreVerticalWord {
                    word := WordOnBoard{
                        firstLetterXIdx: horizontalIdx,
                        firstLetterYIdx: verticalWordFirstLetterYIdx,
                        lastLetterXIdx: horizontalIdx,
                        lastLetterYIdx: verticalWordFirstLetterYIdx + len(verticalWordTiles) - 1,
                        wordTiles: verticalWordTiles,
                    }
                    newWords = append(newWords, word)
                    newWordsIds = append(newWordsIds, GetIdStringForTileSlice(verticalWordTiles))
                }

			}
		}
	}

    if len(newWords) == 0 {
        return newWords, errors.New("No new words found on board.")
    }

    return newWords, nil

}


func IsConnectedToCenterTile(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile, alreadyCheckedTilesMap [][]bool) bool {
	// check whether a given tile on the board is connected through other tiles
	// with letters to the center
	// This is done by recursively following adjacent tiles to the center tile.

	// Requires:
	// - The vertical and horizontal index of the tile to be checked
	//   (parameter 1&2)
	// - The two-dimensional array of all the tiles on the board (parameter 3)
	// - A two dimensional slice with booleans describing which tile
	//   indexes have already been checked by the IsConnectedToCenterTile
	//   recursion. This will prevent infinite recursions.
	//   This slice must have the same dimensions as the tiles slice
	//   and should initially be filled with value "false" throughout.
	//   nil may be passed instead which will let this
	//   function create the slice for itself
	//
	// Guarantees:
	// - Return true if the given tile has a letter on it and at the same time
	//   has a connection to the center tile through a series of other tiles
	//   with letters on them
	// - Return true if the given tile is the center tile
	// - Return false if the given tile is empty (no letter),
	//   if the tile index is invalid or
	//   if there is no possible path of tiles with letters to the center tile.

	if alreadyCheckedTilesMap == nil {
		alreadyCheckedTilesMap = make([][]bool, VERTICAL_TILES_AMOUNT)
		for verticalIdx := 0; verticalIdx < VERTICAL_TILES_AMOUNT; verticalIdx++ {
			alreadyCheckedTilesMap[verticalIdx] = make([]bool, HORIZONTAL_TILES_AMOUNT)
		}
	}

	if ! AreValidBoardCoordinates(verticalTileIdx, horizontalTileIdx) {
		return false
	}

	alreadyCheckedTilesMap[verticalTileIdx][horizontalTileIdx] = true

	if (TileIsCenterTile(verticalTileIdx, horizontalTileIdx)) {
		return true
	}

	_, err := GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
	if err != nil {
		return false
	}

	// Follow all vertically and horizontally adjacent tiles.
	// Return true as soon as a connection has been found
	// through a conneted tile
	if AreValidBoardCoordinates(verticalTileIdx+1, horizontalTileIdx) && !alreadyCheckedTilesMap[verticalTileIdx+1][horizontalTileIdx] && IsConnectedToCenterTile(verticalTileIdx+1, horizontalTileIdx, tiles, alreadyCheckedTilesMap) {
		return true
	}

	if AreValidBoardCoordinates(verticalTileIdx-1, horizontalTileIdx) && !alreadyCheckedTilesMap[verticalTileIdx-1][horizontalTileIdx] && IsConnectedToCenterTile(verticalTileIdx-1, horizontalTileIdx, tiles, alreadyCheckedTilesMap) {
		return true
	}

	if AreValidBoardCoordinates(verticalTileIdx, horizontalTileIdx+1) && !alreadyCheckedTilesMap[verticalTileIdx][horizontalTileIdx+1] && IsConnectedToCenterTile(verticalTileIdx, horizontalTileIdx+1, tiles, alreadyCheckedTilesMap) {
		return true
	}

	if AreValidBoardCoordinates(verticalTileIdx, horizontalTileIdx-1) && !alreadyCheckedTilesMap[verticalTileIdx][horizontalTileIdx-1] && IsConnectedToCenterTile(verticalTileIdx, horizontalTileIdx-1, tiles, alreadyCheckedTilesMap) {
		return true
	}

	return false

}

func IsLegalPlacement(verticalTileIdx int, horizontalTileIdx int, tiles [][]Tile) (bool, string) {
	// Check whether the letter can be placed on a given tile.
	//
	// Guarantees:
	// - Return true if placement is possible, otherwise false
	// - Return a string as second parameter describing why the placement is
	//   illegal, otherwise emptystring.
	// - This function does not consider any letters but only the location
	//   on the board tile

	// If the center tile is empty, it is the only tile
	// onto which a letter can be placed.
	_, err := GetLetterFromTile((VERTICAL_TILES_AMOUNT-1)/2, (HORIZONTAL_TILES_AMOUNT-1)/2, tiles)
	if err != nil {
		if TileIsCenterTile(verticalTileIdx, horizontalTileIdx) {
			return true, ""
		}
		return false, "Center Tile empty. No other placements legal."
	}

	// Make sure the tile on which the letter is to be placed
	// is not already occupied another letter.
	_, err = GetLetterFromTile(verticalTileIdx, horizontalTileIdx, tiles)
	if err == nil {
		return false, "Tile occupied"
	}

	// If there have already been letters placed on the board by the
	// active player, all new letters must be connected to this unlocked tile
	// through one single word.
	hasUnlockedLetters, unlockedLetterVerticalIdx, unlockedLetterHorizontalIdx := HasUnlockedLetters(tiles)
	if hasUnlockedLetters {
		// get the horizontal word
		// if it exists i.e. if it's longer than 1 character
		// new placements are only legal on the same vertical idx
		// (i.e. same horizontal row)
		hasHorizontalWord, _, _ := GetHorizontalWordAtTile(unlockedLetterVerticalIdx, unlockedLetterHorizontalIdx, tiles)
		if hasHorizontalWord && (unlockedLetterVerticalIdx != verticalTileIdx) {
			return false, "Not same vertical index as horizontal word at unlocked tile"
		}

		//Do the same for a vertical adjacent word.
		hasVerticalWord, _, _ := GetVerticalWordAtTile(unlockedLetterVerticalIdx, unlockedLetterHorizontalIdx, tiles)
		if hasVerticalWord && (unlockedLetterHorizontalIdx != horizontalTileIdx) {
			return false, "Not same horizontal index as vertical word at unlocked tile"
		}

	}

	return true, ""

}

func (game *Game) LockLetters() {
	// lock all letters on the board so they cannot be
	// removed by the player anymore
    // Requires
    // - A reference to a game struct as base variable
    // Guarantees:
    // - Set all tiles on the board to isLocked State
    // - Remove tile effects from all letters that are
    //   being locked now so that their effects
    //   will not be accounted for at succeesing words
    //   in which the now locked tiles are involved.
    //   (Each tile effect must only be accounted for the first time
    //    it is part of a new word)

	for verticalIdx, row := range game.Tiles {
		for horizontalIdx, tile := range row {
			if tile.Letter.Character != 0 && ! tile.IsLocked {
                game.Tiles[verticalIdx][horizontalIdx].Effect = NO_TILE_EFFECT
				game.Tiles[verticalIdx][horizontalIdx].IsLocked = true
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
				IsLegalPlacement(verticalIdx, horizontalIdx, game.Tiles)
		}
	}

}

func HasUnlockedLetters(tiles [][]Tile) (bool, int, int) {
	// Check if there are unlocked tiles on the board that have letters on it.
	//
	// Can be uses for placement legality checking when a letter
	// has already been placed by the active player.
	//
	// Requires:
	// - two-dimensional array of board tiles as first parameter
	// Guarantees:
	// - Return true if there is at least one tile on the board
	//   that is not locked and has a letter on it i.e. is not empty and movable
	// - Return false if all of the tiles on the board are locked or if none
	//   of the locked tiles on the board have a letter on them
	// - If true, return also the vertical and horizontal index of the
	//   first found unlocked tile with a letter on it
	// - return -1 for the indexes if false.

	for verticalIdx, tileRow := range tiles {
		for horizontalIdx, tile := range tileRow {
			_, err := GetLetterFromTile(verticalIdx, horizontalIdx, tiles)
			if err == nil && !tile.IsLocked {
				log.Printf("Unlocked letter found. First occurrence at v:%d, h:%d,", verticalIdx, horizontalIdx)
				return true, verticalIdx, horizontalIdx
			} else if err == nil && tile.IsLocked {
				log.Printf("Locked v:%d, h:%d,", verticalIdx, horizontalIdx)
			}
		}
	}

	log.Printf("No unlocked letters on board")
	return false, -1, -1

}
