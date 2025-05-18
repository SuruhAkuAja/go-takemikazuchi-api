package model

import "time"

type WorkerWallet struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	WorkerID      uint64    `gorm:"column:worker_id;"`
	WalletType    string    `gorm:"column:wallet_type"`
	AccountName   string    `gorm:"column:account_name"`
	AccountNumber string    `gorm:"column:account_number"`
	BankName      string    `gorm:"column:bank_name"`
	Balance       float64   `gorm:"column:balance;default:0"`
	IsPrimary     bool      `gorm:"column:is_primary;default:false"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
