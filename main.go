package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/ValidatePayeeRef", ValidatePayeeRef)
	http.HandleFunc("/ValidateBusiness", ValidateBusinessEndpoint)

	http.ListenAndServe(":8080", nil)
}
