package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

func getBalance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, "Invalid ID supplied")
		return
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid ID supplied "+err.Error())
		return
	}

	var userBalance decimal.Decimal
	var balanceExist bool

	userBalance, balanceExist, err = getUserBalance(userID)
	if !balanceExist {
		c.IndentedJSON(http.StatusNotFound, "Balance not found")
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	c.IndentedJSON(
		http.StatusOK,
		map[string]interface{}{
			"balance": userBalance.String(),
			"userID":  userID})
}
