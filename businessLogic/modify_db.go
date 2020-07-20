package businessLogic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sync"
	utils "testTask"
	"testTask/model"
	"time"
)

func UpdateDB(c *gin.Context) {
	var lock sync.Mutex
	duration, err := strconv.Atoi(utils.Duration)
	if err != nil {
		log.Println(err)
	}
	timer := time.NewTicker(time.Second * time.Duration(duration))
	defer timer.Stop()
	for {
		/* run forever */
		select {
		case <-timer.C:
			go func() {
				lock.Lock()
				defer lock.Unlock()
				MakeDBCalculations(c)
			}()
		}
	}
}

func MakeDBCalculations(c *gin.Context) {
	DB, err := GetDB()
	defer utils.CloseDB(DB)
	if err != nil {
		utils.LogCancels(err.Error())
		return
	}
	rowsToCancel, err := strconv.Atoi(utils.CancelRows)
	if err != nil {
		utils.LogCancels(err.Error())
		return
	}
	rows, err := DB.Query("SELECT state, amount, transactionId FROM transaction_queue WHERE cancelled = false AND MOD (id, 2) = 1 ORDER BY id  DESC LIMIT $1", rowsToCancel)
	if err != nil {
		utils.LogCancels(err.Error())
		return
	}
	defer utils.CloseRows(rows)
	transactionsToCancel := make([]model.Transaction, 0, rowsToCancel)
	var state string
	var amount float64
	var transactionId string
	for rows.Next() {
		err = rows.Scan(&state, &amount, &transactionId)
		var transaction = model.Transaction{state, amount, transactionId}
		transactionsToCancel = append(transactionsToCancel, transaction)
	}
	revert := make([]float64, 0, rowsToCancel)
	for _, t := range transactionsToCancel {
		if err != nil {
			utils.LogCancels(err.Error())
		}
		if t.State == "win" {
			revert = append(revert, -t.Amount)
		} else {
			revert = append(revert, t.Amount)
		}
	}
	balance, err := GetBalance(c, DB)
	if err != nil {
		utils.LogCancels(err.Error())
	}
	for _, val := range revert {
		balance += val
		if balance <= 0 {
			balance = 0
		}
	}
	_, err = DB.Exec("UPDATE balance SET balance = $1 WHERE id = $2", balance, BalanceId)
	if err != nil {
		utils.LogCancels(err.Error())
	}
	for _, transaction := range transactionsToCancel {
		_, err := DB.Exec("UPDATE transaction_queue SET cancelled = true WHERE transactionId = $1", transaction.Id)
		if err != nil {
			utils.LogCancels(err.Error())
		}
	}
	timestamp := time.Now().UnixNano()
	if len(transactionsToCancel) != 0 {
		loginfo := fmt.Sprint("timestamp: &v cancelled: &v", timestamp, transactionsToCancel)
		utils.LogCancels(loginfo)
	}
}


