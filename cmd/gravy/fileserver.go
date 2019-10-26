package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d-sparks/gravy/paramfileserver"
	"github.com/spf13/cobra"
)

var fileserverCmd = &cobra.Command{
	Use:   "fileserver",
	Short: "fileserver for params",
	Run:   fileserverFn,
}

var port int
var folder string

func init() {
	rootCmd.AddCommand(fileserverCmd)
	fileserverCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to serve on")
	fileserverCmd.Flags().StringVarP(&folder, "folder", "f", "./data/mock/alphavantage/", "Folder to serve")
}

func fileserverFn(cmd *cobra.Command, args []string) {
	pfs := paramfileserver.ParamFileServer{Folder: folder}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), pfs); err != nil {
		log.Println(err.Error())
	}
}
