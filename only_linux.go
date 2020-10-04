// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "syscall"
    "os"
    "time"
)


func GetDriveType(DriveLetter string) (uint32) {
    fmt.Printf("[!] Drive Listing not relevant to Linux\n")

    return 0
}

func winFreeDisk() (uint64, uint64) {
    var diskAll, diskFree uint64

    fs := syscall.Statfs_t{}
    err := syscall.Statfs("/", &fs)

    if err != nil {
        return 0, 0
    }

    diskAll = fs.Blocks * uint64(fs.Bsize)
    diskFree = fs.Bfree * uint64(fs.Bsize)

    //fmt.Printf("[!] Free Disk Space Listing not implemented for Linux\n")
    //return 0, 0

    return diskAll, diskFree
}


// Windows Console Hide and Show
func winConHideShow(HideOrShow int) {
    fmt.Printf("[+] Console Hide/Show not Implemented for Linux\n")
}


// Windows Get Volume Information
func winGetVolInfo(rootDrive string) (string) {
    fmt.Printf("[!] Drive Listing not Implemeted for Linux\n")
    return "UNKNOWN"
}


// Gets the Modified, Create and Access time of a file
func FTime(FileName string) (time.Time, time.Time, time.Time) {
     var atime, mtime, ctime time.Time

    FileInfo, err_ftime := os.Stat(FileName)
    if err_ftime != nil {
        return time.Time{}, time.Time{}, time.Time{}
    }

    var stat = FileInfo.Sys().(*syscall.Stat_t)
    atime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
    ctime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
    mtime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)

    return atime, mtime, ctime
}


// Gets the Modified, Create and Access time of a file
func GetOSVer() string {
    return "Linux"
}

