package main

import (
	"flag"
	"fmt"
	"github.com/leiyang23/nmon-parser/handler"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var listenPort *int64
var nmonFile, format *string

func webServer(port int64) {
	http.HandleFunc("/assets/{name}", handler.FileHandler)
	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/download", handler.ExportExcelHandler)
	http.HandleFunc("/", handler.MainHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	webSubCmd := flag.NewFlagSet("web", flag.ExitOnError)
	listenPort = webSubCmd.Int64("port", 8081, "default port 8081")

	exportSubCmd := flag.NewFlagSet("export", flag.ExitOnError)
	format = exportSubCmd.String("format", "xlsx", "export file format, supported formats xlsx/json")
	nmonFile = exportSubCmd.String("nmon", "", "nmon file")

	if len(os.Args) < 2 {
		fmt.Println("subcommand is required: web/export")
		os.Exit(1)
	}
	subCmd := os.Args[1]

	switch subCmd {
	case "web":
		err := webSubCmd.Parse(os.Args[2:])
		if err != nil {
			webSubCmd.PrintDefaults()
			panic("parse web flag failed")
		}
		fmt.Printf("Starting web server, visit http://localhost:%v\n", *listenPort)
		webServer(*listenPort)
	case "export":
		err := exportSubCmd.Parse(os.Args[2:])
		if err != nil {
			exportSubCmd.PrintDefaults()
			panic("parse export flag failed")
		}
		if len(*nmonFile) <= 0 {
			fmt.Println("nmon file is required")
			os.Exit(1)
		}
		fmt.Println("nmon file: ", *nmonFile, "export format: ", *format)
		fmt.Println("this feature has not yet been implemented -_-")
	default:
		fmt.Println("not support subcommand")
	}
}
