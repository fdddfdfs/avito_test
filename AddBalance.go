package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
	"net/http"
)

func addBalance(c *gin.Context) {
	userID, amount, err := getAddData(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "invalid id/balance supplied"+err.Error())
		return
	} else if amount.LessThanOrEqual(decimal.Zero) {
		c.IndentedJSON(http.StatusBadRequest, "invalid balance supplied, balance must be greater then 0")
		return
	} else if userID < 0 {
		c.IndentedJSON(http.StatusBadRequest, "invalid id supplied, id must be positive")
		return
	}

	var tx *pgx.Tx
	tx, err = conn.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error")
		return
	}

	defer tx.Rollback()

	var newBalance decimal.Decimal
	newBalance, err = addBalanceToUser(userID, amount)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error with adding balance "+err.Error())
		return
	}

	var commandTag pgx.CommandTag
	commandTag, err = addBalanceAddition(userID, amount)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error with logging balance "+err.Error())
		return
	} else if commandTag.RowsAffected() != 1 {
		c.IndentedJSON(http.StatusInternalServerError, "Error with logging balance")
		return
	}

	err = tx.Commit()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error")
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"balance": newBalance.String(),
		"userID":  userID})
}
