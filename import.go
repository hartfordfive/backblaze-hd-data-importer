package main

import (
	"bytes"
	uuid "code.google.com/p/go-uuid/uuid"
	"flag"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var data_files = []string{
	"https://docs.backblaze.com/public/hard-drive-data/2013_data.zip",
	"https://docs.backblaze.com/public/hard-drive-data/2014_data.zip",
}

var debug, es_port, num_workers, verbose int
var es_host, file_dir, index_prefix string
var csv_files []string
var wg sync.WaitGroup
var csv_files_to_scan = make(chan string, 30000)

// Used for ElasticSearch bulk indexer
var totalBytes = 0
var sets = 0
var pending_docs = 0

func scanFile(fpath string, f os.FileInfo, err error) error {

	fname := filepath.Base(fpath)
	dir, err := filepath.Abs(filepath.Dir(fpath))
	if err != nil {
		fmt.Println("\t", err)
		return err
	}

	// If it's not a file, then return immediately
	file, err := os.Open(fpath)
	defer file.Close()
	if err != nil {
		if debug >= 2 {
			fmt.Println("Error opening file:", err)
		}
		return nil
	}

	finfo, err := file.Stat()
	if err != nil {
		if debug >= 2 {
			fmt.Println("Error getting file stats:", err)
		}
		return nil
	}

	mode := finfo.Mode()

	if err != nil {
		return nil
	}

	// Ensure that it's a valid file type
	ext := path.Ext(fname)

	// Ensure that it's not a directory
	if mode.IsDir() {
		return nil
	} else if ext == ".csv" {
		csv_files = append(csv_files, dir+"/"+fname)
	}

	return nil
}

func loadAndInsertInElasticSearch(indexer *elastigo.BulkIndexer, file_channel chan string, s *sync.WaitGroup) {

	var colum_heading []string

	for {

		select {
		case f := <-file_channel:

			fmt.Printf("Importing %s into ElasticSearch\n", f)

			csv_file, err := NewCsvFile(f)
			defer csv_file.Handle.Close()

			// If error opening file, then continue to next
			if err != nil {
				continue
			}

			for {

				line, read_err := csv_file.ReadLine()

				// If there's a read error, break out and continue to the next file
				if read_err != nil {
					break
				}

				// Reached the end of the file, move on to the next
				if len(line) == 0 && read_err == nil {
					break
				}

				// Use the first line as the column headers, then continue
				if csv_file.CurrLineNum == 0 {
					colum_heading = line
					continue
				}

				date := strings.Split(line[0], "-")
				//t := time.Now()
				capacity_bytes, _ := strconv.ParseInt(line[3], 10, 64)
				failure, _ := strconv.Atoi(line[4])

				obj_normalized := SmartNormalized{}
				obj_raw := SmartRaw{}
				mutable1 := reflect.ValueOf(&obj_normalized).Elem()
				mutable2 := reflect.ValueOf(&obj_raw).Elem()
				i := 0

				for _, col := range colum_heading {

					if i < 5 {
						i += 1
						continue
					}

					parts := strings.Split(col, "_")
					val, _ := strconv.ParseInt(line[i], 10, 64)

					if parts[2] == "normalized" {
						mutable1.FieldByName("Value" + parts[1]).SetInt(val)
					} else if parts[2] == "raw" {
						mutable2.FieldByName("Value" + parts[1]).SetInt(val)
					}

					i += 1
				}

				obj := HDDSpecs{line[0], line[1], line[2], capacity_bytes, failure, obj_normalized, obj_raw}
				_ = indexer.Index(index_prefix+"-"+date[0]+date[1], "hdd-data", uuid.NewUUID().String(), "", nil, obj, true)
				pending_docs += 1
				if pending_docs >= 300 {
					indexer.Flush()
					pending_docs = 0
				}
			}

			if err != nil {
				panic(err)
			}

		}

		fmt.Println("\tDone processing file!")
		s.Done()

	}

}

func queueCsvFileForDbInsert(file string) {
	select {
	case csv_files_to_scan <- file:
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//flag.StringVar(&configFile, "c", "logger.conf", "Load configs from specified config")
	flag.IntVar(&debug, "v", 0, "Start in fmt.Println mode (debug)")
	flag.StringVar(&es_host, "h", "localhost", "ElasticSearch host")
	flag.IntVar(&es_port, "p", 9200, "ElasticSearch port")
	flag.StringVar(&file_dir, "d", "_data", "Directory to scan for files")
	flag.IntVar(&num_workers, "w", -1, "Number of workers")
	flag.StringVar(&index_prefix, "i", "blaze-hd-reliability-data", "Directory to scan for files")
	flag.Parse()

	debug := flag.Lookup("v")
	if debug.Value.String() == "1" {
		verbose = 1
	}

	fmt.Printf("\n---------- Backblaze Hard Drive Data Importer --------\n\n")

	if num_workers == -1 {
		num_workers = runtime.NumCPU()
	}

	if index_prefix == "" {
		fmt.Println("Error: Empty index prefix is invalid!")
		os.Exit(1)
	}

	if es_port < 1 || es_port > 65535 {
		fmt.Println("Error: Port must be between 1 and 65535")
		os.Exit(1)
	}

	// *************** Extract the CSV data files ***************

	// If the following files exist locally, the don't download them

	wg.Add(2)
	for _, f := range data_files {

		remote_file := strings.Split(f, "/")
		local_file := file_dir + "/" + remote_file[len(remote_file)-1]

		if Exists(local_file) == false {
			if verbose == 1 {
				fmt.Println("File doesn't exist, downloading it...")
			}
			Download(f, file_dir+"/")
		}

		// If the extracted directory already exists, then skip the zip file
		if _, err := os.Stat(local_file[0 : len(local_file)-9]); err == nil {
			if verbose == 1 {
				fmt.Printf("Extracted archive %s already exists.  Skipping...\n", f[0:len(f)-9])
			}
			wg.Done()
			continue
		}

		go func(wg *sync.WaitGroup, source_file string, dest_file string) {
			if verbose == 1 {
				fmt.Printf("Extracting file: %s\n", source_file)
			}
			err := Unzip(source_file, dest_file)
			if err != nil {
				fmt.Printf("\tError extracting file %s: %v\n", source_file, err)
			} else {
				fmt.Printf("\tSuccessfully extracted %s to %s\n", source_file, source_file[0:len(source_file)-9])
			}
			wg.Done()
		}(&wg, local_file, file_dir)
	}
	wg.Wait()

	// *************** Now get all the csv files that have been created ***************

	if verbose == 1 {
		fmt.Println("Getting list of CSV files to parse....")
	}

	for _, url := range data_files {

		parts := strings.Split(url, "/")
		fn := parts[len(parts)-1]

		err := filepath.Walk(file_dir+"/"+fn[0:len(fn)-9], scanFile)
		if err != nil {
			fmt.Printf("Error scanning for csv files: %v\n", err)
		}

	}

	fmt.Println("Total Files to Process: ", len(csv_files))

	// ************** Iterate over each CSV file **************

	fmt.Println("Reading CSV files and importing into ES (may take a while)...")
	c := elastigo.NewConn()
	c.Domain = es_host
	c.Port = strconv.Itoa(es_port)
	indexer := c.NewBulkIndexer(4)
	//indexer.BulkMaxDocs = 50
	indexer.Sender = func(buf *bytes.Buffer) error {
		totalBytes += buf.Len()
		sets += 1
		return indexer.Send(buf)
	}
	indexer.Start()

	// Spawn X number of workers to scan open files and insert their data into ElasticSearch
	var wg2 sync.WaitGroup
	wg2.Add(len(csv_files))
	for i := 0; i < num_workers; i++ {
		if verbose == 1 {
			fmt.Printf("Starting queue worker #%d\n", i)
		}
		go loadAndInsertInElasticSearch(indexer, csv_files_to_scan, &wg2)
	}

	num_files_queued_for_insert := 0
	for _, f := range csv_files {
		// Send it on the filesPending channel
		queueCsvFileForDbInsert(f)
		num_files_queued_for_insert += 1
	}

	if verbose == 1 {
		fmt.Printf("Total files queued for insert: %d\n", num_files_queued_for_insert)
	}

	wg2.Wait()
	indexer.Stop()
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
