package main

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"time"
)

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
