package main

import "log"

func main() {
	app := NewApp()

	err := app.Start("7070")
	if err != nil {
		log.Fatalln(err)
	}
}
