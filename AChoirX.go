// ****************************************************************
// "As we benefit from the inventions of others, we should be
//  glad to share our own...freely and gladly."
//  - Benjamin Franklin
//
// AChoirX v0.01 - Convert from C to Go for Xplatform capability
//
//
//
//
// Other Libraries and code I use:
//  Syslog: go get github.com/NextronSystems/simplesyslog
//
// ****************************************************************

package main

import (
    "fmt"
    "time"
    "os"
    "strings"
    "runtime"
    "net/http"
    "io"
    "bufio"
    "crypto/tls"
    syslog "github.com/NextronSystems/simplesyslog"
)


// Global Variable Settings
var Version = "v4.4"                            // AChoir Version
var RunMode = "Run"                             // Character Runmode Flag (Build, Run, Menu)
var ConsOut = "[+] Console Output"              // Console, Log, Syslog strings
var iRunMode = 0                                // Int Runmode Flag (0, 1, 2)
var inFnam = "AChoir.ACQ"                       // Script Name
var ACQName = "ACQ-IR-LocalHost-00000000-0000"  // AChoir Unique Collection Name
var DiskDrive = "C:"                            // Disk Drive (/DRV)
var iCase = 0                                   // Case Information Processing Mode
var consOrFile = 0                              // Console Input instead of File
var opArchit = "AMD64"                          // Architecture
var opSystem = "Windows"                        // Which Operating System are we running on
var iopSystem = 0                               // Operating System Flag (0=win, 1=lin, 2=osx, 3=?)
var slashDelim byte = '\\'                      // Directory Delimiter Win vs. Lin vs. OSX
var WGetURL = "http://Company.com/file"         // URL For HTTP Get
var RootSlash = 0                               // Last Occurance of Slash to find Root URL
var CurrFil = "Current.fil"                     // Current File Name
var inUser = "Joe"                              // UserId
var inPass = "Pa$$w0rd"                         // Password
var Numberz = "0123456789"                      // String to convert from Char to Int
var VarArray[10][256] string                    // Variables Array VR0-VR9
var iVar = -1                                   // Index of the Variable Array
var FullDateTime = "01/01/0001 - 01:01:01"      // Date and Time
var iHtmMode = 0                                // Have we begun writing the HTML Index File

// Input and Output Records
var Conrec = "Console Record"                   // Console Output Record

// Default Case Settings 
var caseNumbr = "ACQ-IR-host-YYYMMDD-HHMM"      // Case Number
var evidNumbr = "001"                           // Evidence Number
var caseDescr = "AChoir Collection"             // Case Description
var caseExmnr = "Unknown"                       // Case Examiner

// Syslog Variables
var Syslogd = "127.0.0.1"                       // Syslog Server 
var Syslogp = "514"                             // Syslog Port
var SyslogTMSG = "AChoir Syslog Started."       // Initialize Syslog Messages 
var SyslogServer = "127.0.0.1:514"              // Syslog Server:Port
var tlsConfig *tls.Config                       // TLS Config

// Message and Log Levels
var iLogOpen = 0                                // Is the LogFile Open Yet
var setMSGLvl = 2                               // Display Message Level - Default=2 (med)
var iSyslogLvl = 0                              // Syslog Level - Default=0 (Off)

// Global File Names
var IniFile = "C:\\AChoir\\AChoir.Acq"          // AChoir Script File
var LogFile = "C:\\AChoir\\LogFile.dat"         // AChoir Log File
var HtmFile = "C:\\AChoir\\Index.htm"           // AChoir HTML Output File
var WGetFile = "C:\\AChoir\\Download.dat"       // Downloaded WGet File
var LstFile = "C:\\AChoir\\Data.Lst"            // List of Data
var ChkFile = "C:\\AChoir\\Data.Chk"            // Check For File Existence
var BACQDir = "C:\\AChoir"                      // Base Acquisition Directory
var CachDir = "C:\\AChoir\\Cache"               // AChoir Caching Directory 
var ForFile = "C:\\AChoir\\Cache\\ForFiles"     // Do action for these Files
var MCpFile = "C:\\AChoir\\Cache\\MCpFiles"     // Do action for Multiple File Copies
var ForDisk = "C:\\AChoir\\Cache\\ForDisk"      // Do Action for Multiple Disk Drives 
var CurrDir = ""                                // Current Directory

// Windows OS Variables
var WinRoot = "NA"                              // Windows Root Directory
var Procesr = "NA"                              // Processor
var TempVar = "NA"                              // Windows Temp Directory
var ProgVar = "NA"                              // Windows Program Files

