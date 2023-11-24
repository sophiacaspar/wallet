package data

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	testDBName = "test_wallet_db"
)

func setupTestDB(t *testing.T) *sql.DB {
	connStr := "postgres://wallets:pa55w0rd@db:5432/wallets_db?sslmode=disable"

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
	
	if err := db.Close(); err != nil {
		t.Fatalf("Error closing database connection: %s", err)
	}

	connStr = "postgres://wallets:pa55w0rd@db:5432/"+ testDBName + "?sslmode=disable"
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

func teardownTestDB(t *testing.T, db *sql.DB) {
    if db != nil {
        // Close the database connection
        if err := db.Close(); err != nil {
            t.Fatalf("Error closing database connection: %s", err)
        }
    }
}

func TestWalletInsert(t *testing.T) {
	db := setupTestDB(t) // Create a test PostgreSQL DB
	defer teardownTestDB(t, db) // Cleanup function to be executed after this test completes

	model := WalletModel{DB: db}

	wallet := &Wallet{
		PlayerId:      "player123",
		Balance:       100,
		LastUpdatedAt: time.Now(),
	}

	err := model.Insert(wallet)
	if err != nil {
		t.Errorf("Insert failed: %s", err)
	}
}

func TestWalletGet(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	model := WalletModel{DB: db}

	walletID := int64(1)
	expectedWallet := &Wallet{
		ID:            walletID,
		PlayerId:      "player123",
		Balance:       100,
	}

	err := model.Insert(expectedWallet)
	if err != nil {
		t.Errorf("Insert failed: %s", err)
	}

	retrievedWallet, err := model.Get(walletID)
	if err != nil {
		t.Errorf("Get failed: %s", err)
	}

	if !reflect.DeepEqual(retrievedWallet, expectedWallet) {
		t.Errorf("Retrieved wallet does not match expected. Got: %+v, Expected: %+v", retrievedWallet, expectedWallet)
	}
}

func TestWalletUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	model := WalletModel{DB: db}

	walletID := int64(1)
	initialBalance := 100
	expectedWallet := &Wallet{
		ID:            walletID,
		PlayerId:      "player123",
		Balance:       initialBalance,
		LastUpdatedAt: time.Now(),
	}

	err := model.Insert(expectedWallet)
	if err != nil {
		t.Errorf("Insert failed: %s", err)
	}

	expectedWallet.Balance = 200 // Updated balance value

	err = model.Update(expectedWallet)
	if err != nil {
		t.Errorf("Update failed: %s", err)
	}

	retrievedWallet, err := model.Get(walletID)
	if err != nil {
		t.Errorf("Get failed: %s", err)
	}

	if retrievedWallet.Balance != expectedWallet.Balance {
		t.Errorf("Retrieved wallet balance does not match expected. Got: %d, Expected: %d", retrievedWallet.Balance, expectedWallet.Balance)
	}
}


func TestWalletGetNotFound(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    model := WalletModel{DB: db}

    walletID := int64(1)

    // Attempt to retrieve a wallet that does not exist in the test database
    retrievedWallet, err := model.Get(walletID)
    if err == nil {
        t.Error("Expected error, but got nil")
    }

    expectedErrMsg := "wallet not found"
    if err != nil && err.Error() != expectedErrMsg {
        t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
    }

    if retrievedWallet != nil {
        t.Error("Expected retrievedWallet to be nil, but it's not")
    }
}
