// ****************************************************************
// * Only Include this file for windows version                   *
// * - Uses the Embed Directive - Requires GoLang 1.16            *
// ****************************************************************

package main

import (
    "fmt"
    "regexp"
    "io"
    "errors"
    _ "embed"
    "syscall"
    "unsafe"
    "os"
    "time"
    "strings"
    "bytes"
    "archive/zip"
    "golang.org/x/sys/windows"
    "golang.org/x/sys/windows/registry"
    "github.com/gonutz/w32/v2"
    ntfs "www.velocidex.com/golang/go-ntfs/parser"
)


// ****************************************************************
// * NTFS Raw Copy Constants                                      *
// ****************************************************************
const (
    NTFSAttrType_Data        = 128
    NTFSAttrID_Normal        = 0
    NTFSAttrID_ADS           = 5
    NTFSWildcard_Stream_Name = ":*:"
    NTFSWildcard_Stream_ID   = uint16(0xffff)
)


// ****************************************************************
// * NTFS Raw Copy Variables                                      *
// ****************************************************************
var (
    ErrReturnedNil             = errors.New("[!] Result returned nil reference")
    ErrInvalidInput            = errors.New("[!] Invalid input")
    ErrDeviceInaccessible      = errors.New("[!] Raw device is not accessible")
    ErrPrivilegedAccountWanted = errors.New("[!] Require privileged token, please UAC elevate")
)


// ****************************************************************
// * Embed Zip File in embdata byte array                         *
// ****************************************************************
//go:embed WinEmbed.zip
var embdata []byte


// ****************************************************************
// * List All Windows Drive Letters                               *
// ****************************************************************
func GetDriveType(DriveLetter string) (uint32) {
    //var drvRoot []uint16
    var uinRoot *uint16

    uinRoot, _ = windows.UTF16PtrFromString(DriveLetter)

    // return windows.GetDriveType(&drvRoot[0])
    return windows.GetDriveType(uinRoot)
}


// ****************************************************************
// Check the Disk Space and return Total and Free Space           *
// ****************************************************************
func winFreeDisk() (uint64, uint64) {

    var uinRoot *uint16
    var uinCallerFreeBytes, uinTotalBytes, uinTotalFree uint64

    // What is our Base Directory Drive Root
    BaseDrive := fmt.Sprintf("%c:\\", BaseDir[0])
    uinRoot, _ = windows.UTF16PtrFromString(BaseDrive)

    //uinRoot, _ = windows.UTF16PtrFromString("C:\\")
    _ = windows.GetDiskFreeSpaceEx(uinRoot, &uinCallerFreeBytes, &uinTotalBytes, &uinTotalFree)

   return uinTotalBytes, uinTotalFree
}


// ****************************************************************
// Windows Console Hide and Show                                  *
// ****************************************************************
func winConHideShow(HideOrShow int) {
    console := w32.GetConsoleWindow()
    if console != 0 {
        _, consoleProcID := w32.GetWindowThreadProcessId(console)
        if w32.GetCurrentProcessId() == consoleProcID {
            if HideOrShow == 0 {
                w32.ShowWindowAsync(console, w32.SW_HIDE)
            } else {
                w32.ShowWindowAsync(console, w32.SW_SHOW)
            }
        }
    }
}


// ****************************************************************
// Windows Get Volume Information                                 *
// ****************************************************************
func winGetVolInfo(rootDrive string) (string) {
    var vol_err error

    rootPathName := windows.StringToUTF16Ptr(rootDrive)
    fileSystemName := make([]uint16, windows.MAX_PATH+1)

    if vol_err = windows.GetVolumeInformation(rootPathName, nil, 0, nil, nil, nil, &fileSystemName[0], uint32(len(fileSystemName))); vol_err != nil {
        volType = "UNKNOWN"
    } else {
        volType = windows.UTF16ToString(fileSystemName)
    }

    return volType
}


