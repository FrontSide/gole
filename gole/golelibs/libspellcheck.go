package golelibs

import (
    "github.com/ddliu/go-dict"
    "log"
)

func IsAValidWord(word string) (bool) {
    // Return whether a word is a legitimate english word.
    // Utilize the go-dict library which makes use of
    // the locally stored dict on the computer
    //
    // requires:
    // - the existance of a file with a list of words at
    //   /usr/share/dict/words
    //   this restricts the use of this software to unix systems

    log.Println("Call spell check dictionsry for word: " + word)
    dict := dict.NewDict()
    dict.Load("/usr/share/dict/words")
    _, wordExists := dict.Get(word)
    return wordExists

}
