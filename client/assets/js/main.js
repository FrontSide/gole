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
                     console.log("No GameID found. Prompt for new game.")
                     promptNewGame(initNewGame)
                     return
             } else {
                     // Get the game id from the storage module
                     // and assign it to the in-memory game id
                     game.id = existingGameID
                     reload()
             }
         })

});

function initNewGame() {
         // to be called AFTER the new player names have been set by the user
         // Ideally as callback for the new game prompt
         // Initiates the new game

         // Ineract with gole server to create new game
         createNewGame()
         reload()
         playsoundGameStart()
}

function gameOver() {
         console.log("Game is over")
         updateScoreBoard()
         promptGameOver(scoreboard)
         playsoundGameOver()
}

function reload() {
         console.log("reload")
         reloadBoard()
         reloadPlayer()
}

function reloadBoard() {
         getBoard()
         if (game.board == null) {
         promptError("There has been a connection issue. The current game could not be restored.", "Start new Game", promptNewGame, initNewGame)
         }
         drawBoard()
}

function reloadPlayer() {
         if (game.board == null) {
         promptError("There has been a connection issue. The current game could not be restored.", "Start new Game", promptNewGame, initNewGame)
         }
         getActivePlayer()
         drawPlayer()
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


         //Check first if a letter is already activated and deactivate all if so
         if (activatedLetter) {
                 deactivateLetter(activatedLetter)
         }

         console.log(tileDiv)
         $(tileDiv).addClass("gole-tile-activated")
         activatedLetter = letter
         if (originYIdx >= 0 && originXIdx >= 0) {
                 removeLetterOrigin.verticalIdx = originYIdx
                 removeLetterOrigin.horizontalIdx = originXIdx
         }
         console.log(activatedLetter)
         console.log("LetterToMove: " + activatedLetter.Character)

         playsoundTilePickUp()
         showLegalPlacements()

}

// deactivate a letter tile
//
// called when a player clicks on a letter tile while it is activated
function deactivateLetter(letter, tileDiv) {

         //If no tile div is given deactivate all visualy
         if (!tileDiv) {
                 $(".gole-tile").removeClass("gole-tile-activated")
         } else {
                 $(tileDiv).removeClass("gole-tile-activated")
         }

         activatedLetter = null
         removeLetterOrigin.verticalIdx = null
         removeLetterOrigin.horizontalIdx = null
         hideLegalPlacements()
         console.log("LetterDeactivated: " + letter.Character)

}

function placeLetterOnTile(xIdx, yIdx, board_tile) {
         // Innitiate the action of trying to place a letter tile
         // on the board
         // Requires:
         // - x and y index of the tile on the board on which the letter tile
         //   is to be placed
         // - the board_tile object onto which the letter tile is to be placed
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
         console.log("Place letter on tile. Tile:")
         console.log(board_tile)

         if (!activatedLetter) {
                 console.error("No letter activated. Cannot place letter.")
                 return
         }

         if (!board_tile.PlacementIsLegal) {
                 console.error("Placement not Legal. Abort.")
                 return
         }

         if (activatedLetter.Character == WILDCARD_CHARACTER.charCodeAt()) {
                 console.log("Letter is wildcard")
                 promptWildcardLetterSelection(placeLetter, activatedLetter.Id, xIdx, yIdx)
         } else {
                 // The initial null argument is for the wildcard replacement
                 // letter code
                 placeLetter(false, activatedLetter.Id, xIdx, yIdx)
         }

         reload()

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
         reload()
         playPointsGainAnimation(gainedPoints)
}

function wordConfirmErrorRoutine(errorMessage) {
         // To be called whenever a player has played an invald word
         // which has been denied by the gole server
         // Requires:
         // - Error message from gole server as 1st argument
         // Guarantees:
         // - Prompt the error message handed over to this method
         // - Play WordInvalid Sound
         promptError(errorMessage)
         playsoundWordInvalid()
}
