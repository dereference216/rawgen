package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "\x1b[38;5;120mINF\x1b[0m \x1b[38;5;239m> ", log.Ltime)
	WarnLogger = log.New(os.Stdout, "\x1b[38;5;203mWRN\x1b[0m \x1b[38;5;239m> ", log.Ltime)
	ErrorLogger = log.New(os.Stdout, "\x1b[38;5;209mERR\x1b[0m \x1b[38;5;239m> ", log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "\x1b[38;5;221mDBG\x1b[0m \x1b[38;5;239m> ", log.Ltime|log.Lshortfile)
}
