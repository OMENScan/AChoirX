package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)
 
func main() {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()
 
	serverReader := bufio.NewReader(con)
 
	for {
 
		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')
 
		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}