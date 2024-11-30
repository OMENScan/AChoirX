// ****************************************************************
// "As we benefit from the inventions of others, we should be
//  glad to share our own...freely and gladly."
//  - Benjamin Franklin
//
// AChCleanr = Clean Up AChoirX After Running
//
// AChCleanr v00.00.01 - Alpha 1

package main

import (
    "fmt"
    "time"
    "os"
    "strings"
    "strconv"
    "encoding/csv"
    "regexp"
    "runtime"
    "path/filepath"
    "bufio"
    "github.com/shirou/gopsutil/cpu")


// Global Variable Settings
var Version = "v00.00.01"                       // Version
var RunMode = "Run"                             // Character Runmode Flag (Build, Run, Menu)
var iRunMode = 0                                // Int Runmode Flag (0, 1, 2)
var Console_Status = 0                          // Was the Console Ever Invoked
var inFnam = "AChCleanr.ACQ"                    // Script Name
var MyProg = "none"                             // My Program Name and Path (os.Args[0])
var MyHash = "none"                             // My Hash
var inUser = "none"                             // UserId
var inPass = "none"                             // Password
var opArchit = "AMD64"                          // Architecture
var opSystem = "Windows"                        // Which Operating System are we running on
var iopSystem = 0                               // Operating System Flag (0=win, 1=lin, 2=osx, 3=?)
var OSVersion = "Windows 10.0.0"                // Operating System Version: OS Major.Minor.Build
var slashDelim byte = '\\'                      // Directory Delimiter Win vs. Lin vs. OSX
var slashDelimS string = "\\"                   // Same Thing but a String
var iNative = 0                                 // Are we Native 64Bit on 64Bit (Native = 1, NonNative = 0)
var sNative = ""                                // Native String (blank or Non-) 
var iIsAdmin = 0                                // Are we an Admin 
var sIsAdmin = ""                               // Are we an Admin String (blank or Non-) 
var yIsAdmin = "No"                             // Are we an Admin (Yes | No)
var LastRC = 0                                  // Last Return Code From External Executable
var RunMe = 0                                   // Used to Track Conditional Run Routines
var Looper = 0                                  // Flag to Keep Looping
var ForMe = 0                                   // Flag to identify &For is being used
var LstMe = 0                                   // Flag to identify &LST is being used
var ifFor = 0                                   // Flag contains FOR, FO0 - FOP
var iFor = 0                                    // Loop Counter FOR, FO0 - FOP
var iMaxCnt = 0                                 // Maximum Record Count (Set by FOR:, LST: &FOR, &LST)
var iLstCnt = 0                                 // LST Record Counter
var iForCnt = 0                                 // FOR Record Counter
var iLst = 0                                    // Loop Counter LST, LS0 - LSP
var ifLst = 0                                   // Flag contains LST, LS0 - LSP
var LoopNum = 0                                 // Loop Counter
var NotFound = 0                                // Flag for ensuring that only one Found Rec Increments RunMe
var YesFound = 0                                // Flag for ensuring that only one Found Rec Increments RunMe
var ForSlash = 0                                // Last Occurance of Slash to find File in Path
var ForFName = "File.txt"                       // Parsed File name from Path
var iFiltype = 0                                // Filter Type Include(1) or Exclude(2)
var FltRecFound = 0                             // Found a filter Record match
var iFiltscope = 0                              // Filter Scope Full Match(1) or Partial Match(2)
var consOrFile = 0                              // Console Input instead of File
var JmpLbl = "LBL:Top"                          // Working Jump Label Build String
var iSleep = 0                                  // Seconds to Sleep
var iOPNisOpen = 0                              // Is the User Defined File Open?
var FullDateTime = "01/01/0001 - 01:01:01"      // Date and Time

// Max CPU for Throttleing
var cpu_max float64 = 999

//Tokenize Records
var tokCount = 0                                // Token Counter
var Delims = ",\\/"                             // Tokenizing Delimiters (Lst, For(Win), For(Lin)
var iParseQuote = 0                             // CSV Quoted String Parsing (0=Strict, 1=Lazy) 

// Global File Names
var IniFile = "C:\\AChoir\\AChCleanr.Acq"       // AChoir Script File
var IncFile = "C:\\AChoir\\AChCleanr.Acq"       // AChoir Include (INC:) Script File
var ForFile = "C:\\AChoir\\Cache\\ForFiles"     // Do action for these Files
var LstFile = "C:\\AChoir\\Data.Lst"            // List of Data
var FltFile = "C:\\AChoir\\Cache\\FltFiles"     // Filter File for &LST and &FOR
var ChkFile = "C:\\AChoir\\Data.Chk"            // Check For File Existence
var CachDir = "C:\\AChoir\\Cache"               // AChoir Caching Directory 
var CurrDir = ""                                // Current Directory
var TempDir = ""                                // Temp Build Directory String
var CurrFil = "Current.fil"                     // Current File Name
var CurrWorkDir = "C:\\AChoir"                  // Current Workin Directory
var OpnHndl *os.File                            // User Defined Output File(s)

