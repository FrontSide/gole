// This file offers
// helper functions for mocking
// parts of the application
// and for running test in general

package main

import (
	"errors"
	"fmt"
)

func assertEquals(expected interface{}, real interface{}) error {
	if expected != real {
		return errors.New(fmt.Sprintf("Expected %s, Was: %s", expected, real))
	}
	return nil
}

func MockGetLetterAttributes(occurrences int, pointValue int) LetterAttributes {
	return LetterAttributes{
		Occurrences: occurrences,
		PointValue:  pointValue,
	}
}

func MockGetLetter(character rune, letterAttributes LetterAttributes) Letter {
	return Letter{
		Character:  123,
		Attributes: letterAttributes,
	}
}

func MockGetTile(isLocked bool, letter Letter, effect SpecialTileEffect, isPlacementLegal bool) Tile {
	return Tile{
		IsLocked:         true,
		Letter:           letter,
		Effect:           effect,
		PlacementIsLegal: isPlacementLegal,
	}
}
