// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "syscall"
    "os"
    "time"
    "golang.org/x/sys/windows"
    "github.com/gonutz/w32"
)


func GetDriveType(DriveLetter string) (uint32) {
    //var drvRoot []uint16
    var uinRoot *uint16

    uinRoot, _ = windows.UTF16PtrFromString(DriveLetter)

    // return windows.GetDriveType(&drvRoot[0])
    return windows.GetDriveType(uinRoot)
}


// Check the Disk Space and return Total and Free Space
func winFreeDisk() (uint64, uint64) {

    var uinRoot *uint16
    var uinCallerFreeBytes, uinTotalBytes, uinTotalFree uint64

    uinRoot, _ = windows.UTF16PtrFromString("C:\\")
    _ = windows.GetDiskFreeSpaceEx(uinRoot, &uinCallerFreeBytes, &uinTotalBytes, &uinTotalFree)

   return uinTotalBytes, uinTotalFree
}


// Windows Console Hide and Show
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


// Windows Get Volume Information
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


// Gets the Modified, Create and Access time of a file
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


// Gets the Windows Version: Major, Minor, and BuildNumber
func GetOSVer() string {
    osVersion := w32.RtlGetVersion()
    str_osVersion := fmt.Sprintf("Windows %d.%d.%d", osVersion.MajorVersion, osVersion.MinorVersion, osVersion.BuildNumber)

    return str_osVersion
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
