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

// The ASCII representation of the wildcard character
// needs to correspond with the wildcard charater defined on the server side
var WILDCARD_CHARACTER = '*'

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
             setGameID(game.id);
         });

         console.log("New Game ID:" + game.id);

}

function rearrangeLettersInHand(letterIds, fSuccessCallback) {
        // Send a request to gole server to rearrange the letters
        // in the active player's hand.
        // Requires:
        // - Nothing, if the letters should be rearranged randomly
        // - An array containing the new desired arrangement of the
        //   letters in the active player's hand represented by their letterIds
        // - A successCallback function to be called if the server returns
        //   with 200 OK.
        // Guarantees:
        // - Sends a request to the gole server to rearrange the active player's
        //   hand either according to the gieven array of letterIds
        //   or randomly if no letterIds are given.
        // - Invoke the successCallback function after the server has
        //   returned with an HTTP 200.

        $.ajax({
            async: false,
            method: "POST",
            url: server.url + "/hand/sort",
            data: JSON.stringify(
                {
                    "LetterIds": letterIds,
                    "GameId": game.id
                }
            ),
        })
        .done(function(id) {
                console.log("Hand Rearranged")
                if (fSuccessCallback) {
                        return fSuccessCallback()
                }
                return null
        })
        .fail(function(response) {
                return response
        })

}

function replaceWildcardLetter(letterCharacterCode, letterId, fSuccessCallback) {
        // Replace a wilcard tile letter with a real letter character
        // Requires:
        // - The ASCII code of the character that is supposed to replace
        //   the wildcard character on the tile
        // - The letter id of the wildcard tile letter which is to be replaced
        // - A callback function to be called when the action was successful
        // - Optionally, arguments to pass to the successCallback function
        //   (not visible in function signature)
        // Guarantees:
        // - Send a wildcardReplace request to the gole server
        // - Run the successCallback function with optional given arguments
        //   if the operation was successul
        // - Return the error message coming in from the server if not successful

        var extraCallbackArguments = Array.prototype.slice.call(arguments, 3)

        $.ajax({
            async: false,
            method: "POST",
            url: server.url + "/wildcard/replace",
            data: JSON.stringify(
                {
                    "LetterId": letterId,
                    "ReplacementLetter": letterCharacterCode,
                    "GameId": game.id
                }
            ),
        })
        .done(function(id) {
            console.log("Wildcard Replaced")
            return fSuccessCallback.apply(null, extraCallbackArguments)
        })
        .fail(function(response) {
            return response
        })

}

function placeLetter(replaceWildcardLetterCode, letterId, tilesXCoordinate, tilesYCoordinate) {
         // Required:
         // - Optionally, if the letter to be placed is a wildcard letter:
         //   The ASCII code of the character that is supposed to replace
         //   the wildcard character on the tile.
         //   This is the first parameter because if this function
         //   is passed to another one as a callback,
         //   the otehr function may pass the replaceWildcardLetterCode
         //   as first argument.
         // - ID (as given by the gole server) of the character to be placed
         // - coordinates of the tile the letter is to be placed on
         // - boolean describing whether or not
         //   the letter to be places origintes from a wildcard tile
         // Guarantees:
         // - Send letterPlacement request to gole server
         // - Return null if operation was successfult and server returned with HTTP ok
         // - Return error message if operation was unsuccesful

         if (replaceWildcardLetterCode) {
                 console.log("Call replace Wildcard letter")
                 return replaceWildcardLetter(replaceWildcardLetterCode, letterId, placeLetter, false, letterId, tilesXCoordinate, tilesYCoordinate)
         }

         $.ajax({
             async: false,
             method: "POST",
             url: server.url + "/place",
             data: JSON.stringify(
                 {
                     "TileXCoordinate": tilesXCoordinate,
                     "TileYCoordinate": tilesYCoordinate,
                     "LetterId": letterId,
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

function removeLetter(fSuccessCallbak) {
         // Requires:
         // - The x and y coordinates of the tile from which the letter
         //   is to be removed to be set in the global
         //   removeLetterOrigin.horizontalIdx and removeLetterOrigin.verticalIdx
         //   variables.
         // - Optional: A callback function that is invoked when the
         //   http response returns a success.
         // Guarantees:
         // - Send request to remove letter to gole server

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
                 if (fSuccessCallbak){
                         fSuccessCallbak()
                 }
         })
         .fail(function(response){
                 console.log("Letter Removal Failed.")
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
         //   (200 response from gole server) with gained points as first argument
         //   and list of words for which points were gaines as second argument
         // - Trigger fGameOverCallback callback function if request succeeded
         //   and game is over (250 response from gole server)
         //   with gained points as first argument and list of words
         //   for which points were gaines as second argument

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
             wordConfirmResponse = JSON.parse(response)
             fSuccessCallback(wordConfirmResponse.GainedPoints)
         })
         .fail(function(response){
             fErrorCallback(response.responseText);
         });


}
