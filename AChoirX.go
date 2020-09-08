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
// Changes from AChoir:
//  Environment Variable Expansion now uses GoLang $Var or ${Var} 
//
//
// ****************************************************************

package main

import (
    "fmt"
    "time"
    "os"
    "strings"
    "text/scanner"
    "encoding/csv"
    "regexp"
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
var ForSlash = 0                                // Last Occurance of Slash to find File in Path
var CurrFil = "Current.fil"                     // Current File Name
var inUser = "Joe"                              // UserId
var inPass = "Pa$$w0rd"                         // Password
var Numberz = "0123456789"                      // String to convert from Char to Int
var VarArray[10][256] string                    // Variables Array VR0-VR9
var iVar = -1                                   // Index of the Variable Array
var FullDateTime = "01/01/0001 - 01:01:01"      // Date and Time
var iHtmMode = 0                                // Have we begun writing the HTML Index File
var RunMe = 0                                   // Used to Track Conditional Run Routines
var Looper = 0                                  // Flag to Keep Looping 
var ForMe = 0                                   // Flag to identify &For is being used
var LstMe = 0                                   // Flag to identify &LST is being used
var DskMe = 0                                   // Flag to identify &DSK is being used
var LoopNum = 0                                 // Loop Counter
var ForFName = "File.txt"                       // Parsed File name from Path
var iNative = 0                                 // Are we Native 64Bit on 64Bit (Native = 1, NonNative = 0)
var sNative = ""                                // Native String (blank or Non-) 
var iIsAdmin = 0                                // Are we an Admin 
var sIsAdmin = ""                               // Are we an Admin String (blank or Non-) 

//Tokenize Records
var tokRec scanner.Scanner                      // Used to Tokenize Records into Slices
var tokCount = 0                                // Token Counter
var Delims = ","                                // Tokenizing Delimiters

// Input and Output Records
var Conrec = "Console Record"                   // Console Output Record
var Tmprec = "Formatted Console Record"         // Console Formatting Record
var Filrec = "File Record"                      // File Record 
var Lstrec = "File Record"                      // List Record 
var Dskrec = "File Record"                      // Disk Record 
var Inprec = "Console Input"                    // Console Input Record 
var o32VarRec = "32 bit Variables"              // 32 Bit Variable Expansion Record
var o64VarRec = "64 bit Variables"              // 64 Bit Variable Expansion Record

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
var ACQDir = ""                                 // Relative Acquisition Directory
var CachDir = "C:\\AChoir\\Cache"               // AChoir Caching Directory 
var ForFile = "C:\\AChoir\\Cache\\ForFiles"     // Do action for these Files
var MCpFile = "C:\\AChoir\\Cache\\MCpFiles"     // Do action for Multiple File Copies
var ForDisk = "C:\\AChoir\\Cache\\ForDisk"      // Do Action for Multiple Disk Drives 
var CurrDir = ""                                // Current Directory
var TempDir = ""                                // Temp Build Directory String
var TempAcq = ""                                // Temp ACQ Directory String

// Windows OS Variables
var WinRoot = "NA"                              // Windows Root Directory
var Procesr = "NA"                              // Processor
var TempVar = "NA"                              // Windows Temp Directory
var ProgVar = "NA"                              // Windows Program Files

// Global File Handles
var IniScan *bufio.Scanner                      // IO Reader for File or Console
var ForScan *bufio.Scanner                      // IO Reader for ForFile
var LstScan *bufio.Scanner                      // IO Reader for LstFile
var DskScan *bufio.Scanner                      // IO Reader for DskFile
var LogHndl *os.File                            // File Handle for the LogFile
var HtmHndl *os.File                            // File Handle for the HtmFile
var IniHndl *os.File                            // File Handle for the IniFile
var ForHndl *os.File                            // File Handle for the ForFile
var LstHndl *os.File                            // File Handle for the LstFile
var DskHndl *os.File                            // File Handle for the DskFile
var log_err error                               // Logging Errors
var htm_err error                               // HTML Writer Errors
var ini_err error                               // Ini File Errors
var for_err error                               // For File Errors
var lst_err error                               // Lst File Errors
var dsk_err error                               // Dsk File Errors
var for_rcd bool                                // Return Code for ForFile Read
var lst_rcd bool                                // Return Code for LstFile Read
var dsk_rcd bool                                // Return Code for DskFile Read

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
    caseDescr = fmt.Sprintf("AChoirX Live Acquisition: %s", ACQName)
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
        if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/HELP") {
            fmt.Printf("\nAChoirX ver: %s, Argument/Options:\n", Version)

            fmt.Printf(" /Help - This Description\n")
            fmt.Printf(" /BLD - Run the Build.ACQ Script (Build the AChoirX Toolkit)\n")
            fmt.Printf(" /MNU - Run the Menu.ACQ Script (A Simple AChoirX Menu)\n")
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
        } else if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/CSE") {
            iCase = 2
        } else if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/BLD") {
            RunMode = "Bld"
            inFnam = "Build.ACQ"
            iRunMode = 0
        } else if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/RUN") {
            RunMode = "Run"
            inFnam = "AChoir.ACQ"
            iRunMode = 1
        } else if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/MNU") {
            RunMode = "Mnu"
            inFnam = "Menu.ACQ"
            iRunMode = 3
        } else if len(os.Args[i]) > 6 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/DRV:") {
            if os.Args[i][6] == ':' {
                DiskDrive = os.Args[i][5:7]
                fmt.Println("[+] Disk Drive Set: ", DiskDrive)
            } else {
                fmt.Println("[!] Invalid Disk Drive Setting: ", os.Args[i][5:])
            }
        } else if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/CON") {
            consOrFile = 1
            RunMode = "Con"
            inFnam = "Console"
            iRunMode = 1
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/INI:") {
            // Check if Input is Console
            if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/INI:Console") {
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
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/GET:") {
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
	} else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/USR:") {
            if (os.Args[i][5] =='?') {
                cons_readr := bufio.NewReader(os.Stdin)
                fmt.Print("[?] Enter Share Mapping UserID > ")
                cons_text, _ := cons_readr.ReadString('\n')

                // Trim whitespace and print
                inUser = strings.TrimSpace(cons_text)
            } else {
                inUser = os.Args[i][5:]
            }
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/PWD:") {
            if (os.Args[i][5] =='?') {
                cons_readr := bufio.NewReader(os.Stdin)
                fmt.Print("[?] Enter Share Mapping Password > ")
                cons_text, _ := cons_readr.ReadString('\n')

                // Trim whitespace and print
                inPass = strings.TrimSpace(cons_text)
            } else {
                inPass = os.Args[i][5:]
            }
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/VR") && (os.Args[i][4] ==':') {
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
  
    ConsOut = fmt.Sprintf("[+] AChoirX ver: %s, Mode: %s, OS: %s, Proc: %s\n", Version, RunMode, opSystem, opArchit)
    ConsLogSys(ConsOut, 1, 1)

    showTime("Start Acquisition")

    fmt.Fprintf(LogHndl, "[+] Directory Has Been Set To: %s%c%s\n", BaseDir, slashDelim, CurrDir)
    fmt.Fprintf(LogHndl, "[+] Input Script Set:\n     %s\n\n", IniFile)


    //****************************************************************
    //* Are we running Non-Native (Sysnative vs. System32)           *
    //****************************************************************
    if iopSystem == 0 {
        TempDir = fmt.Sprintf("%s\\Sysnative", WinRoot)

        if _, fol_err := os.Stat(TempDir); os.IsNotExist(fol_err) {
            sNative = "64Bit "
            iNative = 1
        } else {
            sNative = "32Bit NON-"
            iNative = 0
        }


        //****************************************************************
        //* Check If We are an Admin                                     *
        //****************************************************************
        if (IsUserAdmin()) {
            iIsAdmin = 1
            sIsAdmin = ""
        } else {
            iIsAdmin = 0;
            sIsAdmin = "NON-"
        }

        ConsOut = fmt.Sprintf("[+] Running as Windows %sNative (%sAdmin)\n", sNative, sIsAdmin)
        ConsLogSys(ConsOut, 1, 1)
    }


    //****************************************************************
    // If iRunMode=1 Create the BACQDir - Base Acquisition Dir       *
    //****************************************************************
    if iRunMode == 1  {
        // Have we created the Base Acquisition Directory Yet?
        ConsOut = fmt.Sprintf("[+] Creating Base Acquisition Directory: %s\n", BACQDir)
        ConsLogSys(ConsOut, 1, 3)

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


    if consOrFile == 1 {
        ConsOut = fmt.Sprintf("[+] Switching to Console Input.\n")
        ConsLogSys(ConsOut, 1, 1)
        fmt.Printf(">>>");

        IniScan = bufio.NewScanner(os.Stdin)

    } else {
        IniHndl, ini_err = os.Open(IniFile)
        if ini_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Ini File: %s\n", IniFile)
            ConsLogSys(ConsOut, 1, 2)
        }

        IniScan = bufio.NewScanner(IniHndl)
        
    }

    for IniScan.Scan() {
        RunMe = 0;  // Conditional run Script default is yes

        //Remove any preceding blanks
        Tmprec = strings.TrimSpace(IniScan.Text())

        // Dont Process any Comments
        if strings.HasPrefix(Tmprec, "*") {
            continue
        }


        //****************************************************************
        //* Conditional Execution                                        *
        //****************************************************************
        if RunMe > 0 {
            if strings.HasPrefix(strings.ToUpper(Tmprec), "32B:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "64B:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "VER:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "CKY:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "CKN:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "EQU:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "NEQ:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "RC=:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "RC!:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "RC>:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "RC<:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "N>>:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "N<<:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "N==:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "END:") {
                RunMe--
            }
        } else {
            Looper = 1;

            //****************************************************************
            //* ForFiles Looper Setup                                        *
            //****************************************************************
            if strings.HasPrefix(strings.ToUpper(Tmprec), "&FOR") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FO0") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FO1") ||
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FO2") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FO3") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FO4") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FO5") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FO6") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FO7") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FO8") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FO9") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOA") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOB") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOC") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOD") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOE") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOF") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOG") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOH") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOI") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOJ") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOK") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOL") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOM") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FON") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&FOO") || strings.HasPrefix(strings.ToUpper(Tmprec), "&FOP") {
                ForMe = 1

                ForHndl, for_err = os.Open(ForFile)

                if for_err != nil {
                    ConsOut = fmt.Sprintf("[!] &FOR Directory has not been set with the FOR: command.  Ignoring &FOR Loop...\n")
                    ConsLogSys(ConsOut, 1, 2)

                    Looper = 0
                }

                ForScan = bufio.NewScanner(ForHndl)
            } else {
                ForMe = 0
            }


            //****************************************************************
            //* LstFiles Looper Setup                                        *
            //****************************************************************
            if strings.HasPrefix(strings.ToUpper(Tmprec), "&LST") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LS0") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LS1") ||
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LS2") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LS3") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LS4") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LS5") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LS6") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LS7") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LS8") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LS9") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSA") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSB") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSC") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSD") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSE") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSF") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSG") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSH") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSI") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSJ") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSK") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSL") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSM") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSN") || 
            strings.HasPrefix(strings.ToUpper(Tmprec), "&LSO") || strings.HasPrefix(strings.ToUpper(Tmprec), "&LSP") {
                LstMe = 1;

                LstHndl, lst_err = os.Open(LstFile)

                if lst_err != nil {
                    ConsOut = fmt.Sprintf("[!] &LST File was not found (LST: not set): %s\n", LstFile)
                    ConsLogSys(ConsOut, 1, 2)

                    Looper = 0
                }

                LstScan = bufio.NewScanner(LstHndl)
            } else {
                LstMe = 0
            }


            //****************************************************************
            //* DskFiles Looper Setup                                        *
            //****************************************************************
            if strings.HasPrefix(strings.ToUpper(Tmprec), "&DSK") {
                DskMe = 1

                DskHndl, dsk_err = os.Open(ForDisk)

                if dsk_err != nil {
                    ConsOut = fmt.Sprintf("[!] &DSK Listing was not found (DSK: not set): %s\n", ForDisk)
                    ConsLogSys(ConsOut, 1, 2)

                    Looper = 0
                }

                DskScan = bufio.NewScanner(DskHndl)
            } else {
                DskMe = 0
            }


            //****************************************************************
            //* Loop (FOR: and LST:) until Looper = 1                        *
            //****************************************************************
            LoopNum = 0;
            for Looper > 0 {
                if ForMe == 0 && LstMe == 0 && DskMe == 0 {
                    Looper = 0
                } else if ForMe == 1 && LstMe == 0 && DskMe == 0 {
                    for_rcd = ForScan.Scan()
                    for_err = ForScan.Err()

                    // No Error and no EOF - So Process the Record
                    if for_err == nil && for_rcd == true {
                        Filrec = strings.TrimSpace(ForScan.Text())

                        Looper = 1
                        LoopNum++

                        //****************************************************************
                        //* Get Just the File Name                                       *
                        //****************************************************************
                        ForSlash = strings.LastIndexByte(Filrec, slashDelim)
                        if (ForSlash == -1) {
                            ForFName = Filrec
                        } else if len(Filrec[RootSlash+1:]) < 2 {
                            ForFName = "Unknown"
                        } else {
                            ForFName = Filrec[RootSlash+1:]
                        }

                    } else {
                        break;
                    }
                } else if ForMe == 0 && LstMe == 1 && DskMe == 0 {
                    lst_rcd = LstScan.Scan()
                    lst_err = LstScan.Err()

                    // No Error and no EOF - So Process the Record
                    if lst_err == nil && lst_rcd == true {
                        Lstrec = strings.TrimSpace(LstScan.Text())

                        Looper = 1
                        LoopNum++
                    } else {
                        break
                    }
                } else if ForMe == 0 && LstMe == 0 && DskMe == 1 {
                    dsk_rcd = DskScan.Scan()
                    dsk_err = DskScan.Err()

                    // No Error and no EOF - So Process the Record
                    if dsk_err == nil && dsk_rcd == true {
                        Dskrec = strings.TrimSpace(DskScan.Text())

                        Looper = 1
                        LoopNum++
                    } else {
                        break
                    }
                } else {
                    Looper = 0

                    ConsOut = fmt.Sprintf("[!] AChoirX does not yet support Nested Looping (&LST + &FOR)\n     > %s\n", Tmprec)
                    ConsLogSys(ConsOut, 1, 2)

                    Tmprec = fmt.Sprintf("***: Command Bypassed")
                }


                //****************************************************************
                //* Check for System Variables and Expand them                   *
                //****************************************************************
                //* This function changes in AChoirX - AChoir uses %EnVar%       *
                //* but native GOLang support $Var or ${Var}.  AChoirS now uses  *
                //* the Native GoLang functions to prevent reinventing the wheel *
                //****************************************************************
                varConvert(Tmprec)


                //****************************************************************
                //* Now Further expand o32VarRec for Achoir unique variables     *
                //****************************************************************
                if CaseInsensitiveContains(o32VarRec, "&Dir") {

                    if len(CurrDir) > 0 {
                        TempDir = fmt.Sprintf("%s\\%s", BaseDir, CurrDir)
                    } else {
                        TempDir = BaseDir
                    }

                    repl_Dir := NewCaseInsensitiveReplacer("&Dir", TempDir)
                    o32VarRec = repl_Dir.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Fil") {

                    repl_Fil := NewCaseInsensitiveReplacer("&Fil", CurrFil)
                    o32VarRec = repl_Fil.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Tim") {

                    showTime("&Tim")
                    repl_Tim := NewCaseInsensitiveReplacer("&Tim", FullDateTime)
                    o32VarRec = repl_Tim.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Inp") {

                    repl_Inp := NewCaseInsensitiveReplacer("&Inp", Inprec)
                    o32VarRec = repl_Inp.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Hst") {

                    repl_Hst := NewCaseInsensitiveReplacer("&Hst", cName)
                    o32VarRec = repl_Hst.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Acq") {

                    if len(ACQDir) > 0 {
                        TempAcq = fmt.Sprintf("%s\\%s", BACQDir, ACQDir)
                    } else {
                        TempAcq = BACQDir
                    }

                    repl_Acq := NewCaseInsensitiveReplacer("&Acq", TempAcq)
                    o32VarRec = repl_Acq.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Acn") {

                    repl_Acn := NewCaseInsensitiveReplacer("&Acn", ACQName)
                    o32VarRec = repl_Acn.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Win") {

                    repl_Win := NewCaseInsensitiveReplacer("&Win", WinRoot)
                    o32VarRec = repl_Win.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Fo") {

                    // Split string, we will likely need it split 
                    runeDelims := []rune(Delims)
                    tokRdr := csv.NewReader(strings.NewReader(Filrec))
                    tokRdr.Comma = runeDelims[0]
                    tokRdr.FieldsPerRecord = -1
                    tokFields, tok_err := tokRdr.Read()

                    if tok_err != nil {
                        ConsOut = fmt.Sprintf("[!] Parsing Error for Record(%d): %s\n", LoopNum, tok_err)
                        ConsLogSys(ConsOut, 1, 2)
                        continue
                    }                    
	
                    tokCount = len(tokFields)
                    if tokCount < 25 {
                        for i := tokCount; i < 26; i++ {
                            tokFields = append(tokFields, "")
                        }
                    }

                    if CaseInsensitiveContains(o32VarRec, "&For") {

                        repl_For := NewCaseInsensitiveReplacer("&For", Filrec)
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo0") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo0", tokFields[0])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo1") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo1", tokFields[1])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo2") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo2", tokFields[2])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo3") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo3", tokFields[3])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo4") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo4", tokFields[4])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo5") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo5", tokFields[5])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo6") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo6", tokFields[6])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo7") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo7", tokFields[7])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo8") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo8", tokFields[8])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fo9") {

                        repl_For := NewCaseInsensitiveReplacer("&Fo9", tokFields[9])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foa") {

                        repl_For := NewCaseInsensitiveReplacer("&Foa", tokFields[10])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fob") {

                        repl_For := NewCaseInsensitiveReplacer("&Fob", tokFields[11])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foc") {

                        repl_For := NewCaseInsensitiveReplacer("&Foc", tokFields[12])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fod") {

                        repl_For := NewCaseInsensitiveReplacer("&Fod", tokFields[13])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foe") {

                        repl_For := NewCaseInsensitiveReplacer("&Foe", tokFields[14])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fof") {

                        repl_For := NewCaseInsensitiveReplacer("&Fof", tokFields[15])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fog") {

                        repl_For := NewCaseInsensitiveReplacer("&Fog", tokFields[16])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foh") {

                        repl_For := NewCaseInsensitiveReplacer("&Foh", tokFields[17])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foi") {

                        repl_For := NewCaseInsensitiveReplacer("&Foi", tokFields[18])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foj") {

                        repl_For := NewCaseInsensitiveReplacer("&Foj", tokFields[19])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fok") {

                        repl_For := NewCaseInsensitiveReplacer("&Fok", tokFields[20])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fol") {

                        repl_For := NewCaseInsensitiveReplacer("&Fol", tokFields[21])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fom") {

                        repl_For := NewCaseInsensitiveReplacer("&Fom", tokFields[22])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fon") {

                        repl_For := NewCaseInsensitiveReplacer("&Fon", tokFields[23])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Foo") {

                        repl_For := NewCaseInsensitiveReplacer("&Foo", tokFields[24])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                    if CaseInsensitiveContains(o32VarRec, "&Fop") {

                        repl_For := NewCaseInsensitiveReplacer("&Fop", tokFields[25])
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

                }


































                // Testing...
                fmt.Printf("Out: %s\nOut: %s\n", o32VarRec, o64VarRec)


            }





















        }

        // Testing - Echo Input
        //fmt.Printf(Tmprec)


        if consOrFile == 1 {
            fmt.Printf(">>>")
        }
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
    SendLogMSG = strings.TrimSpace(SendLogMSG)


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

    HtmHndl, htm_err = os.Create(HtmFile)
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

        fmt.Fprintf(HtmHndl, "<html><head><title>AChoirX Artifacts</title></head>\n")
        fmt.Fprintf(HtmHndl, "<body>\n")
        fmt.Fprintf(HtmHndl, "<h2>Welcome to AChoirX %s</h2>\n\n", Version)
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

    // Remove CRLF to LF
    Conrec = strings.TrimSpace(Conrec)

    if conLog == 1 {
        if conHide == 1 {
            ConsOut = fmt.Sprintf("[?] *Redacted*\n")
        } else {
            ConsOut = fmt.Sprintf("[?] %s\n", Conrec);
        }

        ConsLogSys(ConsOut, 1, 1)

    }
}


