// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "syscall"
    "os"
    "time"
    "bytes"
    "archive/zip"
)


// ****************************************************************
// * Windows Drive Letters - Not relevant for Linux               *
// ****************************************************************
func GetDriveType(DriveLetter string) (uint32) {
    fmt.Printf("[!] Drive Listing not relevant to Linux\n")

    return 0
}


// ****************************************************************
// Check the Disk Space and return Total and Free Space           *
// ****************************************************************
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


// ****************************************************************
// Windows Console Hide and Show                                  *
// ****************************************************************
func winConHideShow(HideOrShow int) {
    fmt.Printf("[+] Console Hide/Show not Implemented for Linux\n")
}


// ****************************************************************
// Windows Get Volume Information                                 *
// ****************************************************************
func winGetVolInfo(rootDrive string) (string) {
    fmt.Printf("[!] Drive Listing not Implemeted for Linux\n")
    return "UNKNOWN"
}


// ****************************************************************
// Gets the Modified, Create and Access time of a file            *
// ****************************************************************
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


// ****************************************************************
// Platform Specific - Return OS (Version)                        *
// ****************************************************************
func GetOSVer() string {
    return "Linux"
}


//****************************************************************
// Check for Linux Root by checking Effective UserId             *
//****************************************************************
func IsUserAdmin() bool {
    if os.Geteuid() != 0 {
        return false
    } 
    return true
}


//******************************************************************
// Get Memory Size: Linux Version                                  *
//  Copied from:                                                   *
//  https://github.com/pbnjay/memory/blob/master/memory_windows.go *
//******************************************************************
func sysTotalMemory() uint64 {
    sysin := &syscall.Sysinfo_t{}
    sys_err := syscall.Sysinfo(sysin)
    if sys_err != nil {
        return 0
    }

    // If this is a 32-bit system, then these fields are
    // uint32 instead of uint64.
    // So we always convert to uint64 to match signature.
    return uint64(sysin.Totalram) * uint64(sysin.Unit)
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

    Unzfiles, unz_err := Unzip(ZipRdr.File, BaseDir)
    if unz_err != nil {
        ConsOut = fmt.Sprintf("[!] Unzip Error: %s\n", unz_err)
        ConsLogSys(ConsOut, 1, 1)
        return false
    }

    for iUnz := 0; iUnz < len(Unzfiles); iUnz++ {
        ConsOut = fmt.Sprintf("[*] Unzipped File: %s\n", Unzfiles[iUnz])
        ConsLogSys(ConsOut, 1, 1)
    }

    return true

}