// Global File Handles
var LogHndl *os.File                            // File Handle for the LogFile
var HtmHndl *os.File                            // File Handle for the HtmFile
var log_err error                               // Logging Errors
var htm_err error                               // HTML Writer Errors

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
    cName, host_err := os.Hostname()
    if host_err != nil {
        cName = "LocalHost"
    }


    // Get Operating System and Architecture
    opArchit = runtime.GOARCH

    opSystem = runtime.GOOS
    switch opSystem {
    case "windows":
        iopSystem = 0
        slashDelim = '\\'
        WinRoot = os.Getenv("SYSTEMROOT")
        Procesr = os.Getenv("PROCESSOR_ARCHITECTURE")
        TempVar = os.Getenv("TEMP")
        ProgVar = os.Getenv("PROGRAMFILES")
    case "linux":
        iopSystem = 1
        slashDelim = '/'
    case "darwin":
        iopSystem = 2
        slashDelim = '/'
    default:
        iopSystem = 3
        slashDelim = '/'
    }


    // Initial Settings and Confiiguration
    ACQName = fmt.Sprintf("ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 
    inFnam = "AChoir.ACQ"


    // Default Case Settings 
    caseNumbr = ACQName
    evidNumbr = "001"
    caseDescr = fmt.Sprintf("AChoir Live Acquisition: %s", ACQName)
    caseExmnr = "Unknown"


    // What Directory are we in?
    BaseDir, cwd_err := os.Getwd()
    if cwd_err != nil {
        BaseDir = fmt.Sprintf("\\AChoir\\ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 

    }

    CurrWorkDir := BaseDir

    // Remove any Trailing Slashes.  This happens if CWD is a mapped network drive (since it is at the root directory)
    // Note: slashDelim was set above based on OS
    if last := len(BaseDir) - 1; last >= 0 && BaseDir[last] == slashDelim {
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
            iRunMode = 1
        } else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:5], "/INI:") {
            // Check if Input is Console
            if strings.EqualFold(os.Args[i], "/INI:Console") {
                consOrFile = 1
                RunMode = "Con"
                inFnam = os.Args[i][5:]
                iRunMode = 1
            } else if len(os.Args[i]) < 254 {
                RunMode = "Ini"
                inFnam = os.Args[i][5:]

                // Initially Set iRunmode to 2 (in case we are running remote)
                // Avoids Creating a Local BACQDIR
                iRunMode = 2
            } else {
                fmt.Println("[!] /INI: Too Long - Greater than 254 chars")
            }
        } else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:5], "/GET:") {
            WGetURL = os.Args[i][5:]
            RootSlash = strings.LastIndexByte(WGetURL, '/')
            if (RootSlash == -1) {
                CurrFil = fmt.Sprintf("%s%cAChoir.Acq", slashDelim, CurrWorkDir)
            } else if len(WGetURL[RootSlash+1:]) < 2 {
                CurrFil = fmt.Sprintf("%s%cAChoir.Acq", slashDelim, CurrWorkDir)
            } else {
                CurrFil = fmt.Sprintf("%s%c%s", CurrWorkDir, slashDelim, WGetURL[RootSlash+1:])
            }

            fmt.Println("[+] HTTP GetFile: ", WGetURL, CurrFil)

            http_err := DownloadFile(CurrFil, WGetURL)
            if http_err != nil {
                fmt.Println("[!] Downloaded Failed: " + WGetURL)        
            } else {
                fmt.Println("[+] Downloaded Success: " + WGetURL)        
            }
	} else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:5], "/USR:") {
            if (os.Args[i][5] =='?') {
                cons_readr := bufio.NewReader(os.Stdin)
                fmt.Print("[?] Enter Share Mapping UserID > ")
                cons_text, _ := cons_readr.ReadString('\n')

                // Trim whitespace and print
                inUser = strings.TrimSpace(cons_text)
            } else {
                inUser = os.Args[i][5:]
            }
        } else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:5], "/PWD:") {
            if (os.Args[i][5] =='?') {
                cons_readr := bufio.NewReader(os.Stdin)
                fmt.Print("[?] Enter Share Mapping Password > ")
                cons_text, _ := cons_readr.ReadString('\n')

                // Trim whitespace and print
                inPass = strings.TrimSpace(cons_text)
            } else {
                inPass = os.Args[i][5:]
            }
        } else if len(os.Args[i]) > 5 && strings.EqualFold(os.Args[i][0:3], "/VR") && (os.Args[i][4] ==':') {
            iVar = strings.IndexByte(Numberz, os.Args[i][3])
            if (iVar == -1) {
                fmt.Println("[!] Invalid Variable: ", os.Args[i][1:4])
            } else if len(os.Args[i]) > 250 {
                fmt.Println("[!] Variable Exceeds 250 Bytes: ", os.Args[i][1:4])
            } else {
                VarArray[iVar][0] = os.Args[i][5:]
            }
        } else {
            fmt.Println("[!] Bad Argument: ", os.Args[i])
        }
    }


    //****************************************************************
    // Set Initial File Names (BaseDir needs to be set 1st)          *
    //****************************************************************
    IniFile = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, inFnam)
    WGetFile = fmt.Sprintf("%s%cAChoir.Dat", BaseDir, slashDelim)
    LstFile = fmt.Sprintf("%s%cLstFiles", BaseDir, slashDelim)
    ChkFile = fmt.Sprintf("%s%cAChoir.exe", BaseDir, slashDelim)

    BACQDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, ACQName)
    CachDir = fmt.Sprintf("%s%c%s%cCache", BaseDir, slashDelim, ACQName, slashDelim)

    ForFile = fmt.Sprintf("%s%cForFiles", CachDir, slashDelim)
    MCpFile = fmt.Sprintf("%s%cMCpFiles", CachDir, slashDelim)
    ForDisk = fmt.Sprintf("%s%cForDisks", CachDir, slashDelim)


    //****************************************************************
    // Create Log Dir if it aint there                               *
    //****************************************************************
    LogFile = fmt.Sprintf("%s%cLogs", BaseDir, slashDelim)
    DirAllocErr(LogFile)


    //****************************************************************
    //* Logging!                                                     *
    //****************************************************************
    LogFile = fmt.Sprintf("%s%cLogs%c%s.Log", BaseDir, slashDelim, slashDelim, ACQName)
    LogHndl, log_err = os.Create(LogFile)

    if log_err != nil {
        fmt.Println("[!] Could not Open Log File.")
        os.Exit(3)
    }

    iLogOpen = 1
  
    ConsOut = fmt.Sprintf("[+] AChoir ver: %s, Mode: %s, OS: %s, Proc: %s\n", Version, RunMode, opSystem, opArchit)
    ConsLogSys(ConsOut, 1, 1)

    showTime("Start Acquisition")

    fmt.Fprintf(LogHndl, "[+] Directory Has Been Set To: %s%c%s\n", BaseDir, slashDelim, CurrDir)
    fmt.Fprintf(LogHndl, "[+] Input Script Set:\n     %s\n\n", IniFile)


    //****************************************************************
    // If iRunMode=1 Create the BACQDir - Base Acquisition Dir       *
    //****************************************************************
    if iRunMode == 1  {
        // Have we created the Base Acquisition Directory Yet?
        ConsOut = fmt.Sprintf("[+] Creating Base Acquisition Directory: %s\n", BACQDir)
        ConsLogSys(ConsOut, 1, 0)

        // Check to see if BACQDir Directory Already Exists
        _, exist_err := os.Stat(BACQDir)
        if os.IsNotExist(exist_err) {
            DirAllocErr(BACQDir)
            DirAllocErr(CachDir)
            PreIndex()
        }
    }

    // Should We Gather Case Information (/CSE)
    if iCase == 2 {
        getCaseInfo(1)
    }





    // Print Stuff Cause GoLang makes us use variables 
    fmt.Println("[+] Windows EnVars: ", WinRoot, Procesr, TempVar, ProgVar)








    // Clean-Up Routines for Testing...
    fmt.Fprintf(LogHndl, "[+] Closing LogFile\n")
    LogHndl.Close()
    // End of Cleanup Testing Routines




}


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
    // Get the data via HTTP or HTTPS Get
    http_resp, http_err := http.Get(url)
    if http_err != nil {
        return http_err
    }

    defer http_resp.Body.Close()

    // Create the file
    http_out, http_err := os.Create(filepath)
    if http_err != nil {
        return http_err
    }

    defer http_out.Close()

    // Write the body to file
    _, http_err = io.Copy(http_out, http_resp.Body)
    return http_err
}


