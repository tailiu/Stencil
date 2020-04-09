package migrate_v1

import (
	"fmt"
	"log"
	"stencil/db"

	"github.com/gookit/color"
	"github.com/jlaffaye/ftp"
)

func GetFTPClient() *ftp.ServerConn {
	addr := fmt.Sprintf("%s:%s", db.FTP_SERVER_ADDR, db.FTP_SERVER_PORT)

	color.Info.Println("Dialing FTP Server ", addr)
	client, err := ftp.Dial(addr)
	if err != nil {
		log.Println("Unable to connect to FTP server: ", err)
		panic(err)
	}

	color.Info.Print("Trying to log in...")
	if err := client.Login(db.FTP_USER, db.FTP_PASSWORD); err != nil {
		color.Error.Println(" FTP Login Authentication Error: ", err)
		panic(err)
	}
	color.Info.Println(" FTP Connection Established!")
	return client
}
