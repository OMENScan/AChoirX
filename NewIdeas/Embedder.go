// ****************************************************************
// "As we benefit from the inventions of others, we should be
//  glad to share our own...freely and gladly."
//  - Benjamin Franklin
//
// Embedder - v0.01
//  AChoirX Embedded File Creator
//  This Utility looks for 3 Zip Files:
//   WinEmbed.zip, LinEmbed.zip, and OSXEmbed.zip
//
//  It creates a Byte Array in Output Files:
//   WinEmbed.go, LinEmbed.go, and OSXEmbed.go
//  
// Those files are then compiled into the Windows, Linux, and OSX
//  versions of AChoirX - These embedded files are unzipped by 
//  AChoirX on each platform to allow files to be included with 
//  the AChoirX binary Executable
//
package main

import (
    "os"
    "fmt"
    "io/ioutil"
)

func main() {
    var EmbHndl *os.File

    //**************************************
    //* Write out Windows Embedded Go File *
    //**************************************
    EmbHndl, emb_err := os.Create("WinEmbed.go")
    if emb_err != nil {
        fmt.Printf("[!] Error Creating Windows Embed File: %s\n", emb_err)
    } else {
        embdata, emb_err := ioutil.ReadFile("WinEmbed.zip")
        if emb_err != nil {
            fmt.Printf("[!] Error Creating Windows Embed File: %s\n", emb_err)
        } else {
            fmt.Printf("[+] Now Creating Creating Windows Embed File: WinEmbed.go\n")

            fmt.Fprint(EmbHndl, "package main\n\nvar embdata = []byte{")
            for i, v := range embdata {
                if i > 0 {
                    fmt.Fprintf(EmbHndl, ",")
                }
                fmt.Fprintf(EmbHndl, "%d", v)
            }

        fmt.Fprintf(EmbHndl, "}")
        }
    EmbHndl.Close()
    }

    //**************************************
    //* Write out Linux Embedded Go File   *
    //**************************************
    EmbHndl, emb_err = os.Create("LinEmbed.go")
    if emb_err != nil {
        fmt.Printf("[!] Error Creating Linux Embed File: %s\n", emb_err)
    } else {
        embdata, emb_err := ioutil.ReadFile("LinEmbed.zip")
        if emb_err != nil {
            fmt.Printf("[!] Error Creating Linux Embed File: %s\n", emb_err)
        } else {
            fmt.Printf("[+] Now Creating Creating Linux Embed File: LinEmbed.go\n")

            fmt.Fprint(EmbHndl, "package main\n\nvar embdata = []byte{")
            for i, v := range embdata {
                if i > 0 {
                    fmt.Fprintf(EmbHndl, ",")
                }
                fmt.Fprintf(EmbHndl, "%d", v)
            }

        fmt.Fprintf(EmbHndl, "}")
        }
    EmbHndl.Close()
    }

    //**************************************
    //* Write out OSX Embedded Go File     *
    //**************************************
    EmbHndl, emb_err = os.Create("OSXEmbed.go")
    if emb_err != nil {
        fmt.Printf("[!] Error Creating OSX Embed File: %s\n", emb_err)
    } else {
        embdata, emb_err := ioutil.ReadFile("OSXEmbed.zip")
        if emb_err != nil {
            fmt.Printf("[!] Error Creating OSX Embed File: %s\n", emb_err)
        } else {
            fmt.Printf("[+] Now Creating Creating OSX Embed File: OSXEmbed.go\n")

            fmt.Fprint(EmbHndl, "package main\n\nvar embdata = []byte{")
            for i, v := range embdata {
                if i > 0 {
                    fmt.Fprintf(EmbHndl, ",")
                }
                fmt.Fprintf(EmbHndl, "%d", v)
            }

        fmt.Fprintf(EmbHndl, "}")
        }
    EmbHndl.Close()
    }

}
