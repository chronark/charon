package main

import (
	"net/http"
)

type GeocodingHandler interface {
	Forward( http.ResponseWriter, r *http.Request) 
}