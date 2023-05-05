package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	algorithm "src/backend/src/backend/algorithm"
	sql_connection "src/backend/src/backend/sql"
)

type RequestData struct {
	Message string `json:"message"`
	Kmpbm   bool   `json:"kmpbm"`
}

type ResponseData struct {
	Result string `json:"result"`
}

type HistoryData struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ChatHistoryData struct {
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}

func handlerSaveHistory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		// Handle preflight request
		return
	}
	// Extract the message from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var historyData HistoryData
	if err := json.Unmarshal(body, &historyData); err != nil {
		http.Error(w, "Invalid history data", http.StatusBadRequest)
		return
	}

	sql_connection.Create_History(db, historyData.Question, historyData.Answer)
	// Set the response headers

	fmt.Fprint(w, "Success!")

}

func handlerLoadHistory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		// Handle preflight request
		return
	}

	questions, answers := sql_connection.Read_All_History(db)
	chatHistoryData := ChatHistoryData{
		Questions: questions,
		Answers:   answers,
	}

	// Convert the response object to JSON
	responseJSON, err := json.Marshal(chatHistoryData)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")

	// Send the response back
	w.Write(responseJSON)

	// Send a response back
	// fmt.Fprint(w, "Success")
}

func handlerMessage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		// Handle preflight request
		return
	}
	// Extract the message from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestData RequestData
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Print the message
	// fmt.Println(algorithm.ProcessQuestion(db, requestData.Message, true))
	result := algorithm.ProcessQuestion(db, requestData.Message, requestData.Kmpbm)
	responseData := ResponseData{
		Result: result,
	}
	if requestData.Kmpbm {
		fmt.Println("used KMP")
	} else {
		fmt.Println("used BM")
	}

	// Convert the response object to JSON
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")

	// Send the response back
	w.Write(responseJSON)

	// Send a response back
	// fmt.Fprint(w, "Success")
}

func main() {
	data_source := "root:lololol@/"
	db, err := sql.Open("mysql", data_source)

	if err != nil {
		panic(err.Error())
	}

	sql_connection.Create_Database(db)

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		handlerMessage(w, r, db)
	})
	http.HandleFunc("/api/save/history", func(w http.ResponseWriter, r *http.Request) {
		handlerSaveHistory(w, r, db)
	})
	http.HandleFunc("/api/get/history", func(w http.ResponseWriter, r *http.Request) {
		handlerLoadHistory(w, r, db)
	})

	port := ":3001" // Choose any port number you prefer
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

	// println(algorithm.ProcessQuestion(db, "tambahkan pertanyaan papope dengan jawaban poppop", true))
}
