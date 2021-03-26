package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dacort/forklift/internal/forklift"
)

var templateURI string

func getTemplateURI() string {
	// First we check if there is an environment variable `FORKLIFT_URI`
	// If there is we use that.
	val, ok := os.LookupEnv("FORKLIFT_URI")
	if ok {
		return val
	}

	// Second, we check for a `-w` command-line argument.
	flag.Parse()
	return templateURI

	// If we don't find either, the main routine will just create a passthrough writer
}

func disableLogging() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func runCmd() io.Reader {
	argsWithoutProg := os.Args[1:]
	cmd := exec.Command(argsWithoutProg[0], argsWithoutProg[1:]...)
	// var stdout bytes.Buffer
	out, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	return out
}

func getInputFromPipeOrCmd() io.Reader {
	if isInputFromPipe() {
		return os.Stdin
	}

	return runCmd()
}

func main() {
	var writer forklift.Destination
	u := getTemplateURI()

	reader := bufio.NewScanner(getInputFromPipeOrCmd())

	if strings.HasPrefix(u, "s3") {
		writer = forklift.NewS3Uploader(getTemplateURI()) // ("s3://dacort-east/tmp/{{json \"name\"}}/data.json")
	} else {
		// On stdout we just want to dump the data
		disableLogging()
		writer = forklift.NewStdOutPassthrough()
	}

	for reader.Scan() {
		record := forklift.JSONRecord{Raw: reader.Text()}
		writer.AddRecord(record)
	}

	writer.Close()
}

func init() {
	flag.StringVar(&templateURI, "w", "", "URI template for the output location")

}