// ****************************************************************
// Gets the Modified, Create and Access time of a file            *
// ****************************************************************
func FTime(FileName string) (time.Time, time.Time, time.Time) {
    var atime, mtime, ctime time.Time

    stat, err_ftime := os.Stat(FileName)
    if err_ftime != nil {
        return time.Time{}, time.Time{}, time.Time{}
    }

    wstat := stat.Sys().(*syscall.Win32FileAttributeData)
    atime = time.Unix(0, wstat.LastAccessTime.Nanoseconds())
    mtime = time.Unix(0, wstat.LastWriteTime.Nanoseconds())
    ctime = time.Unix(0, wstat.CreationTime.Nanoseconds())

    return atime, mtime, ctime
}


// ****************************************************************
// Gets the Windows Version: Major, Minor, and BuildNumber        *
// ****************************************************************
func GetOSVer() string {
    osVersion := w32.RtlGetVersion()
    str_osVersion := fmt.Sprintf("Windows %d.%d.%d", osVersion.MajorVersion, osVersion.MinorVersion, osVersion.BuildNumber)

    return str_osVersion
}


//****************************************************************
// Check for Windows Admin Privs by opening Physical Drive0      *
//****************************************************************
func IsUserAdmin() bool {
    _, adm_err := os.Open("\\\\.\\PHYSICALDRIVE0")
    if adm_err != nil {
        return false
    } 
    return true
}


//******************************************************************
// Get Memory Size: Windows Version                                *
//  Copied from:                                                   *
//  https://github.com/pbnjay/memory/blob/master/memory_windows.go *
//******************************************************************
// omitting a few fields for brevity...
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx

func sysTotalMemory() uint64 {
    type memStatusEx struct {
        dwLength     uint32
        dwMemoryLoad uint32
        ullTotalPhys uint64
        unused       [6]uint64
    }

    kernel32, load_err := syscall.LoadDLL("kernel32.dll")
    if load_err != nil {
        return 0
    }

    // GetPhysicallyInstalledSystemMemory is simpler, but broken on
    // older versions of windows (and uses this under the hood anyway).
    globalMemoryStatusEx, find_err := kernel32.FindProc("GlobalMemoryStatusEx")
    if find_err != nil {
        return 0
    }

    ptr_msx := &memStatusEx{
        dwLength: 64,
    }

    mems_msx, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(ptr_msx)))
    if mems_msx == 0 {
        return 0
    }

    return ptr_msx.ullTotalPhys
}

//****************************************************************
// UnEmbed and UnZip Embedded File(s)                            *
//****************************************************************
func UnEmbed(embdata []byte) bool {
    ZipRdr, unz_err := zip.NewReader(bytes.NewReader(embdata), int64(len(embdata)))
    if unz_err != nil {
        ConsOut = fmt.Sprintf("[!] Unzip Error: %s\n", unz_err)
        ConsLogSys(ConsOut, 1, 1)
        return false
    }

    unz_err = Unzip(ZipRdr.File, BaseDir)
    if unz_err != nil {
        ConsOut = fmt.Sprintf("[!] Unzip Error: %s\n", unz_err)
        ConsLogSys(ConsOut, 1, 1)
        return false
    }

    return true

}


