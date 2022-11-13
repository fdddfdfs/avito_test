package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
	"net/http"
)

func reserveMoney(c *gin.Context) {
	var commandTag pgx.CommandTag
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

	var userBalance decimal.Decimal
	var balanceExist bool
	userBalance, balanceExist, err = getUserBalance(userID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	} else if !balanceExist {
		c.IndentedJSON(http.StatusNotFound, "user not exist")
		return
	} else if userBalance.LessThan(price) {
		c.IndentedJSON(http.StatusBadRequest, "not enough balance")
		return
	}

	var statusID int64
	statusID, err = getStatusID("reserved")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	var tx *pgx.Tx
	tx, err = conn.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error")
		return
	}

	defer tx.Rollback()

	commandTag, err = addOrder(orderID, serviceID, price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	} else if commandTag.RowsAffected() != 1 {
		c.IndentedJSON(http.StatusInternalServerError, "invalid rows number")
		return
	}

	var orderTableID int64
	orderTableID, err = getOrderTableID(orderID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if commandTag, err = updateUserBalance(userID, userBalance.Sub(price)); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	} else if commandTag.RowsAffected() != 1 {
		c.IndentedJSON(http.StatusNotFound, errors.New("user doesnt exist"))
		return
	}

	if commandTag, err = addReservation(userID, orderTableID, statusID); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	} else if commandTag.RowsAffected() != 1 {
		c.IndentedJSON(http.StatusBadRequest, "Cannot add reservation")
		return
	}

	err = tx.Commit()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error")
		return
	}

	err = tx.Commit()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error")
		return
	}

	c.IndentedJSON(http.StatusOK, "Successful operation")
}
