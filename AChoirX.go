// ****************************************************************
// "As we benefit from the inventions of others, we should be
//  glad to share our own...freely and gladly."
//  - Benjamin Franklin
//
// AChoirX v0.01 - Convert from C to Go for Xplatform capability
//
// ****************************************************************

package main

import (
    "fmt"
    "time"
    "os"
    "strings"
)


// Global Variable Settings
var Version = "v4.4"
var RunMode = "Run"


// Main Line
func main() {
    // Get Time and Date
    lclTime := time.Now()
    iMonth := int(lclTime.Month())
    iDay := lclTime.Day()
    iYYYY := lclTime.Year()
    iHour := lclTime.Hour()
    iMin := lclTime.Minute()

    // Get Host Name
    cName,err := os.Hostname()
    if err != nil {
        cName = "LocalHost"
    }


    // Build the &ACQ Incident Number
    ACQName := fmt.Sprintf("ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 


    // Default Case Settings 
    caseNumbr := ACQName
    evidNumbr := "001"
    caseDescr := fmt.Sprintf("AChoir Live Acquisition: %s", ACQName)
    caseExmnr := "Unknown"

    for i := 1; i < len(os.Args); i++ {
        if strings.EqualFold(os.Args[i], "/Help") {
            fmt.Printf("\nAChoirX ver: %s, Argument/Options:\n", Version);

            fmt.Printf(" /Help - This Description\n");
            fmt.Printf(" /BLD - Run the Build.ACQ Script (Build the AChoir Toolkit)\n");
            fmt.Printf(" /MNU - Run the Menu.ACQ Script (A Simple AChoir Menu)\n");
            fmt.Printf(" /RUN - Run the AChoir.ACQ Script to do a Live Acquisition\n");
            fmt.Printf(" /DRV:<x:> - Set the &DRV parameter\n");
            fmt.Printf(" /USR:<UserID> - User to Map to Remote Server\n");
            fmt.Printf(" /PWD:<Password> - Password to Map to Remote Server\n");
            fmt.Printf(" /MAP:<Server\\Share> - Map to a Remote Server\n");
            fmt.Printf(" /GET:<URL/File> - Get a File using HTTP.\n");
            fmt.Printf(" /INI:<File Name> - Run the <File Name> script instead of AChoir.ACQ\n");
            fmt.Printf(" /CSE - Ask For Case, Evidence, and Examiner Information\n");
            fmt.Printf(" /CON- Run with Interactive Console Input (Same as /Ini:Console)\n");

            os.Exit(0);
        }
    }



    fmt.Println("AChoirX Version: ", Version)
    fmt.Println("AChoirX RunMode: ", RunMode)

    fmt.Println("Case Number: ", caseNumbr)
    fmt.Println("Evidence Number: ", evidNumbr)
    fmt.Println("Case Description: ", caseDescr)
    fmt.Println("Case Examiner: ", caseExmnr)

}