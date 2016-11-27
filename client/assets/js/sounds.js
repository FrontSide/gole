/*
 * Plays appropriate sounds for a game interaction.
 */

function playsoundGameStart() {
    $('#soundNewGame')[0].play();
}

function playsoundGameOver() {
    $('#soundGameOver')[0].play();
}

function playsoundTilePickUp() {
    $('#soundTilePickUp')[0].play();
}

function playsoundTilePlacementSuccess() {
    $('#soundTilePlacementSuccess')[0].play();
}

function playsoundTilePlacementFailure() {
    $('#soundTilePlacementFailure')[0].play();
}

function playsoundWordConfirmed() {
    $('#soundWordConfirmed')[0].play();
}

function playsoundWordInvalid() {
    $('#soundWordInvalid')[0].play();
}
