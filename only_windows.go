// ****************************************************************
// * Only Include this file for windows version                   *
// ****************************************************************

package main

import (
    "fmt"
    "syscall"
    "unsafe"
    "os"
    "time"
    "bytes"
    "archive/zip"
    "golang.org/x/sys/windows"
    "github.com/gonutz/w32"
)


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

    uinRoot, _ = windows.UTF16PtrFromString("C:\\")
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
func UnEmbed() bool {

var embdata = []byte{80,75,3,4,20,0,0,0,8,0,182,102,125,81,80,158,38,162,243,3,0,0,156,11,0,0,10,0,0,0,65,67,104,111,105,114,46,65,67,81,181,86,95,79,219,48,16,127,175,212,239,112,66,19,15,104,164,18,123,139,180,135,144,164,107,53,90,186,38,192,152,242,98,146,107,107,145,196,193,118,10,249,246,59,39,109,3,52,45,155,180,33,30,92,251,126,127,124,206,157,29,56,247,246,217,191,252,235,247,106,74,152,10,141,54,132,43,174,128,254,245,10,193,207,30,48,73,48,1,15,23,172,76,53,56,238,74,112,249,19,130,88,242,66,91,0,99,109,98,215,40,43,120,96,138,199,159,129,229,9,236,40,153,2,6,75,204,81,178,20,100,153,226,103,80,43,81,166,9,228,66,195,3,66,169,48,177,54,44,133,228,25,147,60,173,136,3,240,133,101,69,138,160,26,165,150,114,12,18,99,145,101,72,58,177,68,166,121,190,36,192,206,90,40,68,250,200,53,44,132,132,76,72,4,77,130,82,148,203,21,196,34,77,49,214,92,228,228,28,12,229,255,200,165,31,218,174,40,170,25,211,171,175,51,38,117,163,226,214,218,181,217,153,196,5,234,120,101,89,86,191,231,184,63,236,104,38,23,253,158,59,187,183,79,78,239,120,30,21,155,128,232,236,44,58,59,129,147,83,39,126,58,121,71,3,223,104,79,5,204,68,202,227,170,102,10,238,3,251,219,108,142,202,28,212,96,14,131,95,112,126,142,47,168,116,34,74,253,245,148,197,121,180,93,183,244,203,190,175,57,46,185,210,178,106,125,209,76,195,43,113,9,138,173,17,70,223,175,38,81,224,187,55,243,113,120,15,198,88,20,96,92,74,174,43,139,162,3,138,233,68,92,15,195,59,103,238,111,16,98,161,159,153,196,163,8,103,82,7,155,193,177,56,154,9,253,109,104,165,52,102,199,162,71,206,220,107,125,140,152,76,222,250,104,207,192,41,10,87,100,5,211,116,56,98,41,89,166,34,39,139,89,188,66,107,181,198,237,161,236,229,208,95,99,174,111,57,62,163,132,43,177,84,109,46,253,181,110,71,198,234,151,139,87,122,27,239,95,46,162,103,158,35,5,24,112,123,248,175,144,83,250,226,215,175,145,170,82,185,153,195,183,208,195,22,117,12,30,151,244,75,200,170,169,215,96,126,51,217,57,53,41,109,71,71,156,122,146,68,165,138,136,176,149,235,140,84,178,140,72,195,187,180,18,166,15,26,27,9,165,21,56,198,15,202,53,143,81,181,62,223,184,59,146,131,143,77,237,66,255,208,149,159,175,185,20,121,70,231,10,183,212,158,216,67,138,234,141,157,230,67,115,39,30,12,98,178,174,247,202,174,246,236,231,132,54,74,251,253,128,81,173,195,21,213,222,174,142,239,38,99,23,158,22,8,41,205,118,243,253,24,250,6,210,205,40,5,101,79,193,56,167,46,152,49,211,241,118,212,33,83,143,53,235,96,221,77,92,7,108,152,223,67,38,135,33,78,154,118,122,113,69,158,55,77,183,211,78,142,90,105,70,123,100,15,185,128,110,250,105,19,83,135,116,106,140,103,70,102,193,151,157,18,227,217,102,113,192,210,180,91,97,27,98,216,223,131,60,174,138,148,85,222,52,56,136,93,44,105,181,211,153,51,159,117,122,98,178,160,45,119,19,58,178,48,144,78,62,74,197,229,248,58,216,227,116,191,223,219,111,107,110,122,25,6,148,52,139,4,250,61,32,222,195,235,112,174,226,252,96,234,47,185,80,141,23,127,234,217,173,84,91,73,251,90,239,2,254,78,236,252,162,149,219,75,64,64,29,56,41,83,122,147,152,143,174,51,183,206,129,138,169,161,6,165,218,99,166,57,109,102,96,240,84,154,71,204,96,33,224,106,28,132,84,29,31,144,108,108,190,55,184,235,93,157,214,104,135,64,137,144,7,28,110,177,173,63,21,67,99,76,87,5,126,5,213,4,128,41,7,250,73,223,243,81,162,218,227,126,21,7,183,238,7,176,47,157,91,251,198,232,81,40,205,232,134,140,28,123,47,208,0,130,230,214,117,111,154,203,118,26,222,4,254,220,242,156,176,223,27,94,207,109,215,142,12,11,221,84,81,174,75,26,25,201,150,39,58,29,94,111,174,157,211,161,144,134,35,223,204,139,139,61,103,159,230,24,87,113,138,214,37,207,155,230,236,5,223,237,33,127,193,100,75,73,43,219,107,192,83,143,118,244,26,210,241,208,2,26,150,92,241,186,113,185,194,188,69,53,189,94,71,76,173,234,186,150,154,47,88,172,27,177,81,48,178,73,165,223,251,13,80,75,1,2,63,0,20,0,0,0,8,0,182,102,125,81,80,158,38,162,243,3,0,0,156,11,0,0,10,0,36,0,0,0,0,0,0,0,32,0,0,0,0,0,0,0,65,67,104,111,105,114,46,65,67,81,10,0,32,0,0,0,0,0,1,0,24,0,149,141,122,184,145,198,214,1,149,141,122,184,145,198,214,1,225,179,45,20,196,197,214,1,80,75,5,6,0,0,0,0,1,0,1,0,92,0,0,0,27,4,0,0,0,0}


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

