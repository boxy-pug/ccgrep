package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	filePointers    []*os.File
	filepaths       []string
	args            []string
	regex           *regexp.Regexp
	recurseDir      bool // -r
	invertSearch    bool // -v
	caseInsensitive bool // -i
	dirToSearch     string
}

func main() {

	cfg := loadConfig()
	fmt.Printf("cfg loaded, args: %v\n", cfg.args)

	fmt.Println("Printing filepaths:")
	for _, file := range cfg.filepaths {
		fmt.Printf("%s\n", file)
	}

	fmt.Printf("Regex loaded: %v\n", cfg.regex)

	processLines(cfg)

}

func loadConfig() Config {
	var err error
	var cfg Config
	flag.BoolVar(&cfg.recurseDir, "r", false, "recurse directory tree")
	flag.BoolVar(&cfg.invertSearch, "v", false, "invert search")
	flag.BoolVar(&cfg.caseInsensitive, "i", false, "case insensitive")
	flag.Parse()

	cfg.args = flag.Args()

	regex := cfg.args[0]
	if cfg.caseInsensitive {
		regex = "(?i)" + regex
	}

	fmt.Println(regex)

	cfg.regex, err = regexp.Compile(regex)
	if err != nil {
		fmt.Println("error parsing regex: ", err)
	}

	if cfg.recurseDir {
		if len(cfg.args) == 1 {
			cfg.dirToSearch = "."
		} else {
			cfg.dirToSearch = cfg.args[1]
		}
		fmt.Printf("cfgdir to search is: %v\n", cfg.dirToSearch)
		recurseDirFetchFilepaths(&cfg)
	} else {
		cfg.filepaths = append(cfg.filepaths, cfg.args[1])
	}

	return cfg

}

func recurseDirFetchFilepaths(cfg *Config) *Config {
	err := filepath.Walk(cfg.dirToSearch, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fmt.Println("Checking path ", path)

		if !info.IsDir() {
			fmt.Printf("Appending: %v\n", path)
			cfg.filepaths = append(cfg.filepaths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return cfg
	}
	return cfg
}

func processLines(cfg Config) {
	var matchesFound int
	for _, filepath := range cfg.filepaths {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("error opening filepath as file: %v", err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if cfg.regex.MatchString(line) {
				matchesFound++
				if !cfg.invertSearch {
					if len(cfg.filepaths) > 1 {
						fmt.Printf("%s: ", filepath)
					}
					fmt.Println(line)
				}

			} else {
				if cfg.invertSearch {
					matchesFound++
					if len(cfg.filepaths) > 1 {
						fmt.Printf("%s: ", filepath)
					}
					fmt.Println(line)
				}
			}
		}
	}
	if matchesFound == 0 {
		fmt.Println("No matches found")
		os.Exit(1)
	}
}
