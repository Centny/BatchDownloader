package main

import (
	"runtime"
	"strconv"
	"strings"
	"net/http"
	"fmt"
	"io"
	"os"
)
var pchain chan int=make(chan int)
var running int=0
func chandl(uri string,tpath string,fname string){
	if running>3{
		<-pchain
	}
	running++
	go download(uri,tpath,fname,&running,pchain)
	// fmt.Println("processing ",running)
}
func download(uri string,tpath string,fname string,ri *int,c chan int) {
	defer func(c chan int){
		running--
		c <- 1
	}(c)
	fmt.Println("downloading:",uri)
	os.MkdirAll(tpath,0777)
	out, err := os.Create(tpath+"/"+fname)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("download success:", tpath, n)
}
func loopArgs(uri string,tpath string,fname string,step int,args []string){
	if step>=len(args){
		chandl(uri,tpath,fname)
		return
	}
	vals:=strings.Split(args[step],"-")
	sval,eval,vsize:=vals[0],vals[1],0
	var err error
	if len(vals)>2{
		vsize,err=strconv.Atoi(vals[2])
		if err!=nil{
			panic(fmt.Sprintf("invalid length:%s",args[step]))
		}
	}
	if vsize<1{
		vsize=1
	}
	tstep:="\\"+strconv.Itoa(step+1)
	if sval[0]>='0' && sval[0]<='9'{
		var snum,enum int
		snum,err= strconv.Atoi(sval)
		if err !=nil{
			panic(fmt.Sprintf("invalid start value:%s",args[step]))
		}
		enum,err= strconv.Atoi(eval)
		if err !=nil{
			panic(fmt.Sprintf("invalid end value:%s",args[step]))
		}
		if enum<snum{
			panic(fmt.Sprintf("end less start value:%s",args[step]))
		}
			var fmts string
		if vsize>1{
			fmts="%0"+strconv.Itoa(vsize)+"d"
		}else{
			fmts="%d"
		}
		for i := snum; i < enum+1; i++ {
			rval:=fmt.Sprintf(fmts,i)
			uri_t:=strings.Replace(uri,tstep,rval,-1)
			tpath_t:=strings.Replace(tpath,tstep,rval,-1)
			fname_t:=strings.Replace(fname,tstep,rval,-1)
			loopArgs(uri_t,tpath_t,fname_t,step+1,args)
		}
	}else if (sval[0]>='a' && sval[0]<='z' && eval[0]>='a' && eval[0]<='z') || (sval[0]>='A' && sval[0]<='Z' && eval[0]>='A' && eval[0]<='Z'){
		count :=(eval[0]-sval[0])+1
		if count<1{
			panic(fmt.Sprintf("end less start value:%s",args[step]))
		}
		for i := (uint8)(0); i < count; i++ {
			rchar:=sval[0]+i
			var rval string=""
			for j := 0; j < vsize; j++ {
				rval+=string(rchar)
			}
			uri_t:=strings.Replace(uri,tstep,rval,-1)
			tpath_t:=strings.Replace(tpath,tstep,rval,-1)
			fname_t:=strings.Replace(fname,tstep,rval,-1)
			loopArgs(uri_t,tpath_t,fname_t,step+1,args)
		}
	}else{
		panic(fmt.Sprintf("invalid argument:%s",args[step]))
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	if len(os.Args)<4{
		fmt.Println("Usage:bdown <URL Pattern> <Folder Pattern> <File Name Pattern> [Batch number or char]")
		fmt.Println("  Example:")
		fmt.Println("    bdown \"http://github.com/\\1/\\2\" \"\\1\" \"\\2.html\" a-z:3 1-100:3")
		fmt.Println("    bdown \"http://github.com/\\1/\\2\" \"\\1\" \"\\2.html\" A-Z:3 1-100:3")
		return
	}
	loopArgs(os.Args[1],os.Args[2],os.Args[3],0,os.Args[4:])
	for ;running>0; {
		<-pchain
	}
}
