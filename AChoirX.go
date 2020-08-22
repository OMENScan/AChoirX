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
    "runtime"
)


// Global Variable Settings
var Version = "v4.4"                            // AChoir Version
var RunMode = "Run"                             // Character Runmode Flag (Build, Run, Menu)
var iRunMode = 0                                // Int Runmode Flag (0, 1, 2)
var inFnam = "AChoir.ACQ"                       // Script Name
var ACQName = "ACQ-IR-LocalHost-00000000-0000"  // AChoir Unique Collection Name
var DiskDrive = "C:"                            // Disk Drive (/DRV)
var iCase = 0                                   // Case Information Processing Mode
var consOrFile = 0                              // Console Input instead of File
var opSystem = "Windows"                        // Which Operating System are we running on
var iopSystem = 0                               // Operating System Flag (0=win, 1=lin, 2=osx, 3=?)

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


    // Get Operating System
    opSystem = runtime.GOOS
    switch opSystem {
    case "windows":
        iopSystem = 0
    case "linux":
        iopSystem = 1
    case "darwin":
        iopSystem = 2
    default:
        iopSystem = 3
    }


    // Initial Settings and Confiiguration
    ACQName = fmt.Sprintf("ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 
    inFnam = "AChoir.ACQ"


    // Default Case Settings 
    caseNumbr := ACQName
    evidNumbr := "001"
    caseDescr := fmt.Sprintf("AChoir Live Acquisition: %s", ACQName)
    caseExmnr := "Unknown"

    // What Directory are we in?
    BaseDir, err := os.Getwd()
    if err != nil {
        BaseDir = fmt.Sprintf("\\AChoir\\ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 

    }

    CurrWorkDir := BaseDir

    // Remove any Trailing Slashes.  This happens if CWD is a
    // mapped network drive (since it is at the root directory)
    if last := len(BaseDir) - 1; last >= 0 && BaseDir[last] == '\\' {
        BaseDir = BaseDir[:last]
    }



    // Start by Parsing any Command Line Parameters
    for i := 1; i < len(os.Args); i++ {
        if strings.EqualFold(os.Args[i], "/Help") {
            fmt.Printf("\nAChoirX ver: %s, Argument/Options:\n", Version)

            fmt.Printf(" /Help - This Description\n")
            fmt.Printf(" /BLD - Run the Build.ACQ Script (Build the AChoir Toolkit)\n")
            fmt.Printf(" /MNU - Run the Menu.ACQ Script (A Simple AChoir Menu)\n")
            fmt.Printf(" /RUN - Run the AChoir.ACQ Script to do a Live Acquisition\n")
            fmt.Printf(" /DRV:<x:> - Set the &DRV parameter\n")
            fmt.Printf(" /USR:<UserID> - User to Map to Remote Server\n")
            fmt.Printf(" /PWD:<Password> - Password to Map to Remote Server\n")
            fmt.Printf(" /MAP:<Server\\Share> - Map to a Remote Server\n")
            fmt.Printf(" /GET:<URL/File> - Get a File using HTTP.\n")
            fmt.Printf(" /INI:<File Name> - Run the <File Name> script instead of AChoir.ACQ\n")
            fmt.Printf(" /CSE - Ask For Case, Evidence, and Examiner Information\n")
            fmt.Printf(" /CON- Run with Interactive Console Input (Same as /Ini:Console)\n")

            os.Exit(0)
        } else if strings.EqualFold(os.Args[i], "/CSE") {
            iCase = 2
        } else if strings.EqualFold(os.Args[i], "/BLD") {
            RunMode = "Bld"
            inFnam = "Build.ACQ"
            iRunMode = 0
        } else if strings.EqualFold(os.Args[i], "/RUN") {
            RunMode = "Run"
            inFnam = "AChoir.ACQ"
            iRunMode = 1
        } else if strings.EqualFold(os.Args[i], "/MNU") {
            RunMode = "Mnu"
            inFnam = "Menu.ACQ"
            iRunMode = 3
        } else if len(os.Args[i]) > 6 && strings.EqualFold(os.Args[i][0:5], "/DRV:") {
            if os.Args[i][6] == ':' {
                DiskDrive = os.Args[i][5:7]
                fmt.Println("[+] Disk Drive Set: ", DiskDrive)
            } else {
                fmt.Println("[!] Invalid Disk Drive Setting: ", os.Args[i][5:])
            }
        } else if strings.EqualFold(os.Args[i], "/CON") {
            consOrFile = 1
            RunMode = "Con"
            inFnam = "Console"
            iRunMode = 1;
        } else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:5], "/INI:") {
            // Check if Input is Console
            if strings.EqualFold(os.Args[i], "/INI:Console") {
                consOrFile = 1;
                RunMode = "Con"
                inFnam = os.Args[i][5:]
                iRunMode = 1;
            } else if len(os.Args[i]) < 254 {
                RunMode = "Ini"
                inFnam = os.Args[i][5:]

                // Initially Set iRunmode to 2 (in case we are running remote)
                // Avoids Creating a Local BACQDIR
                iRunMode = 2;
            } else {
                fmt.Println("[!] /INI: Too Long - Greater than 254 chars")
            }
        }
    }




    // Print Stuff Cause GoLang makes us use variables 
    fmt.Println("[+] AChoirX Version: ", Version)
    fmt.Println("[+] AChoirX RunMode: ", RunMode)
    fmt.Println("[+] Operating System: ", opSystem)


    fmt.Println("[+] Case Number: ", caseNumbr)
    fmt.Println("[+] Evidence Number: ", evidNumbr)
    fmt.Println("[+] Case Description: ", caseDescr)
    fmt.Println("[+] Case Examiner: ", caseExmnr)
    fmt.Println("[+] Script: ", inFnam)
    fmt.Println("[+] Base Dir: ", BaseDir)
    fmt.Println("[+] Curr Work Dir: ", CurrWorkDir)

}