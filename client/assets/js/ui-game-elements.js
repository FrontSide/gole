/*
 * Defines ngameplay ui elements like the board and player hand.
 */

function initControlPanel() {
        // Execute initial configuration
        // for the gole control panel

        $("#gole-start-new-game-button").click(function(){
                console.log("new game button pressed")
                promptNewGame(initNewGame)
        })

}

function drawBoard() {

         //Clear Board container
         $("div.gole-board-container").html("")

         $.each(game.board, function(yIdx, column) {
                 $.each(column, function(xIdx, tile) {

                        tileEffectColorClass = ""
                        tileEffectText = ""
                        switch(tile.Effect) {
                                case 0: tileEffectColorClass =
                                          "gole-board-tile-effect-double-letter";
                                        tileEffectText = "DL"
                                        break;
                                case 1: tileEffectColorClass =
                                          "gole-board-tile-effect-triple-letter";
                                        tileEffectText = "TL"
                                        break;
                                case 2: tileEffectColorClass =
                                          "gole-board-tile-effect-double-word";
                                        tileEffectText = "DW"
                                        break;
                                        case 3: tileEffectColorClass =
                                          "gole-board-tile-effect-triple-word";
                                        tileEffectText = "TW"
                                        break;
                                case 5: tileEffectColorClass =
                                          "gole-board-tile-effect-center";
                                        tileEffectText =
                                          '<i class="fa fa-star" aria-hidden="true"></i>'
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

                                 var tileDiv = $("<div>", {class: "gole-tile gole-tile-margin"})
                                 var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

                                 var characterToDisplay
                                 if (tile.Letter.Character == WILDCARD_CHARACTER.charCodeAt()) {
                                         characterToDisplay = '<i class="fa fa-star-o" aria-hidden="true"></i>'
                                 } else {
                                         characterToDisplay = String.fromCharCode(tile.Letter.Character).toUpperCase()
                                 }

                                 letterDiv.html(characterToDisplay)

                                 var letterValueDiv = $("<div>", {class: "gole-tile-letter-value-container"})
                                 letterValueDiv.html(tile.Letter.Attributes.PointValue)

                                 if (tile.potentialPointsForWord > 0) {
                                     console.log("Draw potential Points Appendix")
                                     var potentialPointsDiv = $("<div>", {class: "gole-tile-potential-points-appendix-container"})

                                     switch (tile.wordAlignment) {
                                         case wordAlignment.HORIZONTAL:
                                            potentialPointsDiv.addClass("gole-tile-potential-points-horizontal-appendix-container");
                                            break;
                                         case wordAlignment.VERTICAL:
                                            potentialPointsDiv.addClass("gole-tile-potential-points-vertical-appendix-container");
                                            break;
                                     }

                                     potentialPointsDiv.html(tile.potentialPointsForWord)
                                     tileDiv.append(potentialPointsDiv)
                                 }

                                 tileDiv.append(letterDiv)
                                 tileDiv.append(letterValueDiv)
                                 boardTileDiv.append(tileDiv)

                             if (!tile.IsLocked) {

                                     tileDiv.addClass('gole-tile-selectable')

                                     tileDiv.draggable({
                                             snap: ".gole-board-tile",
                                             snapMode: "inner",
                                             revert: "invalid",
                                             connectToSortable: ".gole-active-player-hand-container"
                                 })

                                 tileDiv.on("dragstart", function(){
                                         if (activatedLetter != tile.Letter) {
                                                 console.log("enter drag away from board :: " + yIdx + ", " + xIdx)
                                                 activateLetter(tile.Letter, this, yIdx, xIdx)
                                                 console.log("invoke remove letter")
                                                 removeLetter(null)
                                         }
                                 })

                             } else {
                                     tileDiv.addClass('gole-tile-locked')
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
                                             playsoundTilePlacementSuccess()
                                         }

                                 } else if (tile.Letter.Character != 0) {

                                         if (tile.IsLocked) {
                                                 console.log("sorry locked")
                                         } else if (activatedLetter === tile.Letter) {
                                                 deactivateLetter(tileDiv)
                                         } else {
                                                 activateLetter(tile.Letter, tileDiv)
                                         }

                                 }

                         });
                         $("div.gole-board-container").append(boardTileDiv)
                })
                $("div.gole-board-container").append("<div style='clear:both'></div>")
         })

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
                 deactivateLetter()
                 getArrayOfLetterIdsInHand(handContainerDiv, rearrangeLettersInHand ,reload)
             }
         })

         $.each(activePlayer.LettersInHand, function(idx, letter) {

                 var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
                 tileDiv.attr('id', letter.Id)
                 var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

                 tileDiv.draggable({
                        snap: ".gole-board-tile",
                        snapMode: "inner",
                        revert: "invalid",
                        connectToSortable: ".gole-active-player-hand-container"
                  })

                 var characterToDisplay
                 if (letter.Character == WILDCARD_CHARACTER.charCodeAt()) {
                         // If the letter to display is a wildcard tile we will
                         // display a special character on the tile and
                         // add a "isWildcardTile" attribute to the tile Div.
                         characterToDisplay = '<i class="fa fa-star-o" aria-hidden="true"></i>'
                         tileDiv.attr('data-isWildcardTile', true)
                 } else {
                         // Go returns the character of a letter tile as an int8 (rune) code
                         // so we need to convert to string and put it uppercase before printing
                         characterToDisplay = String.fromCharCode(letter.Character).toUpperCase()
                 }

                 letterDiv.html(characterToDisplay)

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

         var confirmWordButton = $("<button>", {class: "gole-gameplay-button gole-confirm-word-button"})
         confirmWordButton.html("Confirm Word")
         confirmWordButton.click(function(){
                 confirmWord(wordConfirmSuccessRoutine, wordConfirmErrorRoutine, gameOver)
         })

         var shuffleHandButton = $("<button>", {class: "gole-gameplay-button gole-shuffle-hand-button"})
         shuffleHandButton.html('<i class="fa fa-random" aria-hidden="true"></i>')
         shuffleHandButton.click(function(){
                 console.log("shuffle hand button pressed")
                 rearrangeLettersInHand(null, reload)
         })

         $("div.gole-active-player-container").append(nameDiv)
         $("div.gole-active-player-container").append(pointsDiv)
         $("div.gole-active-player-container").append(shuffleHandButton)
         $("div.gole-active-player-container").append(handContainerDiv)
         $("div.gole-active-player-container").append(confirmWordButton)
}

function showLegalPlacements() {
         $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").addClass("gole-board-tile-legal-placement-show")
}

function hideLegalPlacements() {
         $(".gole-board-tile-illegal-placement, .gole-board-tile-legal-placement").removeClass("gole-board-tile-legal-placement-show")
}
