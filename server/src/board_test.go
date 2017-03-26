package main

import (
	"testing"
)

var mockLetterAttributes = []LetterAttributes{
	MockGetLetterAttributes(2, 4),
	MockGetLetterAttributes(1, 9),
}

var mockLetters = []Letter{
	MockGetLetter(123, mockLetterAttributes[0]),
	MockGetLetter(125, mockLetterAttributes[1]),
}

var mockTiles = []Tile{
	MockGetTile(true, mockLetters[0], DOUBLE_LETTER_TILE_EFFECT, false),
	MockGetTile(false, mockLetters[1], DOUBLE_LETTER_TILE_EFFECT, true),
}

func TestTileSliceToStringSuccess(t *testing.T) {

	err := assertEquals("[{\"IsLocked\":true,\"Letter\":{\"Character\":123,\"Attributes\":{\"Occurrences\":2,\"PointValue\":4}},\"Effect\":0,\"PlacementIsLegal\":false},{\"IsLocked\":true,\"Letter\":{\"Character\":123,\"Attributes\":{\"Occurrences\":1,\"PointValue\":9}},\"Effect\":0,\"PlacementIsLegal\":true}]", TileSliceToString(mockTiles))
	if err != nil {
		t.Error(err.Error())
	}

}

func TestTileSliceToStringEmptySlice(t *testing.T) {

	err := assertEquals("[]", TileSliceToString([]Tile{}))
	if err != nil {
		t.Error(err.Error())
	}

}

func TestTileSliceToStringNil(t *testing.T) {

	err := assertEquals("null", TileSliceToString(nil))
	if err != nil {
		t.Error(err.Error())
	}

}
