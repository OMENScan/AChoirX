// ****************************************************************
// * Only Include this file for OSX version                       *
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
    "bytes"
    "archive/zip"
)


// ****************************************************************
// * Embed Zip File in embdata byte array                         *
// ****************************************************************
//go:embed OSXEmbed.zip
var embdata []byte


// ****************************************************************
// * Windows Drive Letters - Not relevant for OSX                 *
// ****************************************************************
func GetDriveType(DriveLetter string) (uint32) {
    fmt.Printf("[!] Drive Listing not relevant to OSX\n")

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

    //fmt.Printf("[!] Free Disk Space Listing not implemented for OSX\n")
    //return 0, 0

    return diskAll, diskFree
}


// ****************************************************************
// Windows Console Hide and Show                                  *
// ****************************************************************
func winConHideShow(HideOrShow int) {
    fmt.Printf("[+] Console Hide/Show not Implemented for OSX\n")
}

//****************************************************************
// Accept a Key String and convert it to registry key            *
//****************************************************************
func makeKey(khive string) {
    fmt.Printf("[+] Reistry Parser not Implemented for OSX\n")
}


// ****************************************************************
// Windows Get Volume Information                                 *
// ****************************************************************
func winGetVolInfo(rootDrive string) (string) {
    fmt.Printf("[!] Drive Listing not Implemeted for OSX\n")
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
    //atime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
    //ctime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
    //mtime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
    atime = timespecToTime(stat.Atimespec)
    mtime = timespecToTime(stat.Mtimespec)
    ctime = timespecToTime(stat.Ctimespec)

    return atime, mtime, ctime
}


// ****************************************************************
// Convert Timespec to time.Time                                  *
// ****************************************************************
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}


// ****************************************************************
// Platform Specific - Return OS (Version)                        *
// ****************************************************************
func GetOSVer() string {
    return "OSX"
}


//****************************************************************
// Check for OSX Root by checking Effective UserId               *
//****************************************************************
func IsUserAdmin() bool {
    if os.Geteuid() != 0 {
        return false
    } 
    return true
}



//******************************************************************
// Get Memory Size: OSX Version                                    *
//  Copied from:                                                   *
//  https://github.com/pbnjay/memory/blob/master/memory_windows.go *
//******************************************************************
func sysTotalMemory() uint64 {
    sysmem, mem_err := sysctlUint64("hw.memsize")
    if mem_err != nil {
        return 0
    }
    return sysmem
}


func sysctlUint64(ctlname string) (uint64, error) {
    ctl_data, ctl_err := syscall.Sysctl(ctlname)
    if ctl_err != nil {
        return 0, ctl_err
    }

    // hack because the string conversion above drops a \0
    byt_data := []byte(ctl_data)
    if len(byt_data) < 8 {
        byt_data = append(byt_data, 0)
    }

    return *(*uint64)(unsafe.Pointer(&byt_data[0])), nil

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
