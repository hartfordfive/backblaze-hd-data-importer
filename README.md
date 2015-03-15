
#Importer for Backblaze Hard Drive Reliability Data Sets

## Description

This is a Go application that imports the open-sourced Backblaze data set on hard drive reliability into Elasticsearch.  The data file archives (ZIP files) are automatically downloaded and extracted if they aren't already present in your current working directory.

You can find out more details regarding this data in the following blog post:

https://www.backblaze.com/blog/hard-drive-data-feb2015/

## Considerations

Please note that this application is most definitely not optimized to it's potential.  Improvements could be made although considering it's only used once to import the data, I might make small improvements here and there although I don't feel the need to spend much time improving it.  In this case, good enough is good enough! (words inspired by Alex Martelli's presentation at OSCON 2013) If you see something is really broken or that could be optimized in the code, feel free to issue a pull request and I'll take a look at it when I have a chance.


## Requirments:

- Go 1.3.3 (elastigo library non functional with 1.4.X)
- Google UUID Library (code.google.com/p/go-uuid/uuid)
- Elastigo Libary (github.com/mattbaird/elastigo/lib)

## Installation

- `git clone https://github.com/hartfordfive/backlaze-hd-data-importer.git`
- `cd backblaze-hd-data-importer`
- `go build`

## Basic Usage

`./backblaze-hd-data-importer -h [ELASTIC_SEARCH_HOST]`

## Options

  -d=".": Directory to scan for files
  -h="localhost": ElasticSearch host
  -i="blaze-hdd-test-data": Directory to scan for files
  -p=9200: ElasticSearch port
  -v=0: Start in fmt.Println mode (debug)
  -w=-1: Number of workers

- `-d [ZIP_DATA_FILES_PATH]` : Location that contains the ZIP data files.  If already present in specified location, they will not be downloaded (default = "_data/")
- `-h [ELASTICSEARCH_HOST]` : The hostname/IP of the Elasticsearch server (default = localhost)
- `-i [INDEX_PREFIX]` : The prefix to use for the indices that will be created ("[INDEX_PREFIX]-[YEAR][MONTH]") (default = blaze-hd-reliability-data)
- `-p [PORT]` : Port whicht to use to connect to Elasticsearch (default = 9200)
- `-v [0|1]` : Set verbose debug mode (default = 0)
- `-w [NUM_WORKERS]` : Sets the number of works to spawn.  A worker imports a single file.  Leave to -1 in order to set number of works to number of cores on your machine. (default = -1)


