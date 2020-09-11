// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "golang.org/x/sys/windows"
)

func winListDrives() {
    var drvRoot []uint16 
    drvRoot, _ = windows.UTF16FromString("C:\\")
    fmt.Printf("C Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("H:\\")
    fmt.Printf("H Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("I:\\")
    fmt.Printf("I Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("Z:\\")
    fmt.Printf("Z Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("M:\\")
    fmt.Printf("M Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

}