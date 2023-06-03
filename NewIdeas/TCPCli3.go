package main
 
import (
	"bufio"
	"time"
	"io"
	"log"
	"net"
      "fmt"
	"strings"
)
 
func main() {

      var conRetries = 0
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
			log.Println("[!] Server closed the connection")
		default:
			log.Printf("[!] Server error: %v\n", Rerr)
		}


            // If we got an Error Try to re-connect (10x Max)
		if Rerr != nil {
			conRetries++

			if conRetries > 10 {
				log.Println("[!] Max Connection Retries Exceeded...  Exiting...")
	                  return
			} else {
				log.Printf("[!] Connection to Handler Lost...  Retrying (%d)...\n", conRetries)
				time.Sleep (time.Duration(60) * time.Second)

				con, err = net.Dial("tcp", "0.0.0.0:9999")
				if err != nil {
					log.Printf("[!] Unable to Connect...  Retrying (%d)...\n", conRetries)
				} else {
					defer con.Close()
 					serverReader = bufio.NewReader(con)
				}
			}

		} else {
                  conRetries = 0
			// Should we Quit
			if strings.TrimSpace(serverResponse) == "BYE:" {
				log.Println("[!] Server request to close the connection so closing")
				return
			} else if strings.HasPrefix(strings.ToUpper(serverResponse), "AUTH:") {
				log.Printf("[+] Auth Recieved from Server: %s\n", strings.TrimSpace(serverResponse[5:]))
				log.Printf("[+] Vrfy Sent to Server: %s:%s\n", strings.TrimSpace(serverResponse[5:]), strings.TrimSpace(serverResponse[5:]))

				ConsOut = fmt.Sprintf("Vrfy: %s:%s\n", strings.TrimSpace(serverResponse[5:]), strings.TrimSpace(serverResponse[5:]))
				if _, Werr := con.Write([]byte(ConsOut)); Werr != nil {
					log.Printf("[!] Failed to respond with Vrfy: %v\n", Werr)
				}
			} else {
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
	}
}
