package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	TimeStamp time.Time `json:"timestamp"`
}

type MessageStorage struct {
	messages []Message
	mu       sync.Mutex
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		messages: make([]Message, 0),
	}
}

func (ms *MessageStorage) AddMessage(msg Message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.messages = append(ms.messages, msg)
}

func (ms *MessageStorage) GetMessage() []Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.messages
}

var messageStorage = NewMessageStorage()

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			messages := messageStorage.GetMessage()
			json.NewEncoder(w).Encode(messages)
		} else if r.Method == http.MethodPost {
			var msg Message
			err := json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			msg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
			msg.TimeStamp = time.Now()

			messageStorage.AddMessage(msg)
			w.WriteHeader(http.StatusCreated)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
