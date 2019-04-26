package skp

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"skpconvfile/conver"
	logger "skpconvfile/log"
)

const (
	defaultLogPath = "~/log/skpconvfile.log"
)

/*
type SkpWorker struct {
	indir  string
	outdir string
}

func NewSkpWorker(in string, out string) *SkpWorker {
	return &SkpWorker{
		indir:  in,
		outdir: out,
	}
}
*/
func Run() {
	indir := flag.String("in", "", "in file directory")
	outdir := flag.String("out", "", "out file directory")
	logfile := flag.String("log", "", "log file")
	flag.Parse()
	if *indir == "" {
		panic("in file directory is NULL")
	}
	if *outdir == "" {
		panic("out file directory is NULL")
	}
	if *logfile == "" {
		*logfile = defaultLogPath
	}
	logfp, err := os.OpenFile(*logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		fmt.Println(err)
	}
	defer logfp.Close()
	logger.MyLoger = log.New(logfp, "SKP CONVER FILE:", log.LstdFlags)
	filepath.Walk(*indir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		logger.Println("-------Begin Proc[" + path + "]-------")
		logger.Println("infile:", path)
		outfile := *outdir + "/" + f.Name()
		logger.Println("outfile:", outfile)
		conver := conver.NewConverFile(path, outfile)
		if err := conver.Conver(); err != nil {
			logger.Println(err)
		}
		logger.Println("-------End Proc[" + outfile + "]-------")
		return nil
	})
}
