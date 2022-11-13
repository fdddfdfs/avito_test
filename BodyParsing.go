package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func getAddData(c *gin.Context) (int64, decimal.Decimal, error) {
	body, err := getBodyJSON(c)
	if err != nil {
		return 0, decimal.Zero, err
	}

	var userID int64
	userID, err = parseInt(body["userID"])
	if err != nil {
		return 0, decimal.Zero, err
	}

	var amount decimal.Decimal
	amount, err = parseDecimal(body["balance"])
	if err != nil {
		return 0, decimal.Zero, err
	}

	return userID, amount, nil
}

func getOrderData(c *gin.Context) (int64, int64, int64, decimal.Decimal, error) {
	var err error
	body, err := getBodyJSON(c)
	if err != nil {
		return 0, 0, 0, decimal.Zero, err
	}

	var userID, orderID, serviceID int64
	var price decimal.Decimal
	if userID, err = parseInt(body["userID"]); err != nil {
		return 0, 0, 0, decimal.Zero, err
	}

	if orderID, err = parseInt(body["orderID"]); err != nil {
		return 0, 0, 0, decimal.Zero, err
	}

	if serviceID, err = parseInt(body["serviceID"]); err != nil {
		return 0, 0, 0, decimal.Zero, err
	}

	if price, err = parseDecimal(body["price"]); err != nil {
		return 0, 0, 0, decimal.Zero, err
	}

	return userID, orderID, serviceID, price, err
}
