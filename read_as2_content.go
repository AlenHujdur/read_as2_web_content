package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func fileSize(filename string) int64 {
	println("Getting file size for: " + filename)
	fi, err := os.Stat(filename)
	if err != nil {
		return -1
	}
	return fi.Size()
}

func readObjects(filename string, endline string) {
	println("Reading file: " + filename)
	if filename == "" {
		filename = "msg_as2_att.txt"
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Print("Error opening file: " + filename + "\n")
		panic(err)
	}
	defer f.Close()

	var content [][]string
	var files [][]string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		println("Reading line: " + scanner.Text())
		line := scanner.Text()
		if strings.HasPrefix(line, endline) {
			break
		}
		fmt.Print("Processing line: " + line + "\n")
		content = append(content, strings.Split(line, ";,"))
	}

	for _, line := range content {
		println("Processing line: " + strings.Join(line, ";"))
		for _, element := range line {
			files = append(files, strings.Split(element, ";"))
		}
	}

	i := 0
	tmpFilename := ""
	tmpContent := ""
	tmpFileID := ""
	tmpType := ""

	os.MkdirAll("sha1files", os.ModePerm)
	os.MkdirAll("files", os.ModePerm)

	csvFile, err := os.Create("as2_files.csv")
	if err != nil {
		println("Error creating file")
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"p_msg_as2;type;filename;filesize;sha1hash"})

	for _, file := range files {
		println("Processing file: " + strings.TrimSpace(file[1]) + " " + strings.TrimSpace(file[4]))
		if len(file) > 5 && len(strings.TrimSpace(file[8])) > 76 {
			data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(file[8]))
			if err != nil {
				panic(err)
			}
			filename := filepath.Join("files/", "as2_msg_file_"+strconv.Itoa(i)+"_"+strings.TrimSpace(file[5]))
			err = ioutil.WriteFile(filename, data, 0644)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(filepath.Join("sha1files/", fmt.Sprintf("%x", sha1.Sum(data))), data, 0644)
			if err != nil {
				panic(err)
			}
			writer.Write([]string{strings.TrimSpace(file[1]) + ";" + strings.TrimSpace(file[4]) + ";" + filename + "; " + strconv.FormatInt(fileSize(filename), 10) + "; " + fmt.Sprintf("%x", sha1.Sum(data))})
			tmpFilename = filename
			tmpFileID = strings.TrimSpace(file[1])
			tmpType = strings.TrimSpace(file[4])
			i++
		} else if len(file) > 5 && len(strings.TrimSpace(file[8])) == 76 && strings.TrimSpace(file[1]) != tmpFileID {
			data, err := base64.StdEncoding.DecodeString(tmpContent)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(tmpFilename, data, 0644)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(filepath.Join("sha1files/", fmt.Sprintf("%x", sha1.Sum(data))), data, 0644)
			if err != nil {
				panic(err)
			}
			writer.Write([]string{tmpFileID + ";" + tmpType + ";" + tmpFilename + "; " + strconv.FormatInt(fileSize(tmpFilename), 10) + "; " + fmt.Sprintf("%x", sha1.Sum(data))})
			tmpFilename = filepath.Join("files/", "as2_msg_file_"+strconv.Itoa(i)+"_"+strings.TrimSpace(file[5]))
			tmpContent = strings.TrimSpace(file[8])
			tmpFileID = strings.TrimSpace(file[1])
			tmpType = strings.TrimSpace(file[4])
			i++
		} else if len(strings.TrimSpace(file[0])) <= 76 {
			tmpContent += strings.TrimSpace(file[0])
			fmt.Println("writing content for " + tmpFilename)
		}
	}

	fmt.Println("done")
}

func main() {
	readObjects("msg_as2_att.txt", "#EndOfLine#")
}
