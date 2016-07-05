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

        // Assign a unique id to the tile div
        // the id will consist of the the name of the player
        // who initially owns, the unix timestamp at assignment of the ID
        // + the index of the iteration
        // and the letter code
        tile_id = activePlayer.Name + "-" + (Date.now() + idx) + "-" + letter.Character
        tile = new Tile(tile_id, letter)

        var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
        tileDiv.attr("id", tile_id)

        // map the tile object to the id on a global mapping object
        // this way it will be easy later to determine
        // which tileDiv has which tile object which is necessary for the gameplay
        divTileMapping[tile_id] = tile

        var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

        // Go returns the character of a letter tile as an int8 (rune) code
        // so we need to convert to string and put it uppercase before printing
        letterDiv.html(String.fromCharCode(letter.Character).toUpperCase())

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
        if (divTileMapping[this.id].isActivated) {
            deactivate_tile(divTileMapping[this.id])
        } else {
            activate_tile(divTileMapping[this.id])
        }
    });

}

// Stores the information about which div (identified by its ID)
// has which tile object on it, necessary for UI.
// Structure {"div-id": "tileObject"}
var divTileMapping = {}

// stores information about the currently activated tile
// i.e. the tile that is to be moved
var activatedTile = null

// check whether the activated tile is moveable
// and show the user where it can be moved to and that it is now active
//
// called when a player clicks on a tile while it is deactivated
function activate_tile(tile) {

    //Check first if a letter is already activated and deactivate it if so
    if (activatedTile) {
        deactivate_tile(activatedTile)
    }

    tile.isActivated = true
    $("#" + tile.id).addClass("gole-tile-activated")
    activatedTile = divTileMapping[tile.id]
    console.log(activatedTile)
    console.log("LetterToMove: " + activatedTile.letter.Character)
}

// deactivate a tile visually and on the tile instance
//
// called when a player clicks on a tile while it is activated
// and when a letter is placed on a valid position on the board
function deactivate_tile(tile) {
    tile.isActivated = false
    $("#" + tile.id).removeClass("gole-tile-activated")
    activatedTile = null
    console.log("LetterDEactivated: " + tile.letter.Character)
}
