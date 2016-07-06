// Main module for drawing the board and calling
// the gole js library



// Tiles are coloured red/green (locked/únloked) when player hovers
// with tile.
// No network communication needed for this
// since board has been loaded completely before

$(document).ready(function() {
    createNewGame(["oasch", "babab"])
    getBoard()
    console.log(game.board)
    getActivePlayer()
    drawBoard()
    drawPlayer()
});

function drawBoard() {

    //Clear Board container
    $("div.gole-board-container").html("")

    $.each(game.board, function(yIdx, column) {
        $.each(column, function(xIdx, tile) {

            tileEffectColorClass = ""
            tileEffectText = ""
            switch(tile.Effect) {
                case 0: tileEffectColorClass = "gole-board-tile-effect-double-letter";
                        tileEffectText = "DL"
                        break;
                case 1: tileEffectColorClass = "gole-board-tile-effect-triple-letter";
                        tileEffectText = "TL"
                        break;
                case 2: tileEffectColorClass = "gole-board-tile-effect-double-word";
                        tileEffectText = "DW"
                        break;
                case 3: tileEffectColorClass = "gole-board-tile-effect-triple-word";
                        tileEffectText = "TW"
                        break;
                case 5: tileEffectColorClass = "gole-board-tile-effect-center";
                        tileEffectText = '<i class="fa fa-star" aria-hidden="true"></i>'
                        break;
            }

            tileLegalPlacementColorClass = ""
            if (tile.PlacementIsLegal) {
                tileLegalPlacementColorClass = "gole-board-tile-legal-placement"
            } else {
                tileLegalPlacementColorClass = "gole-board-tile-illegal-placement"
            }

            var tileDiv = $("<div>", {class: "gole-board-tile " + tileEffectColorClass + " " + tileLegalPlacementColorClass})

            tileInscriptionText = ""
            if (tile.Letter.Character == 0 && tileEffectText) {
                tileInscriptionText = tileEffectText
            } else if (tile.Letter.Character != 0) {
                tileInscriptionText = String.fromCharCode(tile.Letter.Character).toUpperCase()
            }

            tileDiv.html(tileInscriptionText)

            //register tile click events
            tileDiv.click(function(){

                if (activatedLetter) {

                    // Check if the tile is activateable. i.e. if it's not already part
                    // of a played word on the board.
                    if (tile.Letter.Character != 0) {
                        console.log("Cannot place letter. Occupied.")
                    } else {
                        placeLetterOnTile(xIdx, yIdx, tile)
                    }

                } else if (tile.Letter.Character != 0) {

                    if (tile.IsLocked) {
                        console.log("sorry locked")
                    } else if (activatedLetter === tile.Letter) {
                        deactivateLetter(tile.Letter, tileDiv)
                    } else {
                        activateLetter(tile.Letter, tileDiv)
                    }

                }

            });
            $("div.gole-board-container").append(tileDiv)
        })
        $("div.gole-board-container").append("<div style='clear:both'></div>")
    })
}

function drawPlayer() {

    var nameDiv = $("<div>", {class: "gole-active-player-name-container"})
    nameDiv.html(activePlayer.Name)

    var pointsDiv = $("<div>", {class: "gole-active-player-points-container"})
    pointsDiv.html(activePlayer.Points)

    var handContainerDiv = $("<div>", {class: "gole-active-player-hand-container"})

    $.each(activePlayer.LettersInHand, function(idx, letter) {

        var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
        var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

        // Go returns the character of a letter tile as an int8 (rune) code
        // so we need to convert to string and put it uppercase before printing
        letterDiv.html(String.fromCharCode(letter.Character).toUpperCase())

        var letterValueDiv = $("<div>", {class: "gole-tile-letter-value-container"})
        letterValueDiv.html(letter.Attributes.PointValue)

        tileDiv.append(letterDiv)
        tileDiv.append(letterValueDiv)
        handContainerDiv.append(tileDiv)

        //register tile click event
        tileDiv.click(function(){
            if (activatedLetter === letter) {
                deactivateLetter(letter, this)
            } else {
                activateLetter(letter, this)
            }
        });

    })

    $("div.gole-active-player-container").append(nameDiv)
    $("div.gole-active-player-container").append(pointsDiv)
    $("div.gole-active-player-container").append(handContainerDiv)

}

// stores information about the currently activated tile
// i.e. the tile that is to be moved
var activatedLetter = null

// check whether the letter tile to be activates is moveable
// and show the user where it can be moved to and that it is now active
//
// called when a player clicks on a tile while it is deactivated
function activateLetter(letter, tileDiv) {

    //Check first if a letter is already activated and deactivate all if so
    if (activatedLetter) {
        deactivateLetter(activatedLetter)
    }

    console.log(tileDiv)
    $(tileDiv).addClass("gole-tile-activated")
    activatedLetter = letter
    console.log(activatedLetter)
    console.log("LetterToMove: " + activatedLetter.Character)

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
    hideLegalPlacements()
    console.log("LetterDEactivated: " + letter.Character)

}

function showLegalPlacements() {
    $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").addClass("gole-board-tile-legal-placement-show")
}

function hideLegalPlacements() {
    $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").removeClass("gole-board-tile-legal-placement-show")
}

function placeLetterOnTile(xIdx, yIdx, board_tile) {
    console.log("Place letter on tile. Tile:")
    console.log(board_tile)

    if (!activatedLetter) {
        console.log("No letter activated. Cannot place letter.")
        return
    }

    if (!board_tile.PlacementIsLegal) {
        console.log("Placement not Legal. Abort.")
        return
    }

    //Call libgole API request
    placeLetter(xIdx, yIdx, activatedLetter.Character)


    getBoard()
    drawBoard()

}
