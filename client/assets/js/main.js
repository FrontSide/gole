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
    $.each(game.board, function(xIdx, column) {
        $.each(column, function(yIdx, tile) {

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
            $.data(tileDiv, "gole-tile-x-idx", xIdx)
            $.data(tileDiv, "gole-tile-y-idx", yIdx)

            tileInscriptionText = ""
            if (tile.Letter.Character == 0 && tileEffectText) {
                tileInscriptionText = tileEffectText
            } else {
                tileInscriptionText = tile.Letter
            }

            tileDiv.html(tileInscriptionText)

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

        // Assign a unique id to the letter div
        // the id will consist of the last part of the game uuid,
        // the name of the player who initially owns it
        // and the letter code
        letterDiv.id = game.id


        var letterValueDiv = $("<div>", {class: "gole-tile-letter-value-container"})
        letterValueDiv.html(letter.Attributes.PointValue)

        tileDiv.append(letterDiv)
        tileDiv.append(letterValueDiv)
        handContainerDiv.append(tileDiv)
    })

    $("div.gole-active-player-container").append(nameDiv)
    $("div.gole-active-player-container").append(pointsDiv)
    $("div.gole-active-player-container").append(handContainerDiv)

    //register tile click events
    $("div.gole-tile").click(function(){
        console.log("c")
    });

}

// Stores the information about which div (identified by its ID)
// has which letter on it, necessary for UI.
// Structure {"div-id": "letterObject"}
var divLetterMapping = {}

// stores information about the current game move
var move = {
    letter: null,
    letterDiv: null
}
