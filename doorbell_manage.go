package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "math/rand"
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
    "bufio"
    //"github.com/stianeikeland/go-rpio"
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


func ListStoredChimes(rawdirpath string) []string {
    dirslice := strings.Split(rawdirpath, "!")
    //var isrecursed bool
    var dirpath string

    switch lenslice:=len(dirslice);lenslice{
    case 1:
        //isrecursed = false
        dirpath = rawdirpath
    case 2:
        //isrecursed = true
        dirpath = dirslice[0] + "/" + dirslice[1]
    case 3,4,5:
        //isrecursed = true
        dirpath = strings.Join(dirslice, "/")
    default:
        reterr := make([]string,2)
        reterr[0] = "filepath error"
        reterr[1] = "Problem splitting path info"
        return reterr

    }

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
        if strings.HasSuffix(this_file.Name(), ".mp3") {
            ret[i]=this_file.Name()
            i++
        }
        if this_file.IsDir() {
            newrawpath := rawdirpath + "!"+ this_file.Name()
            //@@@@ !!!! RECURSIVE CALL !!!! @@@@
            subdirlist := ListStoredChimes(newrawpath)
            //fmt.Printf("%s\n",this_file.Name())

            sublist := make([]string,len(subdirlist))
            //XXXXXXXXXXX need to prepend dir to each filename

            for cc, this_sub_file := range subdirlist{
                sublist[cc] = newrawpath + this_sub_file
                //cc += 1
            }
            ret = append(ret, sublist...)
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


func SendChimeToAll(file string)(err error){
    /*
    Function to send the passed chime to all satellite doorbells
    Will return the last error or nil.
    */
    var this_err error
    for _,this_doorbell := range(SubscribedDoorbells){
        this_err = SendChime(file, fmt.Sprintf("http://%s:%d/putchime",this_doorbell,CONFIG.Satellite_port))
        if this_err != nil {
            err = this_err
        }
    }

    return err
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
    /*
    So user can send a chime to the doorbell server
    via POST
    */
       fmt.Println("method:", r.Method)
       pagedata := PageData{Token:"Big Secret Token",
                            Serverstr:"192.168.178.20",
                            Reply:""}

        reply:=""
        filepath := ""

        lsdir:=ListStoredChimes(MP3path)
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
           //defer file.Close()
           //fmt.Fprintf(w, "%v", handler.Header)
           filepath = MP3path+handler.Filename
           f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }

           // only want the file if it is an MP3
           if strings.HasSuffix(handler.Filename, ".mp3") {
               io.Copy(f, file)
               reply = fmt.Sprintf("Got file called:%s\n",handler.Filename)
               fmt.Printf("%s",reply)
               f.Close()
               file.Close()
               // now send it out to all satellite doorbells
               err := SendChimeToAll(filepath)
               if err != nil{
                   fmt.Printf("At least one error sending the chime to all:%s\n",err)
               }
           }else{
               reply = fmt.Sprintf("File called:%s does not appear to be an MP3\n",handler.Filename)
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
    /*
    List all of the chimes on the server.
    */
       fmt.Println("method:", r.Method)

       lsdir:=ListStoredChimes(MP3path)
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
    /*
    Write the IP of the new doorbell to the dorbells.txt file
    And send back a success message.
    */
    ip_chunks := strings.Split(req.RemoteAddr,":")
    fmt.Println("IP:",ip_chunks[0])

    f, err := os.OpenFile(Doorbells_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    new_ip := fmt.Sprintf("%s\n",ip_chunks[0])
    new_raw_ip := fmt.Sprintf("%s",ip_chunks[0])

    doesnt_exist := true
    for _,this_ip := range(SubscribedDoorbells){

        fmt.Printf("newip=%s|this_ip=%s\n",new_raw_ip,this_ip)
        if this_ip == new_raw_ip {
            doesnt_exist = false
        }
    }

    if doesnt_exist{
        // write the new ip to the file
        if _, err = f.WriteString(new_ip); err != nil {
            panic(err)
        }
        // put it into the subscribed list
        SubscribedDoorbells = append(SubscribedDoorbells,new_raw_ip)
        fmt.Fprintf(w,"Join Success")
    }else{
        fmt.Fprintf(w,"Doorbell already subscribed")
    }


}

func SyncDoorbell(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    /*
    Send all of the current chimes to the passed doorbell
    */
    ip_chunks := strings.Split(req.RemoteAddr,":")
    the_server := ip_chunks[0]
    fmt.Println("sync IP:",ip_chunks[0])
    //
    // fmt.Fprintf(w,"Join Success")
    //

    // wait for doorbell to sort it's shit out
    //time.Sleep(1000 * time.Millisecond)

    // now we need to ensure that the doorbell has all the files it needs
    lsdir:=ListStoredChimes(MP3path)
    //fmt.Printf("the files:%s", lsdir)

    the_url := fmt.Sprintf("http://%s:%d/putchime",the_server,CONFIG.Satellite_port)
    //the_url = fmt.Sprintf("http://localhost:%s/putchime",CONFIG.Satellite_port)
    the_path := ""
    fmt.Printf("the_url:%s\n", the_url)

    for _, this_file := range lsdir {
        the_path = MP3path + this_file
        fmt.Printf("the_path:%s\n", the_path)

        SendChime( the_path, the_url)
    }
}

func WebRingAllDoorbells(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
    /*
    To tell all of the subscribed doorbells to ring after
    choosing which chime randomly from the dir
    */
    fmt.Printf("About to ring all doorbells from web\n")
    reply := RingAllDoorbells()

    fmt.Printf("Rang all doorbells from web with reply:%s\n",reply)



    // redirect back to the main page
    http.Redirect(w,r,"/",302)

    // pagedata := PageData{Token:"Big Secret Token",
    //                     Reply:reply}
    //
    //
    // t, _ := template.ParseFiles("index_central.gtpl")
    // t.Execute(w,pagedata)
}

/***************************************************************
  DEFINE the special sounds that we may want to Run
  *********************************************************/


func RingSpecial(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    special_sounds := map[string]string{
        "face_detected": "special!facedetected01.mp3",
        "test": "data1.mp3",
    }


    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("specific"))
    specific_name := ps.ByName("specific")

    sound_path := special_sounds[specific_name]
    fmt.Fprintf(w, "soundpath=%s\n",sound_path)

    for _,this_doorbell := range(SubscribedDoorbells){
        go RingADoorbell(sound_path,this_doorbell)
    }
}




//*************************************************************************************
//*************************************************************************************


func RingAllDoorbells () string {
    /* choose a chime randomly then send the ring command to all doorbells
    in the subscribed doorbell list
    */
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



//*************************************************************************************
//*************************************************************************************


func PlayAllDoorbells (filepath string) string {
    // play a specific sound on all doorbells in the subscribed doorbell list
    reply:= ""
    lsdir:=ListStoredChimes(MP3path)
    //fmt.Printf("%s\n",lsdir)
    if len(lsdir) < 1 {
       reply = "No files to choose from"
    }else{

        reply = fmt.Sprintf("Playing: %s  on all doorbells", filepath)
        fmt.Printf("%s\n",reply)


        for _,this_doorbell := range(SubscribedDoorbells){
            go RingADoorbell(filepath,this_doorbell)
        }
    }
    return reply
}


//*************************************************************************************
//*************************************************************************************


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
    fmt.Println("Port:3434")


    MP3path = CONFIG.Doorbell_dir + "/" + MP3subpath
    SubscribedDoorbells = getSubscribedDoorbells()


    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)
    router.Handle("GET","/getchime", GetChime)
    router.Handle("POST","/getchime", GetChime)
    router.Handle("GET","/listchimes", ListChimes)
    router.Handle("GET","/testsend", TestSend)
    router.Handle("GET","/join", SubscribeNewDoorbell)
    router.Handle("GET","/sync", SyncDoorbell)
    router.Handle("GET","/RingAllDoorbells",WebRingAllDoorbells)
    router.Handle("GET","/special/:specific", RingSpecial)


    log.Fatal(http.ListenAndServe(":3434", router))
}
