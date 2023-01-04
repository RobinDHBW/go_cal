package configuration

import (
	"flag"
)

var Port int
var Timeout int
var Folder string

func ReadFlags() {
	//flags and configuration of application
	if flag.Lookup("port") == nil {
		flag.IntVar(&Port, "port", 8080, "Define the port for the application. Default: 8080")
	}
	if flag.Lookup("timeout") == nil {
		flag.IntVar(&Timeout, "timeout", 10, "Define a time (in minutes) after the user is logged off due to inactivity. Default: 10 (min)")
	}
	if flag.Lookup("folder") == nil {
		flag.StringVar(&Folder, "folder", "./files", "Define the folder where the user files are stored in the file system, relative to the root directory. Default: ./files")
	}
	flag.Parse()
}
