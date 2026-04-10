package main

import "log"

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
	}

	err := app.run()
	if err != nil {
		log.Fatal(app.run())
	}

}
