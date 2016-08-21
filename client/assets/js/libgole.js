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
    playerNames: new Array()
}

var activePlayer = {
    Name: null,
    Points: null,
    LettersInHand: null
}
// Active Game Scoreboard
// playerName: points
var scoreboard = {}

// stores information about the currently activated tile
// i.e. the tile that is to be moved
var activatedLetter = null

// If a letter is being dragged away from the board,
// the origin indexes of that tile need to be stored
var removeLetterOrigin = {
    verticalIdx: null,
    horizontalIdx: null
}

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
    // Guarantees:
    // - Send letterPlacement request to gole server
    // - Return null if operation was successfult and server returned with HTTP ok
    // - Return error message if operation was unsuccesful

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
        return null
    })
    .fail(function(response) {
        return response
    })

}

function removeLetter(fSuccessCallback, fErrorCallback) {
    // Requires:
    // - The x and y coordinates of the tile from which the letter
    //   is to be removed to be set in the global
    //   removeLetterOrigin.horizontalIdx and removeLetterOrigin.verticalIdx
    //   variables.
    // - A callback function fSuccessCallback
    //   to be called if the request succeeded
    // - A callback function fErrorCallback
    //   to be called if the request failed
    // Guarantees:
    // - Send request to remove letter to gole server
    // - Trigger fErrorCallback callback function if request failed
    //   (non 200 response from gole server) with response text
    //   as first argument
    // - Trigger fSuccessCallbak callback function if request succeeded
    //   (200 response from gole server) without arguments

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/remove",
        data: JSON.stringify(
            {
                "TileXCoordinate": removeLetterOrigin.horizontalIdx,
                "TileYCoordinate": removeLetterOrigin.verticalIdx,
                "GameId": game.id
            }
        ),
    })
    .done(function(id) {
        console.log("Letter Removed")
        fSuccessCallback()

    })
    .fail(function(response){
        fErrorCallback(response.responseText)
    })

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
    // Get the struct that describes the player that has their current turn

    $.ajax({
        async: false,
        method: "GET",
        url: server.url + "/" + game.id + "/player.json",
    })
    .done(function(player) {
        activePlayer = JSON.parse(player)
    });

}

function updateScoreBoard() {
    // Request the current game's scoreboard.

    $.ajax({
        async: false,
        method: "GET",
        url: server.url + "/" + game.id + "/scoreboard.json",
    })
    .done(function(response) {
        scoreboard = JSON.parse(response)
    });

}

function confirmWord(fSuccessCallback, fErrorCallback, fGameOverCallback) {
    // called after a player has placed
    // all tiles for the current turn
    // Requires:
    // - A callback function fSuccessCallback
    //   to be called if the request succeeded and the game can continue
    // - A callback function fErrorCallback
    //   to be called if the request failed
    // - A callback function fGameOverCallback
    //   to be called if the request succeeded and the game is now over
    // Guarantees:
    // - Send confirm word request to gole server
    // - Trigger fErrorCallback callback function if request failed
    //   (500 response from gole server) with response text
    //   as first argument
    // - Trigger fSuccessCallbak callback function if request succeeded
    //   (200 response from gole server) without arguments
    // - Trigger fGameOverCallback callback function if request succeeded
    //   and game is over (250 response from gole server)

    $.ajax({
        async: false,
        method: "POST",
        url: server.url + "/confirm",
        data: JSON.stringify(
            {
                "GameId": game.id
            }
        ),
        statusCode: {
            250: function (response) {
                fGameOverCallback();
            },
        }
    })
    .done(function(response) {
        fSuccessCallback()
    })
    .fail(function(response){
        fErrorCallback(response.responseText);
    });


}
