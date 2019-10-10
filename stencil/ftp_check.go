package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"github.com/secsy/goftp"
	"github.com/jlaffaye/ftp"
)

func testSecsyGoFTP() {

	config := goftp.Config{
		User:               "zain",
		Password:           "Robust_Killer007",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	ip := "10.230.12.76"
	port := "4410"

	// client, err := goftp.Dial(ip+":"+port)
	client, err := goftp.DialConfig(config, ip+":"+port)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Connected!")

	path := "/home/zain/project/resources/1.jpg"

	wd, err := client.Getwd();
	if  err != nil {
		panic(err)
	}

	// fileInfo, err := client.Stat(path);
	// if  err != nil {
	// 	panic(err)
	// }
	// log.Println(fileInfo)
	
	media, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = client.Store(wd, media)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func testJlaffayeFTP() {

	ip := "10.230.12.76"
	port := "4410"

	log.Println("Dialing...")
	client, err := ftp.Dial(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}
	log.Println("Trying to log in...")
	if err := client.Login("zain", "Robust_Killer007"); err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func main() {
	testJlaffayeFTP()
}