// Global File Handles
var ForScan *bufio.Scanner                      // IO Reader for ForFile
var LstScan *bufio.Scanner                      // IO Reader for LstFile
var FltScan *bufio.Scanner                      // IO Reader for FltFile
var ForHndl *os.File                            // File Handle for the ForFile
var LstHndl *os.File                            // File Handle for the LstFile
var FltHndl *os.File                            // File Handle for the FltFile
var ini_err error                               // Ini File Errors
var for_err error                               // For File Errors
var lst_err error                               // Lst File Errors
var opn_err error                               // User Defined File Errors
var flt_err error                               // Flt File Errors
var cwd_err error                               // Current Directory Errors
var for_rcd bool                                // Return Code for ForFile Read
var lst_rcd bool                                // Return Code for LstFile Read

// Arrays
var ForsArray = []string{"&FO0", "&FO1", "&FO2", "&FO3", "&FO4", "&FO5", "&FO6", "&FO7", "&FO8", 
                         "&FO9", "&FOA", "&FOB", "&FOC", "&FOD", "&FOE", "&FOF", "&FOG", "&FOH", 
                         "&FOI", "&FOJ", "&FOK", "&FOL", "&FOM", "&FON", "&FOO", "&FOP"}
var LstsArray = []string{"&LS0", "&LS1", "&LS2", "&LS3", "&LS4", "&LS5", "&LS6", "&LS7", "&LS8", 
                         "&LS9", "&LSA", "&LSB", "&LSC", "&LSD", "&LSE", "&LSF", "&LSG", "&LSH", 
                         "&LSI", "&LSJ", "&LSK", "&LSL", "&LSM", "&LSN", "&LSO", "&LSP"}

// Host Information
var cName = "localhost"                         // Endpoint Host Name
var cNewName = "localhost"                      // From /NAM: command line option
var host_err error                              // HostName Errors
var MyPID = 0                                   // This Programs Process ID

// Windows OS Variables
var WinRoot = "NA"                              // Windows Root Directory
var Procesr = "NA"                              // Processor
var TempVar = "NA"                              // Windows Temp Directory
var ProgVar = "NA"                              // Windows Program Files

//Tokenize Records
var iDblShard = 0                               // Double Glob Shard Pointer
var WalkfileWild = "Wildcard"                   // Wildcard Portion of the WalkFile Routines

// Input and Output Records
var Inrec = "File Input Record"                 // File Input Record
var Tmprec = "Formatted Console Record"         // Console Formatting Record
var Filrec = "File Record"                      // File Record 
var Lstrec = "List Record"                      // List Record 
var Fltrec = "Filter Record"                    // Filter Record 
var o32VarRec = "32 bit Variables"              // 32 Bit Variable Expansion Record
var Inprec = "Console Input"                    // Console Input Record
var Conrec = "Console Record"                   // Console Output Record

// INC: Nested Arrays
var HndlArry [10]*os.File
var ScanArry[10]*bufio.Scanner
var iIniCount = 0

