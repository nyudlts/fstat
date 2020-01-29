package main

import (
	"context"
	"fmt"
	"github.com/google/go-tika/tika"
	"log"
	"os"
	"path/filepath"
)

var output_file *os.File
var cntxt = context.Background()
var tika_client *tika.Client

func createWriter () *os.File {
	output_file, err := os.Create("file-stats.csv")
	if err != nil { panic(err) }
	output_file.WriteString(fmt.Sprintf("'%s','%s','%s','%s','%s' \n", "name","mime","ext","size","path"))
	return output_file
}


func main() {
	log.Println("* Running FileStats")

	output_file = createWriter()
	defer output_file.Close()
	tika_client = tika.NewClient(nil, "http://localhost:9998")
	dir := os.Args[1]
	count := 0
	var byte_size float64 = 0.0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("* failure accessing a path %q: %v\n", path, err)
		} else {
			f, err := os.Open(path)
			if err != nil {
				log.Printf(err.Error())
			} else {
				if !(info.IsDir()) {
					count = count + 1
					size := info.Size()
					byte_size = byte_size + float64(size)
					mime, _ := tika_client.Detect(cntxt, f)
					name :=  info.Name()
					ext := filepath.Ext(name)
					output_file.WriteString(fmt.Sprintf("'%s','%s','%s','%d','%s'\n", name, mime, ext, size, path))
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("* error walking the path %q: %v\n", dir, err)
		return
	}

	bs_gb := ((byte_size / 1024.0) / 1024.0) / 1024.0
	output_file.WriteString(fmt.Sprintf("'','','',%.2f GB,''\n", bs_gb))

	log.Println("* Complete")
	log.Printf("* %d files scanned", count)




}


