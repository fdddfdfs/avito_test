package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
	"net/http"
)

func unreserveMoney(c *gin.Context) {
	userID, orderID, serviceID, price, err := getOrderData(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid arguments supplied "+err.Error())
		return
	} else if price.LessThan(decimal.Zero) {
		c.IndentedJSON(http.StatusBadRequest, "Price must be positive")
		return
	} else if userID < 0 || orderID < 0 || serviceID < 0 {
		c.IndentedJSON(http.StatusBadRequest, "ID must be positive")
		return
	}

	var commandTag pgx.CommandTag
	commandTag, err = changeReservationStatus(userID, orderID, serviceID, price, "canceled")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	} else if commandTag.RowsAffected() != 1 {
		c.IndentedJSON(http.StatusNotFound, "Reservation doesn`t exist")
		return
	}

	var newBalance decimal.Decimal
	newBalance, err = addBalanceToUser(userID, price)
	if err != nil {
		_, _ = changeReservationStatus(userID, orderID, serviceID, price, "reserved") //change status back

		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"balance": newBalance.String(),
		"userID":  userID})
}
