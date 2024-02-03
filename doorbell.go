package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	//"bufio"
)

const (
	//MP3path = "/home/zane/programming/go/webserver/thesounds/"
	//MP3path = "/home/zane/programming/go/webserver/clientsounds/"
	MP3subpath = "thesounds/"
)

func Ls(dirpath string) []string {
	dir, err := os.Open(dirpath)
	if err != nil {
		fmt.Printf("opening erro+r:%v\n", err)
	}
	ls, lerr := dir.Readdir(0)
	if lerr != nil {
		fmt.Printf("listing error:%v\n", lerr)
	}
	//fmt.Printf("files:%v\n",ls)

	ret := make([]string, len(ls))
	i := 0
	for _, this_file := range ls {
		//fmt.Printf("%s\n",this_file.Name())
		if strings.HasSuffix(this_file.Name(), ".mp3") {
			ret[i] = this_file.Name()
			i++
		}

	}
	// we only want to return the slice with the filenames in it
	return ret[0:i]
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pagedata := PageData{Token: "Big Secret Token",
		Serverstr: "192.168.178.20"}
	t, _ := template.ParseFiles("doorbell_index.gtpl")
	t.Execute(w, pagedata)
	//fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Ring(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Fprintf(w, "ringing\n")

	files := Ls(MP3subpath)

	cmd := exec.Command(MP3player, fmt.Sprintf("%s%s", MP3path, files[0]))

	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	// redirect back to the main page
	http.Redirect(w, r, "/", 301)
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
	cmd := exec.Command("ls") // need to init this var at this scope

	if r.Method == "GET" {

		filename := ps.ByName("name")

		if strings.Index(filename, "!") > -1 {
			// then the file we want played is not in the list but is special
			// we will let central handle the sub dirs but only play files under
			// the main mp3 dir so assume the ! = / char is for under MP3path
			pathfilename := strings.Replace(filename, "!", "/", -1)
			cmd = exec.Command(MP3player, fmt.Sprintf("%s%s", MP3path, pathfilename))

		} else {
			// standard door chime to be run and must be in the list
			filelist := Ls(MP3path)
			if contains(filelist, filename) {
				cmd = exec.Command(MP3player, fmt.Sprintf("%s%s", MP3path, filename))
				//cmd.Stdin = strings.NewReader("some input")
			} else {
				log.Fatal("Can't find chime")
				return
			}
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

}

type PageData struct {
	Serverstr string
	Token     string
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func PutChime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))

		pagedata := PageData{Token: "Big Secret Token",
			Serverstr: "192.168.178.20"}
		//token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("mp3upload.gtpl")
		t.Execute(w, pagedata)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(fmt.Sprintf("FormFile error:%s", err))
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		filepath := r.FormValue("path")
		if filepath == "" {
			f, err := os.OpenFile(MP3path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error opening file:%s", err))
				return
			}
			io.Copy(f, file)
			fmt.Printf("Got file:%s\n", handler.Filename)
			defer f.Close()
		} else {
			// then a path has been given so we want to store the file
			// in a subdirectory
			thepath := MP3path + filepath
			// if the path does not lready exist then just create it
			CreateDirIfNotExist(thepath)

			f, err := os.OpenFile(thepath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error opening file:%s", err))
				return
			}
			io.Copy(f, file)
			fmt.Printf("Got file:%s\n", handler.Filename)
			defer f.Close()

		}

	}
}

//	func PutChime(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	    fmt.Fprintln(rw, "post update")
//	}
type Config struct {
	Doorbell_dir   string
	Satellite_port int
}

var CONFIG Config
var MP3path string
var MP3player string

func GetConfig() Config {
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
	fmt.Println("DIR:", CONFIG.Doorbell_dir)
	fmt.Println("Port:", CONFIG.Satellite_port)
	fmt.Println("Player:", CONFIG.mp3player)

	MP3path = CONFIG.Doorbell_dir + "/" + MP3subpath
	MP3player = CONFIG.mp3player

	// no cmd parms so we just run normally
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/ring", Ring)
	router.Handle("GET", "/putchime", PutChime)
	router.Handle("POST", "/putchime", PutChime)
	router.GET("/ringchime/:name", RingChime)
	//router.GET("/playspecial/path=:path", playspecial)
	//router.PUT("/putchime", PutChime)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CONFIG.Satellite_port), router))

}
