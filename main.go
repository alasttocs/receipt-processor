package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	ID            string `json:"id"`
	Retailer      string `json:"retailer"`
	PurchaseDate  string `json:"purchaseDate"`
	PurchaseTime  string `json:"purchaseTime"`
	Items         []Item `json:"items"`
	Total         string `json:"total"`
	Points        int    `json:"points"`
	CalulationErr bool   `json:"calulationErr"` //	Flag to indicate if there was an error in the calculation of the points
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var PointsMap = make(map[string]int)

// command line flags
var debugMode bool
var noAuthMode bool
var logFileName string
var logToFile bool

var logger *log.Logger

func init() {
	logger = log.Default()
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	// handle command line flags
	flag.BoolVar(&debugMode, "debug", false, "Run in debug mode")
	flag.BoolVar(&noAuthMode, "noauth", false, "Run in test mode")
	flag.BoolVar(&logToFile, "log", false, "Enable logging to a file")
	flag.StringVar(&logFileName, "logfile", "logs/logfileAPI.log", "Override log file name")
	flag.Parse()

	if debugMode {
		logger.Println("Running in debug mode, will log debug information")
	}
	if noAuthMode {
		logger.Println("Running in noAuth mode API Keys will not be validated")
	} else {
		// simplified generation of API keys in memory for simplicity of code review
		hashAPIKeys([]string{"key1", "key2", "key3"})
	}
	if logToFile {
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Println("Failed to open log file: ", err)
		}
		defer logFile.Close()
		logger.Println("Logging to file: ", logFileName)
		logger.SetOutput(logFile)
	}
	// create a new router and define routes
	r := mux.NewRouter()
	if !noAuthMode {
		r.Use(validateAPIKey)
	}
	r.HandleFunc("/receipts/process", ProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(BadRoute)
	logger.Println("Server is ready to handle requests.")
	logger.Fatal(http.ListenAndServe(":8080", r))
}

// function to process a reciept generation request
func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		logger.Println("Error decoding request body: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receipt.ID = uuid.New().String()
	receipt.Points = CalculatePoints(receipt)
	PointsMap[receipt.ID] = receipt.Points

	response := struct {
		ID string `json:"id"`
	}{
		ID: receipt.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Println("(Process Receipts) Error encoding response", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// function to look up points for a given receipt
func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recieptID := vars["id"]
	points, recieptFound := PointsMap[recieptID]
	if !recieptFound {
		http.Error(w, "recipet not found", http.StatusNotFound)
		return
	}

	response := struct {
		Points int `json:"points"`
	}{
		Points: points,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Println("(Get Points) Error encoding response", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// function to handle routing errors
func BadRoute(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 not found", http.StatusNotFound)
}

// function to calculate points for given receipt
// Allows for calulation to continue even if there are issues with reciept data,
// receipts is marked as having a calculation error, but the total points are still calculated
func CalculatePoints(receipt Receipt) int {
	points := 0

	points += retailerNamePoints(receipt.Retailer)
	points += receiptTotalPoints(&receipt)
	points += itemPoints(&receipt)
	points += dateAndTimePoints(&receipt)

	if debugMode {
		logger.Println("CalculatePoints: ", points)
	}
	return points
}
