package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type CreateNewGameRequestBody struct {
	PlayerNames []string
}

type SortHandRequestBody struct {
	LetterIds []string
	GameId string
}

type ReplaceWildcardRequestBody struct {
	LetterId string
	ReplacementLetter rune
	GameId string
}

type PlaceLetterRequestBody struct {
	TileXCoordinate int
	TileYCoordinate int
	LetterId        string
	IsWildcard      bool
	GameId          string
}

type RemoveLetterRequestBody struct {
	TileXCoordinate int
	TileYCoordinate int
	GameId          string
}

type ConfirmWordRequestBody struct {
	GameId string
}

type ConfirmWordResponse struct {
	GainedPoints int
	Words        []string
}

func CreateNewGameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	//Requires:
	// - stringified json obect with one key 'PlayerNames'
	//   that has an array of strings as value, with the player names
	// Guarantees:
	// - String response with new game ID

	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody CreateNewGameRequestBody
	var err error
	err = requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	var gameId string
	gameId, err = StartNewGame(requestBody.PlayerNames...)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write([]byte(gameId))

}

func GetBoardHandler(responseWriter http.ResponseWriter, request *http.Request) {

	id := mux.Vars(request)["id"]

	var err error
	var game *Game
	game, err = GetGameByUUID(id)

	if err != nil {
		log.Println("Not a valid GameID: ", id)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	var boardJson []byte
	boardJson, err = json.Marshal(game.Tiles)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write(boardJson)
}


func SortHandHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Sort or shuffle the hand of the active player
	// Requires:
	// - An incoming HTTP Request Body with values to all keys
	//   as they are defined in the SortHandRequestBody struct
	//   (matching key name, valid data type) whereas the LetterIds
	//   fields may have an empty value. If the value is non empty,
	//   it must contain a slice of letterIDs of letters from the
	//   active players hand. The amount of letterIDs must match
	//   the actual number of letters in the active players hand.
	// Guarantees:
	// - If the LetterIDs value is empty, the letters in the active players
	//   hands are shuffled randomly and stored back to the player's hand.
	// - If the LatterIDs value is valid and non-empty, the letters in
	//   the hand of the active player are sorted according to the
	//   array and stored to the player's hand accordingly.
	// - If successful, HTTP 200 and the gameId is returned
	// - If there is an error in either the request handler or the
	//   game logic, HTTP 500 and the error message is returned.

	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody SortHandRequestBody
	err := requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
	}

	var game *Game
	game, err = GetGameByUUID(requestBody.GameId)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	var activePlayer *Player
	activePlayer, err = GetActivePlayer(game)
	if err != nil {
		http.Error(responseWriter, "Error when trying to retrieve player.", 500)
		return
	}

	if requestBody.LetterIds != nil {
		err = activePlayer.SortHand(requestBody.LetterIds)
	} else {
		activePlayer.ShuffleHand()
	}

	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write([]byte(requestBody.GameId))

}


func ReplaceWildcardHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Handle the client request to replace a wildcard letter tile with an
	// actuall letter
	// Requires:
	// - An incoming HTTP Request Body with values to all keys
	//   as they are defined in the ReplaceWildcardRequestBody struct
	//   (matching key name, valid data type)
	// Guarantees:
	// - Call the ReplaceWildcard function that will replace the wildcard
	//   character on a tile with an actual letter
	// - Will return HTTP 200 and the id of the letter struct
	//   if the replacement was successful
	// - Will return HTTP 500 if there has been an error either in the
	//   request handler function or the game loggic.

	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody ReplaceWildcardRequestBody
	err := requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
	}

	var game *Game
	game, err = GetGameByUUID(requestBody.GameId)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	var activePlayer *Player
	activePlayer, err = GetActivePlayer(game)
	if err != nil {
		http.Error(responseWriter, "Error when trying to retrieve player.", 500)
		return
	}

	err = activePlayer.ReplaceWildcard(requestBody.LetterId, requestBody.ReplacementLetter)

	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write([]byte(requestBody.LetterId))

}

func PlaceLetterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Handle the clien request to place a letter on the board
	// Requires:
	// - An incoming HTTP Request Body with values to all keys
	//   as they are defined in the PlaceLetterRequestBody struct
	//   (matching key name, valid data type)
	// Guarantees:
	// - Call the PlaceLetter function that will handle the
	//   game logic of placing a letter from the player
	//   hand on a board tile
	// - Will return with code 200 and the GameID if the letter was
	//   placed successfully
	// - Will respond with code 500 if there has been an error in either
	//   the HTTP request handler function or the game logic function.

	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody PlaceLetterRequestBody
	err := requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
	}

	var game *Game
	game, err = GetGameByUUID(requestBody.GameId)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	err = PlaceLetter(game, requestBody.TileYCoordinate,
		requestBody.TileXCoordinate, requestBody.LetterId)

	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write([]byte(game.Id))

}

func RemoveLetterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody RemoveLetterRequestBody
	err := requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, "Invalid body", 500)
	}

	var game *Game
	game, err = GetGameByUUID(requestBody.GameId)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	err = RemoveLetter(game, requestBody.TileYCoordinate, requestBody.TileXCoordinate)

	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write([]byte(game.Id))

}

func ConfirmWordHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Trigger the function to finish up the round after a player has placed all
	// letters for this round.
	// Requires:
	// - GameId in Request Body
	// Guarantees:
	// - HTTP 500 response if an error occured
	// - HTTP 200 if the turn has been finished successfully and the next player
	//   can continue with the game.
	// - HTTP 250 if the turn has been finished successfully and if the game is
	//   now over.
	// - In the case of a successful word confirmation,
	//   A json structure containing the gained points and
	//   the word(s) with for which the points have been awarded
	//   are returned
	//   Structure:
	//   { GainedPoints: int, Words: []string }

	HTTP_GAME_OVER_CODE := 250
	HTTP_DEFAULT_CODE := 200
	HTTP_ERROR_CODE := 500

	log.Println("Confirming word")

	requestBodyDecoder := json.NewDecoder(request.Body)
	var requestBody ConfirmWordRequestBody
	err := requestBodyDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(responseWriter, "Invalid body", HTTP_ERROR_CODE)
		return
	}

	var game *Game
	game, err = GetGameByUUID(requestBody.GameId)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), HTTP_ERROR_CODE)
		return
	}

	confirmWordResponse := ConfirmWordResponse{}
	confirmWordResponse.GainedPoints, confirmWordResponse.Words, err = FinishTurn(game)

	if err != nil {
		log.Println("Not a valid GameID: ", requestBody.GameId)
		http.Error(responseWriter, err.Error(), HTTP_ERROR_CODE)
		return
	}

	confirmWordResponseJson, err := json.Marshal(confirmWordResponse)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	if game.GameOver {
		responseWriter.WriteHeader(HTTP_GAME_OVER_CODE)
	} else {
		responseWriter.WriteHeader(HTTP_DEFAULT_CODE)
	}

	responseWriter.Write(confirmWordResponseJson)

}

func GetActivePlayerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]

	var err error
	game, err := GetGameByUUID(id)

	if err != nil {
		log.Println("Not a valid GameID: ", id)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	var playerList []byte
	playerList, err = json.Marshal(game.Players)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}
	log.Println(string(playerList))

	var activePlayer *Player
	activePlayer, err = GetActivePlayer(game)
	if err != nil {
		http.Error(responseWriter, "Error when trying to retrieve player.", 500)
		return
	}

	log.Println("The active player's hand::")
	log.Println(activePlayer.LettersInHand)

	var playerJson []byte
	playerJson, err = json.Marshal(activePlayer)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Write(playerJson)

}

func GetScoreBoardHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Return a json object describing the players of the game with the given
	// is and their game points.
	// Keys in the returned object will be the player name
	// with the points as the value

	id := mux.Vars(request)["id"]

	var err error
	game, err := GetGameByUUID(id)

	if err != nil {
		log.Println("Not a valid GameID: ", id)
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	log.Println("ok")

	var scoreBoard []byte
	scoreBoard, err = json.Marshal(game.GetScoreBoard())
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}
	responseWriter.Write(scoreBoard)
}

func StartWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/new", CreateNewGameHandler).Methods("POST")
	r.HandleFunc("/{id}/board.json", GetBoardHandler).Methods("GET")
	r.HandleFunc("/{id}/player.json", GetActivePlayerHandler).Methods("GET")
	r.HandleFunc("/wildcard/replace", ReplaceWildcardHandler).Methods("POST")
	r.HandleFunc("/hand/sort", SortHandHandler).Methods("POST")
	r.HandleFunc("/place", PlaceLetterHandler).Methods("POST")
	r.HandleFunc("/remove", RemoveLetterHandler).Methods("POST")
	r.HandleFunc("/confirm", ConfirmWordHandler).Methods("POST")
	r.HandleFunc("/{id}/scoreboard.json", GetScoreBoardHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
