package fileHandler

import (
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type FileHandler struct {
	dataPath  string
	fileNames []string
}

// NewFH - Initialize structs from disk
//
// dataPath string - Path where files should be stored
//
// return FileHandler - instance of type FileHandler
func NewFH(dataPath string) FileHandler {

	err := os.MkdirAll(dataPath, 777)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	var fN []string
	for _, f := range files {
		fN = append(fN, f.Name())
	}
	return FileHandler{dataPath, fN}
}

// SyncToFile - Write bytes to file
//
// json []byte - ByteArray
//
// id int - file id, e.g.: 1.json
func (fh *FileHandler) SyncToFile(json []byte, id int) {
	fileName := strconv.Itoa(id) + ".json"
	file, err := os.Create(path.Join(fh.dataPath, fileName)) //os.Create --> if already existing, file will be truncated
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fh.fileNames = append(fh.fileNames, fileName)
	_, err = file.Write(json)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadFromFile - Read from a file
//
// id int - file id, e.g.: 1.json
//
// return string - file content as string
func (fh *FileHandler) ReadFromFile(id int) string {
	fP := path.Join(fh.dataPath, strconv.Itoa(id)+".json")
	file, err := os.Open(fP)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	byteVal, _ := io.ReadAll(file)
	return string(byteVal)
}

// ReadAll - Read all files from disk
//
// return []string - array of string representing file contents
func (fh *FileHandler) ReadAll() []string {
	var uStrings []string

	out := make(chan string, 1)
	go func() {
		for _, name := range fh.fileNames {
			id, err := strconv.Atoi(strings.Split(name, ".")[0])
			if err != nil {
				log.Fatal(err)
			}
			out <- fh.ReadFromFile(id)
		}
		close(out)
	}()
	for s := range out {
		uStrings = append(uStrings, s)
	}
	return uStrings
}
