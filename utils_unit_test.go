// main_test.go

package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logFileName := "logs/testlogfile.log"
	logger = log.Default()
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Println("Failed to open log file: ", err)
	}
	defer logFile.Close()
	logger.Println("Logging to file: ", logFileName)
	logger.SetOutput(logFile)

	retCode := m.Run()

	os.Exit(retCode)
}

func TestRetailerNamePoints(t *testing.T) {
	// empty retailer name
	expectedPoints := 0
	resultPoints := retailerNamePoints("")
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// single word retailer name
	expectedPoints = 6
	resultPoints = retailerNamePoints("abc123")
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// multiple word retailer name
	expectedPoints = 6
	resultPoints = retailerNamePoints("abc 123")
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// untrimmed name
	expectedPoints = 6
	resultPoints = retailerNamePoints("   abc 123   ")
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// name with special character
	expectedPoints = 6
	resultPoints = retailerNamePoints("abc!!$123")
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
}

func TestReceiptTotalPoints(t *testing.T) {
	receipt := &Receipt{}
	// zero total
	receipt.Total = "0"
	expectedPoints := 0
	resultPoints := receiptTotalPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// receipt total with zero cents qualify for 50 points & 25 points
	receipt.Total = "1.00"
	expectedPoints = 75
	resultPoints = receiptTotalPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// receipt total qualify for 25 points
	receipt.Total = "1.75"
	expectedPoints = 25
	resultPoints = receiptTotalPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// reciept total no points
	receipt.Total = "1.40"
	expectedPoints = 0
	resultPoints = receiptTotalPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// receipt total with parsing error (empty total)
	receiptErr := &Receipt{}
	receiptErr.Total = ""
	expectedPoints = 0
	resultPoints = receiptTotalPoints(receiptErr)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d -", expectedPoints, resultPoints)
	}
	if !receiptErr.CalulationErr {
		t.Errorf("Expected calculation error = true, got %v", receiptErr.CalulationErr)
	}

}

func TestItemPoints(t *testing.T) {
	receipt := &Receipt{}
	// single item description qualifies for multiple of 3
	expectedPoints := 1
	receipt.Items = []Item{
		{
			ShortDescription: "abc123",
			Price:            "1.00",
		},
	}
	resultPoints := itemPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}

	// single item description does not qualify for multiple of 3
	expectedPoints = 0
	receipt.Items = []Item{
		{
			ShortDescription: "abc1234",
			Price:            "1.00",
		},
	}
	resultPoints = itemPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}

	// multiple item description, qualify for multiple of 3, qualify for a pair
	// check spacing and trimming
	receipt.Items = []Item{
		{
			ShortDescription: "abc 12",
			Price:            "1.00",
		},
		{
			ShortDescription: "   abc123   ",
			Price:            "1.00",
		},
	}
	expectedPoints = 7
	resultPoints = itemPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}

	// multiple item description, only one qualify for multiple of 3, qualify for a pair
	receipt.Items = []Item{
		{
			ShortDescription: "abc123",
			Price:            "1.00",
		},
		{
			ShortDescription: "abc1234",
			Price:            "1.00",
		},
	}
	expectedPoints = 6
	resultPoints = itemPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}

	// three items, qualify for a pair, none for rule of 3
	receipt.Items = []Item{
		{
			ShortDescription: "abc1234",
			Price:            "1.00",
		},
		{
			ShortDescription: "abc1234",
			Price:            "1.00",
		},
		{
			ShortDescription: "abc1234",
			Price:            "1.00",
		},
	}
	expectedPoints = 5
	resultPoints = itemPoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// empty item price
	errReceipt := &Receipt{}
	errReceipt.Items = []Item{
		{
			ShortDescription: "abc123",
			Price:            "",
		},
	}
	expectedPoints = 0
	resultPoints = itemPoints(errReceipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	if !errReceipt.CalulationErr {
		t.Errorf("Expected calculation error = true, got %v", errReceipt.CalulationErr)
	}
}

func TestDateAndTimePoints(t *testing.T) {
	receipt := &Receipt{}
	// odd day, outside time window
	receipt.PurchaseDate = "2024-01-01"
	receipt.PurchaseTime = "11:00"
	expectedPoints := 6
	resultPoints := dateAndTimePoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// even day, outside time window
	receipt.PurchaseDate = "2024-01-02"
	receipt.PurchaseTime = "11:00"
	expectedPoints = 0
	resultPoints = dateAndTimePoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// odd day, inside time window
	receipt.PurchaseDate = "2024-01-01"
	receipt.PurchaseTime = "15:00"
	expectedPoints = 16
	resultPoints = dateAndTimePoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// even day, inside time window
	receipt.PurchaseDate = "2024-01-02"
	receipt.PurchaseTime = "15:00"
	expectedPoints = 10
	resultPoints = dateAndTimePoints(receipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	// parsing error in date
	errReceipt := &Receipt{}
	errReceipt.PurchaseDate = "2024-01-00"
	errReceipt.PurchaseTime = "11:00"
	expectedPoints = 0
	resultPoints = dateAndTimePoints(errReceipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	if !errReceipt.CalulationErr {
		t.Errorf("Expected calculation error = true, got %v", errReceipt.CalulationErr)
	}
	// parsing error in time
	errReceipt.PurchaseDate = "2024-01-02"
	errReceipt.PurchaseTime = "25:00"
	errReceipt.CalulationErr = false
	expectedPoints = 0
	resultPoints = dateAndTimePoints(errReceipt)
	if resultPoints != expectedPoints {
		t.Errorf("Expected %d, got %d", expectedPoints, resultPoints)
	}
	if !errReceipt.CalulationErr {
		t.Errorf("Expected calculation error = true, got %v", errReceipt.CalulationErr)
	}
}
