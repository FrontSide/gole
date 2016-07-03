// Library that interacts with the gole server IPA

var server = {
    url: "http://localhost:8000"
}

var game = {
    id: null,
    board: null,
}

var activePlayer = {
    Name: null,
    Points: null,
    LettersInHand: null
}

function createNewGame(playerNames) {

    // Assign handlers immediately after making the request,
    // and remember the jqxhr object for this request

    //requires:
    // - Array of strings for player Names

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/new",
        data: JSON.stringify({ "PlayerNames": playerNames }),
    })
    .done(function(id) {
        game.id = id
    });

    console.log("New Game ID:" + game.id);


}

function placeLetter(tilesXCoordinate, tilesYCoordinate, letter) {

    letterASCIICode = "a"//get letter ascii code (server expects rune)

}

function getBoard() {
    //Get the complete board with all properties
    //so that the function requesting it in the
    //main module can draw it to the HTML page

    $.ajax({
        async: false,
        method: "GET",
        url: server.url + "/" + game.id + "/board.json",
    })
    .done(function(board) {
        game.board = JSON.parse(board)
    });

}

function getActivePlayer() {
    //Get the struct that describes the player that has their current turn

    $.ajax({
        async: false,
        method: "GET",
        url: server.url + "/" + game.id + "/player.json",
    })
    .done(function(player) {
        activePlayer = JSON.parse(player)
    });

}

function confirmWord() {
    // called after a player has placed
    // all tiles for the current turn
}
