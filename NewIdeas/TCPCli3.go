package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
      "fmt"
	"strings"
)
 
func main() {

	var ConsOut = "[+] Console Output"

	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()
 
	serverReader := bufio.NewReader(con)
 
	for {
 
		// Waiting for the server response
		serverResponse, Rerr := serverReader.ReadString('\n')
 
		switch Rerr {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", Rerr)
			return
		}


		// Should we Quit
		if strings.TrimSpace(serverResponse) == "BYE:" {
			log.Println("[!] Server request to close the connection so closing")
			return
		}


		// Write back to server
		ConsOut = fmt.Sprintf("[!] Command Received: %s", serverResponse)
		if _, Werr := con.Write([]byte(ConsOut)); Werr != nil {
			log.Printf("[!] Failed to respond to client: %v\n", Werr)
		}

		ConsOut = fmt.Sprintf("[!] Simulate Response: %s", serverResponse)
		if _, Werr := con.Write([]byte(ConsOut)); Werr != nil {
			log.Printf("[!] Failed to respond to client: %v\n", Werr)
		}
	}
}
