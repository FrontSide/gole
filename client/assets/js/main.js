// Main module for drawing the board and calling
// the gole js library



// Tiles are coloured red/green (locked/únloked) when player hovers
// with tile.
// No network communication needed for this
// since board has been loaded completely before

$(document).ready(function() {
    createNewGame(["oasch", "babab"])
    getBoard()
    getLegalPlacementMap()
    drawBoard()
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
            if (game.legalPlacements[yIdx][xIdx]) {
                tileLegalPlacementColorClass = "gole-board-tile-legal-placement"
            } else {
                tileLegalPlacementColorClass = "gole-board-tile-illegal-placement"
            }

            var tileDiv = $("<div>", {class: "gole-board-tile " + tileEffectColorClass + " " + tileLegalPlacementColorClass})
            $.data(tileDiv, "gole-tile-x-idx", xIdx)
            $.data(tileDiv, "gole-tile-y-idx", yIdx)

            tileInscriptionText = ""
            if (tile.Letter == 0 && tileEffectText) {
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
