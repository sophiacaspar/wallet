package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"wallet.sophiacaspar/internal/data"
)

func (app *application) createWalletHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		var input struct {
			PlayerId string `json:"playerId"`
			Balance int `json:"balance"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		wallet := &data.Wallet{
			PlayerId: input.PlayerId,
			Balance: input.Balance,
		}

		err = app.models.Wallets.Insert(wallet)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/wallet/%d", wallet.ID))
		
		// Write the JSON response with a 201 Created status code and get the location header set
		err = app.writeJSON(w, http.StatusCreated, envelope{"wallet": wallet}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (app *application) getWalletHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		id := r.URL.Path[len("/v1/wallet/"):]
		app.getWalletById(w, r, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getWalletById(w http.ResponseWriter, r *http.Request, id string){
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil{
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	wallet, err := app.models.Wallets.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("wallet not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	
	if err := app.writeJSON(w, http.StatusOK, envelope{"wallet":wallet}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) subtractFromWalletHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/v1/wallet/subtract/"):]
    app.updateWalletBalance(w, r, Subtract, id)
}

func (app *application) addToWalletHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/v1/wallet/add/"):]
    app.updateWalletBalance(w, r, Add, id)
}

func (app *application) updateWalletBalance(w http.ResponseWriter, r *http.Request, op operationType, id string) {
    idInt, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    wallet, err := app.models.Wallets.Get(idInt)
    if err != nil {
        switch {
        case errors.Is(err, errors.New("wallet not found")):
            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        default:
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }

    var input struct {
        Amount *int `json:"amount"`
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    switch op {
    case Add:
        if input.Amount != nil {
            wallet.Balance += *input.Amount
        }
    case Subtract:
        if input.Amount != nil {
            if wallet.Balance-(*input.Amount) < 0 {
                http.Error(w, "Bad Request: Can't subtract more than available balance", http.StatusBadRequest)
                return
            }
            wallet.Balance -= *input.Amount
        }
    }

    err = app.models.Wallets.Update(wallet)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }

    if err := app.writeJSON(w, http.StatusOK, envelope{"wallet": wallet}, nil); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }
}