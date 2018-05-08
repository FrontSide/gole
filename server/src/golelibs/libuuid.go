package golelibs

import (
        "os/exec"
        "strings"
        "log"
)

func GetNewUUID() string {
	// Generate a new uuid.
	// Requires:
	// - The uuidgen binary installed on the machine
        //   on which this function is executed
	// Guarantees:
	// - Return a new uuid as generated by the called binary
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(uuid))
}