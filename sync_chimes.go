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
    cmd_format_msg := "The correct syntax is sync_chimes <serverIP(or hostname)>\n"

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

            switch buf.String(){
                case "Join Success":
                    fmt.Printf("Server reported success in subscribing this doorbell\n")

                case "Doorbell already subscribed":
                    fmt.Printf("Server reported this doorbell was already subscribed. \nAttempting sync\n")

                //joinflag = true
                default:
                    fmt.Printf("Something went wrong when trying to join %s\n Message:%s\n", parms[0],buf.String())
                    os.Exit(1)
            }
            // if we are here then a sync will be required
            syncreq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:3434/sync", parms[0]),nil)
            //req.Header.Add("Accept", "application/json")
            _, syncerr := client.Do(syncreq)
            if syncerr != nil {
                //log.Print(err)
                fmt.Printf("@@@ ERROR:%s",syncerr)
                os.Exit(1)
            }


        } else {
            fmt.Printf("%s\n",cmd_format_msg)
        }


    }else {
        fmt.Printf("%s\n",cmd_format_msg)
    }
}