// Cleaner Variables
var iClnExe = 0
var iClnDir = 0
var iClnAcn = 0
var ClnExe  = "AChoirX.exe"
var ClnDir  = "C:\\AChoir"
var ClnAcn  = "ACQ-IR-LocalHost-00000000-0000"  // AChoir Unique Collection Name
// Main Line
func main() {
    // Get Host Name
    cName, host_err = os.Hostname()
    if host_err != nil {
        cName = "LocalHost"
    }

    // Get My PID
    MyPID = os.Getpid()

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


    // Initial Settings and Configuration
    slashDelimS = fmt.Sprintf("%c", slashDelim)
    setDirectories(cName)
    inFnam = "AChCleanr.ACQ"


    // Start by Parsing any Command Line Parameters
    for i := 1; i < len(os.Args); i++ {
        if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/HELP") {
            fmt.Printf("\nAChCleanr ver: %s, Argument/Options:\n", Version)

            fmt.Printf(" /Help - This Description\n")
            fmt.Printf(" /EXE:<EXE Name> - The AChoirX EXE File to Delete\n")
            fmt.Printf(" /DIR:<Directory> - The AChoirX Root Directory to Clean\n")
            fmt.Printf(" /ACN:<Acquisition> - The AChoirX Acquisition/directory\n")
            fmt.Printf(" /INI:<File Name> - Run the <File Name> script instead of AChCleanr.ACQ\n\n")

            fmt.Printf(" Usage:\n")
            fmt.Printf(" AChCleanr Cleans up and deletes files from an AChoirX Run in order to DISOLVE the software\n\n")
            fmt.Printf(" AChCleanr /EXE:<AChoirExecutable.exe> /DIR:<AchoirRoot> /INI:<Alternate Script>\n")

            os.Exit(0)
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/INI:") {
            // Check if Input is Console
            if strings.HasPrefix(strings.ToUpper(os.Args[i]), "/INI:CONSOLE") {
                consOrFile = 1
                RunMode = "Con"
                Console_Status = 1
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
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/EXE:") {
            if len(os.Args[i]) < 254 {
                iClnExe = 1
                ClnExe  = os.Args[i][5:]
            } else {
                fmt.Println("[!] /EXE: Too Long - Greater than 254 chars")
            }
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/DIR:") {
            if len(os.Args[i]) < 254 {
                iClnDir = 1
                ClnDir  = os.Args[i][5:]
            } else {
                fmt.Println("[!] /DIR: Too Long - Greater than 254 chars")
            }
        } else if len(os.Args[i]) > 5 && strings.HasPrefix(strings.ToUpper(os.Args[i]), "/ACN:") {
            if len(os.Args[i]) < 254 {
                iClnAcn = 1
                ClnAcn  = os.Args[i][5:]
            } else {
                fmt.Println("[!] /DIR: Too Long - Greater than 254 chars")
            }
        }
    }


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
    // Set Initial File Names (ClnDir needs to be set 1st)           *
    //****************************************************************
    IniFile = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, inFnam)
    CachDir = fmt.Sprintf("%s%cCache", ClnDir, slashDelim)
    ForFile = fmt.Sprintf("%s%cForFiles", CachDir, slashDelim)
    LstFile = fmt.Sprintf("%s%cLstFiles", CachDir, slashDelim)


    //****************************************************************
    // Console or File Mode?                                         *
    //****************************************************************
    if consOrFile == 1 {
        fmt.Printf("[+] Switching to Console Input.\n")
        fmt.Printf(">>>")
        ScanArry[iIniCount] = bufio.NewScanner(os.Stdin)
    } else {
        HndlArry[iIniCount], ini_err = os.Open(IniFile)
        if ini_err != nil {
            fmt.Printf("[!] Error Opening Ini File: %s\n", IniFile)
        }

        ScanArry[iIniCount] = bufio.NewScanner(HndlArry[iIniCount])
    }

    // Do This Forever
    for {
        // Input Scan (File and Console) until Error (EOF)
        if !ScanArry[iIniCount].Scan() {
            //****************************************************************
            // See if this was a child INI (INC:). If so, return to parent   *
            //****************************************************************
            if iIniCount > 0 {
                fmt.Printf("[+] End of INC: Input...  Returning...\n")
                HndlArry[iIniCount].Close()
                iIniCount--
                continue
            } else {
               fmt.Printf("[+] End of Input...  Exiting...\n")
                cleanUp_Exit(LastRC);
                os.Exit(LastRC);
                break
            }
        }

        //Remove any preceding blanks
        Tmprec = strings.TrimSpace(ScanArry[iIniCount].Text())
  
        // Dont Process any Comments
        if strings.HasPrefix(Tmprec, "*") {
            continue
        }


        //****************************************************************
        //* Conditional Execution                                        *
        //****************************************************************
        if RunMe > 0 {
            if strings.HasPrefix(strings.ToUpper(Tmprec), "CKY:") {
                RunMe++
            } else if strings.HasPrefix(strings.ToUpper(Tmprec), "CKN:") {
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
                    fmt.Printf("[!] &FOR Directory has not been set with the FOR: command.  Ignoring &FOR Loop...\n")
                    Looper = 0
                }

                iMaxCnt = 0
                iForCnt = 0
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
                    fmt.Printf("[!] &LST File was not found (LST: not set): %s\n", LstFile)
                    Looper = 0
                }

                iMaxCnt = 0
                iLstCnt = 0
                LstScan = bufio.NewScanner(LstHndl)
            } else {
                LstMe = 0
            }


            //****************************************************************
            //* Loop (FOR: and LST:) until Looper = 1                        *
            //****************************************************************
            LoopNum = 0
            NotFound = 0
            YesFound = 0

            for Looper > 0 {
                if ForMe == 0 && LstMe == 0 {
                    Looper = 0
                } else if ForMe == 1 && LstMe == 0 {
                    for_rcd = ForScan.Scan()
                    for_err = ForScan.Err()

                    // No Error and no EOF - So Process the Record
                    if for_err == nil && for_rcd == true {
                        Filrec = strings.TrimSpace(ForScan.Text())

                        Looper = 1
                        LoopNum++
                        iMaxCnt++
                        iForCnt++

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
                        fmt.Printf("[+] For Records Processed: %d\n", LoopNum)
                        
                        break
                    }
                } else if ForMe == 0 && LstMe == 1 {
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
                        iLstCnt++
                    } else {
                        fmt.Printf("[+] Lst Records Processed: %d\n", LoopNum)
                        break
                    }
                } else {
                    Looper = 0
                    fmt.Printf("[!] AChoirX does not yet support Nested Looping (&LST + &FOR)\n     > %s\n", Tmprec)
                    Tmprec = fmt.Sprintf("***: Command Bypassed")
                }


                //****************************************************************
                //* If Filtering has been confiured, the same Filtering File     *
                //*  will be used for both &FOR and &LST                         *
                //****************************************************************
                if iFiltype > 0 {
                    if LstMe == 1 || ForMe == 1 {
                        //****************************************************************
                        //* Check for &LST (LstMe) or &FOR (ForMe) filtering             *
                        //****************************************************************
                        FltHndl, flt_err := os.Open(FltFile)

                        if flt_err != nil {
                            iFiltype = 0
                            fmt.Printf("[!] Filter file could not be opened: %s\n", FltFile)
                        } else {
                            FltScan = bufio.NewScanner(FltHndl)

                            FltRecFound = 0 
                            for FltScan.Scan() {
                                Fltrec = strings.TrimSpace(FltScan.Text())

                                if LstMe == 1 {
                                    //****************************************************************
                                    //* &LST Look for Match - Filtscope is Full(1) or Part(2)        *
                                    //****************************************************************
                                    //fmt.Printf("[c] -%s-%s-\n", Lstrec, Fltrec)
                                    if iFiltscope == 1 {
                                        if strings.ToUpper(Lstrec) == strings.ToUpper(Fltrec) {
                                            FltRecFound = 1
                                        } 
                                     } else {
                                        if CaseInsensitiveContains(Lstrec, Fltrec) {
                                            FltRecFound = 1
                                        }
                                    }
                                } else if ForMe == 1 {
                                    //****************************************************************
                                    //* &For Look for Match - Filtscope is Full(1) or Part(2)        *
                                    //****************************************************************
                                    if iFiltscope == 1 {
                                        if strings.ToUpper(Filrec) == strings.ToUpper(Fltrec) {
                                            FltRecFound = 1
                                        }
                                    } else {
                                        if CaseInsensitiveContains(Filrec, Fltrec) {
                                            FltRecFound = 1
                                        }
                                    }
                                }
                            }
                            FltHndl.Close()


                            //****************************************************************
                            //* Include Filter Logic(1) - If not found, bypass               *
                            //****************************************************************
                            if iFiltype == 1 && FltRecFound == 0 {
                                continue
                            }

                            //****************************************************************
                            //* Exclude Filter Logic(2) - if found, bypass                   *
                            //****************************************************************
                            if iFiltype == 2 && FltRecFound == 1 {
                                continue
                            }
                        }
                    }
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
                        TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, CurrDir)
                    } else {
                        TempDir = ClnDir
                    }

                    repl_Dir := NewCaseInsensitiveReplacer("&Dir", TempDir)
                    o32VarRec = repl_Dir.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Exe") {

                    repl_Exe := NewCaseInsensitiveReplacer("&Exe", ClnExe)
                    o32VarRec = repl_Exe.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Fil") {

                    repl_Fil := NewCaseInsensitiveReplacer("&Fil", CurrFil)
                    o32VarRec = repl_Fil.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Acn") {

                    repl_Acn := NewCaseInsensitiveReplacer("&Acn", ClnAcn)
                    o32VarRec = repl_Acn.Replace(o32VarRec)
                }

                if CaseInsensitiveContains(o32VarRec, "&Win") {

                    repl_Win := NewCaseInsensitiveReplacer("&Win", WinRoot)
                    o32VarRec = repl_Win.Replace(o32VarRec)
                }


                // First process &FOR - This should prevent unnecessary parsing
                if CaseInsensitiveContains(o32VarRec, "&For") {
                    repl_For := NewCaseInsensitiveReplacer("&For", Filrec)
                    o32VarRec = repl_For.Replace(o32VarRec)
                }


                // If any &FO_ are left, process those 
               if CaseInsensitiveContains(o32VarRec, "&Fo") {

                    // Always allow &For - Even if Parsing Fails
                    if CaseInsensitiveContains(o32VarRec, "&For") {
                        repl_For := NewCaseInsensitiveReplacer("&For", Filrec)
                        o32VarRec = repl_For.Replace(o32VarRec)
                    }

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

                    if iParseQuote == 1 {
                        // Does not enforce strict matching quotes when parsing
                        tokRdr.LazyQuotes = true
                    }

                    tokFields, tok_err := tokRdr.Read()

                    if tok_err != nil {
                        fmt.Printf("[!] Parsing Error for Record(%d): %s\n", LoopNum, tok_err)
                        
                        repl_For := NewCaseInsensitiveReplacer("&Fo", "Parse_Err_")
                        o32VarRec = repl_For.Replace(o32VarRec)
                    } else {
                        tokCount = len(tokFields)
                        if tokCount < 25 {
                            for i := tokCount; i < 26; i++ {
                                tokFields = append(tokFields, "")
                            }
                        }


                        // Look for Replacements &Fo0 - FoP
                        for iFor = 0; iFor < 26; iFor++ {
                            if CaseInsensitiveContains(o32VarRec, ForsArray[iFor]) {
                                repl_For := NewCaseInsensitiveReplacer(ForsArray[iFor], tokFields[iFor])
                                o32VarRec = repl_For.Replace(o32VarRec)
                            }
                        }
                    }
                }


                //Clear out the &LST - This should prevent unnecessary parsing
                if CaseInsensitiveContains(o32VarRec, "&Lst") {
                    repl_Lst := NewCaseInsensitiveReplacer("&Lst", Lstrec)
                    o32VarRec = repl_Lst.Replace(o32VarRec)
                }

                // If any &LS_ are left, process those 
                if CaseInsensitiveContains(o32VarRec, "&Ls") {

                    // Always allow &LST - Even if Parsing Fails
                    if CaseInsensitiveContains(o32VarRec, "&Lst") {
                        repl_Lst := NewCaseInsensitiveReplacer("&Lst", Lstrec)
                        o32VarRec = repl_Lst.Replace(o32VarRec)
                    }

                    // Split string, we will likely need it split 
                    runeDelims := []rune(Delims)
                    tokRdr := csv.NewReader(strings.NewReader(Lstrec))
                    tokRdr.Comma = runeDelims[0]
                    tokRdr.FieldsPerRecord = -1
                    tokRdr.TrimLeadingSpace = true

                    if iParseQuote == 1 {
                        // Does not enforce strict matching quotes when parsing
                        tokRdr.LazyQuotes = true
                    }

                    tokFields, tok_err := tokRdr.Read()

                    if tok_err != nil {
                        fmt.Printf("[!] Parsing Error for Record(%d): %s\n", LoopNum, tok_err)
                        
                        repl_Lst := NewCaseInsensitiveReplacer("&Ls", "Parse_Err_")
                        o32VarRec = repl_Lst.Replace(o32VarRec)
                    } else {
                        tokCount = len(tokFields)
                        if tokCount < 25 {
                            for i := tokCount; i < 26; i++ {
                                tokFields = append(tokFields, "")
                            }
                        }

                        // Look for Replacements &Ls0 - LsP
                        for iLst = 0; iLst < 26; iLst++ {
                            if CaseInsensitiveContains(o32VarRec, LstsArray[iLst]) {
                                repl_Lst := NewCaseInsensitiveReplacer(LstsArray[iLst], tokFields[iLst])
                                o32VarRec = repl_Lst.Replace(o32VarRec)
                            }
                        }
                    }
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

                if CaseInsensitiveContains(o32VarRec, "&Chk") {

                    repl_Chk := NewCaseInsensitiveReplacer("&Chk", ChkFile)
                    o32VarRec = repl_Chk.Replace(o32VarRec)
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
                        fmt.Printf("[*] Jumping Does not make sense in Interactive Mode.  Ignoring...\n")
                        
                    } else {
                        // Rewind File and Jump to a Label (LBL:)
                        // Note: For This to work we have to Reset both the Handle and the Scanner!
                        RunMe = 0
                        HndlArry[iIniCount].Seek(0, 0)
                        ScanArry[iIniCount] = bufio.NewScanner(HndlArry[iIniCount])

                        JmpLbl = fmt.Sprintf("LBL:%s", Inrec[4:])

                        for ScanArry[iIniCount].Scan() {
                            Tmprec = strings.TrimSpace(ScanArry[iIniCount].Text())
                            if Tmprec == JmpLbl {
                                break
                            }
                        }
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
                            TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, CurrDir)
                        } else  {
                            CurrDir = ""
                            TempDir = ClnDir
                        }
                    } else {
                        if len(Inrec) > 4 {
                            //Check to see if it is an append or new &Dir
                            //Dont add // if it's new!
                            if len(CurrDir) > 0 {
                                CurrDir += slashDelimS
                            }

                            CurrDir += Inrec[4:]
                            TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, CurrDir)
                        }
                    }

                    // Have we created this Directory already?
                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        fmt.Printf("[+] Creating Sub-Directory: %s\n", CurrDir)

                        ExpandDirs(TempDir)

                    }

                    // Reset The WorkingDirectory to the new Directory
                    CurrWorkDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, CurrDir)
                    os.Chdir(CurrWorkDir)

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FIL:") {
                    CurrFil = Inrec[4:]
                    TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, CurrDir)

                    if _, CurrDir_err := os.Stat(TempDir); os.IsNotExist(CurrDir_err) {
                        fmt.Printf("[+] Creating Sub-Directory: %s\n", CurrDir)

                        ExpandDirs(TempDir)
                    }

                    fmt.Printf("[+] File has been Set to: %s\n", CurrFil)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "INI:") {
                    IniFile = Inrec[4:]

                    if strings.ToUpper(IniFile) == "CONSOLE" {
                        if consOrFile == 0 {
                            RunMode = "Con"
                            inFnam = "Console"
                            Console_Status = 1
                            iRunMode = 1
                            consOrFile = 1

                            fmt.Printf("[+] Switching to Console (Interactive Mode)\n")

                            HndlArry[iIniCount].Close()
                            ScanArry[iIniCount] = bufio.NewScanner(os.Stdin)
                        }
                    } else {
                        if FileExists(IniFile) {
                            fmt.Printf("[!] Switching to INI File: %s\n", Inrec[4:])

                            // Only close the handle if its not Console. If it is Console Set it back to File
                            if consOrFile == 0 {
                                HndlArry[iIniCount].Close()
                            } else {
                                consOrFile = 0
                            }

                            HndlArry[iIniCount], ini_err = os.Open(IniFile)
                            if ini_err != nil {
                                fmt.Printf("[!] Error Opening Ini File: %s - Exiting.\n", IniFile)

                                cleanUp_Exit(3)
                                os.Exit(3)
                            }

                            RunMode = "Ini"
                            ScanArry[iIniCount] = bufio.NewScanner(HndlArry[iIniCount])
                            RunMe = 0  // Conditional run Script default is yes

                        } else {
                            fmt.Printf("[!] Requested INI File Not Found: %s - Ignored.\n", Inrec[4:])
                            
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "INC:") {
                    // Include an INI file - All variables, Labels, Etc. Remain the same
                    // It is essentially an injected INI file that is inline
                    IncFile = Inrec[4:]

                    if iIniCount > 8 {
                        fmt.Printf("[!] Maximum Nested INC: Files has been exceeded.  Ignoring: %s...\n", IncFile)
                        
                    } else {
                        if FileExists(IncFile) {
                            fmt.Printf("[+] Including INC File: %s\n", Inrec[4:])
                            

                            // Programming note: Here we DO NOT Close the parent file. But bump the array 
                            // to have an additonal open file
                            iIniCount++

                            HndlArry[iIniCount], ini_err = os.Open(IncFile)
                            if ini_err != nil {
                                fmt.Printf("[!] Error Opening Include File: %s - Ignoring.\n", IncFile)
                                
                                iIniCount--
                            } else {
                                RunMode = "Ini"
                                ScanArry[iIniCount] = bufio.NewScanner(HndlArry[iIniCount])
                                RunMe = 0  // Conditional run Script default is yes
                            }
                        } else {
                            fmt.Printf("[!] Requested INC File Not Found: %s - Ignored.\n", Inrec[4:])
                            
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ADM:CHECK") {
                    fmt.Printf("[+] Running as %sAdmin/Root\n", sIsAdmin)
                    

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ADM:FORCE") {
                    if iIsAdmin == 1 {
                        fmt.Printf("[+] Running as Admin/Root - Continuing...\n")
                        
                    } else {
                        fmt.Printf("[+] NOT Running as Admin/Root - Exiting...\n")
                        

                        cleanUp_Exit(3)
                        os.Exit(3)
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:HIDE") {
                    fmt.Printf("[+] Hiding Console...\n")
                    
                    winConHideShow(0)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CON:SHOW") {
                    fmt.Printf("[+] Showing Console...\n")
                    
                    winConHideShow(1)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SLP:") {
                    iSleep, _ = strconv.Atoi(Inrec[4:])
                    time.Sleep (time.Duration(iSleep) * time.Second)
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "INP:") {
                    consInput(Inrec[4:], 1, 0)
                    Inprec = Conrec
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "USR:") {
                    inUser = Inrec[4:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "PWD:") {
                    inPass = Inrec[4:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "VER:") {
                    //****************************************************************
                    //* Check the Input String for Version.  This can be Partial.    *
                    //*  For Example: Windows 10.0.18362 will get a TRUE for         *
                    //*  Ver:Win, Ver:Windows 10, Ver:Windows 10.0.183, etc...       *
                    //****************************************************************
                    if consOrFile == 1 {
                        fmt.Printf("[*] OS Version Detected: %s\n", OSVersion)
                        
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
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CKY:") {
                    ChkFile = strings.Trim(Inrec[4:], "\"")

                    if consOrFile == 1 {
                        if !FileExists(ChkFile) {
                            fmt.Printf("[*] File Does Not Exist: %s\n", ChkFile)
                            
                        } else {
                            fmt.Printf("[*] File Exists: %s\n", ChkFile)
                            
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
                            fmt.Printf("[*] File Does (NOT NOT) Exist: %s\n", ChkFile)
                            
                        } else {
                            fmt.Printf("[*] File Does Not Exist: %s\n", ChkFile)
                            
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
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SAY:") {
                    fmt.Printf("%s\n", Inrec[4:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "ECHO ") {
                    fmt.Printf("%s\n", Inrec[5:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "END:") {
                    if len(Inrec) > 4 {
                        // Is there a parameter after the END: Statement
                        if strings.ToUpper(Inrec[4:]) == "RESET" {
                            // The Parameter is RESET
                            fmt.Printf("[+] Resetting Internal Conditional Execution Flag...\n")
                            
                            RunMe = 0
                        }
                    } else if RunMe > 0 {
                        RunMe--
                    } else if RunMe < 0 {
                        //* Something went wrong and our logic created a negative RunMe - Reset to 0
                        fmt.Printf("[!] Internal Error, Resetting Internal Conditional Execution Flag...\n")
                        
                        RunMe = 0
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "OPN:") {
                    // If we already had a file open, close it now.
                    if iOPNisOpen == 1 {
                        fmt.Printf("[*] Previously Opened File has been Closed.\n")
                        

                        OpnHndl.Close()
                    }

                    OpnHndl, opn_err = os.OpenFile(Inrec[4:], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                    if opn_err != nil {
                        fmt.Printf("[!] File Could not be opened for Append:\n    %s\n", Inrec[4:])
                        

                        iOPNisOpen = 0
                    } else {
                        fmt.Printf("[+] File Opened for Append:\n    %s\n", Inrec[4:])
                        

                        iOPNisOpen = 1
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "OUT:") {
                    if iOPNisOpen == 1 {
                        OpnHndl.WriteString(Inrec[4:]+"\n")
                    } else {
                        fmt.Printf("[!] No File OPN:(ed) for Append:\n")
                        
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "PZZ:") {
                    fmt.Printf("%s\n", Inrec[4:])
                    consInput(Inrec[4:], 1, 0)
                    Inprec = Conrec

                    if len(Inprec) > 0 {
                        if Inprec[0] == 'q' || Inprec[0] == 'Q' {
                            fmt.Printf("[+] You Have Requested AChoirX to Quit.\n")
                            

                            cleanUp_Exit(0)
                            os.Exit(0)
                        }
                    }
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FOR:") {
                    ExpandDirs(CachDir)
                    ForFile = fmt.Sprintf("%s%cForFiles", CachDir, slashDelim)
                    TempDir = Inrec[4:]

                    //*****************************************************************
                    //* If we are using FOR: with a ListFile (FOR:&LST), Append the   *
                    //*  data in ForFile - unless it's our first &LST, then open      *
                    //*  New/Truncate -  Otherwise always open New/Truncate           *
                    //*****************************************************************
                    if (LstMe == 1) && (iLstCnt > 1) {
                      ForHndl, for_err = os.OpenFile(ForFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
                    } else {
                      ForHndl, for_err = os.OpenFile(ForFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
                    }

                    if for_err != nil {
                        fmt.Printf("[!] System File Tracking File Could not be opened.\n")
                        
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

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "DEL:") {
                    TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, Inrec[4:])

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

                            BasicDel := fmt.Sprintf("%s%s", WalkDir, WalkfileWild)
                            DelParser(BasicDel)

                            filepath.Walk(WalkDir, WalkDelGlob)
                        }
                    } else {
                        DelParser(TempDir)
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "CLN:") {

                    TempDir = fmt.Sprintf("%s%c%s", ClnDir, slashDelim, Inrec[4:])

                    remov_err := os.RemoveAll(TempDir)
                    if remov_err != nil {
                        fmt.Printf("[!] Error Cleaning Directory: %s\n", remov_err)
                        
                    }

                } else if strings.HasPrefix(strings.ToUpper(Inrec), "LST:") {
                    LstFile = fmt.Sprintf("%s%c%s", ClnDir,slashDelim, Inrec[4:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "FLT:") {
                    FltFile = fmt.Sprintf("%s%c%s", ClnDir,slashDelim, Inrec[4:])
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "BYE:") {
                    cleanUp_Exit(LastRC);
                    os.Exit(LastRC);
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:DELIMS=") && len(Inrec) > 13 {
                    Delims = Inrec[11:]
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:PARSEQUOTE=STRICT") {
                    iParseQuote=0
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:PARSEQUOTE=LAZY") {
                    iParseQuote=1
                } else if strings.HasPrefix(strings.ToUpper(Inrec), "SET:FILTER=") {
                    //*****************************************************************
                    //* Set Filter to None (0), Include(1), or Exclude(2)             *
                    //*****************************************************************
                    if strings.ToUpper(Inrec[11:]) == "NONE" {
                        iFiltype = 0
                    } else if CaseInsensitiveContains(Inrec[11:], "Incl") {
                        iFiltype = 1
                    } else if CaseInsensitiveContains(Inrec[11:], "Excl") {
                        iFiltype = 2
                    }

                    //*****************************************************************
                    //* Set Filter to Full(1) or Partial(2) Matching                  *
                    //*****************************************************************
                    if CaseInsensitiveContains(Inrec[11:], "Full") {
                        iFiltscope = 1
                    } else if CaseInsensitiveContains(Inrec[11:], "Part") {
                        iFiltscope = 2
                    }

                    //*****************************************************************
                    //* Verify that the Filter Exists, Otherwise set to NONE(0)       *
                    //*****************************************************************
                    FltHndl, flt_err = os.Open(FltFile)

                    if flt_err != nil {
                        fmt.Printf("[!] &FLT File was not found (FLT: not set): %s\n", FltFile)
                        
                        iFiltype = 0
                        break
                    }

                    FltHndl.Close()

                }
            }
        }

        // Show >>> If we are in Console Mode
        if consOrFile == 1 {
            fmt.Printf(">>>")
        }
    }


    //****************************************************************
    //* Cleanup                                                      *
    //****************************************************************
    if RunMe > 0 {
        fmt.Printf("[!] You have and extra END: Hanging! Check your Logic.\n")
        
    }

    cleanUp_Exit(0);
    os.Exit(0);

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
    return cir.toReplace.ReplaceAllLiteralString(str, cir.replaceWith)
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
        fmt.Printf("[!] Parsing Error for: %s\n", SpString)
        
        return SpString, "", 1
    }

    tokCount = len(tokFields)

    if tokCount < 2 {
        fmt.Printf("[!] No Separating Delimiters Found for: %s\n", SpString)
        
        return SpString, "", 1
    } else if tokCount > 2 {
        fmt.Printf("[!] Extra Data Ignored. Tokens: %d\n", tokCount)
        
        return tokFields[0], tokFields[1], 3
    } else {
        return tokFields[0], tokFields[1], 2
    } 
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


func cleanUp_Exit(exitRC int) {
    //****************************************************************
    //* Cleanup and Exit                                             *
    //****************************************************************
    if FileExists(ForFile) {
        os.Remove(ForFile)
    }


    //****************************************************************
    //* All Done with Cleaning                       `               *
    //****************************************************************
    showTime("Cleaning Completed");

    if consOrFile == 0 {
        HndlArry[iIniCount].Close()
        DelParser(IniFile)
    }
}


func WalkForGlob(Walkfilepath string, WalkInfo os.FileInfo, Walk_err error) error {
    //****************************************************************
    //* Walk the filepath looking for DIRECTORIES ONLY. Then Glob    *
    //*  them with the wilcards...  This Approximates ** Globbing    *
    //****************************************************************
    file_stat, stat_err := os.Stat(Walkfilepath)

    if stat_err != nil {
        fmt.Printf("[!] Error Identifying File: %s\n", stat_err)
        
        return stat_err
    }

    if file_stat.IsDir() {
        WalkfileFull := fmt.Sprintf("%s%c*%s", Walkfilepath, slashDelim, WalkfileWild)
        ForParser(WalkfileFull)
    }

    return nil

}


func WalkDelGlob(Walkfilepath string, WalkInfo os.FileInfo, Walk_err error) error {
    //****************************************************************
    //* Walk the filepath looking for DIRECTORIES ONLY. Then Glob    *
    //*  them with the wilcards...  This Approximates ** Globbing    *
    //****************************************************************
    file_stat, stat_err := os.Stat(Walkfilepath)

    if stat_err != nil {
        fmt.Printf("[!] Error Identifying File: %s\n", stat_err)
        
        return stat_err
    }

    if file_stat.IsDir() {
        WalkfileFull := fmt.Sprintf("%s%c*%s", Walkfilepath, slashDelim, WalkfileWild)
        DelParser(WalkfileFull)
    }

    return nil

}


//***************************************************************************
// ForParser: Parses out the FOR: Parameters                                *
//***************************************************************************
func ForParser(ForDir string) {

    files_glob, glob_err := filepath.Glob(ForDir)

    if glob_err != nil {
        fmt.Printf("[!] Error Expanding WildCards: %s\n", glob_err)
        
        return
    }

    for _, file_found := range files_glob {
    //****************************************************************
    //* Ignore Directories - Only Process Files                      *
    //****************************************************************
        file_stat, fstat_err := os.Stat(file_found)

        // v10.00.53 - Add Check for Deleted Files
        if fstat_err != nil {
            fmt.Printf("[!] File Error: %s\n", fstat_err)
            
            continue
        }

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

            fmt.Printf("[*] Non-Native Flag Has Been Detected - Adding Sysnative Redirection: \n %s\n", TempDir)
            

            files_glob, glob_err := filepath.Glob(TempDir)

            if glob_err != nil {
                fmt.Printf("[!] Error Expanding WildCards: %s\n", glob_err)
                
                return
            }

            for _, file_found := range files_glob {
            //****************************************************************
            //* Ignore Directories - Only Process Files                      *
            //****************************************************************
                file_stat, fstat_err := os.Stat(file_found)

                // v10.00.53 - Add Check for Deleted Files
                if fstat_err != nil {
                    fmt.Printf("[!] File Error: %s\n", fstat_err)
                    
                    continue
                }

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
// DelParser: Parses out the DEL: Parameters                                *
//***************************************************************************
func DelParser(DelDir string) {

    files_glob, glob_err := filepath.Glob(DelDir)

    if glob_err != nil {
        fmt.Printf("[!] Error Expanding WildCards: %s\n", glob_err)
        
        return
    }

    for _, file_found := range files_glob {
    //****************************************************************
    //* Ignore Directories - Only Process Files                      *
    //****************************************************************
        file_stat, fstat_err := os.Stat(file_found)

        // v10.00.53 - Add Check for Deleted Files
        if fstat_err != nil {
            fmt.Printf("[!] File Error: %s\n", fstat_err)
            
            continue
        }

        if file_stat.IsDir() {
            continue
        } else {
            os.Remove(file_found)
        }
    }
}


//***************************************************************************
// Check for CPU Throttleing. cpu_max of 999 (default) is bypass checking   *
//***************************************************************************
func cpuThrotl() {
    if cpu_max != 999 {
        for pctloop := 0; pctloop < 10; pctloop++ {
            //Percent calculates the percentage of cpu used either per CPU or combined.
            cpu_percent, _ := cpu.Percent(time.Second,false)
            //fmt.Printf("[+] Total CPU Utilization: %.2f\n", cpu_percent[0])
            fmt.Printf("[+] Total CPU Utilization: %.2f\n", cpu_percent[0])
            
            if cpu_percent[0] > cpu_max {
                fmt.Printf("[!] CPU Throttling Invoked...\n")
                
            } else {
                pctloop = 10
            }
        }
    }
}

//***************************************************************
// Console Input:                                               *
// conLog  - Should we Log This?                                *
// conHide - Should we redact the Input                         *
//***************************************************************
func consInput(consString string, conLog int, conHide int) {
    for {
        fmt.Sprintf("[?] [%s] ", consString)

        con_reader := bufio.NewReader(os.Stdin)
        Conrec, _ = con_reader.ReadString('\n')
    }

    // Remove CRLF to LF
    Conrec = strings.TrimSpace(Conrec)

    if conHide == 1 {
        fmt.Printf("[?] *Redacted*\n")
    } else {
        fmt.Printf("[?] %s\n", Conrec)
    }
}


//***************************************************************************
// ReSet a new Cleaning Directory                                           *
//***************************************************************************
func setDirectories(NewDirName string) {

    // Get Time and Date
    lclTime := time.Now()
    iMonth := int(lclTime.Month())
    iDay := lclTime.Day()
    iYYYY := lclTime.Year()
    iHour := lclTime.Hour()
    iMin := lclTime.Minute()

    // Set Cleaning Directory Name
    ClnAcn = fmt.Sprintf("ACQ-IR-%s-%04d%02d%02d-%02d%02d", NewDirName, iYYYY, iMonth, iDay, iHour, iMin)

    // What Directory are we in? Set our Base and Current Working Directory
    ClnDir, cwd_err = os.Getwd()
    if cwd_err != nil {
        ClnDir = fmt.Sprintf("%cAChoirX%cACQ-IR-%s-%04d%02d%02d-%02d%02d", slashDelim, slashDelim, NewDirName, iYYYY, iMonth, iDay, iHour, iMin) 
    }

    CurrWorkDir = ClnDir

    // Remove any Trailing Slashes.  This happens if CWD is a mapped network drive (since it is at the root directory)
    // Note: slashDelim was set above based on OS
    if last := len(ClnDir) - 1; last >= 0 && ClnDir[last] == slashDelim {
        ClnDir = ClnDir[:last]
    }

    ClnExe = fmt.Sprintf("AChoirX.exe")

}

