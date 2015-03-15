package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CsvFile struct {
	FileName    string
	Handle      *os.File
	CsvData     *csv.Reader
	LinesRead   int64
	CurrLineNum int64
}

func NewCsvFile(file_name string) (f *CsvFile, err error) {

	h, err := os.Open(file_name)
	if err != nil {
		fmt.Println("CSV file open error:", err)
		return nil, err
	}
	//defer h.Close()
	reader := csv.NewReader(h)
	reader.FieldsPerRecord = -1

	f = &CsvFile{
		Handle:      h,
		FileName:    file_name,
		CsvData:     reader,
		CurrLineNum: -1,
	}
	return f, err
}

func (f *CsvFile) ReadLine() (data []string, err error) {

	//reader := csv.NewReader(file)
	/*
	   reader := csv.NewReader(csvfile)
	   reader.FieldsPerRecord = -1 // see the Reader struct information below
	   data, err = reader.ReadAll()

	   if err != nil {
	     fmt.Println(err)
	     os.Exit(1)
	   }
	   return nil, data
	*/

	// read just one record, but we could ReadAll() as well
	entry, err := f.CsvData.Read()

	if err != nil {
		fmt.Println("\tReadLine error: ", err)
	}

	if err == io.EOF {
		return []string{}, nil
	} else if err != nil {
		return []string{}, err
	}

	f.CurrLineNum += 1
	return entry, nil

}
