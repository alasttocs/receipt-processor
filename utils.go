package main

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// function to calculate points based on retailer name
// One point for every alphanumeric character in the retailer name
func retailerNamePoints(retailer string) int {
	points := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}
	if debugMode {
		logger.Println("retailer name points", points)
	}
	return points
}

// function to calculate points based on the total amount of the receipt
// 50 points if the total is a round dollar amount with no cents
// 25 points if the total is a multiple of 0.25
// Assumption: overall value of zero should return 0 points
func receiptTotalPoints(receipt *Receipt) int {
	points := 0
	receiptTotal, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		logger.Println("(receiptTotalPoints) Error parsing total: ", err)
		receipt.CalulationErr = true
	} else {
		if receiptTotal == 0 {
			return points
		}
		if math.Mod(receiptTotal, 1.00) == 0 {
			points += 50
			if debugMode {
				logger.Println("receiptTotalPoints(1):", points)
			}
		}
		if math.Mod(receiptTotal, 0.25) == 0 {
			points += 25
			if debugMode {
				logger.Println("receiptTotalPoints(2)", points)
			}
		}
	}
	return points
}

// function to calculate points based on the items on the receipt
// 5 points for every two items on the receipt
// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
func itemPoints(receipt *Receipt) int {
	points := 0

	points += (len(receipt.Items) / 2) * 5
	if debugMode {
		logger.Println("itemPoints(1):", points)
	}
	for _, item := range receipt.Items {
		trimmedDescLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedDescLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				logger.Println("(itemPoints) Error parsing item price: ", err)
				receipt.CalulationErr = true
			} else {
				points += int(math.Ceil(price * 0.2))
				if debugMode {
					logger.Println("itemPoints(2)", points)
				}
			}
		}
	}
	return points
}

// function to calculate points based on the date and time of the receipt
// 6 points if the day in the purchase date is odd
// 10 points if the time of purchase is after 2:00pm and before 4:00pm
// Assume: UTC time
func dateAndTimePoints(receipt *Receipt) int {
	points := 0

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		logger.Println("(dateAndTimePoints) Error parsing purchase date: ", err)
		receipt.CalulationErr = true
	} else {
		if purchaseDate.Day()%2 != 0 {
			points += 6
		}
		if debugMode {
			logger.Println("dateAndTimePoints(1)", points)
		}
	}
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	startTime := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)
	if err != nil {
		logger.Println("(dateAndTimePoints) Error parsing purchase time: ", err)
		receipt.CalulationErr = true
	} else {
		if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
			points += 10
		}
		if debugMode {
			logger.Println("dateAndTimePoints(2)", points)
		}
	}
	return points
}
