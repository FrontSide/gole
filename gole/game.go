package main

import (
	"errors"
	"fmt"
	"gole/golelibs"
	"log"
	"reflect"
	"strings"
)

type Game struct {
	Id      string
	Players []Player

	// hold index of player who has the current turn
	PlayerIdxWithTurn int

	// 2-dimensional array storing all tiles
	// first index represents vertical, second index horizontal tiles
	Tiles [][]Tile

	// All letters in the backlog
	// i.e. that have not yet been handed to a player
	LetterSet []Letter

	// Flag indicating whether the game is over
	// will be fale until one player has no more letters in hand
	// and the letter backlog is empty
	GameOver bool
}

var MIN_NUMBER_OF_PLAYERS = 2
var MAX_NUMBER_OF_PLAYERS = 4
var MAX_NUMBER_OF_LETTERS_IN_HAND = 7

func (game *Game) GetPlayerByName(playerName string) (Player, error) {
	// get a player from the given game by their name
	for _, existingPlayer := range game.Players {
		if existingPlayer.Name == strings.TrimSpace(playerName) {
			return existingPlayer, nil
		}
	}
	return Player{}, errors.New(
		"Player with name does not exist in the game. " + playerName)
}

func (game *Game) GetScoreBoard() map[string]int {
	// Retrun a ScoreBoard Map with all player's names as Keys
	// and their game points as associated value

	var scoreBoard = make(map[string]int)

	for _, player := range game.Players {
		log.Println(player.Name)
		scoreBoard[player.Name] = player.Points
	}

	return scoreBoard
}

func AddPlayer(playerName string, game *Game) error {
	// Add a player to the list of players for the
	// upcoming game play

	if len(game.Players) >= MAX_NUMBER_OF_PLAYERS {
		return errors.New("No more players can be added to the Game.")
	}

	_, err := game.GetPlayerByName(playerName)
	if err == nil {
		return errors.New("A player with this name already exists.")
	}

	player := Player{Name: playerName}

	for i := 0; i < MAX_NUMBER_OF_LETTERS_IN_HAND; i++ {
		nextLetter, err := PopLetterFromSet(game)
		if err != nil {
			return err
		}
		err = player.AddLetterToHand(nextLetter)
		if err != nil {
			return err
		}
	}

	game.Players = append(game.Players, player)

	return nil

}

func GetActivePlayer(game *Game) (Player, error) {
	// Return the struct of the active player
	// Guarantees:
	// - Return the Player struct of the player with the index
	//   that correlates to the index stored in the PlayerIdxWithTurn
	//   variable in the given game struct

	log.Println("Get Player with index " + string(game.PlayerIdxWithTurn))
	if game.PlayerIdxWithTurn < len(game.Players) && game.PlayerIdxWithTurn >= 0 {
		return game.Players[game.PlayerIdxWithTurn], nil
	}
	return Player{}, errors.New(fmt.Sprintf("Player with index %d is not available.", game.PlayerIdxWithTurn))

}

func PopLetterFromSet(game *Game) (Letter, error) {
	// Pop the last letter (right end) from the
	// letter string of the passed game structure instance.
	//
	// Guarantees:
	// - Next letter from letter backlog string will be returned if exists
	// - Returned letter will be removed from letter backlog string
	// - Returns an empty letter and error
	//   if no letter is left in letter backlog

	if len(game.LetterSet) < 1 {
		return Letter{}, errors.New(
			"Cannot pop letter from set. Empty.")
	}
	var letterToReturn = game.LetterSet[len(game.LetterSet)-1]
	game.LetterSet = game.LetterSet[:len(game.LetterSet)-1]
	return letterToReturn, nil
}

