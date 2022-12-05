package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"errors"
)

type Note struct {
	Id int
	Content string
}

type UpdatedContent struct {
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

func ReadNewContentFromRequest(r *http.Request) (UpdatedContent, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return UpdatedContent{}, err
	}
	var jsonData UpdatedContent
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return UpdatedContent{}, err
	}
	return jsonData, nil
}

func UpdateNotesWithId(todos []Note, id int, content string) error {
	for i, note := range todos {
		if note.Id == id {
			todos[i].Content = content
			return nil
		}
	}
	return errors.New("notes cannot be found")
}

func DeleteNotesWithId(todos *[]Note, id int) error {
	index := -1
	for i, note := range *todos {
		if note.Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("notes cannot be found")
	}
	*todos = append((*todos)[:index], (*todos)[index + 1:]...)
	fmt.Printf("%v\n", todos)
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

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
			if err != nil {
				fmt.Fprintf(w, "Error getting id: %s", err.Error())
			}
			todos, err := ReadData()
			if err != nil {
				fmt.Fprintf(w, "Failed to get existing data: %s", err.Error())
			}
			newContent, err := ReadNewContentFromRequest(r)
			if err != nil {
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}
			err = UpdateNotesWithId(todos, id, newContent.Content)
			if err = WriteToFile(todos); err != nil {
				fmt.Fprintf(w, "Failed to write to data.json: %s", err.Error())
			}
			Response(w, []byte(""), "application/json", http.StatusNoContent)
		case "DELETE":
			id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
			if err != nil {
				fmt.Fprintf(w, "Error getting id: %s", err.Error())
			}
			todos, err := ReadData()
			if err != nil {
				fmt.Fprintf(w, "Failed to get existing data: %s", err.Error())
			}
			err = DeleteNotesWithId(&todos, id)
			fmt.Printf("%v\n", todos)
			if err = WriteToFile(todos); err != nil {
				fmt.Fprintf(w, "Failed to write to data.json: %s", err.Error())
			}
			Response(w, []byte(""), "application/json", http.StatusNoContent)
		default:
			Response(w, []byte(fmt.Sprintf("Method %s not supported", r.Method)), "text", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server crashed: %s", err.Error())
	}
}
