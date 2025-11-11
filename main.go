package main

import "github.com/Awaken1119/assistant-voice/route"

func main() {
	router := route.Router()
	router.Run(":8080")

}
