package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
)

func addOrder(orderID int64, serviceID int64, price decimal.Decimal) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO orders(order_id, service_id, price) values (%d, %d, %s)",
		orderID,
		serviceID,
		price.String())
	return conn.Exec(query)
}

func getOrderTableID(orderID int64) (int64, error) {
	query := fmt.Sprintf("SELECT id FROM orders WHERE order_id = %d", orderID)
	var orderTableID int64
	err := conn.QueryRow(query).Scan(&orderTableID)

	return orderTableID, err
}
