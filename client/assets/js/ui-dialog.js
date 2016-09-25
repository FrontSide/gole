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

function prompt(message, dialogContainerClass, textContainerClass, buttonsToDisplay, extraDivToDisplay) {
        // Open a dialog prompt
        // Requires:
        // Guarantees:

        if (openPromptContainer) {
                console.log("Close existing prompt before opening new one.")
                closePrompt(prompt, message, dialogContainerClass, textContainerClass, buttonsToDisplay, extraDivToDisplay)
                return
        }

        var promptContainerDiv = $("<div>", {class: "gole-prompt-container"})
        openPromptContainer = promptContainerDiv

        var promptDialogContainerDiv = $("<div>", {class: "gole-prompt-dialog-container " + dialogContainerClass})

        if (!buttonsToDisplay || buttonsToDisplay.length < 1) {
                console.error("Cannot create prompt with less than one button.")
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

function closePrompt(afterCloseCallback) {
        // Close the currently open prompt dialog.
        // Requires:
        // - Optional: A callback function to be called after the prompt has
        //             been fully closed
        // - Optional: Arguments to be passed to the callback function
        //             (not visible in function signature)
        // Guarantees:
        // - Check if a prompt is currently open and initiate a close action
        //   is so
        // - Append the css class initiating the close animation to the
        //   open prompt dialog element
        // - Remove the prompt container element
        //   after the animation has stopped, remove the global reference
        //   to the prompt and call the callback function (incl. args) if given

        var callbackArguments = Array.prototype.slice.call(arguments, 1)

        if (openPromptContainer == null) {
                console.log("Nothing to close")
        }
        openPromptContainer.children().addClass("gole-prompt-dialog-container-close")

        openPromptContainer.children().on('animationend webkitAnimationEnd oAnimationEnd MSAnimationEnd',
        function() {
                // As the close function may be called from more than one source,
                // there is the possibility that the reference to the prompt
                // does not exist anymore now after the animation.
                // So it's necessary to check for it again
                // before issuing a remove
                if (openPromptContainer) {
                        openPromptContainer.remove()
                }
                openPromptContainer = null
                if (afterCloseCallback) {
                        afterCloseCallback.apply(null, callbackArguments)
                }
        });

}

function promptError(message, acceptButtonText, onAcceptCallback) {
        // Prompt an error message to the user
        // Requires:
        // - An error message to be prompted
        // - A text fot the accept button on the error dialog prompt
        // - A callback funtion to be called when the accept button is pressed
        // - Optionally: Arguments for the callback function (not visible in
        //   function signature)
        // Guarantees:
        // - show the given acceptButtonText on the neutral button below the text
        // - execute the given onAcceptCallback with optional extra given parameters
        //   once the OK button on the dialog has been pressed

        var callbackArguments = Array.prototype.slice.call(arguments, 3)

        var buttonsToDisplay = new Array()
        var acceptButton = $("<button>", {class: "gole-prompt-dialog-button gole-prompt-dialog-neutral-button"})

        if (!acceptButtonText) {
                acceptButtonText = "OK"
        }

        acceptButton.html(acceptButtonText)
        acceptButton.click(function(){
                closePrompt()
                if (onAcceptCallback) {
                        onAcceptCallback.apply(null, callbackArguments)
                }
        })
        buttonsToDisplay.push(acceptButton)

        prompt(message, "gole-prompt-error-dialog-container", "gole-prompt-text-container gole-prompt-error-text-container", buttonsToDisplay)

}

function promptGameOver(scoreboard) {
        // prompt the game over dialog
        // presenting the winner(s) and the scoreboard
        // Requires:
        // - A scoreboard as first argument
        //   which is a map with player name (key) to points (value) mapping
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

        var buttonsToDisplay = new Array()
        var newGameButton = $("<button>", {class: "gole-prompt-dialog-button gole-prompt-dialog-success-button"})
        newGameButton.html("Start New Game")
        newGameButton.click(function(){
                closePrompt()
                promptNewGame(initNewGame)
        })
        buttonsToDisplay.push(newGameButton)

        prompt(winnerMessage, "gole-prompt-success-dialog-container", "gole-prompt-text-container gole-prompt-neutra-text-container", buttonsToDisplay)

}

function promptWildcardLetterSelection(onSelectCallback) {
        // prompt the letter selection board dialog for when a wildcard is placed
        // Requires:
        // - A reference to a callback function as first argument
        // - Optionally: A number of arguments to be passed to the callback function.
        // Guarantees:
        // - execute given onSelectCallback
        //   once the SelectLetter button on the dialog has been pressed
        // - Pass the ASCII code of the character of the letter
        //   selected by the user as first argument
        //   argument to the onSelectCallback
        // - Pass the optional arguments (not visible in function signature)
        //   as succeeding arguments to the onSelectCallback
        // - Prompt will be closed immediately after letter
        //   or cancel has been selected and a full reload is initiated

        var extraCallbackArguments = Array.prototype.slice.call(arguments, 1)
        var alphabeth_en = "abcdefghijklmnopqrstuvwxyz".split('')

        var alphabethContainerDiv = $("<div>")

        $.each(alphabeth_en, function (_, letterCharacter) {

                var tileDiv = $("<div>", {class: "gole-tile gole-tile-selectable gole-tile-margin"})
                var letterDiv = $("<div>", {class: "gole-tile-letter-character-container"})

                letterDiv.html(letterCharacter.toUpperCase())
                tileDiv.append(letterDiv)

                tileDiv.click(function(){
                        console.log(letterCharacter)
                        var callbackArguments = [letterCharacter.charCodeAt()].concat(extraCallbackArguments)
                        onSelectCallback.apply(null, callbackArguments)
                        closePrompt()
                        reload()
                })

                alphabethContainerDiv.append(tileDiv)

        })

        var dismissButton = $("<button>", {class: "gole-prompt-dialog-button"})
        dismissButton.html("Cancel")
        dismissButton.click(function(){
                closePrompt()
        })

        prompt("Select a letter to be placed on the wildcard tile...", "", "gole-prompt-text-container gole-prompt-neutral-text-container", dismissButton, alphabethContainerDiv)

}

function promptNewGame(onStartCallback) {
        // Prompt the new game dialog
        // Requires:
        // - A callback function to ba called when the startGame button is pressed
        // - Optionally: Arguments to be passed to the callback function
        //   (not visible in the function signature of this promtNewGame function)
        // Guarantees:
        // - Execute the given onStartCallback (with arguments, if given)
        //   once the StartGame button on the dialog has been pressed

        var callbackArguments = Array.prototype.slice.call(arguments, 1)
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

                onStartCallback.apply(null, callbackArguments)
        })

        var dismissButton = $("<button>", {class: "gole-prompt-dialog-button"})
        dismissButton.html("Cancel")
        dismissButton.click(function(){
                closePrompt()
        })

        buttonsToDisplay.push(startButton)
        buttonsToDisplay.push(dismissButton)

        prompt("Please enter your player names (2 - 4 players)...", "", "gole-prompt-text-container gole-prompt-neutral-text-container", buttonsToDisplay, nameTextFieldsContainer)

}

function playPointsGainAnimation(pointsGained, words) {
        // Animation to be played after a player's
        // move has been successfully confirmed
        // and the player has gained points.
        // Reuqires:
        // - pointsGained: amount of points the player
        //   has gained in this turn
        // - words: The new word(s) on the board, for which
        //   the player has gained the points
        // Guarantees:
        // - Plays a short animation that informs about the gained points.
        // The animation keyframes are defined in gole.css

        var animationContainerDiv = $("<div>", {class: "points-gained-animation-container"})
        var pointsContainer = $("<div>", {class: "points-gained-animation-points-container"})

        pointsContainer.html("+" + pointsGained)

        animationContainerDiv.append(pointsContainer)
        $("body").append(animationContainerDiv)

        playsoundWordConfirmed()

}
