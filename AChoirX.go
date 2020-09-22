// ****************************************************************
// "As we benefit from the inventions of others, we should be
//  glad to share our own...freely and gladly."
//  - Benjamin Franklin
//
// AChoirX = AChoir Version 10
//  This is Achoir converted to GoLang. VersionX was chosen to
//  to prevent clash with the C version numbers.  But essentially
//  AChoirX is Achoir converted to GoLang.
//
// AChoirX v10.00.01 - Convert from C to Go for Xplatform capability
//
//
//
//
// Other Libraries and code I use:
//  Syslog: go get github.com/NextronSystems/simplesyslog
//  Sys:    go get golang.org/x/sys
//  w32:    go get github.com/gonutz/w32
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
    "strconv"
    "text/scanner"
    "encoding/csv"
    "regexp"
    "runtime"
    "net/http"
    "path/filepath"
    "io"
    "bufio"
    "crypto/tls"
    syslog "github.com/NextronSystems/simplesyslog"
)


// Global Variable Settings
var Version = "v10.00.01"                       // AChoir Version
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
var slashDelimS string = "\\"                   // Same Thing but a String
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
var iFor = 0                                    // Loop Counter FOR, FO0 - FOP
var iLst = 0                                    // Loop Counter LST, LS0 - LSP
var ifFor = 0                                   // Flag contains FOR, FO0 - FOP
var ifLst = 0                                   // Flag contains LST, LS0 - LSP
var iMaxCnt = 0                                 // Maximum Record Count (Set by FOR:, LST: &FOR, &LST)
var LastRC = 0                                  // Last Return Code From External Executable
var iVrx = 0                                    // Index of VR0 - VR9 Variable array
var iCnx = 0                                    // Index of CN0 - CN9 Variable array
var JmpLbl = "LBL:Top"                          // Working Jump Label Build String
var iSleep = 0                                  // Seconds to Sleep
var volType = ""                                // Volume File System
var isNTFS = 0                                  // Is the Volume NTFS
var iCPS = 0                                    // Copy based on Magic Number (Signature)
var setCPath = 1                                // Output Copy Patch Shard - 0=None, 1=Partial, 2=Full

//Tokenize Records
var tokRec scanner.Scanner                      // Used to Tokenize Records into Slices
var tokCount = 0                                // Token Counter
var Delims = ",\\/"                             // Tokenizing Delimiters (Lst, For(Win), For(Lin)
var CntsTring = ""                              // Convert Cnt Integer Array Variable to String
var MCpFName = "File.Name"                      // Parseing Multiple Files with Glob (Wildcards)
var MCprcO = "FilePath"                         // Build Output File Name from Glob (Wildcards)
var MCpRoot = "FilePath"                        // Multi-Copy Root (Before any WildCards)
var MCpPath = "FilePath"                        // Path minus File Name
var MCpShard = "FilePath"                       // Multi-Copy Expanded Directories Shard (Before FileName)
var iShard =0                                   // Shard Index Pointer
var iAShard =0                                  // Asterisk Shard Index Pointer
var iQShard =0                                  // Question Mark Shard Index Pointer
var iOldLen = 0                                 // Old length of a string
var iNewLen = 0                                 // New length of a string

// Input and Output Records
var Inrec = "File Input Record"                 // File Input Record
var Conrec = "Console Record"                   // Console Output Record
var Tmprec = "Formatted Console Record"         // Console Formatting Record
var Cpyrec = "Copy Record"                      // Used by Copy Routine
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

// Arrays
var ForsArray = []string{"&FO0", "&FO1", "&FO2", "&FO3", "&FO4", "&FO5", "&FO6", "&FO7", "&FO8", 
                         "&FO9", "&FOA", "&FOB", "&FOC", "&FOD", "&FOE", "&FOF", "&FOG", "&FOH", 
                         "&FOI", "&FOJ", "&FOK", "&FOL", "&FOM", "&FON", "&FOO", "&FOP"}
var LstsArray = []string{"&LS0", "&LS1", "&LS2", "&LS3", "&LS4", "&LS5", "&LS6", "&LS7", "&LS8", 
                         "&LS9", "&LSA", "&LSB", "&LSC", "&LSD", "&LSE", "&LSF", "&LSG", "&LSH", 
                         "&LSI", "&LSJ", "&LSK", "&LSL", "&LSM", "&LSN", "&LSO", "&LSP"}
