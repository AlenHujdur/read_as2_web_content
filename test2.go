package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func sanitizeString(input string, replacement rune) (string, error) {
	// Create a decoder for the specific encoding (e.g., Windows-1252)
	decoder := charmap.Windows1252.NewDecoder()

	// Convert the input string to bytes
	inputBytes := []byte(input)

	// Decode the bytes using the specified decoder
	decodedBytes, err := decoder.Bytes(inputBytes)
	if err != nil {
		return "", err
	}

	// Convert the decoded bytes back to string
	sanitizedString := string(decodedBytes)

	// Replace non-UTF-8 characters with the replacement character
	//utf8String, err := encoding.UTF8Validator.Transform([]byte(sanitizedString), replacement) //Replace([]byte(sanitizedString), replacement)
	if err != nil {
		return "", err
	}

	//return string(utf8String), nil
	return sanitizedString, nil
}

func sanitizeString2(input string, replacement rune) (string, error) {
	// Create a decoder for the specific encoding (e.g., Windows-1252)
	decoder := charmap.Windows1252.NewDecoder()

	// Convert the input string to bytes
	inputBytes := []byte(input)

	// Decode the bytes using the specified decoder
	decodedBytes, _, err := transform.Bytes(decoder, inputBytes)
	if err != nil {
		return "", err
	}

	// Replace non-UTF-8 characters with the replacement character
	for i := range decodedBytes {
		if decodedBytes[i] >= 0x80 && decodedBytes[i] <= 0xBF {
			decodedBytes[i] = byte(replacement)
		}
	}

	// Convert the decoded bytes back to string
	sanitizedString := string(decodedBytes)

	return sanitizedString, nil
}

// Check if the line ending is valid
func isValidLineEnding(lineEnding string) bool {
	trimmed := strings.TrimSpace(lineEnding)
	return trimmed == "2#EndOfLine#" || strings.HasPrefix(trimmed, "2#EndOfLine#")
}

func replaceInvalidSymbol(text string) string {
	replacedText := strings.ReplaceAll(text, "�", "_")
	return replacedText
}

func sanitizeFile(filename string) error {
	// Read the file content
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Replace "�" with "_"
	modifiedContent := strings.ReplaceAll(string(content), "�", "_")

	// Write the modified content back to the file
	err = ioutil.WriteFile(filename, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Assuming the input string is stored in the 'data' variable

	if sanitizeFile("msg_as2_att.txt") != nil {
		_, err := sanitizeString("msg_as2_att.txt", '_')
		fmt.Println("Failed to sanitize file:", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile("msg_as2_att.txt")
	if err != nil {
		fmt.Println("Failed to read file:", err)
		os.Exit(1)
	}

	// Convert the data to a string
	dataStr := string(data)

	// Split the data into rows based on the newline character
	rows := strings.Split(dataStr, "\n")
	for i, row := range rows {
		// Split the row into columns using the separator ";," or any whitespace
		columns := regexp.MustCompile(`[;\s]*,[;\s]*`).Split(row, -1)
		if len(columns) != 10 {
			fmt.Println("Invalid number of columns")
			os.Exit(1)
		}

		// Check if the line ending is valid
		lineEnd := columns[9]
		if !isValidLineEnding(lineEnd) {
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
		column10 := columns[9]

		// If line ends with "2#EndOfLine#", skip to the next line
		if strings.HasSuffix(lineEnd, "2#EndOfLine#") {
			continue
		}

		fmt.Println("Row", i+1)
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

		// New attempt

		// Split the data into columns based on the separator " ;,"
		columns := strings.Split(data, ";,")

		// Retrieve the last column value
		lastColumn := columns[len(columns)-1]

		// Decode the Base64 encoded content
		decodedContent, err := base64.StdEncoding.DecodeString(lastColumn)
		if err != nil {
			fmt.Println("Failed to decode Base64 content:", err)
			return
		}
		// end of new attempt

		// Decode the base64 encoded file content
		fileContent, err := base64.StdEncoding.DecodeString(column9)
		if err != nil {
			fmt.Println("Failed to decode file content:", err)
			os.Exit(1)
		}

		// Write the file content to a file
		os.Mkdir("files", 0777)
		sanitized1, err := sanitizeString2(column6, '_')
		if err != nil {
			fmt.Println("Failed to sanitize string:", err)
			return
		}
		sanitized_name, err := replaceInvalidSymbol(sanitized1), nil
		files := "files/" + strconv.Itoa(i) + "_" + sanitized_name
		err = ioutil.WriteFile(files, fileContent, 0644)
		if err != nil {
			fmt.Println("Failed to write file:", err)
			os.Exit(1)
		}

		fmt.Println("File generated successfully.")
	}
}
