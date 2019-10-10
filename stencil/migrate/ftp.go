package migrate
import (
	"os"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"stencil/db"
	"strings"
)


func TransferMedia(filePath string) error {

	addr := fmt.Sprintf("%s:%s", db.FTP_SERVER_ADDR, db.FTP_SERVER_PORT)
	
	log.Println("Dialing FTP Server ", addr)
	client, err := ftp.Dial(addr)
	if err != nil {
		log.Println("Unable to connect to FTP server: ", err)
		return err
	}

	log.Println("Trying to log in...")
	if err := client.Login(db.FTP_USER, db.FTP_PASSWORD); err != nil {
		log.Println("FTP Login Authentication Error: ", err)
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Can't open the file at: ", filePath, err)
		return err
	}

	fpTokens := strings.Split(filePath, "/")
	fileName := fpTokens[len(fpTokens)-1]

	if err := client.Stor("/"+fileName, file); err != nil {
		log.Println("File Transfer Failed: ", err)
		return err
	}

	log.Println("Done!")
	return nil
}