func PlaceLetter(game *Game, verticalTileIdx int,
			horizontalTileIdx int, letterId string) error {
	// Add a letter to the board.
	//
	// Guarantees:
	// - If successful, the affected letter will be moved away from the
	//   active player's hand and put on the specified board tile
	// - Return error if placement is illeal, if the game is over,
	//   if the letter with the given ID is a wildcard letter, that
	//   has not yet been replaced with an actual letter,
	//   or if the active player does not own the letter that is to be placed

	isRawWildcardLetter, err := game.Players[game.PlayerIdxWithTurn].IsRawWildcardLetter(letterId)
	if err != nil {
		return err
	}
	if isRawWildcardLetter {
		return errors.New("Cannot place unsubstituted wildcard letter.")
	}

	if game.GameOver {
		return errors.New("Cannot place letter. Game is over.")
	}

	if isLegal, reason := IsLegalPlacement(
		verticalTileIdx, horizontalTileIdx, game.Tiles); !isLegal {
		return errors.New("Cannot place letter. " + reason)
	}

	var letterStruct Letter
	letterStruct, err = game.Players[game.PlayerIdxWithTurn].PopLetterFromHand(letterId)
	if err != nil {
		return err
	}

	game.Tiles[verticalTileIdx][horizontalTileIdx].Letter = letterStruct
	game.UpdatePlacementLegalityOfAllTiles()

	return nil
}

func RemoveLetter(game *Game, verticalTileIdx int, horizontalTileIdx int) error {
	// Remove one single letter from the board that has
	// not been locked yet

	_, err := GetLetterFromTile(
		verticalTileIdx, horizontalTileIdx, game.Tiles)

	if err != nil {
		return errors.New("Cannot remove letter. Tile empty")
	}

	if game.Tiles[verticalTileIdx][horizontalTileIdx].IsLocked {
		return errors.New("Cannot remove letter. Tile Locked")
	}

	// Hand letter back to player
	err = game.Players[game.PlayerIdxWithTurn].AddLetterToHand(
		game.Tiles[verticalTileIdx][horizontalTileIdx].Letter)

	if err != nil {
		return err
	}

	// Overwrite letter on  tile with empty letter struct
	game.Tiles[verticalTileIdx][horizontalTileIdx].Letter = Letter{}

	// Update placement legality of whole board
	game.UpdatePlacementLegalityOfAllTiles()

	return nil

}

func GetPointsForWord(wordTiles []Tile) (int, string, error) {

	// Calculate the points for a series of
	// tiles, with respect to the point value of a letter
	// and the tile effects
	// Requires:
	// - Slice of tiles with letters on them
	// Guarantees:
	// - Return the points gained with this word (if valid)
	// - Return the word from the tiles as a string
	// - Return an error if the word is invalid

	if len(wordTiles) < 2 {
		return -1, "", errors.New(
			"Can not get points for word. Too short.")
	}

	var word string
	var wordPoints int
	wordPointMultiplicator := 1

	for _, tile := range wordTiles {

		word += string(tile.Letter.Character)

		var letterPoints = tile.Letter.Attributes.PointValue
		if tile.Effect == DOUBLE_LETTER_TILE_EFFECT {
			letterPoints *= 2
		} else if tile.Effect == TRIPLE_LETTER_TILE_EFFECT {
			letterPoints *= 3
		} else if tile.Effect == DOUBLE_WORD_TILE_EFFECT {
			wordPointMultiplicator += 1
		} else if tile.Effect == TRIPLE_WORD_TILE_EFFECT {
			wordPointMultiplicator += 2
		}
		wordPoints += letterPoints

	}

	log.Println("Word to check: " + word)

	wordPoints *= wordPointMultiplicator

	if !golelibs.IsAValidWord(word) {
		return -1, word, errors.New("Not a valid word: " + word)
	}

	return wordPoints, word, nil

}

