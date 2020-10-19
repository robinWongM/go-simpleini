package main

import (
	"fmt"
	"net/http"

	"github.com/robinWongM/go-simpleini/simpleini"
)

func main() {
	// Provide a listener func to deal with changes of .ini
	conf, err := simpleini.Watch("test.ini", func(conf simpleini.Configuration) {
		fmt.Printf("Reloaded. http_port = %s\n", conf.Get("server", "http_port"))
	})
	if err != nil {
		panic(err)
	}

	// Get(section, key)
	fmt.Printf("Loaded. http_port = %s\n", conf.Get("server", "http_port"))

	http.ListenAndServe(":8000", nil)
}