var VarsArray = []string{"&VR0", "&VR1", "&VR2", "&VR3", "&VR4", "&VR5", "&VR6", "&VR7", "&VR8", "&VR9"}
var VarxArray = []string{"VR0:", "VR1:", "VR2:", "VR3:", "VR4:", "VR5:", "VR6:", "VR7:", "VR8:", "VR9:"}
var VardArray = []string{"", "", "", "", "", "", "", "", "", ""}
var CntsArray = []string{"&CN0", "&CN1", "&CN2", "&CN3", "&CN4", "&CN5", "&CN6", "&CN7", "&CN8", "&CN9"}
var CntxArray = []string{"CN0:", "CN1:", "CN2:", "CN3:", "CN4:", "CN5:", "CN6:", "CN7:", "CN8:", "CN9:"}
var CntiArray = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}


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
    slashDelimS = fmt.Sprintf("%c", slashDelim)
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
        BaseDir = fmt.Sprintf("%cAChoir%cACQ-IR-%s-%04d%02d%02d-%02d%02d", slashDelim, slashDelim, cName, iYYYY, iMonth, iDay, iHour, iMin) 

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
        if _, exist_err := os.Stat(BACQDir); os.IsNotExist(exist_err) {
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
            ifFor = 0
            if strings.Contains(strings.ToUpper(Tmprec), "&FOR") {
                ifFor = 1
            }

            // Loop Through to check for &FO0 - &FOP
            for iFor = 0; iFor < 26; iFor++ {
                if strings.Contains(strings.ToUpper(Tmprec), ForsArray[iFor]) {
                    ifFor = 1
                }
            }

            if ifFor == 1 {
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
            ifLst = 0
            if strings.Contains(strings.ToUpper(Tmprec), "&LST") {
                ifLst = 1
            }

            // Loop Through to check for &LS0 - &LSP
            for iLst = 0; iLst < 26; iLst++ {
                if strings.Contains(strings.ToUpper(Tmprec), LstsArray[iLst]) {
                    ifLst = 1
                }
            }

            if ifLst == 1 {
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
                        } else if len(Filrec[ForSlash+1:]) < 2 {
                            ForFName = "Unknown"
                        } else {
                            ForFName = Filrec[ForSlash+1:]
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
                        TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir)
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
                        TempAcq = fmt.Sprintf("%s%c%s", BACQDir, slashDelim, ACQDir)
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

                    // Windows Delim Stored in Delims[1] - Linux Delim Stored in Delims[2]
                    // This is setable with SET:Delims= (Default is: ,\/)
                    if iopSystem == 0 {
                        tokRdr.Comma = runeDelims[1]
                    } else {
                        tokRdr.Comma = runeDelims[2]
                    }

                    tokRdr.FieldsPerRecord = -1
                    tokRdr.TrimLeadingSpace = true
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


                    // Look for Replacements &Fo0 - FoP
                    for iFor = 0; iFor < 26; iFor++ {
                        if CaseInsensitiveContains(o32VarRec, ForsArray[iFor]) {
                            repl_For := NewCaseInsensitiveReplacer(ForsArray[iFor], tokFields[iFor])
                            o32VarRec = repl_For.Replace(o32VarRec)
                        }
                    }
                }


                if CaseInsensitiveContains(o32VarRec, "&Ls") {

                    // Split string, we will likely need it split 
                    runeDelims := []rune(Delims)
                    tokRdr := csv.NewReader(strings.NewReader(Lstrec))
                    tokRdr.Comma = runeDelims[0]
                    tokRdr.FieldsPerRecord = -1
                    tokRdr.TrimLeadingSpace = true
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

                    if CaseInsensitiveContains(o32VarRec, "&Lst") {

                        repl_Lst := NewCaseInsensitiveReplacer("&Lst", Lstrec)
                        o32VarRec = repl_Lst.Replace(o32VarRec)
                    }

                    // Look for Replacements &Ls0 - LsP
                    for iLst = 0; iLst < 26; iLst++ {
                        if CaseInsensitiveContains(o32VarRec, LstsArray[iLst]) {
                            repl_Lst := NewCaseInsensitiveReplacer(LstsArray[iLst], tokFields[iLst])
                            o32VarRec = repl_Lst.Replace(o32VarRec)
                        }
                    }
                }

                if CaseInsensitiveContains(o32VarRec, "&Dsk") {

                    repl_Dsk := NewCaseInsensitiveReplacer("&Dsk", Dskrec)
                    o32VarRec = repl_Dsk.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Num") {

                    repl_Num := NewCaseInsensitiveReplacer("&Num", strconv.Itoa(LoopNum))
                    o32VarRec = repl_Num.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Cnr") {

                    repl_Cnr := NewCaseInsensitiveReplacer("&Cnr", strconv.Itoa(iMaxCnt))
                    o32VarRec = repl_Cnr.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Fnm") {

                    repl_Fnm := NewCaseInsensitiveReplacer("&Fnm", ForFName)
                    o32VarRec = repl_Fnm.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Rcd") {

                    repl_Rcd := NewCaseInsensitiveReplacer("&Rcd", strconv.Itoa(LastRC))
                    o32VarRec = repl_Rcd.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Chk") {

                    repl_Chk := NewCaseInsensitiveReplacer("&Chk", ChkFile)
                    o32VarRec = repl_Chk.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Drv") {

                    repl_Drv := NewCaseInsensitiveReplacer("&Drv", DiskDrive)
                    o32VarRec = repl_Drv.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Prc") {

                    repl_Prc := NewCaseInsensitiveReplacer("&Prc", Procesr)
                    o32VarRec = repl_Prc.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Vck") {
                    
                    repl_Vck := NewCaseInsensitiveReplacer("&Vck", volType)
                    o32VarRec = repl_Vck.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Dsa") {

                    TotBytes, FreeBytes := winFreeDisk()
                    if TotBytes == 0 {
                        ConsOut = fmt.Sprintf("[!] Error retrieving Disk Space stats, or not yet implemented (%s).\n", opSystem)
                        ConsLogSys(ConsOut, 1, 2)
                    }

                    // Even if we got 0 FreeBytes, replace it.
                    repl_Dsa := NewCaseInsensitiveReplacer("&Dsa", strconv.FormatUint(FreeBytes, 10))
                    o32VarRec = repl_Dsa.Replace(o32VarRec)
                }

                // Look for Replacements &VR0 - VR9
                for iVrx = 0; iVrx < 10; iVrx++ {
                    if CaseInsensitiveContains(o32VarRec, VarsArray[iVrx]) {
                        repl_Vrx := NewCaseInsensitiveReplacer(VarsArray[iVrx], VardArray[iVrx])
                        o32VarRec = repl_Vrx.Replace(o32VarRec)
                    }
                }

                // Look for Replacements &CN0 - CN9
                for iCnx = 0; iCnx < 10; iCnx++ {
                    if CaseInsensitiveContains(o32VarRec, CntsArray[iCnx]) {
                        CntsTring = strconv.Itoa(CntiArray[iCnx])
                        repl_Cnx := NewCaseInsensitiveReplacer(CntsArray[iCnx], CntsTring)
                        o32VarRec = repl_Cnx.Replace(o32VarRec)
                    }
                }


                //****************************************************************
                //* Now execute the Actions                                      *
                //****************************************************************
                Inrec = o32VarRec

                if len(Inrec) < 1 {
                    continue
                } else if strings.HasPrefix(Inrec, "*") {
                    continue
                } else if len(Inrec) < 4 {
                    continue
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "LBL:") {
                    continue
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "JMP:") && len(Inrec) > 4 {
                    if consOrFile == 1 {
                        ConsOut = fmt.Sprintf("[*] Jumping Does not make sense in Interactive Mode.  Ignoring...\n")
                        ConsLogSys(ConsOut, 1, 2)
                    } else {
                        // Rewind File and Jump to a Label (LBL:)
                        // Note: For This to work we have to Reset both the Handle and the Scanner!
                        RunMe = 0
                        IniHndl.Seek(0, 0)
                        IniScan = bufio.NewScanner(IniHndl)

                        JmpLbl = fmt.Sprintf("LBL:%s", Inrec[4:])

                        for IniScan.Scan() {
                            Tmprec = strings.TrimSpace(IniScan.Text())
                            if Tmprec == JmpLbl {
                                break
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CSE:") {
                    if strings.HasPrefix(strings.ToUpper(Inrec), "CSE:GET") {
                        getCaseInfo(1)
                    } else {
                        getCaseInfo(0);
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ACQ:") {
                    // Sanity Check - Replace Delimiters base on OS
                    if iopSystem == 0 {
                        Inrec = strings.Replace(Inrec, "/", "\\", -1) 
                    } else {
                        Inrec = strings.Replace(Inrec, "\\", "/", -1) 
                    }


                    // Have we created the Base Acquisition Directory Yet?
                    if _, BACQ_err := os.Stat(BACQDir); os.IsNotExist(BACQ_err) {
                        // Set iRunMode=1 to be sure we post-process the Acquired Artifacts
                        // (In case we had not set it originally due to remote BACQDIR)
                        iRunMode = 1;

                        DirAllocErr(BACQDir);
                        DirAllocErr(CachDir);
                        PreIndex();
                    }

                    // Explicit Path (Dependent upon OS!
                    osACQ := fmt.Sprintf("ACQ:%c", slashDelim)
                    if strings.HasPrefix(strings.ToUpper(Inrec), osACQ) {
                        if len(Inrec) > 5 {
                            ACQDir = fmt.Sprintf("%s", Inrec[5:])
                            TempDir = fmt.Sprintf("%s%c%s", BACQDir, slashDelim, ACQDir)
                        } else  {
                            ACQDir = ""
                            TempDir = BACQDir
                        }
                    } else {
                        if len(Inrec) > 4 {
                            //Check to see if it is an append or new &Acq
                            //Dont add // if it's new!
                            if len(ACQDir) > 0 {
                                ACQDir += slashDelimS
                            }

                            ACQDir += Inrec[4:]
                            TempDir = fmt.Sprintf("%s%c%s", BACQDir, slashDelim, ACQDir);
                        }
                    }

                    // Determine the Level 1 Directory to see if we have it yet
                    // If not, we will want to add it to the HTML file
                    LvlSplit := strings.Split(ACQDir, slashDelimS)
                    LevelOne := fmt.Sprintf("%s%c%s", BACQDir, slashDelim, LvlSplit[0])

                    if _, level_err := os.Stat(LevelOne); os.IsNotExist(level_err) && iHtmMode == 1 {
                        fmt.Fprintf(HtmHndl, "</td><td align=center>\n");
                        fmt.Fprintf(HtmHndl, "<a href=file:%s target=AFrame> %s </a>\n", LvlSplit[0], LvlSplit[0]);
                    }


                    // Have we created this Directory already?
                    if _, ACQDir_err := os.Stat(TempDir); os.IsNotExist(ACQDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Acquisition Sub-Directory: %s\n", ACQDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir);

                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "DIR:") {
                    //****************************************************************
                    //* Set Current Directory                                        *
                    //****************************************************************
                    // Sanity Check - Replace Delimiters base on OS
                    if iopSystem == 0 {
                        Inrec = strings.Replace(Inrec, "/", "\\", -1) 
                    } else {
                        Inrec = strings.Replace(Inrec, "\\", "/", -1) 
                    }

                    // Explicit Path (Dependent upon OS!
                    osDir := fmt.Sprintf("DIR:%c", slashDelim)
                    if strings.HasPrefix(strings.ToUpper(Inrec), osDir) {
                        if len(Inrec) > 5 {
                            CurrDir = fmt.Sprintf("%s", Inrec[5:])
                            TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir)
                        } else  {
                            CurrDir = ""
                            TempDir = BaseDir
                        }
                    } else {
                        if len(Inrec) > 4 {
                            //Check to see if it is an append or new &Dir
                            //Dont add // if it's new!
                            if len(CurrDir) > 0 {
                                CurrDir += slashDelimS
                            }

                            CurrDir += Inrec[4:]
                            TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir);
                        }
                    }

                    // Have we created this Directory already?
                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Sub-Directory: %s\n", CurrDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir);

                    }

                    // Reset The WorkingDirectory to the new Directory
                    CurrWorkDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir);
                    os.Chdir(CurrWorkDir); 

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FIL:") {
                    CurrFil = Inrec[4:]
                    TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir);

                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Sub-Directory: %s\n", CurrDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir);
                    }

                    ConsOut = fmt.Sprintf("[+] File has been Set to: %s\n", CurrFil)
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "VR") {
                    // Look for Replacements VR0: - VR9:
                    for iVrx = 0; iVrx < 10; iVrx++ {
                        if strings.HasPrefix(strings.ToUpper(Inrec), VarxArray[iVrx]) {
                            VardArray[iVrx] = Inrec[4:]
                        }
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CN") {
                    // Look for Replacements CN0: - CN9:
                    for iCnx = 0; iCnx < 10; iCnx++ {
                        if strings.HasPrefix(strings.ToUpper(Inrec), CntxArray[iCnx]) {
                            if Inrec[4:] == "++" {
                                CntiArray[iCnx]++
                            } else if Inrec[4:] == "--" {
                                CntiArray[iCnx]--
                            } else {
                                CntiArray[iCnx], _ = strconv.Atoi(Inrec[4:])
                            }
                        }
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "DRV:") {
                    DiskDrive = Inrec[4:]

                    ConsOut = fmt.Sprintf("[+] Disk Drive Set to: %s\n", DiskDrive)
                    ConsLogSys(ConsOut, 1, 1)

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "INI:") {
                    IniFile = Inrec[4:]

                    if strings.ToUpper(IniFile) == "CONSOLE" {
                        if consOrFile == 0 {
                            RunMode = "Con"
                            inFnam = "Console"
                            iRunMode = 1
                            consOrFile = 1

                            ConsOut = fmt.Sprintf("[+] Switching to Console (Interactive Mode\n")
                            ConsLogSys(ConsOut, 1, 1)

                            IniHndl.Close()
                            fmt.Printf(">>>");
                            IniScan = bufio.NewScanner(os.Stdin)
                        }
                    } else {
                        // _, exist_err := os.Stat(IniFile)
                        // if os.IsNotExist(exist_err) {

                        if FileExists(IniFile) {
                            ConsOut = fmt.Sprintf("[!] Switching to INI File: %s\n", Inrec[4:])
                            ConsLogSys(ConsOut, 1, 1)

                            // Only close the handle if its not Console. If it is Console Set it back to File
                            if consOrFile == 0 {
                                IniHndl.Close()
                            } else {
                                consOrFile = 0
                            }

                            IniHndl, ini_err = os.Open(IniFile)
                            if ini_err != nil {
                                ConsOut = fmt.Sprintf("[!] Error Opening Ini File: %s - Exiting.\n", IniFile)
                                ConsLogSys(ConsOut, 1, 2)

                                //cleanUp_Exit(3)
                                os.Exit(3)
                            }

                            IniScan = bufio.NewScanner(IniHndl)
                            RunMe = 0  // Conditional run Script default is yes

                        } else {
                            ConsOut = fmt.Sprintf("[!] Requested INI File Not Found: %s - Ignored.\n", Inrec[4:])
                            ConsLogSys(ConsOut, 1, 2)
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ADM:CHECK") {
                    if iIsAdmin == 1 {
                        ConsOut = fmt.Sprintf("[+] Running as Admin\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ConsOut = fmt.Sprintf("[+] Running as NON-Admin\n")
                        ConsLogSys(ConsOut, 1, 1)
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ADM:FORCE") {
                    if iIsAdmin == 1 {
                        ConsOut = fmt.Sprintf("[+] Running as Admin - Continuing...\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ConsOut = fmt.Sprintf("[+] NOT Running as Admin - Exiting...\n")
                        ConsLogSys(ConsOut, 1, 1)

                        //cleanUp_Exit(3)
                        os.Exit(3)
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:HIDE") {
                    ConsOut = fmt.Sprintf("[+] Hiding Console...\n")
                    ConsLogSys(ConsOut, 1, 1)
                    winConHideShow(0)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:SHOW") {
                    ConsOut = fmt.Sprintf("[+] Showing Console...\n")
                    ConsLogSys(ConsOut, 1, 1)
                    winConHideShow(1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:MSGLEVEL=MIN") {
                    setMSGLvl = 1
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:MSGLEVEL=STD") {
                    setMSGLvl = 2
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:MSGLEVEL=MAX") {
                    setMSGLvl = 3
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:MSGLEVEL=DEBUG") {
                    setMSGLvl = 4
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SLP:") {
                    iSleep, _ = strconv.Atoi(Inrec[4:])
                    time.Sleep (time.Duration(iSleep) * time.Second)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "INP:") {
                    consInput(Inrec[4:], 1, 0);
                    Inprec = Conrec
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "VCK:") {
                    isNTFS = 0
                    volType = winGetVolInfo(Inrec[4:])

                    // This should only work in Windows - Linux and OSX will be UNKNOWN
                    if volType == "NTFS" {
                        isNTFS = 1
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CPY:") || strings.HasPrefix(strings.ToUpper(Inrec), "CPS:") {
                    //****************************************************************
                    //* Binary Copy From => To                                       *
                    //****************************************************************
                    if strings.HasPrefix(strings.ToUpper(Inrec), "CPS:") {
                        iCPS = 1
                    } else {
                        iCPS = 0
                    }

                    ConsOut = fmt.Sprintf("[+] %s\n", Inrec)
                    ConsLogSys(ConsOut, 1, 1)

                    Cpyrec = Inrec[4:]

                    splitString1, splitString2, SplitRC := twoSplit(Cpyrec)

                    // Remove any duplicate Delimiters - This is necessary to prevent indexing errors when
                    //  The found file does not match the search string (OS ignore duplicated delimiters)
                    oneDelim := fmt.Sprintf("%c", slashDelim)
                    twoDelim := fmt.Sprintf("%c%c", slashDelim, slashDelim)

                    iOldLen = 1
                    iNewLen = 0
                    for iOldLen != iNewLen {
                        iOldLen = len(splitString1)
                        splitString1 = strings.Replace(splitString1, twoDelim, oneDelim, -1)
                        iNewLen = len(splitString1)
                    }

                    iOldLen = 1
                    iNewLen = 0
                    for iOldLen != iNewLen {
                        iOldLen = len(splitString2)
                        splitString2 = strings.Replace(splitString2, twoDelim, oneDelim, -1)
                        iNewLen = len(splitString2)
                    }

                    if SplitRC == 1 {
                        ConsOut = fmt.Sprintf("[!] Copying Requires both a FROM File and a TO Directory\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        fmt.Sprintf("CPY: %s to %s\n", splitString1, splitString2)

                        //****************************************************************
                        //* If we see any wildcards, do search for multiple occurances   *
                        //****************************************************************
                        if strings.Contains(splitString1, "*") || strings.Contains(splitString1, "?") {

                            //****************************************************************
                            // Parse out the Expanded Directory Shard (pre-wildcard)         *
                            //****************************************************************
                            iAShard = strings.IndexByte(splitString1, '*')
                            iQShard = strings.IndexByte(splitString1, '?')

                            if iAShard < 0 {
                                iShard = strings.LastIndexByte(splitString1[:iQShard], slashDelim) + 1
                            } else if iQShard < 0 {
                                iShard = strings.LastIndexByte(splitString1[:iAShard], slashDelim) + 1
                            } else if iAShard < iQShard {
                                iShard = strings.LastIndexByte(splitString1[:iAShard], slashDelim) + 1
                            } else if iQShard < iAShard {
                                iShard = strings.LastIndexByte(splitString1[:iQShard], slashDelim) + 1
                            } else {
                                iShard = 0
                            }

                            if iShard > 1 {
                                MCpRoot = splitString1[:iShard]
                            } else {
                                MCpRoot = ""
                            }


                            //****************************************************************
                            // Do File Search using Glob                                     *
                            //****************************************************************
                            files_glob, err_glob := filepath.Glob(splitString1)

                            if err_glob != nil {
                                ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", err_glob)
                                ConsLogSys(ConsOut, 1, 1)
                                continue
                            }

                            for _, file_found := range files_glob {
                                //****************************************************************
                                //* Ignore Directories - Only Process Files                      *
                                //****************************************************************
                                file_stat, _ := os.Stat(file_found)
                                if file_stat.IsDir() {
                                    continue
                                }


                                //****************************************************************
                                //* Get Just the File Name                                       *
                                //****************************************************************
                                ForSlash = strings.LastIndexByte(file_found, slashDelim)

                                if (ForSlash == -1) {
                                    MCpFName = file_found
                                    MCpShard = ""
                                    MCpPath = ""
                                } else if iShard == 0 {
                                    MCpFName = file_found[ForSlash+1:]
                                    MCpShard = ""
                                    MCpPath = file_found[:ForSlash] 
                                } else if len(file_found[ForSlash+1:]) < 2 {
                                    MCpFName = file_found
                                    MCpShard = file_found[iShard:len(file_found)]
                                    MCpPath = file_found[:ForSlash] 
                                } else {
                                    MCpFName = file_found[ForSlash+1:]
                                    MCpShard = file_found[iShard:len(file_found)-len(MCpFName)]
                                    MCpPath = file_found[:ForSlash] 
                                }

                                //****************************************************************
                                //* Copy to Output File Name                                     *
                                //*  Note: a Shard is any expanded WildCard Directory - Shards   *
                                //*        can be used to logically group duplicate file names   *
                                //****************************************************************
                                if setCPath == 0 {
                                    MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)
                                } else if setCPath == 1 && len(MCpShard) < 1 {
                                    MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)
                                } else if setCPath == 1 {
                                    MCprcO = fmt.Sprintf("%s%c%s%s", splitString2, slashDelim, MCpShard, MCpFName)
                                } else if setCPath == 2 {
                                    MCpPath = strings.Replace(MCpPath, ":", "", -1)
                                    MCprcO = fmt.Sprintf("%s%c%s%s", splitString2, slashDelim, MCpPath, MCpFName)
                                }
                                ConsOut = fmt.Sprintf("[+] Multi-Copy File: %s\n    To: %s\n", file_found, MCprcO)
                                ConsLogSys(ConsOut, 1, 1)

                                nBytes, copy_err := binCopy(file_found, MCprcO)

                                if copy_err != nil {
                                    ConsOut = fmt.Sprintf("[!] Error Copying File: %s\n", copy_err)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    ConsOut = fmt.Sprintf("[+] Copy Complete: %d Bytes Copied\n", nBytes)
                                    ConsLogSys(ConsOut, 1, 1)
                                }
                            }

                            if (iNative == 0) {
                                //****************************************************************
                                //* Replace System32 with Sysnative if we are non-native         *
                                //****************************************************************
                                if CaseInsensitiveContains(splitString1, "\\System32\\") || CaseInsensitiveContains(splitString1, "/System32/") {
                                    TempDir = splitString1

                                    repl_sys := NewCaseInsensitiveReplacer("System32", "sysnative")
                                    TempDir = repl_sys.Replace(TempDir)

                                    ConsOut = fmt.Sprintf("[*] Non-Native Flag Has Been Detected - Adding Sysnative Redirection: \n %s\n", TempDir)
                                    ConsLogSys(ConsOut, 1, 1)

                                    files_glob, err_glob := filepath.Glob(TempDir)

                                    if err_glob != nil {
                                        ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", err_glob)
                                        ConsLogSys(ConsOut, 1, 1)
                                        continue
                                    }

                                    for _, file_found := range files_glob {
                                        //****************************************************************
                                        //* Ignore Directories - Only Process Files                      *
                                        //****************************************************************
                                        file_stat, _ := os.Stat(file_found)
                                        if file_stat.IsDir() {
                                            continue
                                        }

                                        //****************************************************************
                                        //* Get Just the File Name                                       *
                                        //****************************************************************
                                        ForSlash = strings.LastIndexByte(file_found, slashDelim)
                                        if (ForSlash == -1) {
                                            MCpFName = file_found
                                        } else if len(file_found[ForSlash+1:]) < 2 {
                                            MCpFName = file_found
                                        } else {
                                            MCpFName = file_found[ForSlash+1:]
                                        }

                                        //****************************************************************
                                        //* Copy to Output File Name                                     *
                                        //****************************************************************
                                        if setCPath == 0 || len(MCpShard) < 1 {
                                            MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)
                                        } else { 
                                            MCprcO = fmt.Sprintf("%s%c%s%s", splitString2, slashDelim, MCpShard, MCpFName)
                                        }

                                        ConsOut = fmt.Sprintf("[+] Multi-Copy Redir File: %s\n    To: %s\n", file_found, MCprcO)
                                        ConsLogSys(ConsOut, 1, 1)

                                        nBytes, copy_err := binCopy(file_found, MCprcO)

                                        if copy_err != nil {
                                            ConsOut = fmt.Sprintf("[!] Error Copying File: %s\n", copy_err)
                                            ConsLogSys(ConsOut, 1, 1)
                                        } else {
                                            ConsOut = fmt.Sprintf("[+] Copy Complete: %d Bytes Copied\n", nBytes)
                                            ConsLogSys(ConsOut, 1, 1)
                                        }
                                    }
                                }
                            }

                        } else {
                            //****************************************************************
                            //* Get Just the File Name                                       *
                            //****************************************************************
                            ForSlash = strings.LastIndexByte(splitString1, slashDelim)
                            if (ForSlash == -1) {
                                MCpFName = splitString1
                            } else if len(splitString1[ForSlash+1:]) < 2 {
                                MCpFName = splitString1
                            } else {
                                MCpFName = splitString1[ForSlash+1:]
                            }

                            //****************************************************************
                            //* Copy to Output File Name                                     *
                            //****************************************************************
                            MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)

                            ConsOut = fmt.Sprintf("[+] Singl-Copy File: %s\n    To: %s\n", splitString1, MCprcO)
                            ConsLogSys(ConsOut, 1, 1)

                            nBytes, copy_err := binCopy(splitString1, MCprcO)

                            if copy_err != nil {
                                ConsOut = fmt.Sprintf("[!] Error Copying File: %s\n", copy_err)
                                ConsLogSys(ConsOut, 1, 1)
                            } else {
                                ConsOut = fmt.Sprintf("[+] Copy Complete: %d Bytes Copied\n", nBytes)
                                ConsLogSys(ConsOut, 1, 1)
                            }


                            if (iNative == 0) {
                                //****************************************************************
                                //* Replace System32 with Sysnative if we are non-native         *
                                //****************************************************************
                                if CaseInsensitiveContains(splitString1, "\\System32\\") || CaseInsensitiveContains(splitString1, "/System32/") {
                                    TempDir = splitString1
                                    ConsOut = fmt.Sprintf("[*] Non-Native Flag Has Been Detected - Adding Sysnative Redirection: \n %s\n", TempDir)
                                    ConsLogSys(ConsOut, 1, 1)

                                    //****************************************************************
                                    //* Get Just the File Name                                       *
                                    //****************************************************************
                                    ForSlash = strings.LastIndexByte(splitString1, slashDelim)
                                    if (ForSlash == -1) {
                                        MCpFName = splitString1
                                    } else if len(splitString1[ForSlash+1:]) < 2 {
                                        MCpFName = splitString1
                                    } else {
                                        MCpFName = splitString1[ForSlash+1:]
                                    }

                                    //****************************************************************
                                    //* Copy to Output File Name                                     *
                                    //****************************************************************
                                    MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)

                                    ConsOut = fmt.Sprintf("[+] Singl-Copy Redir File: %s\n    To: %s\n", splitString1, MCprcO)
                                    ConsLogSys(ConsOut, 1, 1)

                                    nBytes, copy_err := binCopy(splitString1, MCprcO)

                                    if copy_err != nil {
                                        ConsOut = fmt.Sprintf("[!] Error Copying File: %s\n", copy_err)
                                        ConsLogSys(ConsOut, 1, 1)
                                    } else {
                                        ConsOut = fmt.Sprintf("[+]Copy Complete: %d Bytes Copied\n", nBytes)
                                        ConsLogSys(ConsOut, 1, 1)
                                    }







                                 }
                            }
                        }
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





    // Some Test Code for checking Drive Types for Windows
    // DRIVE_CDROM = 5, DRIVE_FIXED = 3, DRIVE_RAMDISK = 6, DRIVE_REMOTE = 4, DRIVE_REMOVABLE = 2
    if iopSystem == 0 {
        winListDrives()
        winFreeDisk()
    } else {
        fmt.Printf("[!] Bypassing Drive and Memory Routines - We are not running on Windows\n")
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
    if _, exist_err := os.Stat(DirToCreat); os.IsNotExist(exist_err) {
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
    HtmFile = fmt.Sprintf("%s%cIndex.htm", BACQDir, slashDelim)

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


//***********************************************************************
//* Expand out Directories if they do not exist by slashDelim           *
//***********************************************************************
func ExpandDirs(FullDirName string) {
    var iDir = 0
    var TempDirName = ""

    DirSplit := strings.Split(FullDirName, slashDelimS)

    if len(DirSplit) > 1 {
        // Go through Dirs - Start with 1 (Ignore Root)
        TempDirName = DirSplit[0]
        for iDir = 1; iDir < len(DirSplit); iDir++ {
            TempDirName += slashDelimS
            TempDirName += DirSplit[iDir]

            DirAllocErr(TempDirName);
        }
    }
}


func twoSplit(SpString string) (string, string, int) {
    // Split string in two by space - but honor quotes 
    // Success Returns String1 String2 and ReturnCode=2
    // var splitRC = 0 

    tokRdr := csv.NewReader(strings.NewReader(SpString))
    tokRdr.Comma = ' '
    tokRdr.FieldsPerRecord = -1
    tokRdr.TrimLeadingSpace = true
    tokFields, tok_err := tokRdr.Read()


    if tok_err != nil {
        ConsOut = fmt.Sprintf("[!] Parsing Error for: %s\n", SpString)
        ConsLogSys(ConsOut, 1, 2)
        return SpString, "", 1
    }

    tokCount = len(tokFields)

    if tokCount < 2 {
        ConsOut = fmt.Sprintf("[!] No Separating Delimiters Found for: %s\n", SpString)
        ConsLogSys(ConsOut, 1, 2)
        return SpString, "", 1
    } else if tokCount > 2 {
        ConsOut = fmt.Sprintf("[!] Extra Data Ignored. Tokens: %d\n", tokCount)
        ConsLogSys(ConsOut, 1, 2)
        return tokFields[0], tokFields[1], 3
    } else {
        return tokFields[0], tokFields[1], 2
    } 
}


//***********************************************************************
//* Binary Copy Two Files - FromFile, TooFile                           *
//***********************************************************************
func binCopy(FrmFile, TooFile string) (int64, error) {
    var TmpFrmFile = ""
    var TmpTooFile = ""
    var LastSlash = 0
    var iFileCount = 0
    var PathOnly = "/"

    if !FileExists(FrmFile) {
        ConsOut = fmt.Sprintf("[*] Source Copy File Not Found: \n    %s\n", FrmFile)
        ConsLogSys(ConsOut, 1, 2)

        // Check for Sysnative edge case (running 32 bit on 64 bit)
        if (iNative == 0) {
            //****************************************************************
            // Wait... Maybe it's a file Redirect                            *
            //****************************************************************
            //* Replace System32 with Sysnative if we are non-native         *
            //****************************************************************
            if CaseInsensitiveContains(FrmFile, "\\System32\\") || CaseInsensitiveContains(FrmFile, "/System32/") {
                TmpFrmFile = FrmFile

                repl_frm := NewCaseInsensitiveReplacer("System32", "sysnative")
                TmpFrmFile = repl_frm.Replace(TmpFrmFile)

                ConsOut = fmt.Sprintf("[*] Non-Native Flag Has Been Detected - Trying Sysnative Redirection:\n    %s\n", TmpFrmFile)
                ConsLogSys(ConsOut, 1, 1)

                if !FileExists(TmpFrmFile) {
                    ConsOut = fmt.Sprintf("[*] Sysnative Source Copy Also File Not Found:\n    %s\n", TmpFrmFile)
                    ConsLogSys(ConsOut, 1, 1)

                    //No... Sorry... Not Sysnative
                    return 0, fmt.Errorf("[!] Copy Error - File Could Not Be Found: %s", FrmFile)

                } else {
                    ConsOut = fmt.Sprintf("[*] Sysnative Source Copy File Found. Now Substituting:\n    %s\n", TmpFrmFile)
                    ConsLogSys(ConsOut, 1, 1)

                    // Yes... Substitution Successful
                    FrmFile = TmpFrmFile
                }
            }
        }
    }


    //***********************************************************************
    //* Check to make sure its a Regular File                               *
    //***********************************************************************
    FrmFileStat, stat_err := os.Stat(FrmFile)
    if stat_err != nil {
        return 0, stat_err
    }

    if !FrmFileStat.Mode().IsRegular() {
        return 0, fmt.Errorf("[!] Copy Error: %s is not a Regular File", FrmFile)
    }


    //***********************************************************************
    //* Open it up with Defer Close                                         *
    //***********************************************************************
    FrmSource, frm_err := os.Open(FrmFile)
    if frm_err != nil {
        return 0, frm_err
    }
    defer FrmSource.Close()


    //***********************************************************************
    //* Check to make sure Dest Directory Exists                            *
    //***********************************************************************
    LastSlash = strings.LastIndexByte(TooFile, slashDelim)
    PathOnly = TooFile[:LastSlash]
    ExpandDirs(PathOnly)    


    //***********************************************************************
    //* Never OverWrite a File - Rename if it is already there              *
    //***********************************************************************
    if FileExists(TooFile) {
        TmpTooFile = TooFile
        iFileCount = 0
        for FileExists(TmpTooFile) {
            iFileCount++;
            TmpTooFile = fmt.Sprintf("%s(%d)", TooFile, iFileCount)
        }

        TooFile = TmpTooFile
        ConsOut = fmt.Sprintf("[*] Destination File Already Exists.\n    Renamed To: %s\n", TooFile)
        ConsLogSys(ConsOut, 1, 2)
    }



    //***********************************************************************
    //* Now Copy the File                                                   *
    //***********************************************************************
    TooDest, too_err := os.Create(TooFile)
    if too_err != nil {
        return 0, too_err
    }
    defer TooDest.Close()

    nBytes, copy_err := io.Copy(TooDest, FrmSource)
    return nBytes, copy_err
}


//***********************************************************************
// FileExists checks if a file exists and is not a directory before we  *
//  try using it to prevent further errors.                             *
//***********************************************************************
func FileExists(CheckFilename string) bool {
    // See if the File is already there
    fexst_info, fexst_err := os.Stat(CheckFilename)

    if os.IsNotExist(fexst_err) {
        return false
    }

    return !fexst_info.IsDir()

}

