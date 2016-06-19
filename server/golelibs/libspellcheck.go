package golelibs

import (
    "net/http"
    "encoding/json"
    "log"
)

var API_KEY string = "kmfwT9XfBcmshv8f6039EtD4q6rOp1VmG7kjsnicKoIphlUcHA"
var API_ENDPOINT string = "https://montanaflynn-spellcheck.p.mashape.com/check/?text="
//var API_ENDPOINT string = "https://google.com"

func IsAValidWord(word string) bool {
    // Return whether a word is a legitimate english word.
    // Utilized the spellcheck API

    client := &http.Client{}

    request, err := http.NewRequest("GET", API_ENDPOINT + word, nil)
    request.Header.Add("X-Mashape-Key", API_KEY)
    request.Header.Add("Accept", "application/json")
    response, err := client.Do(request)
    defer response.Body.Close()

    if err != nil {
        log.Fatal(err)
    }

    if response.StatusCode != 200 {
        log.Fatal("Request to Dictionary service failed. Response Code: %d", response.StatusCode)
    }

    // Following struct represents the structure, i.e., the keys
    // of the responded JSON. Note that even though the keys
    // in the actual JSON string are lower-case, they need to be
    // defined in upper case in the struct.
    // Otherwise the encoding/json package won't be able to access them
    // since they won't be accessible outside the package they are in.
    type ResponseJsonContainer struct {
        Original, Suggestion, Corrections string
    }

    responseJson := ResponseJsonContainer{}
    json.NewDecoder(response.Body).Decode(&responseJson)

    return responseJson.Suggestion == responseJson.Original

}
