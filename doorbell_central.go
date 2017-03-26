package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    //"os/exec"
    "bytes"
    "html/template"
    "os"
    "io"
    "io/ioutil"
    "time"
    "crypto/md5"
    "strconv"
    "strings"
    "mime/multipart"
    "encoding/json"
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


func SendChime(file , url string)(err error){
    // Function to send the file identified by file in the MP3path dir to the
    // doorbell webserver identified by url
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    // Add your image file
    f, err := os.Open(file)
    if err != nil {
        return
    }
    defer f.Close()

    // only want the name of the file not the path in the form
    path_slice := strings.Split(file,"/")
    filename := path_slice[len(path_slice) -1]
    fmt.Printf("filename : %s\n",filename)

    fw, err := w.CreateFormFile("uploadfile", filename)
    if err != nil {
        return
    }
    if _, err = io.Copy(fw, f); err != nil {
        return
    }
    // Add the other fields
    if fw, err = w.CreateFormField("Token"); err != nil {
        return
    }
    if _, err = fw.Write([]byte("Big Secret Token")); err != nil {
        return
    }
    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
        return
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Submit the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return
    }

    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    }
    return

}







//*************************************************************************************
//*************** WEB SERVER HANDLERS ************************************************
//*************************************************************************************


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    //fmt.Fprint(w, "Welcome!\n")
    t, _ := template.ParseFiles("index_central.gtpl")
    t.Execute(w,"")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}


func GetChime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
       fmt.Println("method:", r.Method)
       pagedata := PageData{Token:"Big Secret Token",
                            Serverstr:"192.168.178.20",
                            Reply:""}

        reply:=""

        lsdir:=Ls(MP3path)
        if len(lsdir) < 1{
            fmt.Printf("no files in dir")
        }

        if r.Method == "POST" {

           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
           defer file.Close()
           //fmt.Fprintf(w, "%v", handler.Header)
           f, err := os.OpenFile(MP3path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           // only want the file if it is an MP3
           if strings.HasSuffix(handler.Filename, ".mp3") {
               io.Copy(f, file)
               reply = fmt.Sprintf("Got file called:%s",handler.Filename)
           }else{
               reply = fmt.Sprintf("File called:%s does not appear to be an MP3",handler.Filename)
           }

       }
       crutime := time.Now().Unix()
       h := md5.New()
       io.WriteString(h, strconv.FormatInt(crutime, 10))

       pagedata = PageData{Token:"Big Secret Token",
                            Serverstr:"192.168.178.20",
                            Reply:reply}
       //token := fmt.Sprintf("%x", h.Sum(nil))

       t, _ := template.ParseFiles("mp3upload_central.gtpl")
       t.Execute(w, pagedata)
}



func ListChimes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
       fmt.Println("method:", r.Method)

       lsdir:=Ls(MP3path)
       if len(lsdir) < 1 {
           emptyls := []string{"NO FILES"}
           lsdir = emptyls
       }

       //fmt.Fprintf(w, "%v", lsdir)
       pagedata := ListPageData{Filelist:lsdir,
                                Ringserver: Doorbell_url}
       fmt.Printf("%v\n",pagedata.Ringserver)
       t, _ := template.ParseFiles("mp3_listing.gtpl")
       t.Execute(w, pagedata)
}

func TestSend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
       fmt.Println("method:", r.Method)
       SendChime(MP3path+"unother1.mp3" , "http://localhost:3400/putchime")



    //    t, _ := template.ParseFiles("mp3_listing.gtpl")
    //    t.Execute(w, pagedata)
}

func SubscribeNewDoorbell(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    ip_chunks := strings.Split(req.RemoteAddr,":")
    fmt.Println("IP:",ip_chunks[0])

    f, err := os.OpenFile(Doorbells_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString(fmt.Sprintf("%s\n",ip_chunks[0]); err != nil {
        panic(err)
    }
    fmt.Fprintf(w,"Join Success")

}

func SyncNewDoorbell(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    ip_chunks := strings.Split(req.RemoteAddr,":")
    the_server := ip_chunks[0]
    fmt.Println("sync IP:",ip_chunks[0])
    //
    // fmt.Fprintf(w,"Join Success")
    //

    // wait for doorbell to sort it's shit out
    //time.Sleep(1000 * time.Millisecond)

    // now we need to ensure that the doorbell has all the files it needs
    lsdir:=Ls(MP3path)
    //fmt.Printf("the files:%s", lsdir)

    the_url := fmt.Sprintf("http://%s:%s/putchime",the_server,CONFIG.Satellite_port)
    the_url = fmt.Sprintf("http://localhost:%s/putchime",CONFIG.Satellite_port)
    the_path := ""
    fmt.Printf("the_url:%s\n", the_url)

    for _, this_file := range lsdir {
        the_path = MP3path + this_file
        fmt.Printf("the_path:%s\n", the_path)

        SendChime( the_path, the_url)
    }


}

func all_doorbells_ring(filename string){
    /*
    Send the command to all doorbells to ring the sent chime
    */


}



//*************************************************************************************
//*************************************************************************************

type Config struct {
    Doorbell_dir string
    Satellite_port  int
}
var CONFIG Config

var MP3path string


func GetConfig() (Config){
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



    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)
    router.Handle("GET","/getchime", GetChime)
    router.Handle("POST","/getchime", GetChime)
    router.Handle("GET","/listchimes", ListChimes)
    router.Handle("GET","/testsend", TestSend)
    router.Handle("GET","/join", SubscribeNewDoorbell)
    router.Handle("GET","/syncnew", SyncNewDoorbell)


    //router.PUT("/putchime", PutChime)

    log.Fatal(http.ListenAndServe(":3434", router))
}
