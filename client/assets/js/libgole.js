// Library that interacts with the gole server IPA

let server = {
    url: 'http://localhost:8000',
};

let game = {
    id: null,
    board: null,
    // The player names will only be needed at the beginning when they are
    // passed over by the user via the new game form
    // Needs to be an Array of strings (2-4 players are allowed)
    playerNames: [],
};

let activePlayer = {
    Name: null,
    Points: null,
    LettersInHand: null,
};

// Active Game Scoreboard
// playerName: points
let scoreboard = {};

// stores information about the currently activated tile
// i.e. the tile that is to be moved
let activatedLetter = null;

// If a letter is being dragged away from the board,
// the origin indexes of that tile need to be stored
let removeLetterOrigin = {
    verticalIdx: null,
    horizontalIdx: null,
};

// After placeing a letter, the potential points
// a player could gain with the word s/he has constructed
// throughout the course of the current turn are returned
// by the gole server and are the stored as an array
// of potentialPointsForWord objects
// (see server documentation)
let potentialPointsForWords = [];

// Enum-like struct holding constants
// defining the alignment of a word.
// This attribute can for example be added to
// the last tile in a word so thet the potentialPointsForWord container
// can be placed accordingly
let wordAlignment = {
    HORIZONTAL: 1,
    VERTICAL: 2,
};

// The ASCII representation of the wildcard character
// needs to correspond with the wildcard charater defined on the server side
let WILDCARD_CHARACTER = '*';

function createNewGame() {
    // Assign handlers immediately after making the request,
    // and remember the jqxhr object for this request

    // Requires:
    // - That the array at game.playerNames is set as a string array
    //   defining the names of the players for the new game
    // Guarantees:
    // - Sets the new game id to a global variable as well a in the store

    console.log('Start Game with players: ' + game.playerNames);

    return $.ajax({
            method: 'POST',
            url: server.url + '/new',
            data: JSON.stringify({
                'PlayerNames': game.playerNames
            }),
        })
        .done(function(id) {
            game.id = id;
            setGameID(game.id);
        });

    console.log('New Game ID:' + game.id);
}

function rearrangeLettersInHand(letterIds) {
    // Send a request to gole server to rearrange the letters
    // in the active player's hand.
    // Requires:
    // - Nothing, if the letters should be rearranged randomly
    // - An array containing the new desired arrangement of the
    //   letters in the active player's hand represented by their letterIds.
    // Guarantees:
    // - Sends a request to the gole server to rearrange the active player's
    //   hand either according to the gieven array of letterIds
    //   or randomly if no letterIds are given.

    return $.ajax({
            method: 'POST',
            url: server.url + '/hand/sort',
            data: JSON.stringify({
                'LetterIds': letterIds,
                'GameId': game.id,
            }),
        })
        .done(function(id) {
            console.log('Hand Rearranged');
            return null;
        })
        .fail(function(response) {
            return response;
        });
}

function replaceWildcardLetter(letterCharacterCode, letterId) {
    // Replace a wilcard tile letter with a real letter character
    // Requires:
    // - The ASCII code of the character that is supposed to replace
    //   the wildcard character on the tile
    // - The letter id of the wildcard tile letter which is to be replaced
    // Guarantees:
    // - Send a wildcardReplace request to the gole server
    // - Return the error message coming in from the server
    //   if not successful
    return $.ajax({
            method: 'POST',
            url: server.url + '/wildcard/replace',
            data: JSON.stringify({
                'LetterId': letterId,
                'ReplacementLetter': letterCharacterCode,
                'GameId': game.id,
            }),
        })
        .done(function(id) {
            console.log('Wildcard Replaced');
        })
        .fail(function(response) {
            return response;
        });
}

function getPotentialPoints() {
    // Call the gole API endpoint that will return
    // the potential points for unplayed but placed words
    // and updated the local data-structure.
    //
    // Requires:
    // -
    //
    // Guarantees:
    // - Assign the potentialPointsForWords array returned by
    //   the server to the client side global array varialbe
    //   potentialPointsForWords i.e. if words were found
    //   and the server returned with OK 200.

    return $.ajax({
            method: 'GET',
            url: server.url + '/' + game.id + '/potentialPoints.json',
        })
        .done(function(httpResponsePotentialPointsForWords) {
            console.log('Successfully retrieved potential points');
            potentialPointsForWords =
                JSON.parse(httpResponsePotentialPointsForWords);
        })
        .fail(function(response) {
            return response;
        });
}

