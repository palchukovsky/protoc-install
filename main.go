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
	protocType = flag.String("type", "protoc",
		`type, "cli" for protoc or "grpc-web" for gRPC Web protoc plugin`)
	version = flag.String("ver", "", "version to install, for ex. 3.10.1")
	outDir  = flag.String("out", "", "output directory")
)

func main() {
	flag.Parse()
	if len(*version) == 0 {
		log.Fatal("Version is not set.\n")
	}
	if len(*outDir) == 0 {
		log.Fatal("Out is not set.\n")
	}
	if *protocType == "cli" {
		installProtoc()
		return
	}
	if *protocType == "grpc-web" {
		installGRPCWebPlugin()
		return
	}
	log.Fatal("Unknown type.\n")
}

func download(source, file string) string {
	path := filepath.Join(*outDir, file)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path
	}

	log.Printf("Downloading from \"%s\"...\n", source)

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		panic(err)
	}

	response, err := http.Get(source)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		panic(err)
	}

	return path
}

func remove(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Failed to remove \"%s\": \"%s\".", path, err)
	}
}

func installProtoc() {

	binFilePath := filepath.Join(*outDir, protocFile)
	if _, err := os.Stat(binFilePath); !os.IsNotExist(err) {
		log.Printf("protoc is already existing (\"%s\").\n", binFilePath)
		return
	}

	archFile := fmt.Sprintf("protoc-%s-%s.zip", *version, protocPlatform)
	url := fmt.Sprintf(
		"https://github.com/protocolbuffers/protobuf/releases/download/v%s/%s",
		*version, archFile)
	archFilePath := download(url, archFile)
	defer remove(archFilePath)

	log.Printf("Unpacking into \"%s\"...\n", binFilePath)
	arch, err := zip.OpenReader(archFilePath)
	if err != nil {
		panic(err)
	}
	defer arch.Close()

	binInsideArchPath := "bin/" + protocFile
	for _, file := range arch.File {
		if file.Name != binInsideArchPath {
			continue
		}
		bin, err := os.OpenFile(binFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(err)
		}
		r, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer r.Close()
		_, err = io.Copy(bin, r)
		if err != nil {
			panic(err)
		}
		bin.Close()
		err = os.Chmod(binFilePath, 0777)
		if err != nil {
			panic(err)
		}
		log.Println("protoc successfully installed.")
		return
	}

	log.Fatalf("Archive \"%s\" does not have protoc.\n", archFilePath)
}

func installGRPCWebPlugin() {
	fileName := "protoc-gen-grpc-web"
	url := fmt.Sprintf(
		"https://github.com/grpc/grpc-web/releases/download/%s/%s-%s-%s%s",
		*version, fileName, *version, grpcWebPlatform, grpcWebFileExt)
	binFilePath := download(url, fileName+grpcWebFileExt)
	err := os.Chmod(binFilePath, 0777)
	if err != nil {
		panic(err)
	}
	log.Println("gRPC Web protoc plugin installed.")
}
