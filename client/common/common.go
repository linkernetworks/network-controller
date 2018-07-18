package common

import "log"

// CheckFatal is a default success handling function, which prints error message and
// makes os.Exit in case of failed.
func CheckFatal(success bool, reason string, msg string) {
	if success {
		log.Printf("%s is sussessful.", msg)
	} else {
		log.Fatalf("%s is not sussessful. The reason is %s.", msg, reason)
	}
}
