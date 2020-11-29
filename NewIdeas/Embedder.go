package main

import (
    "fmt"
    "io/ioutil"
)

func main() {
    embdata, emb_err := ioutil.ReadFile("embedded.zip")
    if emb_err != nil {
        panic(emb_err)
    }

    fmt.Print("var embdata = []byte{")
    for i, v := range embdata {
        if i > 0 {
            fmt.Print(",")
        }
        fmt.Print(v)
    }
fmt.Println("}")
}
