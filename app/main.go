package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
			fmt.Println("database file is corrupted")
		}

		var pageSize uint16
		pageSize, err = iterBuffer(header, 16, 18)
		if err != nil {
			fmt.Println("database file is corrupted")
		}

		buffer := make([]byte, pageSize)
		_, err = databaseFile.ReadAt(buffer, 0)
		if err != nil {
			fmt.Println("database file is corrupted")
		}

		var tableCount uint16
		tableCount, err = iterBuffer(buffer, 103, 105)
		if err != nil {
			fmt.Println("Failed to read value:", err)
		}

		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", tableCount)

		for i := 0; i < int(tableCount); i++ {
			// ptr is the offset of the table pointer in the buffer
			ptr, err := iterBuffer(buffer, 108+i*2, 110+i*2)
			if err != nil {
				fmt.Println("Failed to read value:", err)
			}

			headerStart := ptr + 2
			headerSize := buffer[headerStart]
			firstTypeCode := buffer[headerStart+1]
			secondTypeCode := buffer[headerStart+2]
			start := int(headerStart) + int(headerSize) + int((firstTypeCode-13)/2)
			end := start + int((secondTypeCode-13)/2)
			fmt.Println(string(buffer[start:end]))
		}
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

func iterBuffer(buffer []byte, lower, upper int) (pointer uint16, err error) {
	err = binary.Read(bytes.NewReader(buffer[lower:upper]), binary.BigEndian, &pointer)
	return
}
