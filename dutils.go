package dutils

import (
    "fmt"
    "os"

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
    fmt.Printf("files:%v\n",ls)
    for _, this_file := range ls {
        fmt.Printf("%s\n",this_file.Name())
    }

}
