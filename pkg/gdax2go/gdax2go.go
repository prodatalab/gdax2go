package gdax2go

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/onrik/logrus/filename"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Value blah
type Value struct {
	RootPath string
}

const (
	pkgName = "gdax"
)

var (
	// Val blah
	Val = Value{}
	log = logrus.New()
	err error
)

// Run blah
func Run() {
	log.AddHook(filename.NewHook())
	log.Level = logrus.DebugLevel
	os.Chdir(Val.RootPath)
	err = filepath.Walk(Val.RootPath, func(path string, info os.FileInfo, err error) error {
		fmt.Println()
		if filepath.Ext(path) == ".go" || info.Name() == "odir" {
			return nil
		}
		log.Debugf("STARTING work on: %s\n", path)
		if err != nil {
			return errors.Wrapf(err, "Walk failed for path: %s\n", path)
		}
		// if its a dir, chdir and exit
		if info.IsDir() {
			os.Chdir(path)
			log.Debugf("changing directory: %s\n", path)
			return nil
		}
		// get new file name and struct name
		var structName, nname string
		fname := info.Name()
		log.Debugf("working on file: %s\n", fname)
		fname = strings.Split(fname, ".")[0]
		nname = fname + ".go"
		log.Debugf("the nname is: %s\n", nname)
		if !strings.Contains(fname, "-") {
			structName = strings.Title(fname)
		}
		snameList := strings.Split(fname, "-")
		for _, s := range snameList {
			structName += strings.Title(s)
		}
		// call gojson with new information

		// gojson -name=SubscribeReq -pkg=gdax -o=tmp.go
		// -input=/home/prodatalab/go/src/github.com/prodatalab/msg/gdax/webreturnreturnreturnsocket/websocket-subscribe-req.example.json
		// &&  addzid tmp.go

		// var out bytes.Buffer
		// var curdir string
		cmd := exec.Command("gojson", "-input="+info.Name(), "-name="+structName, "-o="+nname, "-pkg=gdax", "-subStruct")
		cmd.Stdout = os.Stdout
		// fmt.Printf("stdout of gojson:\n %s\n", out.String())
		log.Debugf("executing gojson on %s to create %s\n", info.Name(), nname)
		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "cmd.Run failed")
		}
		// curdir, err = os.Getwd()
		if err != nil {
			return errors.Wrap(err, "Getwd failed")
		}
		cmd = exec.Command("addzid", "-p=gdax", nname)
		cmd.Stdout = os.Stdout
		// fmt.Printf("stdout of addzid:\n %s\n", out.String())
		log.Debugf("executing addzid on %s\n", nname)
		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "cmd.Run failed")
		}
		log.Debugf("COMPLETED work for %s\n", info.Name())
		return nil
	})
}
