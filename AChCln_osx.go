// ****************************************************************
// * Only Include this file for OSX version                       *
// * - Uses the Embed Directive - Requires GoLang 1.16            *
// ****************************************************************

package main

import (
    "fmt"
    "os"
)

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

// ****************************************************************
// Windows Console Hide and Show                                  *
// ****************************************************************
func winConHideShow(HideOrShow int) {
    fmt.Printf("[+] Console Hide/Show not Implemented for OSX\n")
}
