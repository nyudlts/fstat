package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-tika/tika"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var (
	cntxt = context.Background()
	server_uri string
	tika_client *tika.Client
	root_dir string
	output_file_name string
	output_file *os.File
	log_file string
	rootCmd = &cobra.Command{
		Use:   "fstat",
		Short: "A csv generator that provides basic information about files in a file system",
		Run: func(cmd *cobra.Command, args []string) { },
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&root_dir, "dir", "d", "","root directory to walk")
	rootCmd.PersistentFlags().StringVarP(&server_uri, "tika-url", "t", "","url of tika server")
	rootCmd.PersistentFlags().StringVarP(&output_file_name, "output-file", "o", "","path to file to output to")
	rootCmd.PersistentFlags().StringVarP(&log_file, "log-file", "l", "", "path to log file")
}

func initConfig() {
	if root_dir == "" {
		fmt.Println("Location of directory to walk is required")
		os.Exit(1)
	}

	if server_uri == "" {
		server_uri = "http://localhost:9998"
	}

	tika_client = tika.NewClient(nil, server_uri)

	if output_file_name == "" {
		output_file_name = "output.csv"
	}

	output_file = createOutputFile(output_file_name)


	if log_file == "" {
		log_file = "fstat.log"
	}

}

func createOutputFile(output_file_name string) *os.File {
	of, err := os.Create(output_file_name)
	if err != nil {
		log.Printf("* error creating output file %s: %v\n", output_file_name, err)
	}
	output_file.WriteString(fmt.Sprintf("'%s','%s','%s','%s','%s' \n", "name","mime","ext","size","path"))
	return of
}

func Walk() {
	log.Println("* running fstat")
	count := 0
	var byte_size float64 = 0.0

	err := filepath.Walk(root_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("* failure accessing a path %q: %v\n", path, err)
		} else {
			f, err := os.Open(path)
			defer f.Close()
			if err != nil {
				log.Printf("* error walking the path %q: %v\n", root_dir, err)
			}else {
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
		log.Printf("* error walking the path %q: %v\n", root_dir, err)
		return
	}

	bs_gb := ((byte_size / 1024.0) / 1024.0) / 1024.0
	output_file.WriteString(fmt.Sprintf("'','','',%.2f GB,''\n", bs_gb))

	log.Println("* fstat complete")
	log.Printf("* %d files scanned\n", count)
}

func ShutDown() {
	log.Println("* fstat shutting down")
	output_file.Close()
	os.Exit(0)
}
