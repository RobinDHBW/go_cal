package configuration

import (
	"flag"
)

var Port int
var Timeout int
var Folder string
var CertPath string
var KeyPath string

// ReadFlags reads flags from the command line and stores the values in global var
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
	if flag.Lookup("certPath") == nil {
		flag.StringVar(&CertPath, "certPath", "./", "Define path to your ssl cert.pem. Default: ./")
	}
	if flag.Lookup("keyPath") == nil {
		flag.StringVar(&KeyPath, "keyPath", "./", "Define path to your ssl key.pem. Default: ./")
	}
	flag.Parse()
}
