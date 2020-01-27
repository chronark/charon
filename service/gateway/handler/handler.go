package main

import (
	"net/http"
)

type GeocodingHandler interface {
	Forward(w http.ResponseWriter, r *http.Request)
}
