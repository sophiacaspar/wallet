package data

import (
	"database/sql"
	"errors"
	"time"
)

type Wallet struct {
	ID int64 `json:"id"`
	PlayerId string `json:"playerId"`
	Balance int `json:"balance"`
	LastUpdatedAt time.Time `json:"-"`
}

type WalletModel struct {
	DB *sql.DB
}

func (w WalletModel) Insert(wallet *Wallet) error {
	query := `
	INSERT INTO wallets (playerId, balance) 
	VALUES ($1, $2)
	RETURNING id, lastUpdatedAt`

	args := []interface{}{wallet.PlayerId, wallet.Balance}
	return w.DB.QueryRow(query, args...).Scan(&wallet.ID, &wallet.LastUpdatedAt)
}

func (w WalletModel) Get(id int64) (*Wallet, error) {
	if id < 1 {
		return nil, errors.New("wallet not found")
	}

	query := `
		SELECT id, playerId, balance, lastUpdatedAt
		FROM wallets
		WHERE id = $1`

	var wallet Wallet

	err := w.DB.QueryRow(query, id).Scan(
		&wallet.ID,
		&wallet.PlayerId,
		&wallet.Balance,
		&wallet.LastUpdatedAt,
	)
	
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("wallet not found")
		default:
			return nil, err
		}
	}
	return &wallet, nil
}

func (w WalletModel) Update(wallet *Wallet) error {
	query := `
		UPDATE wallets
		SET balance = $1, lastUpdatedAt = $2
		WHERE id = $3 AND playerId = $4
		RETURNING balance, lastUpdatedAt
	`
	args := []interface{}{wallet.Balance, time.Now(), wallet.ID, wallet.PlayerId}
	return w.DB.QueryRow(query, args...).Scan(&wallet.Balance, &wallet.LastUpdatedAt)
}