/***********************************************************/
/* Create a Directory - Err (Exit) if it fails             */
/***********************************************************/
func DirAllocErr(DirToCreat string) {
    // Check to see if Directory Already Exists
    _, exist_err := os.Stat(DirToCreat)

    if os.IsNotExist(exist_err) {
        // Try to Create the Directory
        creat_err := os.MkdirAll(DirToCreat, 0644)
        if creat_err != nil {
            if iLogOpen == 1 {
                // Log Any errors if the Log File is Open
                fmt.Fprintf(LogHndl, "[+] Error Creating Directory\n")
            }

            fmt.Println("[!] Error Creating Directory: ", DirToCreat)
            os.Exit(3)
        }
    }
}


func showTime(showText string) {
    //***************************************************************
    // Show the TIME on console and in log                          *
    //***************************************************************
    // Get Local Time and Date
    showlocal := time.Now()


    if showText == "&Tim" {
        FullDateTime = fmt.Sprintf("%02d/%02d/%04d - %02d:%02d:%02d",
        showlocal.Month(), showlocal.Day(), showlocal.Year(),
        showlocal.Hour(), showlocal.Minute(), showlocal.Second())
    } else {
        fmt.Printf("[+] %s: %02d/%02d/%04d - %02d:%02d:%02d\n", showText,
        showlocal.Month(), showlocal.Day(), showlocal.Year(),
        showlocal.Hour(), showlocal.Minute(), showlocal.Second())
    }

    // Only Log if we have opened the Log File.
    if iLogOpen == 1 {
        fmt.Fprintf(LogHndl, "[+] %s: %02d/%02d/%04d - %02d:%02d:%02d\n", showText,
        showlocal.Month(), showlocal.Day(), showlocal.Year(),
        showlocal.Hour(), showlocal.Minute(), showlocal.Second())

    }
}


