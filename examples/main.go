package main

import (
	"fmt"
	"net/http"
	"simpleini/simpleini"
)

func main() {
	conf, err := simpleini.Watch("test.ini", func(c simpleini.Configuration) {
		fmt.Printf("%v\n", c)
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", conf)

	http.ListenAndServe(":8000", nil)
}
