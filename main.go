package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Note struct {
	Id int
	Content string
}

func ReadData() ([]Note, error) {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		return []Note{}, err
	}
	var todos []Note
	if err = json.Unmarshal(data, &todos); err != nil {
		return []Note{}, err
	}
	return todos, nil
}

func Response(w http.ResponseWriter, responseBody []byte, contentType string, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}

func ReadNoteFromRequest(r *http.Request) (Note, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return Note{}, err
	}
	var jsonData Note
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return Note{}, err
	}
	return jsonData, nil
}

func WriteToFile(content []Note) error {
	output, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("data.json", output, 0666); err != nil {
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			data, err := ReadData()
			if err != nil {
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
			}
			response, _ := json.Marshal(data)
			Response(w, response, "application/json", http.StatusOK)
		case "POST":
			requestBody, err := ReadNoteFromRequest(r)
			if err != nil {
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}
			todos, err := ReadData()
			if err != nil {
				fmt.Fprintf(w, "Failed to get existing data: %s", err.Error())
			}
			todos = append(todos, requestBody)
			if err = WriteToFile(todos); err != nil {
				fmt.Fprintf(w, "Failed to write to data.json: %s", err.Error())
			}
			Response(w, []byte("Successfully wrote data"), "text", http.StatusCreated)
		default:
			Response(w, []byte(fmt.Sprintf("Method %s not supported", r.Method)), "text", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server crashed: %s", err.Error())
	}
}