func AChSyslog(SendLogMSG string) {
    //***************************************************************
    // Send to Syslog                                               *
    //***************************************************************
    // Remove CRLF to prevent Blank Lines in Syslog
    if opSystem == "windows" {
        SendLogMSG = strings.Replace(SendLogMSG, "\r", "", -1)
        SendLogMSG = strings.Replace(SendLogMSG, "\n", "", -1)
    } else {
        SendLogMSG = strings.Replace(SendLogMSG, "\n", "", -1)
    }


    // Not sure why UDP Syslog requires tlsConfig - but it wont compile without it
    SyslogServer = fmt.Sprintf("%s:%s", Syslogd, Syslogp)

    syslog_client, syslog_err := syslog.NewClient(syslog.ConnectionUDP, SyslogServer, tlsConfig)
    if syslog_err != nil {
        // fmt.Println("[!] Problem Defining Syslog Client: ", opSystem)
        return
    }
    defer syslog_client.Close()

    if syslog_err := syslog_client.Send(SendLogMSG, syslog.LOG_LOCAL0|syslog.LOG_NOTICE); syslog_err != nil {
        // fmt.Println("[!] Problem Sending From Syslog Client: ", opSystem)
        return
    }
}


func PreIndex() {
    //***************************************************************
    // Build The Initial Artfact Index.htm                          *
    //***************************************************************
    iHtmMode = 0
    HtmFile = fmt.Sprintf("%s\\Index.htm", BACQDir)

    HtmHndl, htm_err = os.Create(LogFile)
    if htm_err != nil {
        if iLogOpen == 1 {
            fmt.Fprintf(LogHndl, "[!] Could not Create Artifact Index: %s\n", HtmFile);
        }
        if (setMSGLvl > 1) {
            fmt.Printf("[!] Could not Create Artifact Index: %s\n", HtmFile);
        }

        return
    } else {
        iHtmMode = 1

        fmt.Fprintf(HtmHndl, "<html><head><title>AChoir Artifacts</title></head>\n")
        fmt.Fprintf(HtmHndl, "<body>\n")
        fmt.Fprintf(HtmHndl, "<h2>Welcome to AChoir %s</h2>\n\n", Version)
        fmt.Fprintf(HtmHndl, "<p>\n")
        fmt.Fprintf(HtmHndl, "Below is an Index of the Artifacts gathered for Acquisition: <b>%s</b>\n\n", ACQName)
        fmt.Fprintf(HtmHndl, "</p>\n\n")
        fmt.Fprintf(HtmHndl, "<Center><table width=98%%>\n")
        fmt.Fprintf(HtmHndl, "<tr><td align=left>\n")
        fmt.Fprintf(HtmHndl, "<button onclick=\"window.history.back()\">&lt;&lt;</button>\n")
        fmt.Fprintf(HtmHndl, "</td><td align=center>\n")
        fmt.Fprintf(HtmHndl, "<a href=file:./ target=AFrame> Root </a>\n")
    }
}


