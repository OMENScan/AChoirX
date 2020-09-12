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
    var uinRoot *uint16

    uinRoot, _ = windows.UTF16PtrFromString("C:\\")
    fmt.Printf("C Drive Type = %d\n", windows.GetDriveType(uinRoot))

    drvRoot, _ = windows.UTF16FromString("H:\\")
    fmt.Printf("H Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("I:\\")
    fmt.Printf("I Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("Z:\\")
    fmt.Printf("Z Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))

    drvRoot, _ = windows.UTF16FromString("M:\\")
    fmt.Printf("M Drive Type = %d\n", windows.GetDriveType(&drvRoot[0]))
}


// Check the Disk Space and return Total and Free Space
func winFreeDisk() (uint64, uint64) {

    var uinRoot *uint16
    var uinCallerFreeBytes, uinTotalBytes, uinTotalFree uint64

    uinRoot, _ = windows.UTF16PtrFromString("C:\\")
    _ = windows.GetDiskFreeSpaceEx(uinRoot, &uinCallerFreeBytes, &uinTotalBytes, &uinTotalFree)

   return uinTotalBytes, uinTotalFree
}

