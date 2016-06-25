package main

import (
    "github.com/gorilla/mux"
    "log"
)

func init() {
    r := mux.NewRouter()
    r.HandleFunc("/new", CreateNewGame).Methods("POST")
    r.HandleFunc("/board.json", GetBoard).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", r))
}

func CreateNewGameHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func GetBoardHandler(responseWriter http.ResponseWriter, request *http.Request) {
    var game = GetGameByUUID()

    // 1. Get tiles for game id
    // 2. Get board json for tiles
    // 3. return jsons
}
