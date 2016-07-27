// Library that interacts with the gole server IPA

var server = {
    url: "http://localhost:8000"
}

var game = {
    id: null,
    board: null,
    // The player names will only be needed at the beginning when they are
    // passed over by the user via the new game form
    // Needs to be an Array of strings (2-4 players are allowed)
    playerNames: new Array(),
}

var activePlayer = {
    Name: null,
    Points: null,
    LettersInHand: null
}

//function Tile(letter) {
//    this.letter = letter
//    this.isActivated = false
//    this.isLocked = false
//}

function createNewGame() {

    // Assign handlers immediately after making the request,
    // and remember the jqxhr object for this request

    //requires:
    // - That the array at game.playerNames is set as a string array
    //   defining the names of the players for the new game

    console.log("Start Game with players: " + game.playerNames)

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/new",
        data: JSON.stringify({ "PlayerNames": game.playerNames }),
    })
    .done(function(id) {
        game.id = id
        Cookies.set('golegameid', game.id);
    });

    console.log("New Game ID:" + game.id);


}

function placeLetter(tilesXCoordinate, tilesYCoordinate, letter) {

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/place",
        data: JSON.stringify(
            {
                "TileXCoordinate": tilesXCoordinate,
                "TileYCoordinate": tilesYCoordinate,
                "Letter": letter,
                "GameId": game.id
            }
        ),
    })
    .done(function(id) {
        console.log("Letter Placed")
    });

}

function removeLetter(tilesXCoordinate, tilesYCoordinate) {

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/remove",
        data: JSON.stringify(
            {
                "TileXCoordinate": tilesXCoordinate,
                "TileYCoordinate": tilesYCoordinate,
                "GameId": game.id
            }
        ),
    })
    .done(function(id) {
        console.log("Letter Removed")
    });

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

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/confirm",
        data: JSON.stringify(
            {
                "GameId": game.id
            }
        ),
    })
    .done(function(id) {
        console.log("Word confirmed");
        return true;
    });
    return false;

}
