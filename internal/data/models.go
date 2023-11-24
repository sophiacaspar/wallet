package data

import "database/sql"

type Models struct {
	Wallets WalletModel
}

func NewModels(db *sql.DB) Models {
	return Models {
		Wallets: WalletModel{DB: db}, 
	}
}