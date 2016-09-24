package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    "encoding/json"
    "log"
)

type CreateNewGameRequestBody struct {
    PlayerNames []string
}

type PlaceLetterRequestBody struct {
    TileXCoordinate int
    TileYCoordinate int
    Letter rune
    IsWildcard bool
    GameId string
}

type RemoveLetterRequestBody struct {
    TileXCoordinate int
    TileYCoordinate int
    GameId string
}

type ConfirmWordRequestBody struct {
    GameId string
}

type ConfirmWordResponse struct {
    GainedPoints int
    Words []string
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


func GetBoardHandler(responseWriter http.ResponseWriter, request *http.Request){

    id := mux.Vars(request)["id"]

    var err error
    var game *Game
    game, err = GetGameByUUID(id)

    if err != nil {
        log.Println("Available ids are: ", games)
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

func PlaceLetterHandler(responseWriter http.ResponseWriter, request *http.Request){
    requestBodyDecoder := json.NewDecoder(request.Body)
    var requestBody PlaceLetterRequestBody
    err := requestBodyDecoder.Decode(&requestBody)
    if err != nil {
        http.Error(responseWriter, err.Error(), 500)
    }

    var game *Game
    game, err = GetGameByUUID(requestBody.GameId)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    err = PlaceLetter(game, requestBody.TileYCoordinate, requestBody.TileXCoordinate, requestBody.Letter, requestBody.IsWildcard)

    if err != nil {
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    responseWriter.Write([]byte(game.Id))

}

func RemoveLetterHandler(responseWriter http.ResponseWriter, request *http.Request){
    requestBodyDecoder := json.NewDecoder(request.Body)
    var requestBody RemoveLetterRequestBody
    err := requestBodyDecoder.Decode(&requestBody)
    if err != nil {
        http.Error(responseWriter, "Invalid body", 500)
    }

    var game *Game
    game, err = GetGameByUUID(requestBody.GameId)

    if err != nil {
        log.Println("Available ids are: ", games)
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
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), HTTP_ERROR_CODE)
        return
    }

    confirmWordResponse := ConfirmWordResponse{}
    confirmWordResponse.GainedPoints, confirmWordResponse.Words, err = FinishTurn(game)

    if err != nil {
        log.Println("Available ids are: ", games)
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
        log.Println("Available ids are: ", games)
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

    var player Player
    log.Println("Get Player with index " + string(game.PlayerIdxWithTurn))
    if game.PlayerIdxWithTurn < len(game.Players) && game.PlayerIdxWithTurn >= 0 {
        player = game.Players[game.PlayerIdxWithTurn]
    } else {
        log.Println("Player with index " + string(game.PlayerIdxWithTurn) + " is not available.")
        http.Error(responseWriter, "Error when trying to retrieve player.", 500)
        return
    }

    var playerJson []byte
    playerJson, err = json.Marshal(player)
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
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
        return
    }

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
    r.HandleFunc("/place", PlaceLetterHandler).Methods("POST")
    r.HandleFunc("/remove", RemoveLetterHandler).Methods("POST")
    r.HandleFunc("/confirm", ConfirmWordHandler).Methods("POST")
    r.HandleFunc("/{id}/scoreboard.json", GetScoreBoardHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
