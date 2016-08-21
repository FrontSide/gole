/*
 * Defines ngameplay ui elements like the board and player hand.
 */

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

            var boardTileDiv = $("<div>", {class: "gole-board-tile " + tileEffectColorClass})

            tileLegalPlacementColorClass = ""
            if (tile.PlacementIsLegal) {
                boardTileDiv.droppable()
                boardTileDiv.addClass("gole-board-tile-legal-placement")
            } else {
                boardTileDiv.addClass("gole-board-tile-illegal-placement")
            }

            tileInscriptionText = ""
            if (tile.Letter.Character == 0 && tileEffectText) {

                boardTileDiv.html(tileEffectText)
                boardTileDiv.addClass("gole-board-tile-no-tile")

            } else if (tile.Letter.Character != 0) {

                var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
                var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

                letterDiv.html(String.fromCharCode(tile.Letter.Character).toUpperCase())

                var letterValueDiv = $("<div>", {class: "gole-tile-letter-value-container"})
                letterValueDiv.html(tile.Letter.Attributes.PointValue)

                tileDiv.append(letterDiv)
                tileDiv.append(letterValueDiv)
                boardTileDiv.append(tileDiv)

                if (!tile.IsLocked) {

                    tileDiv.draggable({
                        snap: ".gole-board-tile",
                        snapMode: "inner",
                        revert: "invalid",
                        connectToSortable: ".gole-active-player-hand-container"
                    })

                    tileDiv.on("dragstart", function(){
                        console.log("enter drag away from board")
                        activateLetter(tile.Letter, this, yIdx, xIdx)
                    })

                }

            }

            // define actions when a drop on an element
            // occurs on a board tile
            // i.e. player tried to place letter
            boardTileDiv.on("drop", function(){

                console.log("drop action")

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
            $("div.gole-board-container").append(boardTileDiv)
        })
        $("div.gole-board-container").append("<div style='clear:both'></div>")
    })

    var startNewGameButton = $("<button>", {class: "gole-gameplay-button"})
    startNewGameButton.html("Start New Game")
    startNewGameButton.click(function(){
        console.log("new game button pressed")
        promptNewGame(initNewGame)
    })
    $("div.gole-board-container").append(startNewGameButton)
}

// The div that represents the container for the
// players hand i.e. the tiles the player owns and has not yet
// playes. Tiles can be moved from this container to specific
// places on the board and from specific places on the board
// back to the hand.
var handContainerDiv = null

function drawPlayer() {

    $("div.gole-active-player-container").html("")

    var nameDiv = $("<div>", {class: "gole-active-player-name-container"})
    nameDiv.html(activePlayer.Name)

    var pointsDiv = $("<div>", {class: "gole-active-player-points-container"})
    pointsDiv.html(activePlayer.Points +  " Points")

    handContainerDiv = $("<div>", {class: "gole-active-player-hand-container"})
    handContainerDiv.sortable({
        // define actions when a drop of an element
        // occurs on the player hand container
        // i.e. player tried to remove letter
        receive: function(event, ui) {
            if (removeLetterOrigin.verticalIdx && removeLetterOrigin.horizontalIdx) {
                console.log("hand drop action - remove letter")
                removeLetter(reload, promptError)
            }
            deactivateLetter(activatedLetter)
        }
    })

    $.each(activePlayer.LettersInHand, function(idx, letter) {

        var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
        var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

        tileDiv.draggable({
           snap: ".gole-board-tile",
           snapMode: "inner",
           revert: "invalid",
           connectToSortable: ".gole-active-player-hand-container"
        })

        // Go returns the character of a letter tile as an int8 (rune) code
        // so we need to convert to string and put it uppercase before printing
        letterDiv.html(String.fromCharCode(letter.Character).toUpperCase())

        var letterValueDiv = $("<div>", {class: "gole-tile-letter-value-container"})
        letterValueDiv.html(letter.Attributes.PointValue)

        tileDiv.append(letterDiv)
        tileDiv.append(letterValueDiv)
        handContainerDiv.append(tileDiv)

        tileDiv.on("dragstart", function(){
            console.log("enter drag away from hand")
            activateLetter(letter, this)
        })

    })

    var confirmWordButton = $("<button>", {class: "gole-gameplay-button"})
    confirmWordButton.html("Confirm Word")
    confirmWordButton.click(function(){
        confirmWord(reload, promptError, gameOver)
    })

    $("div.gole-active-player-container").append(nameDiv)
    $("div.gole-active-player-container").append(pointsDiv)
    $("div.gole-active-player-container").append(handContainerDiv)
    $("div.gole-active-player-container").append(confirmWordButton)

}

function showLegalPlacements() {
    $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").addClass("gole-board-tile-legal-placement-show")
}

function hideLegalPlacements() {
    $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").removeClass("gole-board-tile-legal-placement-show")
}
