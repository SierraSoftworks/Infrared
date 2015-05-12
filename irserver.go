package main

import (
	"github.com/SierraSoftworks/Infrared/lib"
	"net/http"
)

func main() {
	handler := infrared.Setup()

	http.ListenAndServe(":8080", handler)
}
