package main

import "booky-backend/internal/app"

func main() {
	app := &app.App{}
	if err := app.Run(); err != nil {
		panic(err)
	}
	defer app.Close()
}
