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
// AChoirX v10.00.20 - Alpha 2 - Mostly Feature Complete
// AChoirX v10.00.21 - Add Native S3 File upload capability
// AChoirX v10.00.22 - Multi-Upload S3 File upload capability
// AChoirX v10.00.23 - Multi-Upload S3 Improvements
// AChoirX v10.00.24 - Hash Running Program for non-repudiation
//                     Add &Myp (My Program) and &Myh (My Hash)
// AChoirX v10.00.25 - Add /GXR: - Gets a Zip File, Extracts it, and
//                     runs the script.
// AChoirX v10.00.26 - Copy Running Program to \Progs for non-repudiation
//
// 11/06/2020 - Move to Beta Status:
// AChoirX v10.00.27 - Change CPY: Target File Atime and MTime to match Source
//                     Change FileExists to accept File or Directory & Improve Error handling
//                     Add Quoted Parsing to EXE: and SYS:
// AChoirX v10.00.28 - Add &USR and &PWD - To enable UserID and Password on Command Line
// AChoirX v10.00.29 - Add File Encryption/Decryption - Credit goes to:
//                      https://www.thepolyglotdeveloper.com/2018/02/
//                      encrypt-decrypt-data-golang-application-crypto-packages/
//                     Note: It will use &PWD to Encrypt
// AChoirX v10.00.30 - Fix minor message wording
// AChoirX v10.00.31 - Add Admin/Root Checks - Move subroutine to Platform Specific Files
//                   - Add &Adm Variable = Yes or No  (Running as Admin/Root)
// AChoirX v10.00.32 - Add &Mem (Total Memory) - Copied from: https://github.com/pbnjay/memory
//                      BSD 3-Clause License - Thanks to Jeremy Jay
// AChoirX v10.00.33 - Add support for unzipping an embedded []byte stream
//                   - Embed Platform specific Default Scripts (Win, Lin, OSX)
// AChoirX v10.00.34 - Add Embedder to ToolChain.  Include WinPmem (Memory Dumper) in Embedded Zip
//                     Add TSK Fcat into (Raw NTFS Copy) into Embedded Zip
// AChoirX v10.00.35 - Fix &CNR for FOR: and LST:
//
// AChoirX v10.00.36 - Change Conditional Logic to only count a single occurance of &FOR and &LST comparisons
//                     This prevents the need for multiple END: statements  - Multiple comparisons only get 
//                     a single hit. Add &For and &lst Support to more Actions.
//                   - Add HSH:<Filename> Will put the File hash in the &HSH Variable
//                     (Only supports a single File for now)
//                   - Trim quotes for CKN: and CKY:
//
// Other Libraries and code I use:
//  Syslog: go get github.com/NextronSystems/simplesyslog
//  Sys:    go get golang.org/x/sys
//  w32:    go get github.com/gonutz/w32
//  S3:     go get github.com/aws/aws-sdk-go/...
// Changes from AChoir:
//  Environment Variable Expansion now uses GoLang $Var or ${Var} 
//
// Not Implemented from AChoir (Yet):
//  Raw NTFS - Windows Unique - Use TSK
//  NTP - Not used enough
//  Console Colors - No native cross-platform way to do this
//  Native SMB/CIFS - Windows Unique
//  USB Protection (Registry Key) - Windows Unique
//  Transfer File MetaData on Copy (CTime, Perms, Owner, etc) - Not cross-platform 
// ****************************************************************

package main

import (
    "fmt"
    "time"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "text/scanner"
    "encoding/csv"
    "encoding/hex"
    "archive/zip"
    "regexp"
    "runtime"
    "net/http"
    "path/filepath"
    "io"
    "io/ioutil"
    "bufio"
    "crypto/tls"
    "crypto/md5"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"

    syslog "github.com/NextronSystems/simplesyslog")


// Global Variable Settings
var Version = "v10.00.35"                       // AChoir Version
var RunMode = "Run"                             // Character Runmode Flag (Build, Run, Menu)
var ConsOut = "[+] Console Output"              // Console, Log, Syslog strings
var MyProg = "none"                             // My Program Name and Path (os.Args[0])
var MyHash = "none"                             // My Hash
var iRunMode = 0                                // Int Runmode Flag (0, 1, 2)
var inFnam = "AChoir.ACQ"                       // Script Name
var inEncFile = "AChoir.ECR"                    // Input Encrypted File Name
var ACQName = "ACQ-IR-LocalHost-00000000-0000"  // AChoir Unique Collection Name
var DiskDrive = "C:"                            // Disk Drive (/DRV)
var iCase = 0                                   // Case Information Processing Mode
var consOrFile = 0                              // Console Input instead of File
var opArchit = "AMD64"                          // Architecture
var opSystem = "Windows"                        // Which Operating System are we running on
var iopSystem = 0                               // Operating System Flag (0=win, 1=lin, 2=osx, 3=?)
var OSVersion = "Windows 10.0.0"                // Operating System Version: OS Major.Minor.Build
var slashDelim byte = '\\'                      // Directory Delimiter Win vs. Lin vs. OSX
var slashDelimS string = "\\"                   // Same Thing but a String
var WGetURL = "http://Company.com/file"         // URL For HTTP Get
var WGetDir = "c:\\Achoir"                      // Directory For HTTP Get (Downloaded File)
var RootSlash = 0                               // Last Occurance of Slash to find Root URL
var ForSlash = 0                                // Last Occurance of Slash to find File in Path
var CurrFil = "Current.fil"                     // Current File Name
var inUser = "none"                             // UserId
var inPass = "none"                             // Password
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
var dskTyp uint32 = 3                           // Default Disk Type is Fixed (Type 3)
var LoopNum = 0                                 // Loop Counter
var ForFName = "File.txt"                       // Parsed File name from Path
var iNative = 0                                 // Are we Native 64Bit on 64Bit (Native = 1, NonNative = 0)
var sNative = ""                                // Native String (blank or Non-) 
var iIsAdmin = 0                                // Are we an Admin 
var sIsAdmin = ""                               // Are we an Admin String (blank or Non-) 
var yIsAdmin = "No"                             // Are we an Admin (Yes | No)
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
var iOutOfDiskSpace = 0                         // Did we get any Out Of Disk Space Errors
var iXitCmd = 0                                 // Are we runnning an Exit Command
var XitCmd = "Exit"                             // Exit Command (AChoirX Post Processing)
var iOPNisOpen = 0                              // Is the User Defined File Open?
var setCDepth = 10                              // Set Copy Depth
var LastHash = "none"                           // Last Single File Hash
var NotFound = 0                                // Flag for ensuring that only one Found Rec Increments RunMe
var YesFound = 0                                // Flag for ensuring that only one Found Rec Increments RunMe

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
var iShard = 0                                  // Shard Index Pointer
var iAShard = 0                                 // Asterisk Shard Index Pointer
var iQShard = 0                                 // Question Mark Shard Index Pointer
var iDblShard = 0                               // Double Glob Shard Pointer
var iOldLen = 0                                 // Old length of a string
var iNewLen = 0                                 // New length of a string
var WalkfileWild = "Wildcard"                   // Wildcard Portion of the WalkFile Routines
var WalkfileToo = "TooFiled"                    // File Name Portion of the WalkFile Routines

// Input and Output Records
var Inrec = "File Input Record"                 // File Input Record
var Conrec = "Console Record"                   // Console Output Record
var Tmprec = "Formatted Console Record"         // Console Formatting Record
var Cpyrec = "Copy Record"                      // Used by Copy Routine
var Encrec = "Encrypt Record"                   // Used by Encrypt Routine
var S3Urec = "Copy Record"                      // Used by S3 Upload Routine
var Cmprec = "Compare Record"                   // Used by Compare Routines
var Ziprec = "Zip Record"                       // Used by Unzip Routines
var Getrec = "Get Record"                       // Used by HTTP Get: Routines
var Filrec = "File Record"                      // File Record 
var Lstrec = "List Record"                      // List Record 
var Dskrec = "Disk Record"                      // Disk Record 
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
var Syslogp = "514"                             // Syslog Port string
var iSyslogp = 514                              // Syslog Port Int
var SyslogTMSG = "AChoir Syslog Started."       // Initialize Syslog Messages 
var SyslogServer = "127.0.0.1:514"              // Syslog Server:Port
var tlsConfig *tls.Config                       // TLS Config

// AWS S3 Variables
var S3_REGION = "none"                          // AWS Region
var S3_BUCKET = "none"                          // AWS Bucket
var S3_AWSId = "none"                           // AWS ID
var S3_AWSKey = "none"                          // AWS Secret Key
var S3_Session *session.Session                 // AWS Session
var S3_AWS_SplitRC = 0                          // AWS Split Return Code
var S3_err error                                // S3 Errors
var iS3Login = 0                                // Default is NOT logged in

// Message and Log Levels
var iLogOpen = 0                                // Is the LogFile Open Yet
var setMSGLvl = 2                               // Display Message Level - Default=2 (med)
var iSyslogLvl = 0                              // Syslog Level - Default=0 (Off)

// Global File Names
var IniFile = "C:\\AChoir\\AChoir.Acq"          // AChoir Script File
var LogFile = "C:\\AChoir\\LogFile.dat"         // AChoir Log File
var CpyFile = "C:\\AChoir\\LogFile.dat"         // Copy To this File
var HtmFile = "C:\\AChoir\\Index.htm"           // AChoir HTML Output File
var WGetFile = "C:\\AChoir\\Download.dat"       // Downloaded WGet File
var LstFile = "C:\\AChoir\\Data.Lst"            // List of Data
var ChkFile = "C:\\AChoir\\Data.Chk"            // Check For File Existence
var MD5File = "C:\\AChoir\\Hash.txt"            // Saved Hashes
var BACQDir = "C:\\AChoir"                      // Base Acquisition Directory
var BaseDir = "C:\\AChoir"                      // Base Directory
var ACQDir = ""                                 // Relative Acquisition Directory
var CachDir = "C:\\AChoir\\Cache"               // AChoir Caching Directory 
var ForFile = "C:\\AChoir\\Cache\\ForFiles"     // Do action for these Files
var MCpFile = "C:\\AChoir\\Cache\\MCpFiles"     // Do action for Multiple File Copies
var ForDisk = "C:\\AChoir\\Cache\\ForDisk"      // Do Action for Multiple Disk Drives 
var CurrDir = ""                                // Current Directory
var TempDir = ""                                // Temp Build Directory String
var TempAcq = ""                                // Temp ACQ Directory String
var STDOutF = "C:\\AChoir\\Cache\\STDOut"       // Write exec stdout to this file 
var STDErrF = "C:\\AChoir\\Cache\\STDErr"       // Write exec stderr to this file 
var iSTDOut = 0                                 // Exec STDOut is Console(0) or File(1) 
var iSTDErr = 0                                 // Exec STDErr is Console(0) or File(1) 

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
var OpnHndl *os.File                            // User Defined Output File(s)
var MD5Hndl *os.File                            // Save Hashes of Files
var STDOHndl *os.File                           // STDOut File Handle
var STDEHndl *os.File                           // STDErr File Handle
var log_err error                               // Logging Errors
var htm_err error                               // HTML Writer Errors
var ini_err error                               // Ini File Errors
var for_err error                               // For File Errors
var lst_err error                               // Lst File Errors
var dsk_err error                               // Dsk File Errors
var opn_err error                               // User Defined File Errors
var cwd_err error                               // Current Directory Errors
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
var DrvsArray = []string{"C:\\", "D:\\", "E:\\", "F:\\", "G:\\", "H:\\", "I:\\", "J:\\", "K:\\", "L:\\",
                         "M:\\", "N:\\", "O:\\", "P:\\", "Q:\\", "R:\\", "S:\\", "T:\\", "U:\\", "V:\\",
                         "W:\\", "X:\\", "Y:\\", "Z:\\"}

