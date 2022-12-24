// ****************************************************************
// * Only Include this file for windows version                   *
// * - Uses the Embed Directive - Requires GoLang 1.16            *
// ****************************************************************

package main

import (
    "fmt"
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
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCR\\") {
        subKey, err := registry.OpenKey(registry.CLASSES_ROOT, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCU\\") {
        subKey, err := registry.OpenKey(registry.CURRENT_USER, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKUS\\") {
        subKey, err := registry.OpenKey(registry.USERS, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
    } else if strings.HasPrefix(strings.ToUpper(khive), "HKCC\\") {
        subKey, err := registry.OpenKey(registry.CURRENT_CONFIG, khive[5:], registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Error Opening Key: %s\n", err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
        walkKey(subKey, khive)
    }
}

//****************************************************************
// Enumerate a Registry Key (entries and sub-keys)               *
//****************************************************************
func walkKey(k registry.Key, kname string) {

    names, err := k.ReadValueNames(-1)
    if err != nil {
        ConsOut = fmt.Sprintf("[!] Reading Value Names of %s Failed: %v", kname, err)
        ConsLogSys(ConsOut, 1, 1)
        return
    }

    for _, name := range names {
        _, valtype, err := k.GetValue(name, nil)
        if err != nil {
            ConsOut = fmt.Sprintf("[!] Reading Value Type of %s of %s Failed: %v", name, kname, err)
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
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_SZ: %s\n", kname, name, strVal)
                ConsLogSys(ConsOut, 1, 1)
            }
        case registry.EXPAND_SZ:
            strVal, _, err := k.GetStringValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_SZ: %s\n", kname, name, strVal)
                ConsLogSys(ConsOut, 1, 1)

                expVal, err := registry.ExpandString(strVal)
                if err != nil {
                    ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                    ConsLogSys(ConsOut, 1, 1)
                } else {
                    ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_EXPAND_SZ: %s\n", kname, name, expVal)
                    ConsLogSys(ConsOut, 1, 1)
                }
            }
        case registry.DWORD:
            val64, _, err := k.GetIntegerValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_DWORD: %d\n", kname, name, val64)
                ConsLogSys(ConsOut, 1, 1)
            }
        case registry.QWORD:
            val64, _, err := k.GetIntegerValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_QWORD: %d\n", kname, name, val64)
                ConsLogSys(ConsOut, 1, 1)
            }
        case registry.BINARY:
            binVal, _, err := k.GetBinaryValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_BINARY: %v\n", kname, name, binVal)
                ConsLogSys(ConsOut, 1, 1)
            }
        case registry.MULTI_SZ:
            strVal, _, err := k.GetStringsValue(name)
            if err != nil {
                ConsOut = fmt.Sprintf("[!] Registry Error: %s\n", err)
                ConsLogSys(ConsOut, 1, 1)
            } else {
                ConsOut = fmt.Sprintf("Key: %s\\%s\nREG_MULTI_SZ: %s\n", kname, name, strVal)
                ConsLogSys(ConsOut, 1, 1)
            }
        case registry.FULL_RESOURCE_DESCRIPTOR, registry.RESOURCE_LIST, registry.RESOURCE_REQUIREMENTS_LIST:
            // TODO: not implemented
        default:
            ConsOut = fmt.Sprintf("[!] Value Type %d of %s of %s Failed: %v", valtype, name, kname, err)
            ConsLogSys(ConsOut, 1, 1)
            return
        }
    }

    names, err = k.ReadSubKeyNames(-1)
    if err != nil {
        ConsOut = fmt.Sprintf("[!] Reading Sub-Keys of %s Failed: %v", kname, err)
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

                ConsOut = fmt.Sprintf("[!] Opening Sub-Keys %s of %s Failed: %v", name, kname, err)
                ConsLogSys(ConsOut, 1, 1)
                return
            }

            defer subk.Close()

            walkKey(subk, kname+`\`+name)
        }()
    }
}

