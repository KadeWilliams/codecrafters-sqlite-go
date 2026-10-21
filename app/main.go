package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		header := make([]byte, 100)

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		bTreeBuffer := make([]byte, 12)
		if _, err := databaseFile.Seek(100, io.SeekStart); err != nil {
			fmt.Println("Failed to seek:", err)
			return
		}

		_, err = databaseFile.Read(bTreeBuffer)

		var tableCount uint16
		if err := binary.Read(bytes.NewReader(bTreeBuffer[3:5]), binary.BigEndian, &tableCount); err != nil {
			fmt.Println("Failed to read b tree:", err)
			return
		}

		// You can use print statements as follows for debugging, they'll be visible when running tests.
		//fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", tableCount)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
