package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shyam078/internal-transfer-system/db"

	"github.com/shyam078/internal-transfer-system/models"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc models.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Creating account with ID: %d and balance: %f", acc.AccountID, acc.Balance)
	_, err := db.DB.Exec("INSERT INTO accounts (account_id, balance) VALUES ($1, $2)", acc.AccountID, acc.Balance)
	if err != nil {
		log.Printf("Failed to create account ID %d: %v", acc.AccountID, err)
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Account created successfully with ID: %d", acc.AccountID)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "account_id")
	log.Printf("Fetching account with ID: %s", id)
	var acc models.Account
	err := db.DB.QueryRow("SELECT account_id, balance FROM accounts WHERE account_id = $1", id).Scan(&acc.AccountID, &acc.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Account not found: %s", id)
			http.Error(w, "Account not found", http.StatusNotFound)
		} else {
			log.Printf("DB error fetching account %s: %v", id, err)
			http.Error(w, "DB error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(acc)
	log.Printf("Account fetched successfully: %s", id)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(tx.Amount, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	log.Printf("Performing transaction from %d to %d amount %f", tx.SourceAccountID, tx.DestinationAccountID, amount)
	txErr := performTransaction(tx.SourceAccountID, tx.DestinationAccountID, amount)
	if txErr != nil {
		log.Printf("Transaction failed: %v", txErr)
		http.Error(w, txErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Transaction completed successfully from %d to %d amount %f", tx.SourceAccountID, tx.DestinationAccountID, amount)
}

func performTransaction(sourceID, destID int64, amount float64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	var srcBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE account_id = $1 FOR UPDATE", sourceID).Scan(&srcBalance)
	if err != nil {
		log.Printf("Failed to get source account balance for %d: %v", sourceID, err)
		return err
	}

	if srcBalance < amount {
		log.Printf("Insufficient funds in account %d: balance %f, required %f", sourceID, srcBalance, amount)
		return sql.ErrTxDone
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE account_id = $2", amount, sourceID)
	if err != nil {
		log.Printf("Failed to debit source account %d: %v", sourceID, err)
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE account_id = $2", amount, destID)
	if err != nil {
		log.Printf("Failed to credit destination account %d: %v", destID, err)
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)", sourceID, destID, amount)
	if err != nil {
		log.Printf("Failed to insert transaction record: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
