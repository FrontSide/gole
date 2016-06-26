package main

import (
    "net/http"
    "github.com/gorilla/mux"
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


    requestBodyDecoder := json.NewDecoder(request.Body)
    var requestBody CreateNewGameRequestBody
    err := requestBodyDecoder.Decode(&requestBody)
    if err != nil {
        http.Error(responseWriter, "Invalid body", 500)
    }

    responseWriter.Write([]byte(StartNewGame(requestBody.PlayerNames...)))

}


func GetBoardHandler(responseWriter http.ResponseWriter, request *http.Request){
    id := mux.Vars(request)["id"]
    game, err := GetGameByUUID(id)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
    }

    println(GetBoardAsJson(game))

    responseWriter.Write([]byte(GetBoardAsJson(game)))
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
    }

    err = PlaceLetter(game, requestBody.TileYCoordinate, requestBody.TileXCoordinate, requestBody.Letter)

    if err != nil {
        http.Error(responseWriter, err.Error(), 500)
    }

}

func ConfirmWordHandler(responseWriter http.ResponseWriter, request *http.Request) {

    requestBodyDecoder := json.NewDecoder(request.Body)
    var requestBody ConfirmWordRequestBody
    err := requestBodyDecoder.Decode(&requestBody)
    if err != nil {
        http.Error(responseWriter, "Invalid body", 500)
    }

    var game *Game
    game, err = GetGameByUUID(requestBody.GameId)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
    }

    var points int
    points, err = FinishTurn(game)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
    }
}

func GetPointsOfPlayerHandler(responseWriter http.ResponseWriter, request *http.Request) {
    id := mux.Vars(request)["id"]
    playerName := mux.Vars(request)["playerName"]
    game, err := GetGameByUUID(id)

    if err != nil {
        log.Println("Available ids are: ", games)
        http.Error(responseWriter, err.Error(), 500)
    }

    responseWriter.Write([]byte(GetBoardAsJson(game)))

}

func StartWebServer() {
    r := mux.NewRouter()
    r.HandleFunc("/new", CreateNewGameHandler).Methods("POST")
    r.HandleFunc("/{id}/board.json", GetBoardHandler).Methods("GET")
    r.HandleFunc("/place", PlaceLetterHandler).Methods("POST")
    r.HandleFunc("/confirm", ConfirmWordHandler).Methods("POST")
    r.HandleFunc("/{id}/{playerName}/points", GetPointsOfPlayerHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", r))
}
