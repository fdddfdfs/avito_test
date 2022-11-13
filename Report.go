package main

import "C"
import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"hash/fnv"
	"net/http"
	"time"
)

type ServiceRevenue struct {
	serviceName string
	revenue     string
}

type ServiceIDRevenue struct {
	serviceID int64
	revenue   string
}

func getReport(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")

	if year == "" {
		c.IndentedJSON(http.StatusBadRequest, "Invalid year supplied")
		return
	} else if month == "" {
		c.IndentedJSON(http.StatusBadRequest, "Invalid month supplied")
		return
	}

	var yearInt int64
	var monthInt int64
	var err error
	yearInt, err = parseInt(year)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid year supplied "+err.Error())
		return
	}

	monthInt, err = parseInt(month)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid month supplied "+err.Error())
		return
	}
	var serviceIDRevenues []ServiceIDRevenue
	serviceIDRevenues, err = getReportForPeriod(yearInt, monthInt)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	var serviceRevenues []ServiceRevenue
	serviceRevenues, err = convertToServiceRevenues(serviceIDRevenues)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	var downloadPath string
	var hashedTime string
	var hashInt uint32
	hashInt, err = hash(time.Now().GoString())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
		return
	}

	hashedTime = fmt.Sprintf("%d", hashInt)
	downloadPath = "/downloadReport/" + hashedTime

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"downloadLink": "http://localhost" + downloadPath})

	router.GET(downloadPath, func(c *gin.Context) {

		c.Header("Content-Disposition", "attachment; filename=report_"+year+"_"+month+".csv")
		c.Header("Content-Type", c.GetHeader("Content-Type"))
		c.Header("Content-Length", c.GetHeader("Content-Length"))

		w := csv.NewWriter(c.Writer)
		w.Comma = ';'
		defer w.Flush()

		var records [][]string
		records = append(records, []string{"Service name", "Total revenue"})
		for _, serviceRevenueRow := range serviceRevenues {
			records = append(records, []string{serviceRevenueRow.serviceName, serviceRevenueRow.revenue})
		}

		if err := w.WriteAll(records); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Server error "+err.Error())
			return
		}
	})
}

func convertToServiceRevenues(serviceIDRevenues []ServiceIDRevenue) ([]ServiceRevenue, error) {
	var serviceRevenues []ServiceRevenue
	for _, serviceIDRevenueRow := range serviceIDRevenues {
		var serviceName string
		var err error
		serviceName, err = getServiceName(serviceIDRevenueRow.serviceID)
		if err != nil {
			serviceName = "unknown service name"
		}

		serviceRevenues = append(
			serviceRevenues,
			ServiceRevenue{serviceName, serviceIDRevenueRow.revenue})
	}

	return serviceRevenues, nil
}

func hash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))

	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