//****************************************************************
// convert a record with Environment Variables in it             *
//  - Do manual checks for 64 bit exceptions - Check both 32&64  *
//****************************************************************
// This function changes in AChoirX - AChoir uses %EnVar%        *
// but native GOLang supports $Var or ${Var}.  AChoirX now uses  *
// the Native GoLang functions to prevent reinventing the wheel  *
//****************************************************************
func varConvert(inVarRec string) {
    o32VarRec = os.ExpandEnv(inVarRec)

    //  This doesn't apply to Linux or OSX
    if iopSystem == 0 {
        o64Replac := strings.NewReplacer("\\Program Files\\", "\\Program Files (x86)\\", "System32", "sysnative")
        o64VarRec = o64Replac.Replace(o32VarRec)
    } else {
        o64VarRec = o32VarRec
    }
}


//****************************************************************
// Case Insensitive Search for a SubString in a String           *
//****************************************************************
func CaseInsensitiveContains(CI_String, CI_Substr string) bool {
    CI_String, CI_Substr = strings.ToUpper(CI_String), strings.ToUpper(CI_Substr)
    return strings.Contains(CI_String, CI_Substr)
}


//****************************************************************
// Case Insensitive Replacer Routines                            *
//****************************************************************
type CaseInsensitiveReplacer struct {
    toReplace   *regexp.Regexp
    replaceWith string
}

func NewCaseInsensitiveReplacer(toReplace, with string) *CaseInsensitiveReplacer {
    return &CaseInsensitiveReplacer{
        toReplace:   regexp.MustCompile("(?i)" + toReplace),
        replaceWith: with,
    }
}

func (cir *CaseInsensitiveReplacer) Replace(str string) string {
    return cir.toReplace.ReplaceAllString(str, cir.replaceWith)
}


//****************************************************************
// Check for Windows Admin Privs by opening Physical Drive0      *
//****************************************************************
func IsUserAdmin() bool {
    _, err := os.Open("\\\\.\\PHYSICALDRIVE0")
    if err != nil {
        return false
    }
    return true
}
