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
		var pageSize uint16
		pageSize = readDatabaseFile(databaseFile, header, 16, 18, 0)

		var tableCount uint16
		tableCount = readDatabaseFile(databaseFile, header, 3, 5, 100)

		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", tableCount)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

func readDatabaseFile(databaseFile *os.File, byteArray []byte, lowerBound, upperBound, offset int) (value uint16) {
	if _, err := databaseFile.Seek(int64(offset), io.SeekStart); err != nil {
		fmt.Println("Failed to seek:", err)
		return
	}

	_, err := databaseFile.Read(byteArray)
	if err != nil {
		fmt.Println("Failed to read:", err)
		return
	}

	if err := binary.Read(bytes.NewReader(byteArray[lowerBound:upperBound]), binary.BigEndian, &value); err != nil {
		fmt.Println("Failed to read value:", err)
		return
	}

	return value
}
