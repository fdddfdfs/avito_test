package main

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/shopspring/decimal"
	"time"
)

func getReservationID(userID int64, orderID int64, serviceID int64, price decimal.Decimal) (int64, error) {
	query := fmt.Sprintf(
		"select id "+
			"from reservations "+
			"where user_id = %d and order_id = ("+
			"select id "+
			"from orders "+
			"where orders.order_id = %d and service_id = %d and price = '%s')",
		userID,
		orderID,
		serviceID,
		price.String())
	var reservationID int64
	err := conn.QueryRow(query).Scan(&reservationID)
	return reservationID, err
}

func getStatusID(status string) (int64, error) {
	var statusID int64
	query := fmt.Sprintf("SELECT id FROM statuses WHERE status = '%s' ", status)
	err := conn.QueryRow(query).Scan(&statusID)

	return statusID, err
}

func getReservationStatusID(reservationID int64) (int64, error) {
	var statusID int64
	query := fmt.Sprintf("SELECT status FROM reservations WHERE id = '%d' ", reservationID)
	err := conn.QueryRow(query).Scan(&statusID)

	return statusID, err
}

func checkIsReservationStatusReserved(reservationID int64) (bool, error) {
	var reservedStatusID int64
	var err error
	reservedStatusID, err = getStatusID("reserved")
	if err != nil {
		return false, err
	}

	var reservationStatus int64
	reservationStatus, err = getReservationStatusID(reservationID)
	if err != nil {
		return false, err
	}

	return reservationStatus == reservedStatusID, nil
}

func updateReservationStatus(reservationID int64, statusID int64) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"UPDATE reservations SET status = %d, change_status_date = CURRENT_DATE WHERE id = %d",
		statusID,
		reservationID)
	return conn.Exec(query)
}

func changeReservationStatus(userID int64,
	orderID int64,
	serviceID int64,
	price decimal.Decimal,
	newStatus string) (pgx.CommandTag, error) {
	var reservationID int64
	var err error
	var commandTag pgx.CommandTag
	reservationID, err = getReservationID(userID, orderID, serviceID, price)
	if err != nil {
		return commandTag, err
	}

	var isReserved bool
	isReserved, err = checkIsReservationStatusReserved(reservationID)
	if err != nil {
		return commandTag, err
	} else if !isReserved {
		return commandTag, errors.New("reservation status is not reserved")
	}

	var processedStatusID int64
	processedStatusID, err = getStatusID(newStatus)
	if err != nil {
		return commandTag, err
	}

	if commandTag, err = updateReservationStatus(reservationID, processedStatusID); err != nil {
		return commandTag, err
	}

	return commandTag, nil
}

func getUserBalance(userID int64) (decimal.Decimal, bool, error) {
	query := fmt.Sprintf("SELECT balance FROM users WHERE user_id = %d", userID)
	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return decimal.Zero, false, err
	}

	var userBalanceString string
	var userBalance decimal.Decimal
	var balanceExist bool

	if rows.Next() {
		err = rows.Scan(&userBalanceString)
		if err != nil {
			return decimal.Zero, false, err
		}

		userBalance, err = decimal.NewFromString(userBalanceString[1:])
		if err != nil {
			return decimal.Zero, false, err
		}

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

func addOrder(orderID int64, serviceID int64, price decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO orders(order_id, service_id, price) values (%d, %d, %s)",
		orderID,
		serviceID,
		price.String())
	return conn.Exec(query)
}

func removeOrder(orderID int64) (pgx.CommandTag, error) {
	query := fmt.Sprintf("DELETE from orders where order_id = %d", orderID)
	return conn.Exec(query)
}

func getOrderTableID(orderID int64) (int64, error) {
	query := fmt.Sprintf("SELECT id FROM orders WHERE order_id = %d", orderID)
	var orderTableID int64
	err := conn.QueryRow(query).Scan(&orderTableID)

	return orderTableID, err
}

func addReservation(userID int64, orderTableID int64, statusID int64) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO reservations(user_id, order_id, status,reservation_date) "+
			"VALUES (%d,%d,%d,CURRENT_DATE)",
		userID,
		orderTableID,
		statusID)
	return conn.Exec(query)
}

func getReportForPeriod(year int64, month int64) ([]ServiceIDRevenue, error) {
	var date string
	if month >= 10 {
		date = fmt.Sprintf("%d-%d-01", year, month)
	} else {
		date = fmt.Sprintf("%d-0%d-01", year, month)
	}
	query := fmt.Sprintf(
		"select service_id, sum(price) "+
			"from orders where id in ( "+
			"select order_id "+
			"from reservations "+
			"where status in (SELECT id from statuses where status = 'processed') "+
			"and (change_status_date between '%s' and end_of_month('%s'))) "+
			"group by service_id",
		date,
		date)

	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var serviceID int64
	var serviceRevenueString string
	var report []ServiceIDRevenue
	for rows.Next() {
		err = rows.Scan(&serviceID, &serviceRevenueString)
		if err != nil {
			return nil, err
		}

		report = append(report, ServiceIDRevenue{serviceID, serviceRevenueString[1:]})
	}

	return report, nil

}

