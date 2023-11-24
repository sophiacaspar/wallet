package main

import (
	"net/http"
)

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/wallet", app.createWalletHandler)
	mux.HandleFunc("/v1/wallet/", app.getWalletHandler)
	mux.HandleFunc("/v1/wallet/subtract/", app.subtractFromWalletHandler)
	mux.HandleFunc("/v1/wallet/add/", app.addToWalletHandler)

	return mux
}