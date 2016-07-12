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
    GameId string
}

type ConfirmWordRequestBody struct {
    GameId string
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
        http.Error(responseWriter, "Invalid body", 500)
    }

    var game *Game
    game, err = GetGameByUUID(requestBody.GameId)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    err = PlaceLetter(game, requestBody.TileYCoordinate, requestBody.TileXCoordinate, requestBody.Letter)

    if err != nil {
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    responseWriter.Write([]byte(game.Id))

}

func ConfirmWordHandler(responseWriter http.ResponseWriter, request *http.Request) {

    log.Println("Confirming word")

    requestBodyDecoder := json.NewDecoder(request.Body)
    var requestBody ConfirmWordRequestBody
    err := requestBodyDecoder.Decode(&requestBody)
    if err != nil {
        http.Error(responseWriter, "Invalid body", 500)
        return
    }

    var game *Game
    game, err = GetGameByUUID(requestBody.GameId)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    err = FinishTurn(game)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
        return
    }

    responseWriter.Write([]byte("OK"))

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

func StartWebServer() {
    r := mux.NewRouter()
    r.HandleFunc("/new", CreateNewGameHandler).Methods("POST")
    r.HandleFunc("/{id}/board.json", GetBoardHandler).Methods("GET")
    r.HandleFunc("/{id}/player.json", GetActivePlayerHandler).Methods("GET")
    r.HandleFunc("/place", PlaceLetterHandler).Methods("POST")
    r.HandleFunc("/confirm", ConfirmWordHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