func getServiceName(serviceID int64) (string, error) {
	query := fmt.Sprintf("SELECT service_name FROM services WHERE service_id = %d", serviceID)
	var serviceName string
	err := conn.QueryRow(query).Scan(&serviceName)

	return serviceName, err
}

func addBalanceAddition(userID int64, amount decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO balance_additions(user_id, Amount, addition_date) values (%d, %s, CURRENT_DATE)",
		userID,
		amount.String())
	return conn.Exec(query)
}

func getAllBalanceTransactions(userID int64) ([]Transaction, error) {
	var transactions []Transaction

	additionTransactions, err := getAllBalanceAdditions(userID)
	if err != nil {
		return nil, err
	}

	transactions = append(transactions, additionTransactions...)

	var serviceIDs []int64
	var amounts []string
	var dates []string
	serviceIDs, amounts, dates, err = getAllBalanceReservations(userID, []string{"processed", "reserved"})
	if err != nil {
		return nil, err
	}

	reservationTransactions := getTransactions(serviceIDs, amounts, dates)
	transactions = append(transactions, reservationTransactions...)

	serviceIDs, amounts, dates, err = getAllBalanceReservations(userID, []string{"canceled"})
	if err != nil {
		return nil, err
	}

	canceledTransactions := getTransactions(serviceIDs, amounts, dates)
	transactions = append(transactions, canceledTransactions...)

	return transactions, nil
}

func getAllBalanceAdditions(userID int64) ([]Transaction, error) {
	query := fmt.Sprintf("SELECT Amount, addition_date FROM balance_additions WHERE user_id = %d", userID)
	rows, err := conn.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var amount string
	var date pgtype.Date
	var year, day int
	var month time.Month
	var transactions []Transaction
	for rows.Next() {
		err = rows.Scan(&amount, &date)
		if err != nil {
			return nil, err
		}

		year, month, day = date.Time.Date()

		transactions = append(
			transactions,
			Transaction{
				Amount:          amount[1:],
				Date:            fmt.Sprintf("%d-%d-%d", year, int(month), day),
				TransactionType: "top-up",
				Commentary:      "balance top-up"})
	}

	return transactions, nil
}

func getAllBalanceReservations(userID int64, statuses []string) ([]int64, []string, []string, error) {
	if len(statuses) == 0 {
		return nil, nil, nil, errors.New("statuses len cannot be 0")
	}

	statusesConditions := fmt.Sprintf("where status = '%s'", statuses[0])
	for _, status := range statuses[1:] {
		statusesConditions += fmt.Sprintf(" or status = '%s'", status)
	}
	statusesConditions += "))"
	query := fmt.Sprintf(
		"SELECT orders.service_id,orders.price,reservations.reservation_date "+
			"FROM reservations "+
			"INNER JOIN orders on reservations.order_id = orders.id "+
			"where reservations.order_id in ("+
			"SELECT order_id "+
			"FROM reservations "+
			"WHERE user_id = %d and status in ("+
			"SELECT id "+
			"from statuses "+
			statusesConditions,
		userID)
	rows, err := conn.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, nil, nil, err
	}

	var serviceID int64
	var amount string
	var date pgtype.Date

	var serviceIDs []int64
	var amounts []string
	var dates []string

	var year, day int
	var month time.Month
	for rows.Next() {
		err = rows.Scan(&serviceID, &amount, &date)
		if err != nil {
			return nil, nil, nil, err
		}
		serviceIDs = append(serviceIDs, serviceID)
		amounts = append(amounts, amount[1:])
		year, month, day = date.Time.Date()
		dates = append(dates, fmt.Sprintf("%d-%d-%d", year, int(month), day))
	}

	return serviceIDs, amounts, dates, nil
}

func getTransactions(serviceIDs []int64, amounts []string, dates []string) []Transaction {
	var transactions []Transaction
	var serviceName string
	var err error
	for i := range serviceIDs {
		serviceName, err = getServiceName(serviceIDs[i])
		if err != nil {
			serviceName = "unknown service name"
		}

		transactions = append(
			transactions,
			Transaction{
				Amount:          amounts[i],
				Date:            dates[i],
				TransactionType: "write-off",
				Commentary:      "write off for service named- " + serviceName})
	}
	return transactions
}
