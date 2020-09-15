// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "syscall"
)

func winListDrives() {
    fmt.Printf("[!] Drive Listing not relevant to Linux\n")

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

