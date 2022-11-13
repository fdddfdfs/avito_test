package main

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/shopspring/decimal"
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

func addReservation(userID int64, orderTableID int64, statusID int64) (pgx.CommandTag, error) {
	query := fmt.Sprintf(
		"INSERT INTO reservations(user_id, order_id, status,reservation_date) "+
			"VALUES (%d,%d,%d,CURRENT_DATE)",
		userID,
		orderTableID,
		statusID)
	return conn.Exec(query)
}
