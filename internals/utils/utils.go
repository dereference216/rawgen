package utils

import (
	"io"
	"log"
	"os"
	"sync"
)

func Write(filename string, content string, thread *sync.Mutex) {
	thread.Lock()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if _, err := io.WriteString(file, content+"\n"); err != nil {
		log.Printf("Error writing to file: %v\n", err)
	}
	thread.Unlock()
}
