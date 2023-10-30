/*
zbc_unpack unpacks Gamewave .zbc files into .zbc_unpacked files
*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/namgo/GameWaveFans/pkg/zbc"
	"github.com/spf13/pflag"
)

// flags
var (
	outputName string
)

func parseFlags() {
	pflag.StringVarP(&outputName, "output", "o", "", "name of the output file")
	pflag.Parse()
}

func usage() {
	fmt.Println("Unpacks .zbc format used by Gamewave console")
	fmt.Println("Flags:")
	pflag.PrintDefaults()
}

func main() {
	failed := false
	parseFlags()
	args := pflag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	if outputName != "" && len(args) > 1 {
		fmt.Println("Output name can only be used with one input file")
		usage()
		os.Exit(1)
	}

	for _, inputName := range args {
		f, err := os.Stat(inputName)
		if err != nil {
			fmt.Printf("Failed to get info about %s: %s", inputName, err)
			failed = true
		}
		if f.IsDir() {
			walkFunc := getWalkFunc(&failed)
			err := filepath.Walk(inputName, walkFunc)
			if err != nil {
				fmt.Printf("Failed to unpack dir %s: %s", inputName, err)
				failed = true
			}
		} else {
			if outputName == "" || len(args) > 1 {
				outputName = strings.TrimSuffix(inputName, filepath.Ext(inputName)) + ".zbc_unpacked"
			}
			err := unpackBytecode(inputName, outputName)
			if err != nil {
				fmt.Printf("Failed to unpack %s: %s", inputName, err)
				failed = true
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}

func getWalkFunc(failed *bool) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.ToLower(filepath.Ext(path)) == ".zbc" {
				outputName = strings.TrimSuffix(path, filepath.Ext(path)) + ".zbc_unpacked"
				err := unpackBytecode(path, outputName)
				if err != nil {
					fail := true
					failed = &fail
				}
			}
		}
		return nil
	}
}

func unpackBytecode(inputName, outputName string) error {
	// file deepcode ignore PT: This is CLI tool, this is intended to be traversable
	file, err := os.Open(inputName)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %s", inputName, err)
	}

	packed, err := zbc.IsPacked(file)
	if err != nil {
		return err
	}

	if !packed {
		// file is not packed, skip it
		fmt.Printf("Skipping  %s: is already unpacked\n", inputName)
		return nil
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	fmt.Printf("Unpacking %s\n", inputName)
	unpacked, err := zbc.Unpack(file)
	if err != nil {
		return fmt.Errorf("couldn't parse input file %s: %s", inputName, err)
	}

	outputFile, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("couldn't create output file %s: %s", outputName, err)
	}

	_, err = outputFile.Write(unpacked)
	if err != nil {
		return fmt.Errorf("couldn't write to output file %s: %s", outputName, err)
	}

	err = outputFile.Close()
	if err != nil {
		return fmt.Errorf("couldn't close output file %s: %s", outputName, err)
	}

	return nil
}
