/*
 * Defines non-gameplay ui elements like dialogs and prompts.
 */

var promptType = {
    DEV_DEBUG: 0,
    INFO: 1,
    QUESTION: 2,
    ERROR: 3,
    GAME_OVER: 4
}

// Reference to the container that holds all prompt elements
var openPromptContainer = null

function prompt(message, textContainerClass, buttonsToDisplay, extraDivToDisplay) {

    var promptContainerDiv = $("<div>", {class: "gole-prompt-container"})
    openPromptContainer = promptContainerDiv

    var promptDialogContainerDiv = $("<div>", {class: "gole-prompt-dialog-container"})

    if (buttonsToDisplay.length < 1) {
        console.log("Cannot create prompt with less than one button.")
        return 1
    }

    var promptDialogTextContainer = $("<div>", {class: textContainerClass})
    promptDialogTextContainer.html(message)

    promptDialogContainerDiv.append(promptDialogTextContainer)
    promptDialogContainerDiv.append(extraDivToDisplay)

    $.each(buttonsToDisplay, function(idx, button) {
        promptDialogContainerDiv.append(button)
    })

    promptContainerDiv.append(promptDialogContainerDiv)
    $("body").append(promptContainerDiv)
}

function closePrompt() {
    if (openPromptContainer == null) {
        console.log("Nothing to close")
    }
    openPromptContainer.remove()
    openPromptContainer = null
}

function promptError(message, acceptButtonText, onAcceptCallback, onAcceptCallbackArguments) {
    //prompt an error message to the user
    //show the given acceptButtonText on the neutral button below the text
    //execute the given onAcceptCallback with the onAcceptCallbackArguments
    //once the OK button on the dialog has been pressed

    var buttonsToDisplay = new Array()
    var acceptButton = $("<button>", {class: "gole-prompt-dialog-button gole-prompt-dialog-neutral-button"})

    if (!acceptButtonText) {
        acceptButtonText = "OK"
    }

    acceptButton.html(acceptButtonText)
    acceptButton.click(function(){
        closePrompt()
        if (onAcceptCallback) {
            onAcceptCallback(onAcceptCallbackArguments)
        }
    })
    buttonsToDisplay.push(acceptButton)

    prompt(message, "gole-prompt-text-container gole-prompt-error-text-container", buttonsToDisplay)

}

function promptGameOver(scoreboard) {
    // prompt the game over dialog
    // presenting the winner(s) and the scoreboard
    var winnerPlayerNames = null;
    $.each(scoreboard, function(playerName, playerPoints) {
        if ((winnerPlayerNames == null) || (playerPoints > scoreboard[winnerPlayerNames[0]])) {
            winnerPlayerNames = Array()
            winnerPlayerNames.push(playerName);
            return true; // this is the js equivalent of what is usually "continue"
        }
        // If there is a player with the same amount of points as the one
        // who at the moment is the player with the most points
        // there will be multiple winners.
        else if (playerPoints == scoreboard[winnerPlayerNames[0]]) {
            winnerPlayerNames.push(playerName);
        }
    })

    var winnerMessage = ""
    if (winnerPlayerNames.length > 1) {
        winnerMessage = "The winners are <b>"
        winnerMessage += winnerPlayerNames.join("</b> and <b>")
    } else {
        winnerMessage = "The winner is <b>"
        winnerMessage += winnerPlayerNames[0]
    }
    winnerMessage += "</b> with <b>" + scoreboard[winnerPlayerNames[0]] + "</b> points."

    prompt(message, "gole-prompt-text-container gole-prompt-neutral-text-container")

}

function promptWildcardLetterSelection(onSelectCallback, onSelectCallbackArguments) {
    //prompt the letter selection board dialog for when a wildcard is placed
    //execute the given onSelectCallback with arguments onSelectCallbackArguments
    //once the SelectLetter button on the dialog has been pressed
}

function promptNewGame(onStartCallback, onStartCallbackArguments) {
    //prompt the new game dialog
    //execute the given onStartCallback with arguments onStartCallbackArguments
    //once the StartGame button on the dialog has been pressed

    var nameTextFieldsContainer = $("<div>")

    TEXT_FIELDS_TO_DISPLAY = 4
    var nameTextFields = new Array()
    for (var inputTextFieldCounter = 0; inputTextFieldCounter < TEXT_FIELDS_TO_DISPLAY; inputTextFieldCounter++) {
        var nameTextField = $("<input>", {class: "gole-prompt-dialog-text-field", type: "text", placeholder: "Player " + (inputTextFieldCounter + 1)})
        nameTextFields.push(nameTextField)
        nameTextFieldsContainer.append(nameTextField)
    }

    var buttonsToDisplay = new Array()
    var startButton = $("<button>", {class: "gole-prompt-dialog-button gole-prompt-dialog-success-button"})
    startButton.html("Start Game")
    startButton.click(function(){
        closePrompt()

        // Push the values of the text fields for player names
        // to the in memory game.playerNames array that is used
        // to create the game.
        // Empty text fields are ignored.
        game.playerNames = new Array()
        $.each(nameTextFields, function(idx, textField) {
            if (textField.val().trim() != "") {
                game.playerNames.push(textField.val().trim())
            }
        })

        onStartCallback(onStartCallbackArguments)
    })

    var dismissButton = $("<button>", {class: "gole-prompt-dialog-button"})
    dismissButton.html("Cancel")
    dismissButton.click(function(){
        closePrompt()
    })

    buttonsToDisplay.push(startButton)
    buttonsToDisplay.push(dismissButton)

    prompt("Please enter your player names (2 - 4 players)...", "gole-prompt-text-container gole-prompt-neutral-text-container", buttonsToDisplay, nameTextFieldsContainer)

}