//****************************************************************
// Accept a Key String and convert it to registry key            *
//****************************************************************
func makeKey(khive string) {
    if strings.HasPrefix(strings.ToUpper(khive), "HKLM\\") {
        subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            fmt.Fprintf(RegHndl, "Error Opening Key: %s\n", err)
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
        subKey.Close()
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCR\\") {
        subKey, err := registry.OpenKey(registry.CLASSES_ROOT, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            fmt.Fprintf(RegHndl, "Error Opening Key: %s\n", err)
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
        subKey.Close()
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCU\\") {
        subKey, err := registry.OpenKey(registry.CURRENT_USER, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            fmt.Fprintf(RegHndl, "Error Opening Key: %s\n", err)
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
        subKey.Close()
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKUS\\") {
        subKey, err := registry.OpenKey(registry.USERS, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            fmt.Fprintf(RegHndl, "Error Opening Key: %s\n", err)
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
        subKey.Close()
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCC\\") {
        subKey, err := registry.OpenKey(registry.CURRENT_CONFIG, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            fmt.Fprintf(RegHndl, "Error Opening Key: %s\n", err)
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
        subKey.Close()
    }
}

//****************************************************************
// Enumerate a Registry Key (entries and sub-keys)               *
//  Note: Since the REgistry Entry can contain almost anything   *
//        this routing does TAB Separated instead of Comma       *
//        Separated - which helps to ensure better parsing       *
//****************************************************************
func walkKey(k registry.Key, kname string) {

    names, err := k.ReadValueNames(-1)
    if err != nil {
        ConsOut = fmt.Sprintf("[!] Reading Value Names of %s Failed: %v\n", kname, err)
        ConsLogSys(ConsOut, 1, 1)
        return
    }

    for _, name := range names {
        _, valtype, err := k.GetValue(name, nil)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Reading Value Type of %s of %s Failed: %v\n", name, kname, err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }

        switch valtype {
        case registry.NONE:
        case registry.SZ:
            strVal, _, err := k.GetStringValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                fmt.Fprintf(RegHndl, "REG_SZ\t%s\\%s\t%s\n", kname, name, strVal)
            }
        case registry.EXPAND_SZ:
            strVal, _, err := k.GetStringValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                fmt.Fprintf(RegHndl, "REG_SZ\t%s\\%s\t%s\n", kname, name, strVal)
                expVal, err := registry.ExpandString(strVal)
                if err != nil {
                    ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                    ConsLogSys(ConsOut, 1, 1)
                } else {
                    //ConsOut = fmt.Sprintf("REG_EXPAND_SZ, %s\\%s, `%s`\n", kname, name, expVal)
                    //ConsLogSys(ConsOut, 1, 1)
                    //fmt.Fprintf(RegHndl, "REG_EXPAND_SZ, %s\\%s, `%s`\n", kname, name, expVal)
                    fmt.Fprintf(RegHndl, "REG_EXPAND_SZ\t%s\\%s\t%s\n", kname, name, expVal)
                }
            }
        case registry.DWORD:
            val64, _, err := k.GetIntegerValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_DWORD, %s\\%s, %d\n", kname, name, val64)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_DWORD, %s\\%s, %d\n", kname, name, val64)
                fmt.Fprintf(RegHndl, "REG_DWORD\t%s\\%s\t%d\n", kname, name, val64)
            }
        case registry.QWORD:
            val64, _, err := k.GetIntegerValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_QWORD, %s\\%s, %d\n", kname, name, val64)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_QWORD, %s\\%s, %d\n", kname, name, val64)
                fmt.Fprintf(RegHndl, "REG_QWORD\t%s\\%s\t%d\n", kname, name, val64)
            }
        case registry.BINARY:
            binVal, _, err := k.GetBinaryValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_BINARY, %s\\%s, %v\n", kname, name, binVal)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_BINARY, %s\\%s, %v\n", kname, name, binVal)
                fmt.Fprintf(RegHndl, "REG_BINARY\t%s\\%s\t%v\n", kname, name, binVal)
            }
        case registry.MULTI_SZ:
            strVal, _, err := k.GetStringsValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                //ConsOut = fmt.Sprintf("REG_MULTI_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                //ConsLogSys(ConsOut, 1, 1)
                //fmt.Fprintf(RegHndl, "REG_MULTI_SZ, %s\\%s, `%s`\n", kname, name, strVal)
                fmt.Fprintf(RegHndl, "REG_MULTI_SZ\t%s\\%s\t%s\n", kname, name, strVal)
            }
        case registry.FULL_RESOURCE_DESCRIPTOR, registry.RESOURCE_LIST, registry.RESOURCE_REQUIREMENTS_LIST:
            // TODO: not implemented
        default:
            ConsOut = fmt.Sprintf("[!] Value Type %d of %s of %s Failed: %v\n", valtype, name, kname, err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
    }

    names, err = k.ReadSubKeyNames(-1)
    if err != nil {
        ConsOut = fmt.Sprintf("[!] Reading Sub-Keys of %s Failed: %v\n", kname, err)
        ConsLogSys(ConsOut, 1, 1)
    }

    for _, name := range names {
        func() {
            subk, err := registry.OpenKey(k, name, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
            if err != nil {
                if err == syscall.ERROR_ACCESS_DENIED {
                    // ignore error, if we are not allowed to access this key
                    return
                }

                ConsOut = fmt.Sprintf("[!] Opening Sub-Keys %s of %s Failed: %v\n", name, kname, err)
                ConsLogSys(ConsOut, 1, 1)
                return
            }

            defer subk.Close()

            walkKey(subk, kname+`\`+name)
        }()
    }
}


//****************************************************************
// Raw NTFS Copy Routine.  Uses the Velocidex NTFS library       *
//  Note: Implementation mostly copied from:                     *
//        https://github.com/kmahyyg/go-rawcopy                  *
//        With my own spin on it.                                *
//****************************************************************
func NTFSRawCopy(NTFSInFile string, NTFSOutFile string) {
    if iIsAdmin != 1 {
        ConsOut = fmt.Sprintf("[!] Raw NTFS Access MUST run as Elevated... Bypassing Raw Copy...\n")
        ConsLogSys(ConsOut, 1, 1)
        return
    }

    npath := EnsureNTFSPath(NTFSInFile)

    //****************************************************************
    // fullpath can leave with prefixing backslash, and this library *
    // require file path in slash (*nix format)                      *
    //****************************************************************
    npathRela := strings.Join(npath[1:], "//")

    err := TryRetrieveFile(npath[0], npathRela, NTFSOutFile)
    if err != nil {
        ConsOut = fmt.Sprintf("[!] Error Accessing FileName: %v\n", err)
        ConsLogSys(ConsOut, 1, 1)
        return
    }
}


//****************************************************************
// Raw NTFS Copy Routine.  Check Input File Presence             *
//****************************************************************
func EnsureNTFSPath(volFilePath string) []string {
    return strings.Split(volFilePath, "\\")
}


//****************************************************************
// TryRetrieveFile create a NTFSContext for specific volume and  *
// find the corresponding file, after finding the corresponding  *
// MFT Entry, read the (ATTR(Type=16)-$STANDARD_INFORMATION)     *
// to retrieve file metadata. Then use OpenStream to try extract *
// file from (ATTR(Type=128)-$DATA), read data via raw device    *
// require pagedReader, each read operation must fit a cluster   *
// size, which by default, is 4096 bytes.                        *
//****************************************************************
func TryRetrieveFile(volDiskLetter string, filePath string, NTFSOutFile string) error {
    ConsOut = fmt.Sprintf("[+] Checking Drive Letter...\n")
    ConsLogSys(ConsOut, 3, 3)

    // check user input
    var IsDiskLetter = regexp.MustCompile(`^[a-zA-Z]:$`).MatchString
    if !IsDiskLetter(volDiskLetter) {
        return ErrInvalidInput
    }


    //****************************************************************
    // use UNC path to access raw device to bypass any file lock     *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Opening Raw Device Handle...\n")
    ConsLogSys(ConsOut, 3, 3)
    volFd, err := os.Open("\\\\.\\" + volDiskLetter)
    if err != nil {
        return ErrDeviceInaccessible
    }


    //****************************************************************
    // build a pagedReader for raw device to feed the NTFSContext initializor
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Creating PagedReader with page 4096, cache size 65536...\n")
    ConsLogSys(ConsOut, 3, 3)
    ntfsPagedReader, err := ntfs.NewPagedReader(volFd, 0x1000, 0x10000)
    if err != nil {
        return err
    }


    //****************************************************************
    // build NTFS context for root device                            *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Creating NTFSContext...\n")
    ConsLogSys(ConsOut, 3, 3)
    ntfsVolCtx, err := ntfs.GetNTFSContext(ntfsPagedReader, 0)
    if err != nil {
        return err
    }

    //****************************************************************
    // Get volume root                                               *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Locating Volume Root Directory...\n")
    ConsLogSys(ConsOut, 3, 3)
    ntfsVolRoot, err := ntfsVolCtx.GetMFT(5)
    if err != nil {
        return err
    }

    //****************************************************************
    // Locate MFT_Entry Location                                     *
    // ref: https://www.andreafortuna.org/2017/07/18/how-to-extract- *
    //      data-and-timeline-from-master-file-table-on-ntfs-        *  
    //      filesystem/                                              *
    // a resident file might contain multiple VCN attributes and     *
    // sub-streams, use OpenStream to retrieve data                  *
    // here, we found corresponding MFT Entry first.                 *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Locating file MFT_Entry Location...\n")
    ConsLogSys(ConsOut, 3, 3)
    corrFileEntry, err := ntfsVolRoot.Open(ntfsVolCtx, filePath)
    if err != nil {
        return err
    }

    //****************************************************************
    // After locating MFT_ENTRY, retrieve file metadata information  * 
    // located in corresponding data-stream attribute                *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Identifying File Metadata...\n")
    ConsLogSys(ConsOut, 3, 3)

    corrFileInfo, err := corrFileEntry.StandardInformation(ntfsVolCtx)
    if err != nil {
        return err
    }

    fulPath := ntfs.GetFullPath(ntfsVolCtx, corrFileEntry)
    if err != nil {
        return err
    }

    err = PrintFileMetadata(corrFileInfo, volDiskLetter+"/"+fulPath)
    if err != nil {
        return err
    }


    //****************************************************************
    // retrieve file data by read $DATA                              *
    // check handwritten.go inside NTFS library for more details     *
    // ref: www.velocidex.com/golang/go-ntfs@v0.1.2-0.20221022134455 *
    //      -f91169ef8a39/parser/handwritten.go:63                   *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Retrieving data streaming from attr...\n")
    ConsLogSys(ConsOut, 3, 3)
    corrFileReader, err := ntfs.OpenStream(ntfsVolCtx, corrFileEntry, NTFSAttrType_Data, NTFSWildcard_Stream_ID, NTFSWildcard_Stream_Name)
    if err != nil {
        return err
    }


    //****************************************************************
    // Begin Copy                                                    *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Begin Raw Copy...\n")
    ConsLogSys(ConsOut, 3, 3)
    err = CopyToDestinationFile(corrFileReader, NTFSOutFile)
    if err != nil {
        return err
    }

    //****************************************************************
    // Copy Completed - Apply Original File TimeStamps               *
    //****************************************************************
    ConsOut = fmt.Sprintf("[+] Raw Copy Complete. Applying original file time...\n")
    ConsLogSys(ConsOut, 1, 1)
    err = ApplyOriginalMetadata(volDiskLetter+"/"+fulPath, corrFileInfo, NTFSOutFile)
    if err != nil {
        return err
    }

    return nil
}


//****************************************************************
// Apply Original File TimeStamps                                *
// golang official os package does not support Creation Time     *
// change due to non-POSIX compliant,                            *
// so use windows specific API only.                             *
//****************************************************************
func ApplyOriginalMetadata(path string, info *ntfs.STANDARD_INFORMATION, dst string) error {
    winFileHd, err := windows.Open(dst, windows.O_RDWR, 0)
    defer windows.CloseHandle(winFileHd)
    if err != nil {
        return err
    }

    cTime4Win := windows.NsecToFiletime(info.Create_time().UnixNano())
    aTime4Win := windows.NsecToFiletime(info.File_accessed_time().UnixNano())
    mTime4Win := windows.NsecToFiletime(info.File_altered_time().UnixNano())
    err = windows.SetFileTime(winFileHd, &cTime4Win, &aTime4Win, &mTime4Win)
    if err != nil {
        return err
    }
    return nil
}

//****************************************************************
// Display Metadata to Screen, Log, Etc...                       *   
//****************************************************************
func PrintFileMetadata(stdinfo *ntfs.STANDARD_INFORMATION, fullpath string) error {
    if stdinfo == nil {
        return ErrReturnedNil
    }
    ConsOut = fmt.Sprintf("    File Path: %s\n    File ATime: %s\n    File CTime: %s\n    File MTime: %s\n    MFT MTime: %s\n", 
              fullpath, stdinfo.File_accessed_time().String(), stdinfo.Create_time().String(), stdinfo.File_altered_time().String(),
              stdinfo.Mft_altered_time().String())
    ConsLogSys(ConsOut, 1, 1)
    return nil
}


//****************************************************************
// Copy to Destination File                                      *   
//****************************************************************
func CopyToDestinationFile(src ntfs.RangeReaderAt, dstfile string) error {
    var LastSlash = 0
    var PathOnly = "/"
    var TmpTooFile = ""
    var iFileCount = 0

    if src == nil {
        return ErrReturnedNil
    }

    //***********************************************************************
    //* Check to make sure Dest Directory Exists                            *
    //***********************************************************************
    LastSlash = strings.LastIndexByte(dstfile, slashDelim)
    PathOnly = dstfile[:LastSlash]
    ExpandDirs(PathOnly)


    //***********************************************************************
    //* Never OverWrite a File - Rename if it is already there              *
    //***********************************************************************
    if FileExists(dstfile) {
        TmpTooFile = dstfile
        iFileCount = 0
        for FileExists(TmpTooFile) {
            iFileCount++
            TmpTooFile = fmt.Sprintf("%s(%d)", dstfile, iFileCount)
        }

        dstfile = TmpTooFile
        ConsOut = fmt.Sprintf("[*] Destination File Already Exists.\n    Renamed To: %s\n", dstfile)
        ConsLogSys(ConsOut, 1, 2)
    }

    ConsOut = fmt.Sprintf("[+] Raw Copy to: %s\n", dstfile)
    ConsLogSys(ConsOut, 1, 1)
    dstFd, err := os.Create(dstfile)
    defer dstFd.Sync()
    defer dstFd.Close()
    if err != nil {
        return err
    }

    //****************************************************************
    // Output Cluster Run Info - For Debug only                      *   
    //****************************************************************
    for idx, rn := range src.Ranges() {
        ConsOut = fmt.Sprintf("[+] Split Run %03d : Range Start From %v - Length: %v , IsSparse %v \n", idx, rn.Offset, rn.Length, rn.IsSparse)
        ConsLogSys(ConsOut, 3, 3)
    }

    convertedReader := ConvertFromReaderAtToReader(src, 0)

    wBytes, err := io.Copy(dstFd, convertedReader)
    if err != nil {
        return err
    }

    TooMD5 := GetMD5File(dstfile)
    ConsOut = fmt.Sprintf("[+] Written %d Bytes to Destination Done.\n[+] File Hash: %s\n", wBytes, TooMD5)
    ConsLogSys(ConsOut, 1, 1)

    return nil
}


//****************************************************************
// Reader Definitions                                            *   
//****************************************************************
type readerFromRangedReaderAt struct {
    r      io.ReaderAt
    offset int64
}

func ConvertFromReaderAtToReader(r io.ReaderAt, o int64) *readerFromRangedReaderAt {
    return &readerFromRangedReaderAt{r: r, offset: o}
}

func (rd *readerFromRangedReaderAt) Read(b []byte) (n int, err error) {
    n, err = rd.r.ReadAt(b, rd.offset)
    if n > 0 {
        rd.offset += int64(n)
    }

    return
}





