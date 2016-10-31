package main

import (
	"errors"
	"fmt"
	"gole/golelibs"
	"log"
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

func GetActivePlayer(game *Game) (*Player, error) {
	// Return the struct of the active player
	// Guarantees:
	// - Return the Player struct of the player with the index
	//   that correlates to the index stored in the PlayerIdxWithTurn
	//   variable in the given game struct

	log.Println("Get Player with index " + string(game.PlayerIdxWithTurn))
	if game.PlayerIdxWithTurn < len(game.Players) && game.PlayerIdxWithTurn >= 0 {
		return &game.Players[game.PlayerIdxWithTurn], nil
	}
	return &Player{}, errors.New(fmt.Sprintf("Player with index %d is not available.", game.PlayerIdxWithTurn))

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

type PotentialPointsForWord struct {
    FirstLetterXIdx  int
    FirstLetterYIdx  int
    LastLetterXIdx  int
    LastLetterYIdx  int
    PotentialPoints int
}

type PotentialPointsForWords []*PotentialPointsForWord

func GetPotentialPoints(game *Game) (PotentialPointsForWords, error) {
    // Get the potential points a player could gain for an unplayed but
    // placed word on the board.
    //
    // Guarantees:
    // - Return a list/slice of PotentialPointsForWord structs.
    //   This struct contains the last indexes
    //   (i.e. the coordinates of the word's last tile/letter)
    //   of each unconfirmed word
    //   (i.e. a word that has unlocked tiles/letteres and for which a player
    //    has not redeemed any points)
    //   that has been found on the board and the points that a player
    //   would gain for this word if s/he ended the turn now.
    //   Note that potential points do not consider the validity of
    //   the word nor the legality of the placements of all letters.
    // - Return nil and an error in case of failure

    newWordsOnBoard, err := game.GetNewWordsFromBoard()
    if err != nil {
        return nil, err
    }

    log.Printf("\nNr. Of new words found: %d\n", len(newWordsOnBoard))

    var potentialPointsForWords PotentialPointsForWords

    for _, wordOnBoard := range newWordsOnBoard {
        pointsForWord, _, err := GetPointsForWord(wordOnBoard, false)
        if err != nil {
            return nil, err
        }

        potentialPointsForWord := PotentialPointsForWord{
            LastLetterYIdx: wordOnBoard.lastLetterYIdx,
            LastLetterXIdx: wordOnBoard.lastLetterXIdx,
            FirstLetterXIdx: wordOnBoard.firstLetterXIdx,
            FirstLetterYIdx: wordOnBoard.firstLetterYIdx,
            PotentialPoints: pointsForWord,
        }

        potentialPointsForWords = append(potentialPointsForWords, &potentialPointsForWord)
    }

    return potentialPointsForWords, nil

}

func PlaceLetter(game *Game, verticalTileIdx int,
			horizontalTileIdx int, letterId string) error {
	// Add a letter to the board.
	//
	// Guarantees:
	// - If successful, the affected letter will be moved away from the
	//   active player's hand and put on the specified board tile
	// - Return nil and error if:
    //   -- placement is illeal
    //   -- the game is over
	//   -- the letter with the given ID is a wildcard letter, that
	//      has not yet been replaced with an actual letter
	//   -- the active player does not own the letter that is to be placed

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

func GetPointsForWord(wordOnBoard WordOnBoard, doCheckVailidity bool) (int, string, error) {

	// Calculate the points for a series of
	// tiles, with respect to the point value of a letter
	// and the tile effects
	// Requires:
	// - Slice of tiles with letters on them
    // - A boolean describing whether the word should be
    //   checked for vailidity against a dictionary
	// Guarantees:
	// - Return the points gained with this word (if valid)
	// - Return the word from the tiles as a string
	// - Return an error if the word is invalid if doCheckVailidity
    //   is set to true

	if len(wordOnBoard.wordTiles) < 2 {
		return -1, "", errors.New(
			"Can not get points for word. Too short.")
	}

	var word string
	var wordPoints int
	wordPointMultiplicator := 1

	for _, tile := range wordOnBoard.wordTiles {

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

	if doCheckVailidity && !golelibs.IsAValidWord(word) {
		return -1, word, errors.New("No t a valid word: " + word)
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
    //   This inclused the case that no new words were found on the board.

	// Stores the words that have been successfully confirmed
	// in this round
	var confirmedWords []string

    // Stores the points that have been gaines in this round.
    // After all words have been confirmed,
    // these points will be accounted to the player.
	var points int

    newWordsOnBoard, err := game.GetNewWordsFromBoard()

    if err != nil {
        return -1, nil, err
    }

    if len(newWordsOnBoard) == 0 {
        return -1, nil, errors.New("No new words found on board.")
    }

    for _, wordOnBoard := range newWordsOnBoard {
        pointsForWord, newConfirmdWord, err := GetPointsForWord(wordOnBoard, true)
        if err != nil {
            return -1, nil, err
        }
        points += pointsForWord
        confirmedWords = append(confirmedWords, newConfirmdWord)
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
