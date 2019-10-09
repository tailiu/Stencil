package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"github.com/secsy/goftp"
)


func main() {

	config := goftp.Config{
		User:               "zain",
		Password:           "Robust_Killer007",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	ip := "10.230.12.75"
	port := "21"//"4410"

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