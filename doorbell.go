package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    "os/exec"
    "bytes"
    "html/template"
    "os"
    "io"
    "time"
    "crypto/md5"
    "strconv"
    "strings"
    "io/ioutil"
    "encoding/json"
    //"bufio"
)

const (
    //MP3path = "/home/zane/programming/go/webserver/thesounds/"
    //MP3path = "/home/zane/programming/go/webserver/clientsounds/"
    MP3subpath = "thesounds/"

)


func Ls(dirpath string) []string {
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



func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Ring(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "ringing\n")

    cmd := exec.Command("mpg123", fmt.Sprintf("%s%s", MP3path, "hailed.mp3"))
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
func RingChime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //fmt.Fprintf(w, "ringing\n")
    if r.Method == "GET" {

        filename := ps.ByName("name")
        filelist := Ls(MP3path)

        if contains(filelist, filename){
            cmd := exec.Command("mpg123", fmt.Sprintf("%s%s", MP3path, filename))
        	//cmd.Stdin = strings.NewReader("some input")
        	var out bytes.Buffer
        	cmd.Stdout = &out
        	err := cmd.Run()
        	if err != nil {
        		log.Fatal(err)
        	}
        }
    }
}




type PageData struct {
    Serverstr string
    Token string
}


func PutChime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
       fmt.Println("method:", r.Method)
       if r.Method == "GET" {
           crutime := time.Now().Unix()
           h := md5.New()
           io.WriteString(h, strconv.FormatInt(crutime, 10))

           pagedata := PageData{Token:"Big Secret Token",
                                Serverstr:"192.168.178.20"}
           //token := fmt.Sprintf("%x", h.Sum(nil))

           t, _ := template.ParseFiles("mp3upload.gtpl")
           t.Execute(w, pagedata)
       } else {
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(fmt.Sprintf("FormFile error:%s",err))
               return
           }
           defer file.Close()
           fmt.Fprintf(w, "%v", handler.Header)
           f, err := os.OpenFile(MP3path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(fmt.Sprintf("Error opening file:%s",err))
               return
           }
           defer f.Close()
           io.Copy(f, file)
           fmt.Printf("Got file:%s\n", handler.Filename)
       }
}


// func PutChime(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
//     fmt.Fprintln(rw, "post update")
// }
type Config struct {
    Doorbell_dir string
    Satellite_port  int
}
var CONFIG Config
var MP3path string
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




func main() {

    CONFIG = GetConfig()
    fmt.Println("DIR:",CONFIG.Doorbell_dir)
    fmt.Println("Port:",CONFIG.Satellite_port)

    MP3path = CONFIG.Doorbell_dir + "/" + MP3subpath

    // no cmd parms so we just run normally
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)
    router.GET("/ring", Ring)
    router.Handle("GET","/putchime", PutChime)
    router.Handle("POST","/putchime", PutChime)
    router.GET("/ringchime/:name", RingChime)
    //router.PUT("/putchime", PutChime)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",CONFIG.Satellite_port), router))

}
