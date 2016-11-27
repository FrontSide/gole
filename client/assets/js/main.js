/*
 * Defines the client-side game API
 */

// Tiles are coloured red/green (locked/únloked) when player hovers
// with tile.
// No network communication needed for this
// since board has been loaded completely before

$(document).ready(function() {
    getExistingGameID(function(existingGameID) {
        if (existingGameID == null) {
            console.log('No GameID found. Prompt for new game.');
            promptNewGame(initNewGame);
            return;
        } else {
            // Get the game id from the storage module
            // and assign it to the in-memory game id
            game.id = existingGameID;
            reload();
        }
    });
    initControlPanel();
});

function initNewGame() {
    // to be called AFTER the new player names have been set by the user
    // Ideally as callback for the new game prompt
    // Initiates the new game

    // Ineract with gole server to create new game
    $.when(createNewGame()).done(function() {
        reload();
        playsoundGameStart();
    });
}

function gameOver() {
    console.log('Game is over');
    updateScoreBoard();
    promptGameOver(scoreboard);
    playsoundGameOver();
}

function reload() {
    console.log('reload');
    reloadBoard();
    reloadPlayer();
}

function reloadBoard() {
    console.log('call getBoard');
    $.when(getBoard()).done(function() {
        console.log('getBoard returned');
        if (game.board == null) {
            console.error("no board found after getting board.")
            promptError('There has been a connection issue.' +
                'The current game could not be restored.',
                'Start new Game', promptNewGame, initNewGame);
        }
        $.when(getPotentialPoints()).done(function() {
            updatePotentialPointsForWords();
            drawBoard();
        });
    });
}

function reloadPlayer() {
    console.log('call getActivePlayer');
    $.when(getActivePlayer()).done(function() {
        console.log('getActivePlayer returned');
        if (game.board == null) {
            console.error("no board found after getting active player.")
            promptError('There has been a connection issue.' +
                ' The current game could not be restored.',
                'Start new Game', promptNewGame, initNewGame);
        }
        drawPlayer();
    });
}

function getArrayOfLetterIdsInHand(handContainerDiv) {
    // Get an array containing the ids of the letters in the players hand.
    // Requires:
    // - The div object that represents the player's hand i.e.
    //   the div that holds the player's hand's letter divs
    // Guarantees:
    // - Builds an array of the ids of letters in the player's hand
    //   in the order as they are visually arranged in the player's
    //   hand's div container.
    // - Returns the produces array.

    let AMOUNT_OF_LETTERS = handContainerDiv.children().length;
    let letterIds = [];

    $.each(handContainerDiv.children(), function(idx, tileDiv) {
        letterIds.push(tileDiv.id);
    });

    return letterIds;
}

function activateLetter(letter, tileDiv, originYIdx, originXIdx) {
    // check whether the letter tile to be activates is moveable
    // and show the user where it can be moved to and that it is now active
    //
    // called when a player clicks on a tile while it is deactivated
    // can be used for moving a letter from the player hand to the board
    // or vice versa.
    // if a letter tile is being removed from the board, the origin indexes
    // (i.e. where the letter was originally placed on the board)
    // need to be passed (originYIdx, originXIdx)


    // Check first if a letter is already activated and deactivate all if so
    if (activatedLetter) {
        deactivateLetter();
    }

    console.log('Activate Letter Tile :: ' + tileDiv);
    $(tileDiv).addClass('gole-tile-activated');
    activatedLetter = letter;
    if (originYIdx >= 0 && originXIdx >= 0) {
        removeLetterOrigin.verticalIdx = originYIdx;
        removeLetterOrigin.horizontalIdx = originXIdx;
    }

    playsoundTilePickUp();
    showLegalPlacements()
}

// deactivate a letter tile
//
// called when a player clicks on a letter tile while it is activated
function deactivateLetter(tileDiv) {
    console.log('Deactivate Letter Tile :: ' + tileDiv);

    // If no tile div is given deactivate all visualy
    if (!tileDiv) {
        $('.gole-tile').removeClass('gole-tile-activated');
    } else {
        $(tileDiv).removeClass('gole-tile-activated');
    }

    activatedLetter = null;
    removeLetterOrigin.verticalIdx = null;
    removeLetterOrigin.horizontalIdx = null;
    hideLegalPlacements()
}

