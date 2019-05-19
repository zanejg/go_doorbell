package main

import (
    "fmt"
    "math/rand"
    "net/http"
    //"log"
    //"os/exec"
    //"bytes"
    //"html/template"
    "os"
    //"io"
    "io/ioutil"
    "time"
    //"crypto/md5"
    //"strconv"
    "strings"
    //"mime/multipart"
    "encoding/json"
    "bufio"
    "github.com/stianeikeland/go-rpio"
)

const (
    MP3subpath = "thesounds/"
    //MP3path = "/home/zane/programming/go/webserver/emptysounds/"

    Ring_url = "ringchime/"
    Doorbell_url = "http://localhost:3400/ringchime/"
    Doorbells_file = "./doorbells.txt"

)



type PageData struct {
    Serverstr string
    Token string
    Reply string
}



type ListPageData struct {
    Filelist []string
    Ringserver string
}

//*************************************************************************************
//*************** UTILITY FUNCTIONS ***************************************************
//*************************************************************************************


func ListStoredChimes(dirpath string) []string {
    dir, err:= os.Open(dirpath)
    if err != nil{
        fmt.Printf("opening error:%v\n",err)
    }
    ls, lerr:=dir.Readdir(0)
    if lerr != nil{
        fmt.Printf("listing error:%v\n",lerr)
    }
    //fmt.Printf("files:%v\n",ls)

    ret:= make([]string,len(ls))
    i:=0
    for _, this_file := range ls {
        //fmt.Printf("%s\n",this_file.Name())
        if strings.HasSuffix(this_file.Name(), ".mp3") {
            ret[i]=this_file.Name()
            i++
        }

    }
    // we only want to return the slice with the filenames in it
    return ret[0:i]
}




//*************************************************************************************
//*************************************************************************************


func RingAllDoorbells () string {
    // send the ring command to all doorbells in the subscribed doorbell list
    reply:= ""
    lsdir:=ListStoredChimes(MP3path)
    //fmt.Printf("%s\n",lsdir)
    if len(lsdir) < 1 {
       reply = "No files to choose from"
    }else{
        rand.Seed(time.Now().UnixNano())
        chosen := lsdir[rand.Intn(len(lsdir))]
        reply = fmt.Sprintf("Chose: %s", chosen)
        fmt.Printf("%s\n",reply)


        for _,this_doorbell := range(SubscribedDoorbells){
            go RingADoorbell(chosen,this_doorbell)
        }
    }
    return reply
}



func RingADoorbell(filename string, url string){
    /*
    Send a ring command to a satellite doorbell
    */
    client := &http.Client{}
    this_url := fmt.Sprintf("http://%s:%d/%s%s",
                                        url, CONFIG.Satellite_port, Ring_url,filename)
    fmt.Printf("this_url: %s\n",this_url)
    req, _ := http.NewRequest("GET", this_url, nil)
    //req.Header.Add("Accept", "application/json")
    resp, err := client.Do(req)
    if err != nil {
            //log.Print(err)
            fmt.Printf(fmt.Sprintf("@@@ ERROR:%s for doorbell:%s with response:%s",err,url,resp))
    }else{
        fmt.Printf(fmt.Sprintf("Rang doorbell:%s\n",url))
    }

}



func GetConfig() (Config){
    /*
    Read the config file and set the global vars
    */
    var ret Config
    raw, err := ioutil.ReadFile("./config.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    json.Unmarshal(raw, &ret)
    return ret
}

func getSubscribedDoorbells() []string {
    /*
    Get the list of all of the IPs of the subscribed doorbells.
    Reads them from the doorbells.txt file
    */
    f, err := os.Open(Doorbells_file)
    if err != nil {
          panic(err)
    }
    defer f.Close()

    var ret []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
          ret = append(ret, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
          fmt.Fprintln(os.Stderr, err)
    }

    return ret
}


/**************************************************************************
     MAIN BUTTON EVENT LOOP
***************************************************************************/
func WaitForDoorbellButton(){
    doorbell_reply := ""
    // infinite loop for waiting for button presses
    for {
        // sit here waiting for the pin to go high
        for gpio_pin.Read() == 0 {
            // CPU was maxxing out hopefully this will fix it
            // checking 50 times a second should be ample
            time.Sleep(20 * time.Millisecond)
        }
        // when the the pin goes high then ring the bells
        doorbell_reply = RingAllDoorbells()
        fmt.Printf("Replies: %s",doorbell_reply)
        // wait for 1 sec for debounce etc
        time.Sleep(1000 * time.Millisecond)
        // now wait for the pin to go high before looping
        // to wait for button again
        for gpio_pin.Read() == 1 { }
    }

}
/**************************************************************************
***************************************************************************/


var gpio_pin rpio.Pin

type Config struct {
    Doorbell_dir string
    Satellite_port  int
}
var CONFIG Config
var MP3path string
var SubscribedDoorbells []string






func main() {
    CONFIG = GetConfig()
    fmt.Println("DIR:",CONFIG.Doorbell_dir)
    fmt.Println("Port:",CONFIG.Satellite_port)

    gpio_err := rpio.Open()
    if gpio_err == nil{
        gpio_pin = rpio.Pin(17)
        gpio_pin.Input()
    }

    if gpio_err != nil{
        fmt.Println(fmt.Sprintf("GPIO Error:%s\n",gpio_err.Error()))
    }

    MP3path = CONFIG.Doorbell_dir + "/" + MP3subpath
    SubscribedDoorbells = getSubscribedDoorbells()


    //router.PUT("/putchime", PutChime)
    if gpio_err == nil{
        fmt.Printf("Waiting for Doorbell Button")
        WaitForDoorbellButton()

    }

}
