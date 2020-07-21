package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	utils "testTask"
	"testTask/model"
	"time"
)

type SourceType int

type StateType int

const (
	gameType SourceType = iota + 1
	serverType
	paymentType
)

//In future this will help to scale the app, for now it helps track balance separate from transactions
const BalanceId = 1

var Internal = gin.H{"message": "Internal Error"}

func MakeTransaction(c *gin.Context) {
	insertTransaction(c)
}

func processSourceType(sourceType string) SourceType {
	switch sourceType {
	case "game":
		return gameType
	case "server":
		return serverType
	case "payment":
		return paymentType
	default:
		return 0
	}
}

func RequestIsValid(c *gin.Context, r model.Request, source SourceType) bool {
	convertedAmount, err := strconv.ParseFloat(r.Amount, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Invalid Amount"})
		return false
	}
	if convertedAmount < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Negative Amount"})
		return false
	}
	if r.State != "win" && r.State != "lost" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Wrong Transaction Status"})
		return false
	}
	if r.Id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Wrong Request ID"})
		return false
	}
	if source == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Wrong Source"})
		return false
	}
	return true
}

func insertTransaction(c *gin.Context) {
	sourceType := processSourceType(c.Request.Header.Get("Source-Type"))
	var request model.Request
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
		return
	}
	if RequestIsValid(c, request, sourceType) {
		DB, err := GetDB()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "No DB connection"})
			return
		}
		defer utils.CloseDB(DB)
		IDExist := checkTransactionIsUniq(c, request, DB)
		if IDExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"message": "Transaction already exists"})
			return
		}

		balance, err := GetBalance(c, DB)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
			return
		}
		convertedAmount, err := strconv.ParseFloat(request.Amount, 64)
		var transaction model.Transaction
		transaction.State = request.State
		transaction.Id = request.Id
		transaction.Amount = convertedAmount

		_, err = DB.Exec(
			"INSERT INTO transaction_queue (source, state,  amount, unix_timestamp, transactionId, cancelled, balance_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			sourceType, transaction.State, transaction.Amount, time.Now().UnixNano(), transaction.Id, false, BalanceId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
			return
		}
		newBalance := UpdateBalance(balance, transaction)
		_, err = DB.Exec("UPDATE balance SET balance = $1 WHERE id = $2", newBalance, BalanceId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
			return
		}
		c.JSON(200, gin.H{
			"Message": "Transaction Complete",
		})
	}
}

func GetBalance(c *gin.Context, DB *sql.DB) (float64, error) {
	var balance float64
	rows, err := DB.Query("SELECT balance FROM balance WHERE id = $1", BalanceId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
		return 0, err
	}
	defer utils.CloseRows(rows)
	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Internal)
			return 0, err
		}
	}
	return balance, nil
}

func UpdateBalance(balance float64, t model.Transaction) (newBalance float64) {

	var res float64
	if t.State == "win" {
		res = balance + t.Amount
	} else {
		res = balance - t.Amount
	}
	if res <= 0 {
		return 0
	}
	return res

}

func checkTransactionIsUniq(c *gin.Context, r model.Request, DB *sql.DB) bool {
	var unixTimestamp int
	rows, err := DB.Query("SELECT unix_timestamp FROM transaction_queue WHERE transactionId = $1", r.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return true
	}
	for rows.Next() {
		err = rows.Scan(&unixTimestamp)
		if unixTimestamp != 0 {
			return true
		}
	}
	return false
}