function placeLetterOnTile(xIdx, yIdx, boardTileDiv) {
    // Innitiate the action of trying to place a letter tile
    // on the board
    // Requires:
    // - x and y index of the tile on the board on which the letter tile
    //   is to be placed
    // - the boardTileDiv object onto which the letter tile is to be placed
    // - A activatedLetter needs to be set globally.
    //   This is the letter object of the tile
    //   that is to be placed on the board.
    // Guarantees:
    // - Return nothing and log an error to the console
    //   if no activatedLetter is set
    // - Return nothing and log an error to the console
    //   if placement onto the chosen board tile is not legal
    // - Call the wildcard letter selection prompt if the
    //   character to be placed is a wildcard.
    //   The gole server API library-function and according arguments
    //   will be handed over as callback
    // - In all other cased, call the gole server placement API
    //   with the appropriate arguements
    // - Reload all if no error is returned
    console.log('Place letter on tile. Tile:');

    if (!activatedLetter) {
        console.error('No letter activated. Cannot place letter.');
        return;
    }

    if (!boardTileDiv.PlacementIsLegal) {
        console.error('Placement not Legal. Abort.');
        return;
    }

    // Get the activated tile and check if it has a isWildardTile arribute
    let activatedTiles = $('.gole-tile-activated');

    if (activatedTiles == null) {
        console.error('No gole-tile-activated div found.');
        return;
    }

    // There shoule be only one single activated tile, so the element
    // with index 0 from the array of activated tiles should never be a problem.
    // In case there are more activated tiles, if will always be the first one
    // that is taken for placement here.
    let toBePlacedTileDiv = activatedTiles[0];

    console.log('toBePlacesTileDiv');
    console.log(toBePlacedTileDiv);
    console.log(toBePlacedTileDiv.getAttribute('data-isWildcardTile'));

    if (toBePlacedTileDiv.getAttribute('data-isWildcardTile') == 'true') {
        console.log('Letter is wildcard');
        $.when(
            promptWildcardLetterSelection(placeLetter, activatedLetter.Id, xIdx, yIdx)
        ).then(function() {
            reload();
        });
    } else {
        // The initial null argument is for the wildcard replacement
        // letter code
        $.when(
            placeLetter(false, activatedLetter.Id, xIdx, yIdx)
        ).then(function() {
            reload();
        });
    }
}

function updatePotentialPointsForWords() {
    // To be called after a letter has been placed
    // and after the newest board has been retrieved from the server
    // but before it is drawn.
    // Requires:
    // -
    // Guarantees:
    // - Iterate through the potentialPointsForWords array
    //   to find potential points wor words that the player
    //   could gain for placed but unconfirmed words.
    // - Add the points to a potentialPointsForWord attribute
    //   to the tile with the according index.
    //   The tile that holds this value is always the last
    //   tile of the related word and should be displayed accordingly.
    // - Add the word alignment of the word the potentialPoints are
    //   related to to the wordAlignment attribute of the same tile.
    // You may access the potential points for a word as well as
    // the alignment of the word, the potential points are refering to,
    // by looking for the potentialPointsForWord attribute o
    // and the wordAlignment attribute of the last tile of a given word.
    //
    // For possible values of the wordAlignment, see the
    // enum-like struct with the same name in the libgole module.
    //
    // tile.potentialPointsForWord
    // tile.wordAlignment

    if (potentialPointsForWords) {
        $.each(potentialPointsForWords, function(idx, potentialPointsStruct) {
            console.log('Unconfirmed word starting at y::' + potentialPointsStruct.FirstLetterYIdx + '::x::' + potentialPointsStruct.FirstLetterXIdx);
            console.log('Letter with index y::' + potentialPointsStruct.LastLetterYIdx + '::x::' + potentialPointsStruct.LastLetterXIdx + ':: holds potential points ::' + potentialPointsStruct.PotentialPoints);
            game.board[potentialPointsStruct.LastLetterYIdx][potentialPointsStruct.LastLetterXIdx].potentialPointsForWord = potentialPointsStruct.PotentialPoints;

            if (potentialPointsStruct.LastLetterYIdx == potentialPointsStruct.FirstLetterYIdx) {
                game.board[potentialPointsStruct.LastLetterYIdx][potentialPointsStruct.LastLetterXIdx].wordAlignment = wordAlignment.HORIZONTAL;
            } else {
                game.board[potentialPointsStruct.LastLetterYIdx][potentialPointsStruct.LastLetterXIdx].wordAlignment = wordAlignment.VERTICAL;
            }
        })
    }
}

function wordConfirmSuccessRoutine(gainedPoints) {
    // To be called whenever a player has played a vald word
    // which has been confirmed by the gole server
    // Requires:
    // - The number of gained points in the now finished turn
    // Guarantees:
    // - Reload the board
    // - Reload player hand (new player will be displayed -
    //   has to be givevn by server)
    // - Play WordSuccess Animation
    reload();
    playPointsGainAnimation(gainedPoints);
}

function wordConfirmErrorRoutine(errorMessage) {
    // To be called whenever a player has played an invald word
    // which has been denied by the gole server
    // Requires:
    // - Error message from gole server as 1st argument
    // Guarantees:
    // - Prompt the error message handed over to this method
    // - Play WordInvalid Sound
    promptError(errorMessage);
    playsoundWordInvalid();
}
