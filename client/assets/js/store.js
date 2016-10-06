/*
 * This module offers an API for
 * storing gole user data utilizing the
 * electron-json-storage localStorage library.
 *
 * The library only stores full JSON objects under storage keys,
 * so we need to first store a JSON object at a storage key location,
 * and then, the id under a JSON Object Key location as its value.
 * The stored structure will eventually look as follows:
 * GAME_ID_STORAGE_KEY: {GAME_ID_OBJECT_KEY: GAME_ID_OBJECT_OBJECT_VALUE}
 */

const golestore = require('electron-json-storage');

GAME_ID_STORAGE_KEY = "golegameid"
GAME_ID_OBJECT_KEY = "id"

function getExistingGameID(afterFetchCallback) {
        // Get the Game ID of a previously
        // or currently played game if existing
        // in the local storage.
        // Requires:
        // - A afterFetchCallback function that will be called
        //   after the storage engine has returned the gameID object
        // Guarantees:
        // - Retrieved the GameID stored for gole and passed the result
        //   (if not empty) to the given callback function.
        // - Passes null to the given Callback function in case of an error
        //   or if the retrieved object from the storage is empty,
        //   i.e., no GameID stored.

        console.log(afterFetchCallback)

        golestore.get(GAME_ID_STORAGE_KEY, function(error, data) {
                if (error) {
                        console.err("Error when fetching GameID")
                        afterFetchCallback(null)
                        return
                }

                // Return null if object empty
                if (Object.keys(data).length === 0) {
                        console.log("No GameID found.")
                        afterFetchCallback(null)
                        return
                }

                console.log("GameID found :: " + data[GAME_ID_OBJECT_KEY]);
                afterFetchCallback(data[GAME_ID_OBJECT_KEY])
                return
        });

}

function setGameID(gameID){
        // Store a gameID in localstorage
        // Requires:
        // - Valid Gole GameID as received by gole server
        // Guarantees:
        // - Store the given ID in localstorage
        // - Previously stored gameIDs will be overwritten.

        gameIDObject = {}
        gameIDObject[GAME_ID_OBJECT_KEY] = gameID

        golestore.set(GAME_ID_STORAGE_KEY, gameIDObject, function(error) {
                if (error) {
                        console.err("Failed to store GameID :: " + gameID)
                        throw error;
                }
                console.log("Stored ID to localStorage :: " + gameID)
        });

}

function unsetGameID() {
        // Delete a stored Game ID
        // Guarantees:
        // - Remove the gameID stored in the one available
        //   local store spot.
        // - Do nothing if none exists.

        golestore.remove(GAME_ID_STORAGE_KEY, function(error) {
                if (error) {
                        console.err("Failed to remove GameID")
                        throw error;
                }
        });

}
