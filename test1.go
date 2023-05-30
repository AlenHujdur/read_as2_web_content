package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Assuming the input string is stored in the 'data' variable
	// Read data from the file msg_as2_att.txt
	data, err := ioutil.ReadFile("msg_as2_att.txt")
	if err != nil {
		fmt.Println("Failed to read file:", err)
		os.Exit(1)
	}

	// Convert the data to a string
	dataStr := string(data)

	// Split the data into rows based on the newline character
	rows := strings.Split(dataStr, "\n")
	for _, row := range rows {
		// Split the row into columns using the separator ";," or any whitespace
		columns := regexp.MustCompile(`[;\s]*,[;\s]*`).Split(row, -1)
		if len(columns) != 10 {
			fmt.Println("Invalid number of columns")
			os.Exit(1)
		}

		// Check if the line ends with "2#EndOfLine#"
		lineEnd := columns[9]
		if !strings.HasSuffix(lineEnd, "2#EndOfLine#") {
			fmt.Println("Invalid line ending")
			os.Exit(1)
		}

		// Retrieve column values
		column1 := columns[0]
		column2 := columns[1]
		column3 := columns[2]
		column4 := columns[3]
		column5 := columns[4]
		column6 := columns[5]
		column7 := columns[6]
		column8 := columns[7]
		column9 := columns[8]
		column10 := strings.TrimSuffix(lineEnd, "2#EndOfLine#")

		fmt.Println("Column 1:", column1)
		fmt.Println("Column 2:", column2)
		fmt.Println("Column 3:", column3)
		fmt.Println("Column 4:", column4)
		fmt.Println("Column 5:", column5)
		fmt.Println("Column 6:", column6)
		fmt.Println("Column 7:", column7)
		fmt.Println("Column 8:", column8)
		fmt.Println("Column 9:", column9)
		fmt.Println("Column 10:", column10)

		// Decode the base64 encoded file content
		fileContent, err := base64.StdEncoding.DecodeString(column9)
		if err != nil {
			fmt.Println("Failed to decode file content:", err)
			os.Exit(1)
		}

		//make a folder for the files
		os.Mkdir("files", 0777)

		// Write the file content to a file
		err = ioutil.WriteFile(column6, fileContent, 0644)
		if err != nil {
			fmt.Println("Failed to write file:", err)
			os.Exit(1)
		}

		fmt.Println("File generated successfully.")
	}
}
