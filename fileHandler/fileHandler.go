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

// Initialize structs from disk
func NewFH(dataPath string) FileHandler {
	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	fN := []string{}
	for _, f := range files {
		fN = append(fN, f.Name())
	}
	return FileHandler{dataPath, fN}
}

func (fh FileHandler) SyncToFile(json []byte, id int) {
	fN := strconv.Itoa(id) + ".json"
	file, err := os.Create(path.Join(fh.dataPath, fN)) //os.Create --> if already existing, file will be truncated
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	file.Write(json)
}

func (fh FileHandler) ReadFromFile(id int) string {
	fP := path.Join(fh.dataPath, strconv.Itoa(id)+".json")
	file, err := os.Open(fP)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteVal, _ := io.ReadAll(file)
	return string(byteVal)
}

func (fh FileHandler) ReadAll() []string {
	var uStrings []string

	//@TODO Make parallel
	for _, name := range fh.fileNames {
		id, err := strconv.Atoi(strings.Split(name, ".")[0])
		if err != nil {
			log.Fatal(err)
		}
		uStrings = append(uStrings, fh.ReadFromFile(id))
	}
	return uStrings
}
