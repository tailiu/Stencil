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
		User:               "admin",
		Password:           "123456",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		// Logger:             os.Stderr,
	}

	ip := "127.0.0.1"
	port := "2121"

	// client, err := goftp.Dial(ip+":"+port)
	client, err := goftp.DialConfig(config, ip+":"+port)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Connected!")

	path := "/home/user/Downloads/res/1.jpg"

	if  _, err = client.Getwd(); err != nil {
		panic(err)
	}

	// fileInfo, err := client.Stat("/");
	// if  err != nil {
	// 	panic(err)
	// }
	// log.Println(fileInfo)
	
	media, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = client.Store("/", media)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func testJlaffayeFTP() {

	ip := "10.230.12.75"
	port := "21"

	log.Println("Dialing...")
	client, err := ftp.Dial(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}
	log.Println("Trying to log in...")
	if err := client.Login("cowftp", "Big1Fat2Cow3"); err != nil {
		log.Println("Login Authentication Error!")
		panic(err)
	}

	filePath := "/home/user/Downloads/res/1.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Can't open the file at: ",filePath)
		panic(err)
	}

	if err := client.Stor("/1.jpg", file); err != nil {
		log.Println("File Transfer Failed!")
		panic(err)
	}

	log.Println("Done!")
}

func main() {
	testJlaffayeFTP()
}