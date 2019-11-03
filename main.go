package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	version = flag.String("ver", "", "protoc version, for ex. 3.10.1")
	outDir  = flag.String("out", "", "output directory")
)

func main() {
	flag.Parse()
	if len(*version) == 0 {
		log.Fatal("Version is not set.\n")
	}
	if len(*outDir) == 0 {
		log.Fatal("Destination is not set.\n")
	}

	binFilePath := filepath.Join(*outDir, binFile)
	if _, err := os.Stat(binFilePath); !os.IsNotExist(err) {
		log.Printf("protoc is already existing (\"%s\").\n", binFilePath)
		return
	}

	archFile := fmt.Sprintf("protoc-%s-%s.zip", *version, platform)
	archFilePath := filepath.Join(*outDir, archFile)
	if _, err := os.Stat(binFilePath); os.IsNotExist(err) {

		url := fmt.Sprintf("https://github.com/protocolbuffers/protobuf/releases/download/v%s/%s", *version, archFile)
		log.Printf("Downloading protoc archive from \"%s\"...\n", url)

		if err = os.MkdirAll(filepath.Dir(archFilePath), os.ModePerm); err != nil {
			panic(err)
		}

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		out, err := os.Create(archFilePath)
		if err != nil {
			panic(err)
		}
		defer func() {
			out.Close()
			if err := os.Remove(archFilePath); err != nil {
				log.Printf("Failed to remove archive file: \"%s\".", err)
			}
		}()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Unpacking into \"%s\"...\n", binFilePath)
	arch, err := zip.OpenReader(archFilePath)
	if err != nil {
		panic(err)
	}
	defer arch.Close()

	binInsideArchPath := "bin/" + binFile
	for _, file := range arch.File {
		if file.Name != binInsideArchPath {
			continue
		}
		bin, err := os.OpenFile(binFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(err)
		}
		defer bin.Close()
		r, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer r.Close()
		_, err = io.Copy(bin, r)
		if err != nil {
			panic(err)
		}
		log.Println("Successfully installed.")
		return
	}

	log.Fatalf("Archive \"%s\" does not have protoc.\n", archFilePath)
}