// File Signature Copy Table Variables
var iSigCount = 0
var iSigTMax = 100
var SigTabl [100]string
var TypTabl [100]string


// Main Line
func main() {
    // Get My Name and Path
    MyProg, _ = os.Executable()
    MyHash = GetMD5File(MyProg)

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
    Procesr = opArchit

    opSystem = runtime.GOOS
    switch opSystem {
    case "windows":
        iopSystem = 0
        slashDelim = '\\'
        WinRoot = os.Getenv("SYSTEMROOT")
        //Procesr = os.Getenv("PROCESSOR_ARCHITECTURE")
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

    // If Windows we get Major.Minor.Build - Linux and OSX not implemented yet
    OSVersion = GetOSVer()


    // Initial Settings and Confiiguration
    slashDelimS = fmt.Sprintf("%c", slashDelim)
    ACQName = fmt.Sprintf("ACQ-IR-%s-%04d%02d%02d-%02d%02d", cName, iYYYY, iMonth, iDay, iHour, iMin) 
    inFnam = "AChoir.ACQ"
    iOutOfDiskSpace = 0

    // Default Case Settings 
    caseNumbr = ACQName
    evidNumbr = "001"
    caseDescr = fmt.Sprintf("AChoirX Live Acquisition: %s", ACQName)
    caseExmnr = "Unknown"


    // What Directory are we in?
    BaseDir, cwd_err = os.Getwd()
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
            fmt.Printf(" /GXR:<URL/File> - Get a Zip File using HTTP, Extract the Files, and Run the Script.\n")
            fmt.Printf(" /INI:<File Name> - Run the <File Name> script instead of AChoir.ACQ\n")
            fmt.Printf(" /DEC:<File Name> - Decrypt File using &PWD - Output File Name: Decrypted.dat\n")
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
            if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/INI:CONSOLE") {
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
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/DEC:") {
            inEncFile = os.Args[i][5:]

            if FileExists(inEncFile) {
                plaindata := decryptFile(inEncFile, inPass)

                decFileName, _ := os.Create("Decrypted.dat")
                defer decFileName.Close()
                decFileName.Write(plaindata)
            } else {
                ConsOut = fmt.Sprintf("[!] File to Decrypt Not Found: %s\n", inEncFile)
                ConsLogSys(ConsOut, 1, 1)
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
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/GXR:") {
            WGetURL = os.Args[i][5:]
            CurrFil = fmt.Sprintf("%s%cGXR.Zip", CurrWorkDir, slashDelim)

            fmt.Println("[+] HTTP GXR GetFile: ", WGetURL, CurrFil)

            http_err := DownloadFile(CurrFil, WGetURL)
            if http_err != nil {
                fmt.Println("[!] GXR Downloaded Failed: " + WGetURL)        
            } else {
                fmt.Println("[+] GXR Downloaded Success: " + WGetURL)        
                fmt.Println("[+] Now Expanding GXR Zip File...")        

                ZipRdr, zip_err := zip.OpenReader(CurrFil)
                if zip_err != nil {
                    fmt.Println("[!] GXR Unzip File Open Error: " + CurrFil)        
                } else {
                    defer ZipRdr.Close()

                    Unzfiles, unz_err := Unzip(ZipRdr.File, CurrWorkDir)
                    if unz_err != nil {
                        fmt.Printf("[!] GXR Unzip Error: %s\n", unz_err)
                    }

                    for iUnz := 0; iUnz < len(Unzfiles); iUnz++ {
                        fmt.Printf("[*] Unzipped File: %s\n", Unzfiles[iUnz])
                    }
                }
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
  
    ConsOut = fmt.Sprintf("[+] AChoirX Ver: %s, Mode: %s, OS: %s, Proc: %s\n", Version, RunMode, OSVersion, opArchit)
    ConsLogSys(ConsOut, 1, 1)

    ConsOut = fmt.Sprintf("[+] Path: %s, MD5: %s\n", MyProg, MyHash)
    ConsLogSys(ConsOut, 1, 1)

    showTime("Start Acquisition")

    fmt.Fprintf(LogHndl, "[+] Directory Has Been Set To: %s%c%s\n", BaseDir, slashDelim, CurrDir)
    fmt.Fprintf(LogHndl, "[+] Input Script Set:\n     %s\n\n", IniFile)


    //****************************************************************
    //* Check If We are an Admin/Root                                *
    //****************************************************************
    if (IsUserAdmin()) {
        iIsAdmin = 1
        sIsAdmin = ""
        yIsAdmin = "Yes"
    } else {
        iIsAdmin = 0
        sIsAdmin = "NON-"
        yIsAdmin = "No"
    }


    //****************************************************************
    //* Are we running Non-Native (Sysnative vs. System32)           *
    //****************************************************************
    if iopSystem == 0 {
        //TempDir = fmt.Sprintf("%s\\Sysnative", WinRoot)
        //if _, fol_err := os.Stat(TempDir); os.IsNotExist(fol_err) {
        if strings.ToUpper(Procesr) == "AMD64" {
            sNative = "64Bit "
            iNative = 1
        } else {
            sNative = "32Bit NON-"
            iNative = 0
        }

        ConsOut = fmt.Sprintf("[+] Running as Windows %sNative (%sAdmin)\n", sNative, sIsAdmin)
        ConsLogSys(ConsOut, 1, 1)
    } else {
        ConsOut = fmt.Sprintf("[+] Running as %sRoot\n", sIsAdmin)
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
        fmt.Printf(">>>")

        IniScan = bufio.NewScanner(os.Stdin)

    } else {
        // If the Default IniFile does not exist, let's extract the Default one
        if !FileExists(IniFile) && strings.HasSuffix(strings.ToUpper(IniFile), "ACHOIR.ACQ") {
            ConsOut = fmt.Sprintf("[*] Ini File Does Not Exist UnEmbedding the Default one: %s\n", IniFile)
            ConsLogSys(ConsOut, 1, 2)
            UnEmbed(embdata)
        }

        IniHndl, ini_err = os.Open(IniFile)
        if ini_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Ini File: %s\n", IniFile)
            ConsLogSys(ConsOut, 1, 2)
        }

        IniScan = bufio.NewScanner(IniHndl)
    }

    RunMe = 0  // Conditional run Script default is yes

    for IniScan.Scan() {
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
                if RunMe > 0 {
                    RunMe--
                }
            }
        } else {
            Looper = 1

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

                iMaxCnt = 0
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
                LstMe = 1

                LstHndl, lst_err = os.Open(LstFile)

                if lst_err != nil {
                    ConsOut = fmt.Sprintf("[!] &LST File was not found (LST: not set): %s\n", LstFile)
                    ConsLogSys(ConsOut, 1, 2)

                    Looper = 0
                }

                iMaxCnt = 0
                LstScan = bufio.NewScanner(LstHndl)
            } else {
                LstMe = 0
            }


            //****************************************************************
            //* DskFiles Looper Setup                                        *
            //****************************************************************
            if strings.Contains(strings.ToUpper(Tmprec), "&DSK") {

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
            LoopNum = 0
            NotFound = 0
            YesFound = 0

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
                        iMaxCnt++

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
                        ConsOut = fmt.Sprintf("[+] For Records Processed: %d\n", LoopNum)
                        ConsLogSys(ConsOut, 1, 2)
                        break
                    }
                } else if ForMe == 0 && LstMe == 1 && DskMe == 0 {
                    lst_rcd = LstScan.Scan()
                    lst_err = LstScan.Err()

                    // No Error and no EOF - So Process the Record
                    if lst_err == nil && lst_rcd == true {
                        Lstrec = strings.TrimSpace(LstScan.Text())

                        // Simple way to ignore UTF-16 - Not Great, but it works
                        Lstrec = strings.Replace(Lstrec, "\x00", "", -1) 
                        Lstrec = strings.Replace(Lstrec, "\xfe", "", -1) 
                        Lstrec = strings.Replace(Lstrec, "\xff", "", -1) 
                        //Lstrec = strings.ToValidUTF8(Lstrec, "")   // Force it to be UTF-8

                        Looper = 1
                        LoopNum++
                        iMaxCnt++
                    } else {
                        ConsOut = fmt.Sprintf("[+] Lst Records Processed: %d\n", LoopNum)
                        ConsLogSys(ConsOut, 1, 2)
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
                if strings.Contains(Tmprec, "${") {
                    varConvert(Tmprec)
                 } else {
                     o32VarRec = Tmprec
                 }


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

                if CaseInsensitiveContains(o32VarRec, "&Hsh") {

                    repl_Hsh := NewCaseInsensitiveReplacer("&Hsh", LastHash)
                    o32VarRec = repl_Hsh.Replace(o32VarRec)
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
                    //tokRdr.LazyQuotes = true
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

                if CaseInsensitiveContains(o32VarRec, "&Ver") {

                    repl_Ver := NewCaseInsensitiveReplacer("&Ver", OSVersion)
                    o32VarRec = repl_Ver.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Myp") {

                    repl_Myp := NewCaseInsensitiveReplacer("&Myp", MyProg)
                    o32VarRec = repl_Myp.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Myh") {

                    repl_Myh := NewCaseInsensitiveReplacer("&Myh", MyHash)
                    o32VarRec = repl_Myh.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Vck") {
                    
                    repl_Vck := NewCaseInsensitiveReplacer("&Vck", volType)
                    o32VarRec = repl_Vck.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Adm") {
                    
                    repl_Adm := NewCaseInsensitiveReplacer("&Adm", yIsAdmin)
                    o32VarRec = repl_Adm.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Usr") {
                    
                    repl_Usr := NewCaseInsensitiveReplacer("&Usr", inUser)
                    o32VarRec = repl_Usr.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Pwd") {
                    
                    repl_Pwd := NewCaseInsensitiveReplacer("&Pwd", inPass)
                    o32VarRec = repl_Pwd.Replace(o32VarRec)
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


                if CaseInsensitiveContains(o32VarRec, "&Mem") {

                    TotMemry := sysTotalMemory()
                    if TotMemry == 0 {
                        ConsOut = fmt.Sprintf("[!] Error retrieving Memory stats, or not yet implemented (%s).\n", opSystem)
                        ConsLogSys(ConsOut, 1, 2)
                    }

                    // Even if we got 0 MemoryBytes, replace it.
                    repl_Mem := NewCaseInsensitiveReplacer("&Mem", strconv.FormatUint(TotMemry, 10))
                    o32VarRec = repl_Mem.Replace(o32VarRec)
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
                        getCaseInfo(0)
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
                        iRunMode = 1

                        DirAllocErr(BACQDir)
                        DirAllocErr(CachDir)
                        PreIndex()
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
                            TempDir = fmt.Sprintf("%s%c%s", BACQDir, slashDelim, ACQDir)
                        }
                    }

                    // Determine the Level 1 Directory to see if we have it yet
                    // If not, we will want to add it to the HTML file
                    LvlSplit := strings.Split(ACQDir, slashDelimS)
                    LevelOne := fmt.Sprintf("%s%c%s", BACQDir, slashDelim, LvlSplit[0])

                    if _, level_err := os.Stat(LevelOne); os.IsNotExist(level_err) && iHtmMode == 1 {
                        fmt.Fprintf(HtmHndl, "</td><td align=center>\n")
                        fmt.Fprintf(HtmHndl, "<a href=file:%s target=AFrame> %s </a>\n", LvlSplit[0], LvlSplit[0])
                    }


                    // Have we created this Directory already?
                    if _, ACQDir_err := os.Stat(TempDir); os.IsNotExist(ACQDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Acquisition Sub-Directory: %s\n", ACQDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir)

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
                            TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir)
                        }
                    }

                    // Have we created this Directory already?
                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Sub-Directory: %s\n", CurrDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir)

                    }

                    // Reset The WorkingDirectory to the new Directory
                    CurrWorkDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir)
                    os.Chdir(CurrWorkDir)

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FIL:") {
                    CurrFil = Inrec[4:]
                    TempDir = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, CurrDir)

                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        ConsOut = fmt.Sprintf("[+] Creating Sub-Directory: %s\n", CurrDir)
                        ConsLogSys(ConsOut, 1, 1)

                        ExpandDirs(TempDir)
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

                            ConsOut = fmt.Sprintf("[+] Switching to Console (Interactive Mode)\n")
                            ConsLogSys(ConsOut, 1, 1)

                            IniHndl.Close()
                            //fmt.Printf(">>>")
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

                                cleanUp_Exit(3)
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
                    ConsOut = fmt.Sprintf("[+] Running as %sAdmin/Root\n", sIsAdmin)
                    ConsLogSys(ConsOut, 1, 1)

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ADM:FORCE") {
                    if iIsAdmin == 1 {
                        ConsOut = fmt.Sprintf("[+] Running as Admin/Root - Continuing...\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ConsOut = fmt.Sprintf("[+] NOT Running as Admin/Root - Exiting...\n")
                        ConsLogSys(ConsOut, 1, 1)

                        cleanUp_Exit(3)
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
                    consInput(Inrec[4:], 1, 0)
                    Inprec = Conrec
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "VCK:") {
                    isNTFS = 0
                    volType = winGetVolInfo(Inrec[4:])

                    // This should only work in Windows - Linux and OSX will be UNKNOWN
                    if volType == "NTFS" {
                        isNTFS = 1
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "USR:") {
                    inUser = Inrec[4:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "PWD:") {
                    inPass = Inrec[4:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ENC:") {
                    //****************************************************************
                    //* Encrypt File From => To                                      *
                    //****************************************************************
                    ConsOut = fmt.Sprintf("[+] %s\n", Inrec)
                    ConsLogSys(ConsOut, 1, 1)

                    Encrec = Inrec[4:]

                    splitString1, splitString2, SplitRC := twoSplit(Encrec)

                    if len(splitString1) < 1 {
                        ConsOut = fmt.Sprintf("[!] No File Specified to Encrypt\n")
                        ConsLogSys(ConsOut, 1, 2)
                    } else {
                        if SplitRC == 1 {
                            // Set Output file by appending .ECR to it if no output specified
                            splitString2 = fmt.Sprintf("%s.ECR", splitString1)

                            ConsOut = fmt.Sprintf("[*] Generating Encryption File Name: %s\n", splitString2)
                            ConsLogSys(ConsOut, 1, 1)
                        } 

                        if FileExists(splitString1) {
                            plaindata, _ := ioutil.ReadFile(splitString1)
                            encryptFile(splitString2, plaindata, inPass)

                            if inPass == "none" {
                                ConsOut = fmt.Sprintf("[*] Warning: You are Encrypting with the DEFAULT PASSWORD. This is not recommended.\n")
                                ConsLogSys(ConsOut, 1, 1)
                            }
                        } else {
                            ConsOut = fmt.Sprintf("[!] File to Encrypt Not Found: %s\n", splitString1)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "DEC:") {
                    //****************************************************************
                    //* Decrypt File From => To                                      *
                    //****************************************************************
                    ConsOut = fmt.Sprintf("[+] %s\n", Inrec)
                    ConsLogSys(ConsOut, 1, 1)

                    Encrec = Inrec[4:]

                    splitString1, splitString2, SplitRC := twoSplit(Encrec)

                    if len(splitString1) < 1 {
                        ConsOut = fmt.Sprintf("[!] No File Specified to Decrypt\n")
                        ConsLogSys(ConsOut, 1, 2)
                    } else {

                        if SplitRC == 1 {
                            // Set Output file by appending .DCR to it if no output specified
                            splitString2 = fmt.Sprintf("%s.DCR", splitString1)

                            ConsOut = fmt.Sprintf("[*] Generating Decryption File Name: %s\n", splitString2)
                            ConsLogSys(ConsOut, 1, 1)
                        } 

                        if FileExists(splitString1) {
                            plaindata := decryptFile(splitString1, inPass)

                            decFileName, _ := os.Create(splitString2)
                            defer decFileName.Close()
                            decFileName.Write(plaindata)
                        } else {
                            ConsOut = fmt.Sprintf("[!] File to Decrypt Not Found: %s\n", splitString1)
                            ConsLogSys(ConsOut, 1, 1)
                        }
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
                        //ConsOut = fmt.Sprintf("CPY: %s to %s\n", splitString1, splitString2)
                        //ConsLogSys(ConsOut, 1, 1)


                        //*****************************************************************
                        //* Golang does not support ** - So this code approximates it     *
                        //*  using filepath.Walk.  The limitation is that the string cant *
                        //*  contain another * BEFORE the ** because filepath.Walk does   *
                        //*  not support wildcards. This code is a decent compromise.     *
                        //*****************************************************************
                        DubGlob := fmt.Sprintf("%c**%c", slashDelim, slashDelim)
                        if strings.Contains(splitString1, DubGlob) {
                            iDblShard = strings.Index(splitString1, DubGlob)
                            if iDblShard > 0 {
                                WalkDir := splitString1[:iDblShard]
                                WalkfileWild = splitString1[iDblShard+3:]
                                WalkfileToo = splitString2

                                BasicCopy := fmt.Sprintf("%s%s", WalkDir, WalkfileWild)
                                CopyParser(BasicCopy, splitString2)

                                filepath.Walk(WalkDir, WalkCopyGlob)
                            }
                        } else {
                            CopyParser(splitString1, splitString2)
                        }
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SIG:") {
                    /****************************************************************/
                    /* Clear the File Signature Table, or Load a signature          */
                    /****************************************************************/
                    if strings.HasPrefix(strings.ToUpper(Inrec), "SIG:CLEAR") {
                        iSigCount = 0
                    } else {
                        if strings.Contains(Inrec, "=") {
                            SigSplit := strings.Split(Inrec, "=")
                            if len(SigSplit[0]) > 4 && len(SigSplit[1]) > 0 {
                                TypTabl[iSigCount] = SigSplit[0][4:]
                                SigTabl[iSigCount] = SigSplit[1]

                                // Max Signatures?
                                if iSigCount < iSigTMax {
                                    ConsOut = fmt.Sprintf("[+] Signature Added. Type: %s, Sig: %s, Count: %d/100\n", TypTabl[iSigCount], SigTabl[iSigCount], iSigCount+1)
                                    ConsLogSys(ConsOut, 1, 1)
                                    iSigCount++
                                } else {
                                    ConsOut = fmt.Sprintf("[+] Signature Not Added. Maximum Signature Count is 100\n")
                                    ConsLogSys(ConsOut, 1, 1)
                                }
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "EQU:") {
                    Cmprec = Inrec[4:]
                    splitString1, splitString2, SplitRC := twoSplit(Cmprec)

                    if SplitRC == 1 {
                        ConsOut = fmt.Sprintf("[!] Comparing Requires TWO Strings\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        if(consOrFile == 1) {
                            if splitString1 != splitString2 {
                                ConsOut = fmt.Sprintf("[*] Strings Are NOT Equal: %s != %s\n", splitString1, splitString2)
                                ConsLogSys(ConsOut, 1, 1)
                            } else {
                                ConsOut = fmt.Sprintf("[*] Strings ARE Equal: %s == %s\n", splitString1, splitString2)
                                ConsLogSys(ConsOut, 1, 1)
                            }
                        } else {
                            if splitString1 == splitString2 {
                               // Yes on First Match Only
                                if YesFound == 0 {
                                    YesFound = 1
                                }
                            } else {
                                // No On First Not Match Only
                                if NotFound == 0 {
                                    NotFound = 1
                                }
                            }

                            if NotFound == 1 && YesFound == 0 {
                                // Not Found, Increment Just Once
                                 RunMe++
                                 NotFound = 2
                            } else if YesFound == 1 && NotFound == 2 {
                                // undo the Previous Runme++ and make sure we dont do it again.
                                RunMe--
                                YesFound = 2
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "N>>:") || strings.HasPrefix(strings.ToUpper(Inrec), "N<<:") || strings.HasPrefix(strings.ToUpper(Inrec), "N==:") {
                    Cmprec = Inrec[4:]
                    splitString1, splitString2, SplitRC := twoSplit(Cmprec)

                    if SplitRC == 1 {
                        ConsOut = fmt.Sprintf("[!] Comparing Requires TWO Numbers\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        longString1, _ := strconv.Atoi(splitString1)
                        longString2, _ := strconv.Atoi(splitString2)

                        if strings.HasPrefix(strings.ToUpper(Inrec), "N>>:") {
                            if longString1 > longString2 {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is Greater Than %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // Yes on First Match Only
                                    if YesFound == 0 {
                                        YesFound = 1
                                    }
                                }
                            } else {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is NOT Greater Than %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // No On First Not Match Only
                                    if NotFound == 0 {
                                        NotFound = 1
                                    }
                                }
                            }

                            if(consOrFile != 1) {
                                if NotFound == 1 && YesFound == 0 {
                                     // Not Found, Increment Just Once
                                     RunMe++
                                     NotFound = 2
                                } else if YesFound == 1 && NotFound == 2 {
                                    // undo the Previous Runme++ and make sure we dont do it again.
                                    RunMe--
                                    YesFound = 2
                                }
                            }

                        } else if strings.HasPrefix(strings.ToUpper(Inrec), "N<<:") {
                            if longString1 < longString2 {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is Less Than %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // Yes on First Match Only
                                    if YesFound == 0 {
                                        YesFound = 1
                                    }
                                }
                            } else {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is NOT Less Than %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // No On First Not Match Only
                                    if NotFound == 0 {
                                        NotFound = 1
                                    }
                                }
                            }

                            if(consOrFile != 1) {
                                if NotFound == 1 && YesFound == 0 {
                                     // Not Found, Increment Just Once
                                     RunMe++
                                     NotFound = 2
                                } else if YesFound == 1 && NotFound == 2 {
                                    // undo the Previous Runme++ and make sure we dont do it again.
                                    RunMe--
                                    YesFound = 2
                                }
                            }

                        } else if strings.HasPrefix(strings.ToUpper(Inrec), "N==:") {
                            if longString1 == longString2 {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is Equal To %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // Yes on First Match Only
                                    if YesFound == 0 {
                                        YesFound = 1
                                    }
                                }
                            } else {
                                if(consOrFile == 1) {
                                    ConsOut = fmt.Sprintf("[*] %d Is NOT Equal To %d\n", longString1, longString2)
                                    ConsLogSys(ConsOut, 1, 1)
                                } else {
                                    // No On First Not Match Only
                                    if NotFound == 0 {
                                        NotFound = 1
                                    }
                                }
                            }

                            if(consOrFile != 1) {
                                if NotFound == 1 && YesFound == 0 {
                                     // Not Found, Increment Just Once
                                     RunMe++
                                     NotFound = 2
                                } else if YesFound == 1 && NotFound == 2 {
                                    // undo the Previous Runme++ and make sure we dont do it again.
                                    RunMe--
                                    YesFound = 2
                                }
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "NEQ:") {
                    //****************************************************************
                    //* Check for NOT EQUAL                                          *
                    //****************************************************************
                    Cmprec = Inrec[4:]
                    splitString1, splitString2, SplitRC := twoSplit(Cmprec)

                    if SplitRC == 1 {
                        ConsOut = fmt.Sprintf("[!] Comparing Requires TWO Strings\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        if(consOrFile == 1) {
                            if splitString1 == splitString2 {
                                ConsOut = fmt.Sprintf("[*] Strings are (NOT NOT) Equal: %s == %s\n", splitString1, splitString2)
                                ConsLogSys(ConsOut, 1, 1)
                            } else {
                                ConsOut = fmt.Sprintf("[*] Strings are NOT Equal: %s != %s\n", splitString1, splitString2)
                                ConsLogSys(ConsOut, 1, 1)
                            }
                        } else {
                            if splitString1 != splitString2 {
                               // Yes on First Match Only
                                if YesFound == 0 {
                                    YesFound = 1
                                }
                            } else {
                                // No On First Not Match Only
                                if NotFound == 0 {
                                    NotFound = 1
                                }
                            }

                            if NotFound == 1 && YesFound == 0 {
                                // Not Found, Increment Just Once
                                 RunMe++
                                 NotFound = 2
                            } else if YesFound == 1 && NotFound == 2 {
                                // undo the Previous Runme++ and make sure we dont do it again.
                                RunMe--
                                YesFound = 2
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "VER:") {
                    //****************************************************************
                    //* Check the Input String for Version.  This can be Partial.    *
                    //*  For Example: Windows 10.0.18362 will get a TRUE for         *
                    //*  Ver:Win, Ver:Windows 10, Ver:Windows 10.0.183, etc...       *
                    //****************************************************************
                    if consOrFile == 1 {
                        ConsOut = fmt.Sprintf("[*] OS Version Detected: %s\n", OSVersion)
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        if strings.HasPrefix(strings.ToUpper(OSVersion), strings.ToUpper(Inrec[4:])) {
                            // Yes on First Match Only
                            if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "RC=:") {
                    ChkRC, _ := strconv.Atoi(Inrec[4:]);
                    if consOrFile == 1 {
                        if LastRC != ChkRC {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was not: %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was: %d\n", LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if LastRC == ChkRC {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "RC!:") {
                    ChkRC, _ := strconv.Atoi(Inrec[4:]);
                    if consOrFile == 1 {
                        if LastRC == ChkRC {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was (NOT NOT): %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was not: %d - It was %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if LastRC != ChkRC {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "RC<:") {
                    ChkRC, _ := strconv.Atoi(Inrec[4:]);
                    if consOrFile == 1 {
                        if LastRC >= ChkRC {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was not Less Than: %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was Less Than: %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if LastRC < ChkRC {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "RC>:") {
                    ChkRC, _ := strconv.Atoi(Inrec[4:]);
                    if consOrFile == 1 {
                        if LastRC <= ChkRC {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was not Greater Than: %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Last Return Code was Greater Than: %d - It was: %d\n", ChkRC, LastRC)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if LastRC > ChkRC {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CKY:") {
                    ChkFile = strings.Trim(Inrec[4:], "\"")

                    if consOrFile == 1 {
                        if !FileExists(ChkFile) {
                            ConsOut = fmt.Sprintf("[*] File Does Not Exist: %s\n", ChkFile)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] File Exists: %s\n", ChkFile)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if FileExists(ChkFile) {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CKN:") {
                    ChkFile = strings.Trim(Inrec[4:], "\"")
                    if consOrFile == 1 {
                        if FileExists(ChkFile) {
                            ConsOut = fmt.Sprintf("[*] File Does (NOT NOT) Exist: %s\n", ChkFile)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] File Does Not Exist: %s\n", ChkFile)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if !FileExists(ChkFile) {
                           // Yes on First Match Only
                           if YesFound == 0 {
                                YesFound = 1
                            }
                        } else {
                            // No On First Not Match Only
                            if NotFound == 0 {
                                NotFound = 1
                            }
                        }

                        if NotFound == 1 && YesFound == 0 {
                            // Not Found, Increment Just Once
                             RunMe++
                             NotFound = 2
                        } else if YesFound == 1 && NotFound == 2 {
                            // undo the Previous Runme++ and make sure we dont do it again.
                            RunMe--
                            YesFound = 2
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "64B:") {
                    if consOrFile == 1 {
                        if strings.ToUpper(Procesr) != "AMD64" {
                            ConsOut = fmt.Sprintf("[*] Not running in 64Bit. Processor: %s\n", Procesr)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Running in 64Bit. Processor: %s\n", Procesr)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if strings.ToUpper(Procesr) != "AMD64" {
                            RunMe++
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "32B:") {
                    if consOrFile == 1 {
                        if strings.ToUpper(Procesr) != "386" {
                            ConsOut = fmt.Sprintf("[*] Not running in 32Bit. Processor: %s\n", Procesr)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            ConsOut = fmt.Sprintf("[*] Running in 32Bit. Processor: %s\n", Procesr)
                            ConsLogSys(ConsOut, 1, 1)
                        }
                    } else {
                        if strings.ToUpper(Procesr) != "386" {
                            RunMe++
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "REQ:") {
                    if FileExists(Inrec[4:]) {
                        ConsOut = fmt.Sprintf("[*] [!] Required File Found: %s\n", Inrec[4:])
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ConsOut = fmt.Sprintf("[*] [*] Required File Not Found: %s - Exiting!\n", Inrec[4:])
                        ConsLogSys(ConsOut, 1, 1)

                        cleanUp_Exit(3)
                        os.Exit(3)
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SAY:") {
                    ConsOut = fmt.Sprintf("%s\n", Inrec[4:])
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "END:") {
                    if RunMe > 0 {
                        RunMe--
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "OPN:") {
                    // If we already had a file open, close it now.
                    if iOPNisOpen == 1 {
                        ConsOut = fmt.Sprintf("[*] Previously Opened File has been Closed.\n")
                        ConsLogSys(ConsOut, 1, 1)

                        OpnHndl.Close()
                    }

                    OpnHndl, opn_err = os.OpenFile(Inrec[4:], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                    if opn_err != nil {
                        ConsOut = fmt.Sprintf("[!] File Could not be opened for Append:\n    %s\n", Inrec[4:])
                        ConsLogSys(ConsOut, 1, 1)

                        iOPNisOpen = 0
                    } else {
                        ConsOut = fmt.Sprintf("[+] File Opened for Append:\n    %s\n", Inrec[4:])
                        ConsLogSys(ConsOut, 1, 1)

                        iOPNisOpen = 1
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "OUT:") {
                    if iOPNisOpen == 1 {
                        OpnHndl.WriteString(Inrec[4:]+"\n")
                    } else {
                        ConsOut = fmt.Sprintf("[!] No File OPN:(ed) for Append:\n")
                        ConsLogSys(ConsOut, 1, 1)
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "PZZ:") {
                    ConsOut = fmt.Sprintf("%s\n", Inrec[4:])
                    ConsLogSys(ConsOut, 1, 1)

                    consInput(Inrec[4:], 1, 0)
                    Inprec = Conrec

                    if len(Inprec) > 0 {
                        if Inprec[0] == 'q' || Inprec[0] == 'Q' {
                            ConsOut = fmt.Sprintf("[+] You Have Requested AChoirX to Quit.\n")
                            ConsLogSys(ConsOut, 1, 1)

                            cleanUp_Exit(0)
                            os.Exit(0)
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "HSH:") {
                    if len(Inrec) > 4 {
                        if strings.ToUpper(Inrec[4:]) == "ACQ" {
                            ConsOut = fmt.Sprintf("[+] Now Hashing  Acquisition Files.\n")
                            ConsLogSys(ConsOut, 1, 1)
                            MD5File = fmt.Sprintf("%s%cACQHash.txt", BACQDir, slashDelim)
                            TempDir = BACQDir
                        } else if strings.ToUpper(Inrec[4:]) == "DIR" {
                            ConsOut = fmt.Sprintf("[+] Now Hashing  Entire AChoirX Directory.\n")
                            ConsLogSys(ConsOut, 1, 1)
                            MD5File = fmt.Sprintf("%s%cDirHash.txt", BaseDir, slashDelim)
                            TempDir = BaseDir
                        } else {
                            LastHash = GetMD5File(strings.Trim(Inrec[4:], "\""))
                            break
                        }
                    } else {
                        ConsOut = fmt.Sprintf("[!] No Hashing Parameter Specified.\n")
                        ConsLogSys(ConsOut, 1, 1)
                        break
                    }
                    //****************************************************************
                    // Do File Search using Walk because Glob does not support **    *
                    //****************************************************************
                    MD5Hndl, opn_err = os.OpenFile(MD5File, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
                    if opn_err != nil {
                        ConsOut = fmt.Sprintf("[!] Error Openning Hash Log: %s\n", opn_err)
                        ConsLogSys(ConsOut, 1, 1)
                        break
                    }

                    filepath.Walk(TempDir, MD5FileOut)

                    MD5Hndl.Close()
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "DSK:") {
                    // Checking Drive Types for Windows only
                    // DRIVE_CDROM = 5, DRIVE_FIXED = 3, DRIVE_RAMDISK = 6, DRIVE_REMOTE = 4, DRIVE_REMOVABLE = 2
                    if iopSystem == 0 {
                        ExpandDirs(CachDir)
                        ForDisk = fmt.Sprintf("%s%cForDisks", CachDir, slashDelim)

                        if strings.HasPrefix(strings.ToUpper(Inrec[4:]), "REMOV") {
                            dskTyp = 2
                        } else if strings.HasPrefix(strings.ToUpper(Inrec[4:]), "FIXED") {
                            dskTyp = 3
                        } else if strings.HasPrefix(strings.ToUpper(Inrec[4:]), "REMOT") {
                            dskTyp = 4
                        } else if strings.HasPrefix(strings.ToUpper(Inrec[4:]), "CDROM") {
                            dskTyp = 5
                        } else if strings.HasPrefix(strings.ToUpper(Inrec[4:]), "RAMDI") {
                            dskTyp = 6
                        } else {
                            dskTyp = 3
                        }

                        DskHndl, dsk_err = os.OpenFile(ForDisk, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
                        if dsk_err != nil {
                            ConsOut = fmt.Sprintf("[!] System Disk Tracking File Could not be opened.\n")
                            ConsLogSys(ConsOut, 1, 1)
                            break
                        }


                        for iDrv := 0; iDrv < 24; iDrv++ {
                            GotDriveType := GetDriveType(DrvsArray[iDrv])

                            if GotDriveType == dskTyp {
                                //DskHndl.WriteString(DrvsArray[iDrv] + "\n")

                                DrvLtr := fmt.Sprintf("%c\n", DrvsArray[iDrv][0])
                                DskHndl.WriteString(DrvLtr)
                            }
                        }

                        DskHndl.Close()

                    } else {
                        ConsOut = fmt.Sprintf("[!] Bypassing Drive and Memory Routines - We are not running on Windows\n")
                        ConsLogSys(ConsOut, 1, 1)
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FOR:") {
                    ExpandDirs(CachDir)
                    ForFile = fmt.Sprintf("%s%cForFiles", CachDir, slashDelim)
                    TempDir = Inrec[4:]

                    ForHndl, for_err = os.OpenFile(ForFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
                    if for_err != nil {
                        ConsOut = fmt.Sprintf("[!] System File Tracking File Could not be opened.\n")
                        ConsLogSys(ConsOut, 1, 1)
                        break
                    }


                    //*****************************************************************
                    //* Golang does not support ** - So this code approximates it     *
                    //*  using filepath.Walk.  The limitation is that the string cant *
                    //*  contain another * BEFORE the ** because filepath.Walk does   *
                    //*  not support wildcards. This code is a decent compromise.     *
                    //*****************************************************************
                    DubGlob := fmt.Sprintf("%c**%c", slashDelim, slashDelim)
                    if strings.Contains(TempDir, DubGlob) {
                        iDblShard = strings.Index(TempDir, DubGlob)
                        if iDblShard > 0 {
                            WalkDir := TempDir[:iDblShard]
                            WalkfileWild = TempDir[iDblShard+3:]

                            BasicFor := fmt.Sprintf("%s%s", WalkDir, WalkfileWild)
                            ForParser(BasicFor)

                            filepath.Walk(WalkDir, WalkForGlob)
                        }
                    } else {
                        ForParser(TempDir)
                    }

                    ForHndl.Close()

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "LST:") {
                    LstFile = fmt.Sprintf("%s%c%s", BaseDir,slashDelim, Inrec[4:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "END:") {
                    RunMe--
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "BYE:") {
                    cleanUp_Exit(LastRC);
                    os.Exit(LastRC);
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:DELIMS=") {
                    Delims = Inrec[11:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:COPYPATH=NONE") {
                    setCPath = 0
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:COPYPATH=PART") {
                    setCPath = 1
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:COPYPATH=FULL") {
                    setCPath = 2
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:COPYDEPTH=") {
                    setCDepth, _ = strconv.Atoi(Inrec[14:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:EXESTDOUT=") {
                    if strings.HasPrefix(strings.ToUpper(Inrec[14:]), "CONS") {
                        iSTDOut = 0
                    } else {
                        iSTDOut = 1
                        STDOutF = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, Inrec[14:])
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:EXESTDERR=") {
                    if strings.HasPrefix(strings.ToUpper(Inrec[14:]), "CONS") {
                        iSTDErr = 0
                    } else {
                        iSTDErr = 1
                        STDErrF = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, Inrec[14:])
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:SYSLOGS=") {
                    Syslogd = Inrec[12:]
                    if iSyslogLvl < 1 {
                        iSyslogLvl = 1
                    }

                    ConsOut = fmt.Sprintf("[*] AChoir Version: %s Syslogging Started.  Level: %d  ACQ: %s\n", Version, iSyslogLvl, ACQName)
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:SYSLOGP=") {
                    Syslogp = Inrec[12:]
                    iSyslogp, _ = strconv.Atoi(Syslogp)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:SYSLOGL=NONE") {
                    if iSyslogLvl > 0 {
                        ConsOut = fmt.Sprintf("[*] AChoir Version: %s Syslogging Stopped.  Old Level = %d\n", Version, iSyslogLvl)
                        ConsLogSys(ConsOut, 1, 1)
                    }

                    iSyslogLvl = 0; 
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:SYSLOGL=MIN") {
                    ConsOut = fmt.Sprintf("[*] Syslog Level Set to Min: 1\n")
                    ConsLogSys(ConsOut, 1, 1)

                    iSyslogLvl = 1 
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:SYSLOGL=MAX") {
                    ConsOut = fmt.Sprintf("[*] Syslog Level Set to Max: 2\n")
                    ConsLogSys(ConsOut, 1, 1)

                    iSyslogLvl = 2
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:S3REGION=") {
                    S3_REGION = Inrec[13:]
                    iS3Login = 0  // Reset Login to force a New Session with this Region

                    ConsOut = fmt.Sprintf("[*] S3 Region Set: %s\n", S3_REGION)
                    ConsLogSys(ConsOut, 1, 1)

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:S3BUCKET=") {
                    S3_BUCKET = Inrec[13:]
                    iS3Login = 0  // Reset Login to force a New Session with this Bucket

                    ConsOut = fmt.Sprintf("[*] S3 Bucket Set: %s\n", S3_BUCKET)
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:S3AWSID=") {
                    S3_AWSId = Inrec[12:]
                    iS3Login = 0  // Reset Login to force a New Session with this AWSId

                    ConsOut = fmt.Sprintf("[*] S3 AWS ID Set: %s\n", S3_AWSId)
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:S3AWSKEY=") {
                    S3_AWSKey = Inrec[13:]
                    iS3Login = 0  // Reset Login to force a New Session with this Key

                    ConsOut = fmt.Sprintf("[*] S3 AWS Key Set: <Redacted>\n")
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "XIT:") && len(Inrec) > 4 {
                    iXitCmd = 1

                    if Inrec[4] == '\\' {
                        XitCmd = fmt.Sprintf("%s%s", BaseDir, Inrec[4:])
                    } else {
                        XitCmd = fmt.Sprintf("%s", Inrec[4:])
                    }

                    ConsOut = fmt.Sprintf("[*] Exit Program Set: %s\n", XitCmd)
                    ConsLogSys(ConsOut, 1, 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "EXE:") {
                    RunCommand(Inrec[4:], 1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "EXA:") {
                    RunCommand(Inrec[4:], 2)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SYS:") {
                    RunCommand(Inrec[4:], 3)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "GET:") {
                    Getrec = Inrec[4:]
                    splitString1, splitString2, SplitRC := twoSplit(Getrec)
                    WGetURL = splitString1
                    CurrFil = splitString2

                    if SplitRC == 1 {
                        WGetDir = fmt.Sprintf("%s%c%s", BACQDir, slashDelim, ACQDir)
                        CurrFil = fmt.Sprintf("%s%c%s%cDownload.dat", BACQDir, slashDelim, ACQDir, slashDelim)
                        ExpandDirs(WGetDir)

                        ConsOut = fmt.Sprintf("[!] Target Download File Name Not Specified.  Using: %s\n", CurrFil)
                        ConsLogSys(ConsOut, 1, 1)
                    }

                    ConsOut = fmt.Sprintf("[+] HTTP GetFile: %s (%s)\n", WGetURL, CurrFil)
                    ConsLogSys(ConsOut, 1, 1)

                    http_err := DownloadFile(CurrFil, WGetURL)
                    if http_err != nil {
                        ConsOut = fmt.Sprintf("[+] Download Failed: %s\n", WGetURL)
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ConsOut = fmt.Sprintf("[+] Download Success: %s\n", WGetURL)
                        ConsLogSys(ConsOut, 1, 1)
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "UNZ:") {
                    Ziprec = Inrec[4:]
                    splitString1, splitString2, SplitRC := twoSplit(Ziprec)

                    if SplitRC == 1 {
                        ConsOut = fmt.Sprintf("[!] Unzip Requires both a Zip File Name and Destination Directory\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        ZipRdr, zip_err := zip.OpenReader(splitString1)
                        if zip_err != nil {
                            ConsOut = fmt.Sprintf("[!] Cannot Open Zip File: %s\n", splitString1)
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            defer ZipRdr.Close()

                            Unzfiles, unz_err := Unzip(ZipRdr.File, splitString2)
                            if unz_err != nil {
                                ConsOut = fmt.Sprintf("[!] Unzip Error: %s\n", unz_err)
                                ConsLogSys(ConsOut, 1, 1)
                            }

                            for iUnz := 0; iUnz < len(Unzfiles); iUnz++ {
                                ConsOut = fmt.Sprintf("[*] Unzipped File: %s\n", Unzfiles[iUnz])
                                ConsLogSys(ConsOut, 1, 1)
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "S3S:") {
                    //****************************************************************
                    //* Have we set the Region and Bucket Name?                      *
                    //****************************************************************
                    if S3_REGION == "none" {
                        ConsOut = fmt.Sprintf("[!] Please Set the AWS S3 Bucket REGION Before Starting an AWS Session.\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else if S3_BUCKET == "none" {
                        ConsOut = fmt.Sprintf("[!] Please Set the AWS S3 BUCKET Name Before Starting an AWS Session.\n")
                        ConsLogSys(ConsOut, 1, 1)
                    } else {
                        //****************************************************************
                        //* AWS S3 Bucket Login Parm1:UserID  Parm2: Secret Key          *
                        //****************************************************************
                        ConsOut = fmt.Sprintf("[+] Starting Session with AWS Key and Secret...\n")
                        ConsLogSys(ConsOut, 1, 1)

                        AWSrec := Inrec[4:]
                        S3_AWSId, S3_AWSKey, AWS_SplitRC := twoSplit(AWSrec)

                        if AWS_SplitRC == 1 {
                            ConsOut = fmt.Sprintf("[!] AWS Session Requires Both an ID and a Key\n")
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            S3_Session, S3_err = session.NewSession(&aws.Config {
                                Region: aws.String(S3_REGION),
                                Credentials: credentials.NewStaticCredentials(
                                S3_AWSId, S3_AWSKey, ""),
                            })

                            if S3_err != nil {
                                ConsOut = fmt.Sprintf("[!] Error Starting AWS Session for S3: %s\n", S3_err)
                                ConsLogSys(ConsOut, 1, 1)
                                iS3Login = 0
                            } else {
                                ConsOut = fmt.Sprintf("[*] AWS S3 Session Started...\n")
                                ConsLogSys(ConsOut, 1, 1)
                                iS3Login = 1
                            }
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "S3U:") {
                    //****************************************************************
                    //* See if we have a Sesion.  If not, See if we can start one    *
                    //****************************************************************
                    if iS3Login == 0 {
                        ConsOut = fmt.Sprintf("[+] Checking for AWS Bucket, Region, ID, and Key...\n")
                        ConsLogSys(ConsOut, 1, 1)

                        if S3_AWSId != "none" && S3_AWSKey != "none" && S3_REGION != "none" && S3_BUCKET != "none" {
                            ConsOut = fmt.Sprintf("[+] Starting Session with AWS Key and Secret...\n")
                            ConsLogSys(ConsOut, 1, 1)

                            S3_Session, S3_err = session.NewSession(&aws.Config {
                                Region: aws.String(S3_REGION),
                                Credentials: credentials.NewStaticCredentials(
                                S3_AWSId, S3_AWSKey, ""),
                            })

                            if S3_err != nil {
                                ConsOut = fmt.Sprintf("[!] Error Starting AWS Session for S3: %s\n", S3_err)
                                ConsLogSys(ConsOut, 1, 1)
                                iS3Login = 0
                            } else {
                                ConsOut = fmt.Sprintf("[*] AWS S3 Session Started...\n")
                                ConsLogSys(ConsOut, 1, 1)
                                iS3Login = 1
                            }
                        }
                    }


                    //****************************************************************
                    //* Upload File to S3                                            *
                    //****************************************************************
                    if iS3Login == 1 {
                        ConsOut = fmt.Sprintf("[+] %s\n", Inrec)
                        ConsLogSys(ConsOut, 1, 1)

                        S3Urec = Inrec[4:]

                        splitString1, splitString2, SplitRC := twoSplit(S3Urec)

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
                            ConsOut = fmt.Sprintf("[!] S3 Upload Requires both a FROM File and a TO Directory\n")
                            ConsLogSys(ConsOut, 1, 1)
                        } else {
                            //ConsOut = fmt.Sprintf("S3U: %s to %s\n", splitString1, splitString2)
                            //ConsLogSys(ConsOut, 1, 1)

                            //*****************************************************************
                            //* Golang does not support ** - So this code approximates it     *
                            //*  using filepath.Walk.  The limitation is that the string cant *
                            //*  contain another * BEFORE the ** because filepath.Walk does   *
                            //*  not support wildcards. This code is a decent compromise.     *
                            //*****************************************************************
                            DubGlob := fmt.Sprintf("%c**%c", slashDelim, slashDelim)
                            if strings.Contains(splitString1, DubGlob) {
                                iDblShard = strings.Index(splitString1, DubGlob)
                                if iDblShard > 0 {
                                    WalkDir := splitString1[:iDblShard]
                                    WalkfileWild = splitString1[iDblShard+3:]
                                    WalkfileToo = splitString2

                                    BasicS3Up := fmt.Sprintf("%s%s", WalkDir, WalkfileWild)
                                    S3UpParser(BasicS3Up, splitString2)

                                    filepath.Walk(WalkDir, WalkS3UpGlob)
                                
                                }
                            } else {
                                S3UpParser(splitString1, splitString2)
                            }
                        }
                    } else {
                        ConsOut = fmt.Sprintf("[!] Please Start an AWS Session (Using S3S:) Before Attempting S3 Uploads\n")
                        ConsLogSys(ConsOut, 1, 1)
                    }
                }
            }
        }

        if consOrFile == 1 {
            fmt.Printf(">>>")
        }
    }


    //****************************************************************
    //* Cleanup                                                      *
    //****************************************************************
    if RunMe > 0 {
        ConsOut = fmt.Sprintf("[!] You have and extra END: Hanging! Check your Logic.\n")
        ConsLogSys(ConsOut, 1, 1)
    }

    cleanUp_Exit(0);
    os.Exit(0);

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
        ConsOut = fmt.Sprintf("[+] %s: %02d/%02d/%04d - %02d:%02d:%02d\n", showText,
        showlocal.Month(), showlocal.Day(), showlocal.Year(),
        showlocal.Hour(), showlocal.Minute(), showlocal.Second())

        ConsLogSys(ConsOut, 1, 1)
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
            fmt.Fprintf(LogHndl, "[!] Could not Create Artifact Index: %s\n", HtmFile)
        }
        if (setMSGLvl > 1) {
            fmt.Printf("[!] Could not Create Artifact Index: %s\n", HtmFile)
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
        fmt.Fprintf(LogHndl, ConLogMSG)
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
    iCase = 1
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
            ConsOut = fmt.Sprintf("[?] %s\n", Conrec)
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
    // o32VarRec = os.ExpandEnv(inVarRec)  // Deprecated
    // New code below ensures that only the ${} type expansion is done
    //  this is to avoid expansion of items like $MFT as a variable
    o32VarRec = inVarRec
    VaRegex := regexp.MustCompile(`\$\{[a-zA-Z0-9_]+\}`)
    VaArray := VaRegex.FindAllString(o32VarRec, -1)

    for iRgx := 0; iRgx < len(VaArray); iRgx++ {
        o32VarRec = strings.Replace(o32VarRec, VaArray[iRgx], os.ExpandEnv(VaArray[iRgx]), -1) 
    }


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

            DirAllocErr(TempDirName)
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
    var iCPSFound = 0
    var iSig = 0

    var FrmATime, FrmMTime, FrmCTime time.Time
    var TooATime, TooMTime, TooCTime time.Time

    var FrmMD5, TooMD5 string

    //***********************************************************************
    //* Does the From File Exist?                                           *
    //***********************************************************************
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
    //* Make sure we have enough free disk space                            *
    //***********************************************************************
    _, FreeBytes := winFreeDisk()
    FrmFileSize := FrmFileStat.Size()
    if (int64(FreeBytes) < FrmFileSize) {
        iOutOfDiskSpace = 1
        return 0, fmt.Errorf("[!] Copy Error: Not Enough Disk Space Available: %d", FreeBytes)
    }


    //***********************************************************************
    //* MD5 the Input File                                                  *
    //***********************************************************************
    FrmMD5 = GetMD5File(FrmFile)


    //***********************************************************************
    //* Open it up with Defer Close                                         *
    //***********************************************************************
    FrmSource, frm_err := os.Open(FrmFile)
    if frm_err != nil {
        return 0, fmt.Errorf("[!] Copy Error: %s", frm_err)
    }
    defer FrmSource.Close()


    //***********************************************************************
    //* Get FrmSource File MetaData                                         *
    //***********************************************************************
    FrmATime, FrmMTime, FrmCTime = FTime(FrmFile)

    ConsOut = fmt.Sprintf("[+] From File Meta Data: \n    Bytes: %d\n    ATime: %v\n    MTime: %v\n    CTime: %v\n", FrmFileSize, FrmATime, FrmMTime, FrmCTime)
    ConsLogSys(ConsOut, 1, 1)


    //***********************************************************************
    //* If we are doing CPS (Copy by Signature) check for Magic Numbers     *
    //***********************************************************************
    iCPSFound = 0     // Default to NOT Found

    // For CPY: it's ALWAYS found, for CPS: do the compare
    if iCPS == 0 {
        iCPSFound = 1
    } else if iCPS == 1 {
        /****************************************************************/
        /* If we are doing an CPS - Read the first 32 Bytes and compare */
        /*  to the Signature Table entries                              */
        /****************************************************************/
        // Read in signature bytes
        tmpSigBytes := make([]byte, 32)
        inSize, sigb_err := FrmSource.Read(tmpSigBytes)
        tmpSig := hex.EncodeToString(tmpSigBytes)

        // Parse Out the FileType for Signature Checking
        filetype := filepath.Ext(FrmFile)
        if len(filetype) > 1 {
            filetype = filetype[1:]
        } else {
            filetype = ""
        }


        // Compare with the Signature and FileType Tables
        for iSig = 0; iSig < iSigCount; iSig++ {
            // Make Sure we got some data from the File to Check
            if inSize > 0 && sigb_err == nil {
                if strings.HasPrefix(tmpSig, SigTabl[iSig]) {
                    iCPSFound = 1
                    ConsOut = fmt.Sprintf("[*] Header Signature Match Found in File: %s\n", SigTabl[iSig])
                    ConsLogSys(ConsOut, 2, 2)
                    break
                }
            }

            if filetype == TypTabl[iSig] {
                iCPSFound = 1
                ConsOut = fmt.Sprintf("[*] File Extension Match Found: %s\n", filetype)
                ConsLogSys(ConsOut, 2, 2)
                break
            }
        }
    }

    if iCPSFound == 1 {
        FrmSource.Seek(0, 0)
    } else {
        return 0, fmt.Errorf("[!] No Signature Match in File. File Bypassed\n")
    }


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
            iFileCount++
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

    // File Copy
    nBytes, copy_err := io.Copy(TooDest, FrmSource)


    //***********************************************************************
    //* If Copy worked, Collect Output file Meta Data                       *
    //***********************************************************************
    if copy_err == nil {
        os.Chtimes(TooFile, FrmATime, FrmMTime)

        TooMD5 = GetMD5File(TooFile)

        TooFileStat, stat_err := os.Stat(TooFile)
        if stat_err == nil {
            TooFileSize := TooFileStat.Size()
            TooATime, TooMTime, TooCTime = FTime(TooFile)

            ConsOut = fmt.Sprintf("[+] Copy Complete: %d Bytes Copied\n", nBytes)
            ConsLogSys(ConsOut, 1, 1)

            ConsOut = fmt.Sprintf("[+] To File Meta Data: \n    Bytes: %d\n    ATime: %v\n    MTime: %v\n    CTime: %v\n", TooFileSize, TooATime, TooMTime, TooCTime)
            ConsLogSys(ConsOut, 1, 1)

            if FrmMD5 == TooMD5 {
                ConsOut = fmt.Sprintf("[+] From MD5: %s and To MD5: %s Match\n", FrmMD5, TooMD5)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("[!] From MD5: %s and To MD5: %s DO NOT MATCH\n", FrmMD5, TooMD5)
                ConsLogSys(ConsOut, 1, 1)
            }

            if FrmFileSize == TooFileSize {
                ConsOut = fmt.Sprintf("[+] From File Size: %d and To File Size: %d Match\n", FrmFileSize, TooFileSize)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("[+] From File Size: %d and To File Size: %d DO NOT MATCH\n", FrmFileSize, TooFileSize)
                ConsLogSys(ConsOut, 1, 1)
            }

        }
    }


    return nBytes, copy_err
}


//***********************************************************************
// FileExists checks if a file/directory exists before we try using it  *
//  to prevent further errors.                                          *
//***********************************************************************
func FileExists(CheckFilename string) bool {
    // See if the File is already there
    fexst_info, fexst_err := os.Stat(CheckFilename)

    if fexst_err != nil { return false }
    if os.IsNotExist(fexst_err) { return false }
    if fexst_info.Mode().IsRegular() { return true }
    if fexst_info.Mode().IsDir() { return true }

    // Fell Through - This is not a File or Directory
    return false
}


//***********************************************************************
// MD5 A File                                                           *
//***********************************************************************
func GetMD5File(InFileName string) string {
    MD5Hndl, open_err := os.Open(InFileName)
    if open_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Generating MD5 for: %s\n", InFileName)
        ConsLogSys(ConsOut, 1, 1)
        return "00000000000000000000000000000000"  //No such MD5
    }
    defer MD5Hndl.Close()

    MD5New := md5.New()
    if _, MD5_err := io.Copy(MD5New, MD5Hndl); MD5_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Generating MD5 for: %s\n", InFileName)
        ConsLogSys(ConsOut, 1, 1)
        return "00000000000000000000000000000000"  //No such MD5
    }

    return hex.EncodeToString(MD5New.Sum(nil))

}


func cleanUp_Exit(exitRC int) {
    //****************************************************************
    //* Cleanup and Exit                                             *
    //****************************************************************
    if FileExists(ForFile) {
        os.Remove(ForFile)
    }

    if FileExists(ForDisk) {
        os.Remove(ForDisk)
    }

    if FileExists(MCpFile) {
        os.Remove(MCpFile)
    }


    if iHtmMode == 1 {
        fmt.Fprintf(HtmHndl, "</td><td align=right>\n");
        fmt.Fprintf(HtmHndl, "<button onclick=\"window.history.forward()\">&gt;&gt;</button>\n");
        fmt.Fprintf(HtmHndl, "</td></tr></table></Center>\n<p>\n");
        fmt.Fprintf(HtmHndl, "<iframe name=AFrame style=\"padding:2px;border:3px Lavender solid;\"  height=75%% width=98%% scrolling=auto src=file:./></iframe>\n");
        fmt.Fprintf(HtmHndl, "</p>\n</body></html>\n");

        HtmHndl.Close();
    }


    //****************************************************************
    //* Cleanup - Did we run out of Disk Space?                      *
    //****************************************************************
    if iOutOfDiskSpace == 1 {
        ConsOut = fmt.Sprintf("[!] WARNING: Files Copied During this Acquisition RAN OUT OF DISK SPACE\n")
        ConsLogSys(ConsOut, 1, 1)
    }


    //****************************************************************
    //* Cleanup - Normal Run (RunMode 1)                             *
    //****************************************************************
    if iRunMode == 1 {
        ConsOut = fmt.Sprintf("[+] Setting All Artifacts to Read-Only.\n")
        ConsLogSys(ConsOut, 1, 1)

        filepath.Walk(BACQDir, SetReadOnly)
    }


    //****************************************************************
    //* All Done with Acquisition                    `               *
    //****************************************************************
    showTime("Acquisition Completed");

    if iXitCmd == 1 {
        ConsOut = fmt.Sprintf("[+] Queuing Exit Program: %s\n", XitCmd)
        ConsLogSys(ConsOut, 1, 1)
    }


    //****************************************************************
    //* Make a Copy of the Logfile in the ACQ Directory              *
    //****************************************************************
    if _, BACQDir_err := os.Stat(BACQDir); os.IsNotExist(BACQDir_err) {
        ConsOut = fmt.Sprintf("[+] Base Acquisition Directory Not Found: %s\n", BACQDir)
        ConsLogSys(ConsOut, 1, 1)
    } else {
        iCPS = 0; //ALWAYS Copy LogFile

        ConsOut = fmt.Sprintf("[+] Copying Log File...\n")
        ConsLogSys(ConsOut, 1, 1)

        //Very Last Log Entry - Close Log now, and copy WITHOUT LOGGING
        LogHndl.Close();

        //Reset setCPath to Relative and copy Log
        setCPath = 0

        CpyFile = fmt.Sprintf("%s%c%s.Log", BACQDir, slashDelim, ACQName)
        nBytes, copy_err := binCopy(LogFile, CpyFile)

        if copy_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Copying File: %s\n", copy_err)
            ConsLogSys(ConsOut, 1, 1)

            if nBytes < 1 {
                ConsOut = fmt.Sprintf("[!] File Copy was: %d Bytes\n", nBytes)
                ConsLogSys(ConsOut, 1, 1)

            }
        }
    }


    //****************************************************************
    //* Run Final Exit Program - This will not be logged             *
    //****************************************************************
    if iXitCmd == 1 {
        RunCommand(XitCmd, 1)
    }


    /****************************************************************/
    /* Final Syslog Out                                             */
    /****************************************************************/
    if iSyslogLvl > 0 {
        ConsOut = fmt.Sprintf("[+] AChoir Version: %s Acquisition Completed.  Return Code: %d  ACQ: %s", Version, exitRC, ACQName)
        ConsLogSys(ConsOut, 1, 1)
    }

}


func SetReadOnly(Setfilepath string, SetInfo os.FileInfo, Set_err error) error {
    //****************************************************************
    //* Set File to Read Only                                        *
    //****************************************************************
    ROS_err := os.Chmod(Setfilepath, 0444)
    return ROS_err
 }


func MD5FileOut(MD5filepath string, MD5Info os.FileInfo, MD5_err error) error {
    //****************************************************************
    //* Hash the Files and Write Them to a Log                       *
    //****************************************************************
    //* Ignore Directories - Only MD5 Files                          *
    //****************************************************************
    file_stat, _ := os.Stat(MD5filepath)
    if file_stat.IsDir() {
        MD5Hndl.WriteString("Directory: " + MD5filepath + "\n")
    } else {
        //***********************************************************************
        //* MD5 the Input File                                                  *
        //***********************************************************************
        FileMD5 := GetMD5File(MD5filepath)
        MD5Hndl.WriteString("File: " + MD5filepath + " - MD5: " + FileMD5 +"\n")
    }

    return nil

}


func RunCommand(Commandstring string, Commandtype int) error {
    var Fullpath = BaseDir
    var FullCopyCommand = "none"
    var execMD5 = "00000000000000000000000000000000"
    var cmdSplit []string = nil
    var tmpSTDRedir = 0

    //****************************************************************
    //* Run an External Command, either blocked or unblocked         *
    //****************************************************************
    //tmpSplit := strings.Split(Commandstring, " ")

    // Support Quotes - as of v10.00.27 
    CmdTokRdr := csv.NewReader(strings.NewReader(Commandstring))
    CmdTokRdr.Comma = ' '
    CmdTokRdr.FieldsPerRecord = -1
    CmdTokRdr.TrimLeadingSpace = true
    tmpSplit, CmdTok_err := CmdTokRdr.Read()

    if CmdTok_err != nil {
        ConsOut = fmt.Sprintf("[!] Parsing Error for Command: %s    \n%s\n", Commandstring, CmdTok_err)
        ConsLogSys(ConsOut, 1, 2)
        return CmdTok_err
    }


    //****************************************************************
    //* Look for AChoirX Command Enhancements. Process them, and     *
    //*  pull them out                                               *
    //*  --EXESTDOUT=  Sets STDOut                                   *
    //*  --EXESTDERR=  Sets STDErr                                   *
    //****************************************************************
    for iCmd := 0; iCmd < len(tmpSplit); iCmd++ {
        if strings.HasPrefix(strings.ToUpper(tmpSplit[iCmd]), "--EXESTDOUT=") {
            tmpSTDRedir = 1
            if strings.HasPrefix(strings.ToUpper(tmpSplit[iCmd][12:]), "CONS") {
                iSTDOut = 0
            } else {
                iSTDOut = 1
                STDOutF = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, tmpSplit[iCmd][12:])
            }
        } else if strings.HasPrefix(strings.ToUpper(tmpSplit[iCmd]), "--EXESTDERR=") {
            tmpSTDRedir = 1
            if strings.HasPrefix(strings.ToUpper(tmpSplit[iCmd][12:]), "CONS") {
                iSTDErr = 0
            } else {
                iSTDErr = 1
                STDErrF = fmt.Sprintf("%s%c%s", BaseDir, slashDelim, tmpSplit[iCmd][12:])
            }
        } else {
            cmdSplit = append(cmdSplit, tmpSplit[iCmd])
        }
    }

    //****************************************************************
    //* Make sure the command is there, then hash it.                *
    //****************************************************************
    if Commandtype < 3 {
        // EXE: or EXA: runs from AChoir Dirs

        // Are we requesting an explicit path?
        if len(cmdSplit[0]) > 0 {
            if cmdSplit[0][0] == slashDelim {
                Fullpath = fmt.Sprintf("%s%s", BaseDir, cmdSplit[0])
            } else {
                Fullpath = fmt.Sprintf("%s", cmdSplit[0])
            }

            if FileExists(Fullpath) {
                execMD5 = GetMD5File(Fullpath)
                ConsOut = fmt.Sprintf("[+] Command: %s\n    Command MD5: %s\n", Commandstring, execMD5)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("[!] Command could not be Located: %s\n", Commandstring)
                ConsLogSys(ConsOut, 1, 1)

                return fmt.Errorf("[!] Command Could not be Located: %s", Commandstring)
            }
        } else {
            ConsOut = fmt.Sprintf("[!] Blank command Entered.\n")
            ConsLogSys(ConsOut, 1, 1)

            return fmt.Errorf("[!] Blank Command Entered.\n")
        }
    } else {
        // SYS: Go Search the System Path
        RunPath, path_err := exec.LookPath(cmdSplit[0])
        Fullpath = fmt.Sprintf("%s", RunPath)

        if path_err != nil {
            ConsOut = fmt.Sprintf("[!] Command could not be Located: %s\n    %s\n", Commandstring, path_err)
            ConsLogSys(ConsOut, 1, 1)

            return path_err
        } else {
            execMD5 = GetMD5File(Fullpath)
            ConsOut = fmt.Sprintf("[+] Command: %s\n    Command MD5: %s\n", Commandstring, execMD5)
            ConsLogSys(ConsOut, 1, 1)
        }
    }


    //****************************************************************
    //* Copy the Command to the ACQ Collection for non-repudiation   *
    //****************************************************************
    TempCopyCommand := strings.Replace(Fullpath, slashDelimS, "-", -1)
    TempCopyCommand = strings.Replace(TempCopyCommand, ":", "", -1)
    FullCopyCommand = fmt.Sprintf("%s%cRanProg%c%s-%s", BACQDir, slashDelim, slashDelim, TempCopyCommand, execMD5)

    if !FileExists(FullCopyCommand) {
        ConsOut = fmt.Sprintf("[+] Saving Called Program As: %s\n", FullCopyCommand)
        ConsLogSys(ConsOut, 1, 1)
        _, _ = binCopy(Fullpath, FullCopyCommand)
    }


    //****************************************************************
    //* Setup the command to run                                     *
    //****************************************************************
    run_cmd := exec.Command(Fullpath, cmdSplit[1:]...)
    //run_cmd.Stdout = os.Stdout
    //run_cmd.Stderr = os.Stderr

    if iSTDOut == 0 {
        run_cmd.Stdout = os.Stdout
    } else {
        STDOHndl, _ := os.Create(STDOutF)
        run_cmd.Stdout = STDOHndl
        defer STDOHndl.Close()
    }

    if iSTDErr == 0 {
        run_cmd.Stderr = os.Stderr
    } else {
        STDEHndl, _ := os.Create(STDErrF)
        run_cmd.Stderr = STDEHndl
        defer STDEHndl.Close()
    }
        

    if Commandtype == 1 || Commandtype == 3 {
        // Reset STDout and STDErr if it was Temporary
        if tmpSTDRedir == 1 {
            iSTDOut = 0
            iSTDErr = 0
        }


        // Blocked (EXE: or SYS:)
        run_err := run_cmd.Run()
        if run_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Running Command: %s\n    %s\n", Commandstring, run_err)
            ConsLogSys(ConsOut, 1, 1)

            return run_err
        }

    } else {
        // Reset STDout and STDErr if it was Temporary
        if tmpSTDRedir == 1 {
            iSTDOut = 0
            iSTDErr = 0
        }

        // Non-Blocked/Asynchronous (EXA:)
        strt_err := run_cmd.Start()
        if strt_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Starting Command: %s\n    %s\n", Commandstring, strt_err)
            ConsLogSys(ConsOut, 1, 1)

            return strt_err
        }
    }

    return nil

}


//***************************************************************************
// Unzip will decompress a zip archive, moving all files and folders        *
// within the zip file (parameter 1) to an output directory (parameter 2).  *
//   https://golangcode.com/unzip-files-in-go/                              *
//***************************************************************************
func Unzip(ZipRdrFile []*zip.File, ZipDest string) ([]string, error) {

    var ZipFNames []string

    for _, ZipFile := range ZipRdrFile {

        // Store filename/path for returning and using later on
        ZipFpath := filepath.Join(ZipDest, ZipFile.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(ZipFpath, filepath.Clean(ZipDest)+string(os.PathSeparator)) {
            return ZipFNames, fmt.Errorf("%s: illegal file path", ZipFpath)
        }

        ZipFNames = append(ZipFNames, ZipFpath)

        if ZipFile.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(ZipFpath, os.ModePerm)
            continue
        }

        // Make File
        if zip_err := os.MkdirAll(filepath.Dir(ZipFpath), os.ModePerm); zip_err != nil {
            return ZipFNames, zip_err
        }

        ZipOutFile, zip_err := os.OpenFile(ZipFpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, ZipFile.Mode())
        if zip_err != nil {
            return ZipFNames, zip_err
        }

        Ziprc, zip_err := ZipFile.Open()
        if zip_err != nil {
            return ZipFNames, zip_err
        }

        _, zip_err = io.Copy(ZipOutFile, Ziprc)

        // Close the file without defer to close before next iteration of loop
        ZipOutFile.Close()
        Ziprc.Close()

        if zip_err != nil {
            return ZipFNames, zip_err
        }
    }
    return ZipFNames, nil
}


//***************************************************************************
// CopyParser: Parses out the Copy Parameters for Multi or Single Copy      *
//***************************************************************************
func CopyParser(splitString1 string, splitString2 string) {
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
        files_glob, glob_err := filepath.Glob(splitString1)

        if glob_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
            ConsLogSys(ConsOut, 1, 1)
            return
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

                if nBytes < 1 {
                    ConsOut = fmt.Sprintf("[!] File Copy was: %d Bytes\n", nBytes)
                    ConsLogSys(ConsOut, 1, 1)
                }
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

                files_glob, glob_err := filepath.Glob(TempDir)

                if glob_err != nil {
                    ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
                    ConsLogSys(ConsOut, 1, 1)
                    return
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

                        if nBytes < 1 {
                            ConsOut = fmt.Sprintf("[!] File Copy was: %d Bytes\n", nBytes)
                            ConsLogSys(ConsOut, 1, 1)
                        }
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

            if nBytes < 1 {
                ConsOut = fmt.Sprintf("[!] File Copy was: %d Bytes\n", nBytes)
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

                    if nBytes < 1 {
                        ConsOut = fmt.Sprintf("[!] File Copy was: %d Bytes\n", nBytes)
                        ConsLogSys(ConsOut, 1, 1)
                    }
                }
            }
        }
    }
}


func WalkCopyGlob(Walkfilepath string, WalkInfo os.FileInfo, Walk_err error) error {
    //****************************************************************
    //* Walk the filepath looking for DIRECTORIES ONLY. Then Glob    *
    //*  them with the wilcards...  This Approximates ** Globbing    *
    //****************************************************************
    file_stat, stat_err := os.Stat(Walkfilepath)

    if stat_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Identifying File: %s\n", stat_err)
        ConsLogSys(ConsOut, 1, 1)
        return stat_err
    }

    if file_stat.IsDir() {
        WalkfileFull := fmt.Sprintf("%s%c*%s", Walkfilepath, slashDelim, WalkfileWild)
        CopyParser(WalkfileFull, WalkfileToo)
    }

    return nil

}


func WalkForGlob(Walkfilepath string, WalkInfo os.FileInfo, Walk_err error) error {
    //****************************************************************
    //* Walk the filepath looking for DIRECTORIES ONLY. Then Glob    *
    //*  them with the wilcards...  This Approximates ** Globbing    *
    //****************************************************************
    file_stat, stat_err := os.Stat(Walkfilepath)

    if stat_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Identifying File: %s\n", stat_err)
        ConsLogSys(ConsOut, 1, 1)
        return stat_err
    }

    if file_stat.IsDir() {
        WalkfileFull := fmt.Sprintf("%s%c*%s", Walkfilepath, slashDelim, WalkfileWild)
        ForParser(WalkfileFull)
    }

    return nil

}


func WalkS3UpGlob(Walkfilepath string, WalkInfo os.FileInfo, Walk_err error) error {
    //****************************************************************
    //* Walk the filepath looking for DIRECTORIES ONLY. Then Glob    *
    //*  them with the wilcards...  This Approximates ** Globbing    *
    //****************************************************************
    file_stat, stat_err := os.Stat(Walkfilepath)

    if stat_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Identifying File: %s\n", stat_err)
        ConsLogSys(ConsOut, 1, 1)
        return stat_err
    }

    if file_stat.IsDir() {
        WalkfileFull := fmt.Sprintf("%s%c*%s", Walkfilepath, slashDelim, WalkfileWild)
        S3UpParser(WalkfileFull, WalkfileToo)
    }

    return nil

}


//***************************************************************************
// ForParser: Parses out the FOR: Parameters                                *
//***************************************************************************
func ForParser(ForDir string) {

    files_glob, glob_err := filepath.Glob(ForDir)

    if glob_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
        ConsLogSys(ConsOut, 1, 1)
        return
    }

    for _, file_found := range files_glob {
    //****************************************************************
    //* Ignore Directories - Only Process Files                      *
    //****************************************************************
        file_stat, _ := os.Stat(file_found)
        if file_stat.IsDir() {
            continue
        } else {
            ForHndl.WriteString(file_found + "\n")
        }
    }

    if (iNative == 0) {
        //****************************************************************
        //* Replace System32 with Sysnative if we are non-native         *
        //****************************************************************
        if CaseInsensitiveContains(TempDir, "\\System32\\") || CaseInsensitiveContains(TempDir, "/System32/") {

            repl_sys := NewCaseInsensitiveReplacer("System32", "sysnative")
            TempDir = repl_sys.Replace(TempDir)

            ConsOut = fmt.Sprintf("[*] Non-Native Flag Has Been Detected - Adding Sysnative Redirection: \n %s\n", TempDir)
            ConsLogSys(ConsOut, 1, 1)

            files_glob, glob_err := filepath.Glob(TempDir)

            if glob_err != nil {
                ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
                ConsLogSys(ConsOut, 1, 1)
                return
            }

            for _, file_found := range files_glob {
            //****************************************************************
            //* Ignore Directories - Only Process Files                      *
            //****************************************************************
                file_stat, _ := os.Stat(file_found)
                if file_stat.IsDir() {
                continue
                } else {
                    ForHndl.WriteString(file_found + "\n")
                }
            }
        }
    }
}


//***************************************************************************
// Upload a File to S3 - Assumes we are logged in and have a Session        *
//  https://paulbradley.org/gos3/                                           *
//  https://golangcode.com/uploading-a-file-to-s3/                          *
//***************************************************************************
func uploadFileToS3(S3Session *session.Session, S3FileName string, S3UpldName string) error {
    
    // Open the file
    S3UpFile, S3Up_err := os.Open(S3FileName)
    if S3Up_err != nil {
        return S3Up_err
    }
    defer S3UpFile.Close()

    // get file size and read the file content into a buffer
    S3FileInfo, _ := S3UpFile.Stat()
    var S3FileSize = S3FileInfo.Size()

    //S3Buffer := make([]byte, size)
    //S3UpFile.Read(S3Buffer)

    // Config settings: this is where you choose the bucket, filename, content-type and 
    //  storage class of the file being uploaded
    _, s3err := s3.New(S3Session).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET),
        Key:                  aws.String(S3UpldName),
        ACL:                  aws.String("private"),

        // Body:              bytes.NewReader(S3Buffer),
        // *os.File is an io.Reader
        Body:                 S3UpFile,

        ContentDisposition:   aws.String("attachment"),
        ContentLength:        aws.Int64(S3FileSize),
        //ContentType:          aws.String(http.DetectContentType(S3Buffer)),
        ServerSideEncryption: aws.String("AES256"),
        StorageClass:         aws.String("INTELLIGENT_TIERING"),
    })

    return s3err
}


//***************************************************************************
// S3UpParser: Parses out the Copy Parameters for Multi or Single Copy      *
//***************************************************************************
func S3UpParser(splitString1 string, splitString2 string) {
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
        files_glob, glob_err := filepath.Glob(splitString1)

        if glob_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
            ConsLogSys(ConsOut, 1, 1)
            return
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
            //* Upload to Output File Name                                    *
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


            //****************************************************************
            //* For S3 we will want to change all \ to a /                   *
            //****************************************************************
            MCpShard  = strings.Replace(MCpShard , "\\", "/", -1) 
            MCpPath  = strings.Replace(MCpPath , "\\", "/", -1) 
            MCprcO  = strings.Replace(MCprcO , "\\", "/", -1) 

            ConsOut = fmt.Sprintf("[+] S3 Multi-Upload File: %s\n    To: %s\n", file_found, MCprcO)
            ConsLogSys(ConsOut, 1, 1)

            S3_err = uploadFileToS3(S3_Session, file_found, MCprcO)
            if S3_err != nil {
                ConsOut = fmt.Sprintf("[!] Error Uploading File to S3: %s\n    %s\n", Inrec[4:], S3_err)
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

                files_glob, glob_err := filepath.Glob(TempDir)

                if glob_err != nil {
                    ConsOut = fmt.Sprintf("[!] Error Expanding WildCards: %s\n", glob_err)
                    ConsLogSys(ConsOut, 1, 1)
                    return
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
                    //* Upload to Output File Name                                   *
                    //****************************************************************
                    if setCPath == 0 || len(MCpShard) < 1 {
                        MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)
                    } else { 
                        MCprcO = fmt.Sprintf("%s%c%s%s", splitString2, slashDelim, MCpShard, MCpFName)
                    }


                    //****************************************************************
                    //* For S3 we will want to change all \ to a /                   *
                    //****************************************************************
                    MCpShard  = strings.Replace(MCpShard , "\\", "/", -1) 
                    MCpPath  = strings.Replace(MCpPath , "\\", "/", -1) 
                    MCprcO  = strings.Replace(MCprcO , "\\", "/", -1) 

                    ConsOut = fmt.Sprintf("[+] S3 Multi-Upload Redir File: %s\n    To: %s\n", file_found, MCprcO)
                    ConsLogSys(ConsOut, 1, 1)

                    S3_err = uploadFileToS3(S3_Session, file_found, MCprcO)
                    if S3_err != nil {
                        ConsOut = fmt.Sprintf("[!] Error Uploading File to S3: %s\n    %s\n", Inrec[4:], S3_err)
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
        //* Upload to Output File Name                                   *
        //****************************************************************
        MCprcO = fmt.Sprintf("%s%c%s", splitString2, slashDelim, MCpFName)


        //****************************************************************
        //* For S3 we will want to change all \ to a /                   *
        //****************************************************************
        MCprcO  = strings.Replace(MCprcO , "\\", "/", -1) 

        ConsOut = fmt.Sprintf("[+] S3 Singl-Upload File: %s\n    To: %s\n", splitString1, MCprcO)
        ConsLogSys(ConsOut, 1, 1)

        S3_err = uploadFileToS3(S3_Session, splitString1, MCprcO)
        if S3_err != nil {
            ConsOut = fmt.Sprintf("[!] Error Uploading File to S3: %s\n    %s\n", splitString1, S3_err)
            ConsLogSys(ConsOut, 1, 1)
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


                //****************************************************************
                //* For S3 we will want to change all \ to a /                   *
                //****************************************************************
                MCprcO  = strings.Replace(MCprcO , "\\", "/", -1) 

                ConsOut = fmt.Sprintf("[+] S3 Singl-Upload Redir File: %s\n    To: %s\n", splitString1, MCprcO)
                ConsLogSys(ConsOut, 1, 1)

                S3_err = uploadFileToS3(S3_Session, splitString1, MCprcO)
                if S3_err != nil {
                    ConsOut = fmt.Sprintf("[!] Error Uploading File to S3: %s\n    %s\n", splitString1, S3_err)
                    ConsLogSys(ConsOut, 1, 1)
                }

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
    block, _ := aes.NewCipher([]byte(createHash(passphrase)))
    gcm, enc_err := cipher.NewGCM(block)
    if enc_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Encrypting Data: %s\n", enc_err)
        ConsLogSys(ConsOut, 1, 1)
        return []byte(ConsOut)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, enc_err = io.ReadFull(rand.Reader, nonce); enc_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Encrypting Data: %s\n", enc_err)
        ConsLogSys(ConsOut, 1, 1)
        return []byte(ConsOut)
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext
}


//***************************************************************************
// Encryption Routines: Decrypt a stream of bytes                           *
//***************************************************************************
func decrypt(data []byte, passphrase string) []byte {
    key := []byte(createHash(passphrase))
    block, ciph_err := aes.NewCipher(key)
    if ciph_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data (Cipher): %s\n", ciph_err)
        ConsLogSys(ConsOut, 1, 1)
        return []byte(ConsOut)
    }

    gcm, gcm_err := cipher.NewGCM(block)
    if gcm_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data (GCM): %s\n", gcm_err)
        ConsLogSys(ConsOut, 1, 1)
        return []byte(ConsOut)
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, dec_err := gcm.Open(nil, nonce, ciphertext, nil)
    if dec_err != nil {
        ConsOut = fmt.Sprintf("[!] Error Decrypting Data: %s\n", dec_err)
        ConsLogSys(ConsOut, 1, 1)
        return []byte(ConsOut)
    }

    return plaintext
}


//***************************************************************************
// Encryption Routines: Encrypt a File                                      *
//***************************************************************************
func encryptFile(encFileName string, plainData []byte, passphrase string) {
    encFile, _ := os.Create(encFileName)
    defer encFile.Close()
    encFile.Write(encrypt(plainData, passphrase))
}


//***************************************************************************
// Encryption Routines: Decrypt a File                                      *
//***************************************************************************
func decryptFile(encFileName string, passphrase string) []byte {
    encData, _ := ioutil.ReadFile(encFileName)
    return decrypt(encData, passphrase)
}

