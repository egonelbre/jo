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

	jopath = flag.String("jopath", os.Getenv("JOPATH"), "path to the packages root (default $GOPATH)")
)

func build(args []string) {
	output := "bundle.js"
	pkglist := []string{""}

	if len(args) > 0 {
		if len(args) > 1 {
			pkglist = args[:len(args)-1]
		}
		output = args[len(args)-1]
	}

	pkgs := packages.New(*jopath)
	if err := pkgs.Load(pkglist...); err != nil {
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

func main() {
	flag.Parse()

	if *jopath == "" {
		fmt.Println("JOPATH has not been defined, either add it as an argument or define it as an environment variable.")
		flag.Usage()
		os.Exit(1)
	}

	*jopath = filepath.FromSlash(*jopath)
	var err error
	*jopath, err = filepath.Abs(*jopath)
	if err != nil {
		log.Fatal(err)
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "build":
		build(flag.Args()[1:])
	default:
		fmt.Printf(`Unknown command "%v".`, cmd)
		flag.Usage()
		os.Exit(1)
	}
}
