package main

import (
    "net/http"
    "log"
)

http.HandleFunc("/board.json", GetBoard)


func GetBoard(responseWriter http.ResponseWriter, request *http.Request)

func BoardAsJson(tiles [][]Tile) string {
    // Converts a given two-dimensional
    // slice of tiles into a json representation
    // of the board with all the information the tile
    // struct has as well.
    boardJson, err := json.Marshal(tiles)
    if err != nil {
        log.Fatal(err)
    }
    return boardJson
}
