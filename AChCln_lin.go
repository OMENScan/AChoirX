// ****************************************************************
// * Only Include this file for windows version                   *
// * - Uses the Embed Directive - Requires GoLang 1.16            *
// ****************************************************************

package main

import (
    "fmt"
    "os"
    "syscall"
)

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


// ****************************************************************
// Windows Console Hide and Show                                  *
// ****************************************************************
func winConHideShow(HideOrShow int) {
    fmt.Printf("[+] Console Hide/Show not Implemented for Linux\n")
}
