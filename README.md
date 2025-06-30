# internal-transfer-system

## Setup

1. Install Go & PostgreSQL 17
2. Clone repo
4. Update The .env file
5. Run:
go mod tidy
go run main.go

## API

### POST /accounts
curl --location 'http://127.0.0.1:8080/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "account_id":1212,
    "balance":"96986"
}'
{ "account_id": 123, "initial_balance": "100.00" }

curl --location 'http://127.0.0.1:8080/accounts/1212'
GET /accounts/{account_id}
{ "account_id": 123, "balance": "100.00" }


curl --location 'http://127.0.0.1:8080/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "source_account_id":1212,
    "destination_account_id":245234,
    "amount":"434"
}'
POST /transactions
{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "50.00"
}