function placeLetter(replaceWildcardLetterCode, letterId,
    tilesXCoordinate, tilesYCoordinate) {
    // Required:
    // - Optionally, if the letter to be placed is a wildcard letter:
    //   The ASCII code of the character that is supposed to replace
    //   the wildcard character on the tile.
    //   This is the first parameter because if this function
    //   is passed to another one as a callback,
    //   the other function may pass the replaceWildcardLetterCode
    //   as first argument.
    // - ID (as given by the gole server) of the character to be placed
    // - coordinates of the tile the letter is to be placed on
    // - boolean describing whether or not
    //   the letter to be places origintes from a wildcard tile
    //
    // Guarantees:
    // - Send letterPlacement request to gole server
    // - Return null if operation was successful
    //   and server returned with HTTP ok
    // - Return error message if operation was unsuccesful
    //   (The placement may still have been successful even though
    //    the server returned an error).

    if (replaceWildcardLetterCode) {
        console.log('Call replace Wildcard letter');
        $.when(
            replaceWildcardLetter(replaceWildcardLetterCode, letterId)
        ).done(function() {
            return placeLetter(
                null, letterId, tilesXCoordinate,
                tilesYCoordinate);
        });
    }

    return $.ajax({
            method: 'POST',
            url: server.url + '/place',
            data: JSON.stringify({
                'TileXCoordinate': tilesXCoordinate,
                'TileYCoordinate': tilesYCoordinate,
                'LetterId': letterId,
                'GameId': game.id,
            }),
        })
        .done(function(gameId) {
            console.log('Letter Placed');
            return null;
        })
        .fail(function(response) {
            return response;
        });
}

function removeLetter() {
    // Requires:
    // - The x and y coordinates of the tile from which the letter
    //   is to be removed to be set in the global
    //   removeLetterOrigin.horizontalIdx
    //   and removeLetterOrigin.verticalIdx
    //   variables.
    // Guarantees:
    // - Send request to remove letter to gole server

    return $.ajax({
            method: 'POST',
            url: server.url + '/remove',
            data: JSON.stringify({
                'TileXCoordinate': removeLetterOrigin.horizontalIdx,
                'TileYCoordinate': removeLetterOrigin.verticalIdx,
                'GameId': game.id,
            }),
        })
        .done(function(id) {
            console.log('Letter Removed');
        })
        .fail(function(response) {
            console.log('Letter Removal Failed.');
        });
}

function getBoard() {
    // Get the complete board with all properties
    // so that the function requesting it in the
    // main module can draw it to the HTML page

    // Requires:
    // -
    // Guarantees:
    // - Requests a jsonified board object from the gole server
    //   for the current game session
    // - If the server returns with HTTP 200:
    //   Sets the returned board json to the global game.board variable.

    return $.ajax({
            method: 'GET',
            url: server.url + '/' + game.id + '/board.json',
        })
        .done(function(board) {
            game.board = JSON.parse(board);
        });
}

function getActivePlayer() {
    // Get the struct that describes the player that has their current turn

    // Requires:
    // -
    // Guarantees:
    // - Requests a jsonified player object from the gole server
    //   for the active player of the current game session
    // - If the server returns with HTTP 200:
    //   Sets the returned player json to the global activePlayer variable.

    return $.ajax({
            method: 'GET',
            url: server.url + '/' + game.id + '/player.json',
        })
        .done(function(player) {
            activePlayer = JSON.parse(player);
        });
}

function updateScoreBoard() {
    // Request the current game's scoreboard.

    // Requires:
    // -
    // Guarantees:
    // - Requests a jsonified scoreboard object from the gole server
    //   for the current game session
    // - If the server returns with HTTP 200:
    //   Sets the returned scoreboard json to the global
    //   scoreboard variable.

    return $.ajax({
            method: 'GET',
            url: server.url + '/' + game.id + '/scoreboard.json',
        })
        .done(function(response) {
            scoreboard = JSON.parse(response);
        });
}

function confirmWord() {
    // called after a player has placed
    // all tiles for the current turn
    // Requires:
    // -
    // Guarantees:
    // - Send confirm word request to gole server

    return $.ajax({
        method: 'POST',
        url: server.url + '/confirm',
        data: JSON.stringify({
            'GameId': game.id,
        }),
        statusCode: {
            250: function(response) {
                console.log('Game over.');
            },
        },
    }).fail(function(response) {
        console.log("500 confirm")
        return response;
    });
}
