package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"io"
	"strconv"
)

func parseInt(value interface{}) (int64, error) {
	var parsedValue int64
	var err error

	switch value.(type) {
	case string:
		parsedValue, err = strconv.ParseInt(value.(string), 10, 64)
	case float64:
		parsedValue = int64(value.(float64))
	default:
		return 0, errors.New("cannot get int value")
	}

	return parsedValue, err
}

func parseFloat(value interface{}) (float64, error) {
	var parsedValue float64
	var err error

	switch value.(type) {
	case string:
		parsedValue, err = strconv.ParseFloat(value.(string), 64)
	case float64:
		parsedValue = value.(float64)
	}

	return parsedValue, err
}

func parseDecimal(value interface{}) (decimal.Decimal, error) {
	var parsedValue decimal.Decimal
	var err error = nil

	switch value.(type) {
	case string:
		parsedValue, err = decimal.NewFromString(value.(string))
		parsedValue = parsedValue.RoundBank(2)
	case float64:
		parsedValue = decimal.NewFromFloat(value.(float64)).RoundBank(2)
	default:
		return decimal.Zero, errors.New("cannot get decimal value")
	}

	return parsedValue, err
}

func getBodyJSON(c *gin.Context) (map[string]interface{}, error) {
	body := c.Request.Body

	byteBody, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(byteBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func dateToInts(date string) (int64, int64, int64, error) {
	year, err := parseInt(date[0:4])
	if err != nil {
		return 0, 0, 0, err
	}

	month, err := parseInt(date[0:4])
	if err != nil {
		return 0, 0, 0, err
	}

	day, err := parseInt(date[0:4])
	if err != nil {
		return 0, 0, 0, err
	}

	return year, month, day, nil
}
