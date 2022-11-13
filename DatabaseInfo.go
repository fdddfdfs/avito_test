package main

import (
	"fmt"
)

func getStatusID(status string) (int64, error) {
	var statusID int64
	query := fmt.Sprintf("SELECT id FROM statuses WHERE status = '%s' ", status)
	err := conn.QueryRow(query).Scan(&statusID)

	return statusID, err
}

func getServiceName(serviceID int64) (string, error) {
	query := fmt.Sprintf("SELECT service_name FROM services WHERE service_id = %d", serviceID)
	var serviceName string
	err := conn.QueryRow(query).Scan(&serviceName)

	return serviceName, err
}
