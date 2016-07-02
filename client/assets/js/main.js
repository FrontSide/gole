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
            switch(tile.Effect) {
                case 0: tileEffectColorClass = "gole-board-tile-effect-double-letter"; break;
                case 1: tileEffectColorClass = "gole-board-tile-effect-triple-letter"; break;
                case 2: tileEffectColorClass = "gole-board-tile-effect-double-word"; break;
                case 3: tileEffectColorClass = "gole-board-tile-effect-triple-word"; break;
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
            tileDiv.html(tile.Letter)
            $("div.gole-board-container").append(tileDiv)
        })
        $("div.gole-board-container").append("<div style='clear:both'></div>")
    })
}
