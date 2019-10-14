package migrate
import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"stencil/db"
)

func GetFTPClient() *ftp.ServerConn {
	addr := fmt.Sprintf("%s:%s", db.FTP_SERVER_ADDR, db.FTP_SERVER_PORT)
	
	log.Println("Dialing FTP Server ", addr)
	client, err := ftp.Dial(addr)
	if err != nil {
		log.Println("Unable to connect to FTP server: ", err)
		panic(err)
	}

	log.Println("Trying to log in...")
	if err := client.Login(db.FTP_USER, db.FTP_PASSWORD); err != nil {
		log.Println("FTP Login Authentication Error: ", err)
		panic(err)
	}
	log.Println("FTP Connection Established!")
	return client
}