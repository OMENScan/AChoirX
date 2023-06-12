package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"fmt"
	"strings"
	"strconv"
	mrand "math/rand"
	"time"
	"encoding/hex"
	"crypto/md5"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
)

var SessArry []int
var SessStat []string
var SessKeys []string
var SessHndl []*bufio.Reader
var SessConn []net.Conn
var SessIPV4 []string
var SessCount = 0 
var CurrSess = -1

func main() {

	go ListenForRequest()

	clientReader := bufio.NewReader(os.Stdin)
 
	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		if err != nil {
			log.Println(err)
			continue
		}

		clientRequestrim := strings.TrimSpace(clientRequest)

		if strings.HasPrefix(strings.ToUpper(clientRequestrim), "SESS:NONE") {
			log.Printf("Session Set To: None\n")
			CurrSess = -1
		} else if strings.HasPrefix(strings.ToUpper(clientRequestrim), "SESS:LIST") {
			for i := 0; i < SessCount; i++ {
				log.Printf("Session: %d - IP Address: %s Status: %s\n", SessArry[i], SessIPV4[i], SessStat[i])
			}
		} else if strings.HasPrefix(strings.ToUpper(clientRequestrim), "SESS:") {
			CurrSess, _ = strconv.Atoi(clientRequestrim[5:])
                  if CurrSess > SessCount-1 {
				log.Printf("No Such Session: %d\n", CurrSess)
				CurrSess = -1
			} else if SessStat[CurrSess] != "Active" {
				log.Printf("Session is not Active: %d\n", CurrSess)
				CurrSess = -1
			} else {
				log.Printf("Session Set To: %d (%s)\n", CurrSess, SessIPV4[CurrSess])
			}
		} else if strings.HasPrefix(strings.ToUpper(clientRequestrim), "KILL:") {
			log.Printf("[!] Exiting AChoirX Multi-Handler... All Remote Sessions Will Terminate\n")
			os.Exit(0)
		} else {

			//log.Println(clientRequestrim)

			// Test to see of the Session Array works...
			if CurrSess > -1 {
				if _, serr := SessConn[CurrSess].Write([]byte(clientRequest)); serr != nil {
					log.Printf("Could Not Connect to Session: %d - %v\n", CurrSess, serr)
				}
			} else {
				log.Printf("No Session Selected - Command Not Sent.\n")

			}
		}
	}
}

func ListenForRequest() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
 
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
 
		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)

	}
}
 
func handleClientRequest(con net.Conn) {
	defer con.Close()

	var ConsOut = "[+] Console Output"

	// Generate Auth String (HostName-Randint)
	cName, host_err := os.Hostname()
	if host_err != nil {
		cName = "LocalHost"
	}

	clientReader := bufio.NewReader(con)
	MyCount := SessCount
        authRand := fmt.Sprintf("%s-%d", cName, mrand.Int())

	SessArry = append(SessArry, SessCount)
	SessHndl = append(SessHndl, clientReader)
	SessConn = append(SessConn, con)
	SessStat = append(SessStat, "Conn")
	SessKeys = append(SessKeys, authRand)
	SessIPV4 = append(SessIPV4, con.RemoteAddr().String())
	
	fmt.Printf("[+] Adding Session: %d For IP Address: %s\n", SessCount, SessIPV4[SessCount])
	//log.Printf("Array Length: %d \n", len(SessArry))
	//log.Printf("Session Count: %d\n", SessCount)
	//log.Printf("Session Item: %d\n", SessArry[SessCount])

	// After Connecting - Send an Auth Request
	fmt.Printf("[+] Sending an Auth Request to the newly connected client...\n")

	// Seed the Random Number Generator
    	mrand.Seed(int64(time.Now().Nanosecond()))

	Srv_Auth := fmt.Sprintf("Auth:%s\n", authRand)
	if _, auth_err := con.Write([]byte(Srv_Auth)); auth_err != nil {
		log.Printf("[!] Failed to Initiate Authorization: %v\n", auth_err)
	}

	// Bump the Session Table for the next Session
	SessCount++ 

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		clientRequest = strings.TrimSpace(clientRequest)

		if strings.HasPrefix(strings.ToUpper(clientRequest), "VRFY:") {
                        AuthVrfy := fmt.Sprintf("%s:%s", authRand, authRand)
			if AuthVrfy == strings.TrimSpace(clientRequest[5:]) {
				log.Printf("[+] Auth Verification Passed!\n")
				SessStat[MyCount] = "Active"
			} else {
				// Invalid Auth String. Terminate!
				log.Printf("[+] Invalid AuthVrfy String: %s\n", strings.TrimSpace(clientRequest[5:]))
				SessStat[MyCount] = "Closed"

				if _, auth_err := con.Write([]byte("BYE:\n")); auth_err != nil {
					log.Printf("[!] Failed to Terminate Remote Session: %v\n", auth_err)
				}
				return
			}
			log.Printf("[+] AuthVrfy String: %s\n", AuthVrfy)
			log.Printf("[+] Vrfy Recieved from Client: %s\n", strings.TrimSpace(clientRequest[5:]))
		} else {
			ConsOut = fmt.Sprintf("%d>>> %s", MyCount, clientRequest)
 
			switch err {
			case nil:
				log.Println(ConsOut)
			case io.EOF:
				log.Println("[!] Client closed the connection by terminating the process")
				SessStat[MyCount] = "Closed"
				return
			default:
				log.Printf("[!] Connection Error: %v\n", err)
            	      SessStat[MyCount] = "Closed"
				return
			}
		}
	}
}

//***************************************************************************
// Encryption Routines: Turn a Password into a 32-bit Key                   *
//***************************************************************************
func createHash(key string) string {
    hasher := md5.New()
    hasher.Write([]byte(key))
    return hex.EncodeToString(hasher.Sum(nil))
}


//***************************************************************************
// Encryption Routines: Encrypt a stream of bytes                           *
//***************************************************************************
func encrypt(data []byte, passphrase string) []byte {

    var ConsOut = "[+] Console Output"

    block, _ := aes.NewCipher([]byte(createHash(passphrase)))
    gcm, enc_err := cipher.NewGCM(block)
    if enc_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Encrypting Data: %s\n", enc_err)
        fmt.Printf(ConsOut)
        return []byte(ConsOut)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, enc_err = io.ReadFull(crand.Reader, nonce); enc_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Encrypting Data: %s\n", enc_err)
        fmt.Printf(ConsOut)
        return []byte(ConsOut)
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext
}


//***************************************************************************
// Encryption Routines: Decrypt a stream of bytes                           *
//***************************************************************************
func decrypt(data []byte, passphrase string) []byte {

    var ConsOut = "[+] Console Output"

    key := []byte(createHash(passphrase))
    block, ciph_err := aes.NewCipher(key)
    if ciph_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data (Cipher): %s\n", ciph_err)
        fmt.Printf(ConsOut)
        return []byte(ConsOut)
    }

    gcm, gcm_err := cipher.NewGCM(block)
    if gcm_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data (GCM): %s\n", gcm_err)
        fmt.Printf(ConsOut)
        return []byte(ConsOut)
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, dec_err := gcm.Open(nil, nonce, ciphertext, nil)
    if dec_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data: %s\n", dec_err)
        fmt.Printf(ConsOut)
        return []byte(ConsOut)
    }

    return plaintext
}
