package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"wallet.sophiacaspar/internal/data"
)

func newTestApplication(t *testing.T) *application {
    db := setupTestDB(t)
    models := data.NewModels(db)

    app := &application{
        config:  config{},
        models:  models,
        logger:  log.New(os.Stdout, "", 0),
    }

    app.route()

    return app
}

const (
	testDBName = "test_wallet_db"
)

func teardownTestDB(t *testing.T, db *sql.DB) {
    if db != nil {
        // Close the database connection
        if err := db.Close(); err != nil {
            t.Fatalf("Error closing database connection: %s", err)
        }
    }
}

func setupTestDB(t *testing.T) *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=secretpassword dbname=wallets_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("error connecting to PostgreSQL: %s", err)
	}

	// Drop the test database if it already exists
    _, err = db.Exec("DROP DATABASE IF EXISTS " + testDBName)
	if err != nil {
        t.Fatalf("Error dropping test database: %s", err)
    }

	_, err = db.Exec("CREATE DATABASE " + testDBName) 
	if err != nil {
		t.Fatalf("error creating test database: %s", err)
	}

	connStr = connStr + " dbname=" + testDBName
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("error connecting to test database: %s", err)
	}

	_, err = db.Exec(`
		CREATE TABLE wallets (
			id SERIAL PRIMARY KEY,
			playerId TEXT,
			balance INT,
			lastUpdatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("error creating table: %s", err)
	}

	return db
}

func TestGetWalletHandler(t *testing.T) {
 // TODO: implement this
}


func TestSubtractFromWalletHandler(t *testing.T) {
    app := newTestApplication(t)
    defer teardownTestDB(t, app.models.Wallets.DB)

    // Create a new wallet with a balance to subtract from
    initialBalance := 100
    newWallet := &data.Wallet{
        PlayerId: "player123",
        Balance:  initialBalance,
    }
    if err := app.models.Wallets.Insert(newWallet); err != nil {
        t.Fatalf("failed to insert wallet: %s", err)
    }

    // JSON payload to subtract an amount from the wallet
    payload := map[string]int{"amount": 50}
    jsonStr, err := json.Marshal(payload)
    if err != nil {
        t.Fatalf("could not marshal JSON: %s", err)
    }

    // Build a request to subtract from the wallet
    req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/wallet/subtract/%d", newWallet.ID), bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatalf("could not create request: %s", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a ResponseRecorder to record the response
    responseRecorder := httptest.NewRecorder()

    // Serve the HTTP request and record the response
    app.subtractFromWalletHandler(responseRecorder, req)

    // Check the status code is what we expect
    if status := responseRecorder.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Fetch the updated wallet after subtraction
    updatedWallet, err := app.models.Wallets.Get(newWallet.ID)
    if err != nil {
        t.Fatalf("failed to get updated wallet: %s", err)
    }

    // Check if the balance has been updated correctly after subtraction
    expectedBalance := initialBalance - payload["amount"]
    if updatedWallet.Balance != expectedBalance {
        t.Errorf("incorrect wallet balance after subtraction: got %v want %v",
            updatedWallet.Balance, expectedBalance)
    }
}

func TestAddToWalletHandler(t *testing.T) {
    // Setting up the test application and deferring cleanup
    app := newTestApplication(t)
    defer teardownTestDB(t, app.models.Wallets.DB)

    // Creating a new wallet with an initial balance
    initialBalance := 100
    newWallet := &data.Wallet{
        PlayerId: "player123",
        Balance:  initialBalance,
    }
    if err := app.models.Wallets.Insert(newWallet); err != nil {
        t.Fatalf("failed to insert wallet: %s", err)
    }

    // JSON payload to add an amount to the wallet
    payload := map[string]int{"amount": 50}
    jsonStr, err := json.Marshal(payload)
    if err != nil {
        t.Fatalf("could not marshal JSON: %s", err)
    }

    // Building a request to add to the wallet
    req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/wallet/add/%d", newWallet.ID), bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatalf("could not create request: %s", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Creating a ResponseRecorder to record the response
    responseRecorder := httptest.NewRecorder()

    // Serving the HTTP request and recording the response
    app.addToWalletHandler(responseRecorder, req)

    // Checking the status code is as expected
    if status := responseRecorder.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Fetching the updated wallet after addition
    updatedWallet, err := app.models.Wallets.Get(newWallet.ID)
    if err != nil {
        t.Fatalf("failed to get updated wallet: %s", err)
    }

    // Checking if the balance has been updated correctly after addition
    expectedBalance := initialBalance + payload["amount"]
    if updatedWallet.Balance != expectedBalance {
        t.Errorf("incorrect wallet balance after addition: got %v want %v",
            updatedWallet.Balance, expectedBalance)
    }
}