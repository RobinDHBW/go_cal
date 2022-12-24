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
	fN := []string{} //make([]string, len(files))
	for _, f := range files {
		fN = append(fN, f.Name())
	}
	return FileHandler{dataPath, fN}
}

func (fh FileHandler) SyncToFile(json []byte, id int) {
	fN := strconv.Itoa(id) + ".json"
	file, err := os.Create(path.Join(fh.dataPath, fN))
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
	//var user User

	byteVal, _ := io.ReadAll(file)
	return string(byteVal)
	//json.Unmarshal(byteVal, &user)
	//return user
}

func (fh FileHandler) ReadAll() []string {
	var uStrings []string

	for _, name := range fh.fileNames {
		id, err := strconv.Atoi(strings.Split(name, ".")[0])
		if err != nil {
			log.Fatal(err)
		}
		uStrings = append(uStrings, fh.ReadFromFile(id))
	}
	return uStrings
}
