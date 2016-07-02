// Library that interacts with the gole server IPA

var server = {
    url: "http://localhost:8000"
}

var game = {
    id: null,
    board: null,
    legalPlacements: null,
}

var activePlayer = {
    name: null,
    points: null,
    hand: null
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

    console.log("Game Board:" + game.board);

}

function getLegalPlacementMap() {
    //Get the two dimentional array
    //that contains bools to describe whether
    //a letter can be placed at the tile with given inices tile

        $.ajax({
            async: false,
            method: "GET",
            url: server.url + "/" + game.id + "/legalplacements.json",
        })
        .done(function(legalPlacements) {
            game.legalPlacements = JSON.parse(legalPlacements)
        });

        console.log("Placements Map:" + game.legalPlacements);

}

function getPlayerHand() {
    //Get the two dimentional array
    //that contains bools to describe whether
    //a letter can be placed at the tile with given inices tile

    $.ajax({
        async: false,
        method: "GET",
        url: server.url + "/" + game.id + "/" + activePlayer.name + "/hand",
    })
    .done(function(hand) {
        activePlayer.hand = hand
    });

}

function confirmWord() {
    // called after a player has placed
    // all tiles for the current turn
}
