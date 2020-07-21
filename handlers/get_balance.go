package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	utils "testTask"
)

func GetDB() (*sql.DB, error) {
	DB, err := sql.Open("postgres", utils.FormConnString(
		utils.PGHostname, utils.PGPort, utils.PGUsername, utils.PGPassword, utils.PGDB))
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func CheckBalance(c *gin.Context) {
	DB, err := GetDB()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	defer utils.CloseDB(DB)
	rows, err := DB.Query("SELECT balance FROM balance WHERE id = $1", BalanceId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	defer utils.CloseRows(rows)
	var balance float64
	for rows.Next() {
		err = rows.Scan(&balance)
	}
	c.JSON(200, gin.H{
		"balance": balance,
	})
}

func NullBalance(c *gin.Context) {
	DB, _ := GetDB()
	_, err := DB.Exec("UPDATE balance SET balance = 0 WHERE id = $1", BalanceId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return

	}
	c.JSON(200, gin.H{
		"balance": 0})
}
