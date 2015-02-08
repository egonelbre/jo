package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/egonelbre/jo/packages"
)

func Print(pkgs *packages.Packages) {
	if data, err := json.MarshalIndent(pkgs, "", "  "); err == nil {
		fmt.Println(string(data))
	}
}

var (
	verbose = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()

	root := flag.Arg(0)
	output := flag.Arg(1)

	if root == "" || output == "" {
		fmt.Println("USAGE:")
		fmt.Println("\tjo root/directory bundle.js")
		os.Exit(1)
	}

	root = filepath.FromSlash(root)
	root, _ = filepath.Abs(root)

	pkgs := packages.New(root)
	if err := pkgs.Load(); err != nil {
		log.Fatal(err)
	}

	if err := pkgs.Sort(); err != nil {
		log.Fatal(err)
	}

	if *verbose {
		Print(pkgs)
	}

	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}

	if err := pkgs.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}
