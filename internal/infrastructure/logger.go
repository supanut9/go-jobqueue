package infrastructure

import "log"

func LogInfo(msg string) {
	log.Println("[INFO]", msg)
}

func LogError(err error) {
	log.Println("[ERROR]", err)
}
