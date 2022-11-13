package main

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
)

func getUserBalance(userID int64) (decimal.Decimal, bool, error) {
	query := fmt.Sprintf("SELECT balance FROM users WHERE user_id = %d", userID)
	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return decimal.Zero, false, err
	}

	var userBalanceFloat float64
	var userBalance decimal.Decimal
	var balanceExist bool

	if rows.Next() {
		err = rows.Scan(&userBalanceFloat)
		if err != nil {
			return decimal.Zero, false, err
		}

		userBalance = decimal.NewFromFloat(userBalanceFloat).RoundBank(2)

		balanceExist = true
	} else {
		userBalance = decimal.Zero
		balanceExist = false
	}

	return userBalance, balanceExist, nil
}

func updateUserBalance(userID int64, newBalance decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"UPDATE users "+
			"SET balance = %s "+
			"WHERE user_id = %d",
		newBalance.String(),
		userID)
	return conn.Exec(query)
}

func createNewUser(userID int64, balance decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf("INSERT INTO users(user_id,balance) values (%d,%s)", userID, balance.String())
	return conn.Exec(query)
}

func addBalanceToUser(userID int64, amount decimal.Decimal) (decimal.Decimal, error) {
	var commandTag pgx.CommandTag
	var newBalance decimal.Decimal
	var balanceExist bool
	var err error
	var userBalanceValue decimal.Decimal
	userBalanceValue, balanceExist, err = getUserBalance(userID)
	if err != nil {
		return decimal.Zero, err
	}

	if balanceExist {
		commandTag, err = updateUserBalance(userID, userBalanceValue.Add(amount))
		if err != nil {
			return decimal.Zero, err
		} else if commandTag.RowsAffected() != 1 {
			return decimal.Zero, err
		}

		newBalance = userBalanceValue.Add(amount)
	} else {
		commandTag, err = createNewUser(userID, amount)
		if err != nil {
			return decimal.Zero, err
		} else if commandTag.RowsAffected() != 1 {
			return decimal.Zero, errors.New("invalid number of rows")
		}

		newBalance = amount
	}
	return newBalance, nil
}

func addBalanceAddition(userID int64, amount decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO balance_additions(user_id, amount, addition_date) values (%d, %s, CURRENT_DATE)",
		userID,
		amount.String())
	return conn.Exec(query)
}
