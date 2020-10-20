package main

import (
	"fmt"
	"net/http"

	"github.com/robinWongM/go-simpleini/simpleini"
)

type myListener struct{}

func (m *myListener) Listen(conf simpleini.Configuration) {
	fmt.Printf("Reloaded. http_port = %s\n", conf.Get("server", "http_port"))
}

func main() {
	// Provide a listener func to deal with changes of .ini
	conf, err := simpleini.Watch("test.ini", &myListener{})
	if err != nil {
		panic(err)
	}

	// Get(section, key)
	fmt.Printf("Loaded. http_port = %s\n", conf.Get("server", "http_port"))

	http.ListenAndServe(":8000", nil)
}