func FinishTurn(game *Game) (int, []string, error) {

	// Tiles that have already been respected for point calculation
	// Tiles that already were were locked before this current turn
	// may only be taken into account for rating once.
	// Tiles that the player placed in this turn (i.e. currently unlocked)
	// may be counted twice, if they were connected to two different tiles
	// Guarantees:
	// - If turn was successful, return the points gained,
	//   an array with the word(s) for which the points were awarded
	//   and nil for error
	// - If turn was unsuccessful, return -1, nil and the error

	var confirmedWordTiles [][]Tile
	var confirmedWords []string
	var points int

	// Get all unlocked tiles
	for verticalIdx, column := range game.Tiles {
		for horizontalIdx, tile := range column {
			if tile.Letter != (Letter{}) && !tile.IsLocked {

				if !IsConnectedToCenterTile(verticalIdx, horizontalIdx, game.Tiles, nil) {
					return -1, nil, errors.New(fmt.Sprintf("Tile v:%d,h:%d is isolated from the center tile.", verticalIdx, horizontalIdx))
				}

				hasHorizontalWord, horizontalWordTiles, _ := GetHorizontalWordAtTile(verticalIdx, horizontalIdx, game.Tiles)
				hasVerticalWord, verticalWordTiles, _ := GetVerticalWordAtTile(verticalIdx, horizontalIdx, game.Tiles)

				log.Println("horizontailTiles: " + TileSliceToString(horizontalWordTiles))
				log.Println("verticalTiles: " + TileSliceToString(verticalWordTiles))

				//Check if words have alreaddy been confirmed i.e. rated
				var ignoreHorizontalWord, ignoreVerticalWord bool
				for _, wordTiles := range confirmedWordTiles {
					if hasHorizontalWord && reflect.DeepEqual(wordTiles, horizontalWordTiles) {
						ignoreHorizontalWord = true
					}
					if hasVerticalWord && reflect.DeepEqual(wordTiles, verticalWordTiles) {
						ignoreVerticalWord = true
					}
				}

				if !ignoreHorizontalWord && hasHorizontalWord {
					horizontalWordPoints, word, err := GetPointsForWord(horizontalWordTiles)
					if err != nil {
						return -1, nil, err
					}
					points += horizontalWordPoints
					confirmedWords = append(confirmedWords, word)
					confirmedWordTiles = append(confirmedWordTiles, horizontalWordTiles)
				}

				if !ignoreVerticalWord && hasVerticalWord {
					verticalWordPoints, word, err := GetPointsForWord(verticalWordTiles)
					if err != nil {
						return -1, nil, err
					}
					points += verticalWordPoints
					confirmedWords = append(confirmedWords, word)
					confirmedWordTiles = append(confirmedWordTiles, verticalWordTiles)
				}

				// If the center tile is unlocked
				// make sure that there is a horizontal
				// or vertical word around it.
				// Submitting a single letter as first word
				// is not permissive.
				if TileIsCenterTile(verticalIdx, horizontalIdx) && !hasVerticalWord && !hasHorizontalWord {
					return -1, nil, errors.New("You need to place at least one more letter on the board.")
				}

			}
		}
	}

	// Add earned points to current player
	game.Players[game.PlayerIdxWithTurn].Points += points

	// Fill up player hand with new letters
	numberOflettersToAdd := MAX_NUMBER_OF_LETTERS_IN_HAND - len(game.Players[game.PlayerIdxWithTurn].LettersInHand)
	for addLetterCounter := 0; addLetterCounter < numberOflettersToAdd; addLetterCounter++ {
		newLetter, err := PopLetterFromSet(game)
		if err != nil {
			break
		}

		err = game.Players[game.PlayerIdxWithTurn].AddLetterToHand(newLetter)
		if err != nil {
			return -1, nil, err
		}
	}

	// If the player hand is empty at this stage.
	// The game is considered over as at least one player has no letters left
	// anymore.
	if len(game.Players[game.PlayerIdxWithTurn].LettersInHand) < 1 {
		game.GameOver = true
	}

	game.LockLetters()
	game.UpdatePlacementLegalityOfAllTiles()

	// Give turn to next player
	game.PlayerIdxWithTurn = (game.PlayerIdxWithTurn + 1) % len(game.Players)
	log.Printf("Index of player with turn is now: %d", game.PlayerIdxWithTurn)
	return points, confirmedWords, nil
}
