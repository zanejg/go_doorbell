package main

import (
    "fmt"
    "os"
    "net/http"
    "bufio"
    "bytes"
)

func main() {
        // look for whther there were any commandline args
    parms := os.Args[1:]
    //joinflag := false

    if len(parms) > 0 {
        // then there were cmd line parms so we need to deal with them
        if len(parms) == 1 {
            // then we have the correct number of parms so we will attempt to
            // use them and send the join request to the central doorbell
            client := &http.Client{}
            req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:3434/join", parms[0]),nil)
            //req.Header.Add("Accept", "application/json")
            resp, err := client.Do(req)
            if err != nil {
                //log.Print(err)
                fmt.Printf("@@@ ERROR:%s",err)
                os.Exit(1)
            }

            scanner := bufio.NewScanner(resp.Body)
            scanner.Split(bufio.ScanRunes)
            var buf bytes.Buffer
            for scanner.Scan() {
                buf.WriteString(scanner.Text())
            }
            //fmt.Println(buf.String())

            if buf.String() == "Join Success"{
                fmt.Printf("Server reported success in subscribing this doorbell\n")
                syncreq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:3434/syncnew", parms[0]),nil)
                //req.Header.Add("Accept", "application/json")
                _, syncerr := client.Do(syncreq)
                if syncerr != nil {
                    //log.Print(err)
                    fmt.Printf("@@@ ERROR:%s",syncerr)
                    os.Exit(1)
                }



                //joinflag = true
            } else {
                fmt.Printf("Something went wrong when trying to join %s\n", parms[0])
                os.Exit(1)
            }


        } else {
            fmt.Printf("The correct syntax is doorbell join serverIP(or hostname)\n")
        }


    }
}