func ConsLogSys(ConLogMSG string, thisMSGLvl int, thisSyslog int) {
    //***************************************************************
    // Send to Console, Log, and Syslog                             *
    // thisMSGLvl == The message Level of this message              *
    //  0==None, 1==Min, 2==Standard, 3==Max, 4==Debug              *
    // thisSyslog == Should we send to Syslog                       *
    //  0==None, 1==Min, 2==Standard, 3==Max, 4==Debug              *
    //***************************************************************
    if setMSGLvl >= thisMSGLvl && setMSGLvl > 0 {
        fmt.Printf (ConLogMSG)
    }

    if iLogOpen == 1 {
        fmt.Fprintf(LogHndl, ConLogMSG);
    }
    
    if iSyslogLvl >= thisSyslog && iSyslogLvl > 0 {
        AChSyslog(ConLogMSG) 
    }
}


func getCaseInfo(SayOrGet int) {
    // Display or Get Case information
    // Say = 0, Get = 1
    if SayOrGet == 1 {
        // Enter New Case Information
        if iCase == 1 {
            // We ran this routine already once.
            // Avoid confusing multiple Case Names by running only once!
            ConsOut = fmt.Sprintf("[!] Case Information Can Only Be Entered Once.\n")
            ConsLogSys(ConsOut, 1, 3)
        } else {
            ConsOut = fmt.Sprintf("[*] Default Case Number: %s\n", caseNumbr)
            ConsLogSys(ConsOut, 1, 3)

            consInput("Enter New Case Number (Or Enter To Accept Default): ", 1, 0)
            if len(Conrec) > 0 {
                caseNumbr = Conrec
            }

            ConsOut = fmt.Sprintf("[*] Default Case Description: %s\n", caseDescr)
            ConsLogSys(ConsOut, 1, 3)

            consInput("Enter New Case Description (Or Enter to Accept Default: ", 1, 0)
            if len(Conrec) > 0 {
                caseDescr = Conrec
            }

            ConsOut = fmt.Sprintf("[*] Default Evidence Number: %s\n", evidNumbr)
            ConsLogSys(ConsOut, 1, 3)

            consInput("Enter New Evidence Number (Or Enter to Accept Default): ", 1, 0)
            if len(Conrec) > 0 {
                evidNumbr = Conrec
            }

            ConsOut = fmt.Sprintf("[*] Default Examiner: %s\n", caseExmnr)
            ConsLogSys(ConsOut, 1, 3)

            consInput("Enter New Examiner (Or Enter to Accept Default): ", 1, 0)
            if len(Conrec) > 0 {
                caseExmnr = Conrec
            }
        }
    }

    //****************************************************************
    //* Display Case Information                                     *
    //****************************************************************
    ConsOut = fmt.Sprintf("[*] Case Number: %s\n", caseNumbr)
    ConsLogSys(ConsOut, 1, 1)

    ConsOut = fmt.Sprintf("[*] Case Description: %s\n", caseDescr)
    ConsLogSys(ConsOut, 1, 1)

    ConsOut = fmt.Sprintf("[*] Evidence Number: %s\n", evidNumbr)
    ConsLogSys(ConsOut, 1, 1)

    ConsOut = fmt.Sprintf("[*] Examiner: %s\n", caseExmnr)
    ConsLogSys(ConsOut, 1, 1)

    // Run This Routine ONLY ONCE to avoid ambiguity
    iCase = 1;
}


//***************************************************************
// Console Input:                                               *
// conLog  - Should we Log This?                                *
// conHide - Should we redact the Input                         *
//***************************************************************
func consInput(consString string, conLog int, conHide int) {
    ConsOut = fmt.Sprintf("[?] [%s] ", consString)

    if conLog == 1 {
        // Log it Normal
        ConsLogSys(ConsOut, 1, 1)
    } else {
        // Only Log in Debug Mode
        ConsLogSys(ConsOut, 4, 4)
    }

    con_reader := bufio.NewReader(os.Stdin)
    Conrec, _ = con_reader.ReadString('\n')

    // convert CRLF to LF
    if opSystem == "windows" {
        Conrec = strings.Replace(Conrec, "\r", "", -1)
        Conrec = strings.Replace(Conrec, "\n", "", -1)
    } else {
        Conrec = strings.Replace(Conrec, "\n", "", -1)
    }


    if conLog == 1 {
        if conHide == 1 {
            ConsOut = fmt.Sprintf("[?] *Redacted*\n")
        } else {
            ConsOut = fmt.Sprintf("[?] %s\n", Conrec);
        }

        ConsLogSys(ConsOut, 1, 1)

    }
}

