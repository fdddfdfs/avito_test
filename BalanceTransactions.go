package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"math"
	"net/http"
	"sort"
)

const TransactionsPerPage int64 = 10

type Transaction struct {
	Amount          string `json:"amount"`
	Date            string `json:"date"`
	TransactionType string `json:"transactionType"`
	Commentary      string `json:"commentary"`
}

type ByDate []Transaction

func (a ByDate) Len() int { return len(a) }
func (a ByDate) Less(i, j int) bool {
	var err error
	var year, month, day int64
	year, month, day, err = dateToInts(a[i].Date)
	if err != nil {
		return false
	}

	var year2, month2, day2 int64
	year2, month2, day2, err = dateToInts(a[j].Date)
	if err != nil {
		return false
	}

	if year != year2 {
		return year < year2
	} else if month != month {
		return month < month2
	} else {
		return day < day2
	}
}
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByPrice []Transaction

func (a ByPrice) Len() int { return len(a) }
func (a ByPrice) Less(i, j int) bool {
	var price, price2 decimal.Decimal
	var err error

	price, err = decimal.NewFromString(a[i].Amount)
	if err != nil {
		return false
	}

	price2, err = decimal.NewFromString(a[j].Amount)
	if err != nil {
		return false
	}

	return price.LessThan(price2)
}
func (a ByPrice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func getBalanceTransactions(c *gin.Context) {
	userIDString := c.Param("userID")
	pageString := c.Param("page")
	sortBy := c.Param("sort")

	if userIDString == "" {
		c.IndentedJSON(http.StatusBadRequest, "invalid userID supplied")
		return
	} else if pageString == "" {
		pageString = "1"
	} else if sortBy == "" {
		sortBy = "pd"
	}

	var userID int64
	var err error
	userID, err = parseInt(userIDString)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "invalid userID supplied"+err.Error())
		return
	} else if userID < 0 {
		c.IndentedJSON(http.StatusBadRequest, "userID must be positive")
		return
	}

	var page int64
	page, err = parseInt(pageString)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "invalid page supplied"+err.Error())
		return
	} else if page <= 0 {
		c.IndentedJSON(http.StatusBadRequest, "page must be greater then 0")
		return
	}

	page -= 1

	var allTransactions []Transaction
	allTransactions, err = getAllBalanceTransactions(userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	if int64(len(allTransactions)) < page*TransactionsPerPage {
		c.IndentedJSON(
			http.StatusBadRequest,
			"Page doest exist. Number of pages: "+
				fmt.Sprintf("%d", int64(math.Ceil(float64(len(allTransactions))/float64(TransactionsPerPage)))))
		return
	}

	if sortBy == "d" || sortBy == "da" {
		sort.Sort(ByDate(allTransactions))
	} else if sortBy == "dd" {
		sort.Sort(sort.Reverse(ByDate(allTransactions)))
	} else if sortBy == "p" || sortBy == "pa" {
		sort.Sort(ByPrice(allTransactions))
	} else {
		sort.Sort(sort.Reverse(ByPrice(allTransactions)))
	}

	if int64(len(allTransactions[page:])) < (page+1)*TransactionsPerPage {
		c.IndentedJSON(http.StatusOK, map[string]interface{}{
			"transactions": allTransactions[page:],
			"page":         page + 1})
	} else {
		c.IndentedJSON(http.StatusOK, map[string]interface{}{
			"transactions": allTransactions[page*TransactionsPerPage : (page+1)*TransactionsPerPage],
			"page":         page + 1})
	}
}
