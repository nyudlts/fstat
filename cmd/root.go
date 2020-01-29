package cmd

import (
	"context"
	"github.com/google/go-tika/tika"
	"github.com/spf13/cobra"
	"fmt"
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
	rootCmd.PersistentFlags().StringVarP(&root_dir, "root-dir", "r", "","root directory to walk")
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

	fmt.Println("** running fstat")
}

func createOutputFile(output_file_name string) *os.File {
	of, err := os.Create(output_file_name)
	if err != nil { panic(err) }
	output_file.WriteString(fmt.Sprintf("'%s','%s','%s','%s','%s' \n", "name","mime","ext","size","path"))
	return of
}

func Walk() {

	err := filepath.Walk(root_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("* failure accessing a path %q: %v\n", path, err)
		} else {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("* error walking the path %q: %v\n", root_dir, err)
		return
	}
}

func ShutDown() {
	fmt.Println("** fstat shutting down")
	output_file.Close()
}
