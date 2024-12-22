/*
zbm_pack converts popular image formats to Gamewave .zbm images.
*/
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/iafan/cwalk"
	"github.com/namgo/GameWaveFans/pkg/zbm"
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
	fmt.Println("Packs image to .zbm texture format used by Gamewave console")
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
			walkFunc := getWalkFunc(inputName)
			err := cwalk.Walk(inputName, walkFunc)
			if err != nil {
				fmt.Printf("Failed to unpack dir %s: %s\n", inputName, err)
				// for _, errors := range err.(cwalk.WalkerError).ErrorList {
				// 	fmt.Println(errors)
				// }
				failed = true
			}
		} else {
			if outputName == "" || len(args) > 1 {
				outputName = strings.TrimSuffix(inputName, filepath.Ext(inputName)) + ".zbm"
			}
			err := packTexture(inputName, outputName)
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

func getWalkFunc(basePath string) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			// check if file is an image
			// file deepcode ignore PT: This is CLI tool, this is intended to be traversable
			file, err := os.Open(filepath.Join(basePath, path))
			if err != nil {
				return err
			}
			_, format, err := image.DecodeConfig(file)
			if err != nil {
				return err
			}

			err = file.Close()
			if err != nil {
				return err
			}

			if err != nil && format != zbm.FormatName {
				outputName = filepath.Join(basePath, strings.TrimSuffix(path, filepath.Ext(path))+".zbm")
				return packTexture(filepath.Join(basePath, path), outputName)
			}
		}
		return nil
	}
}

func packTexture(inputName, outputName string) error {
	// file deepcode ignore PT: This is CLI tool, this is intended to be traversable
	file, err := os.Open(inputName)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %s", inputName, err)
	}
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("couldn't read image file config %s: %s", inputName, err)
	}

	fmt.Printf("Packing %s (detected %s): %dx%d\n", inputName, format, config.Width, config.Height)

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("couldn't seek in image file %s: %s", inputName, err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("couldn't read image file %s: %s", inputName, err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("couldn't close image file %s: %s", inputName, err)
	}

	outputFile, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("couldn't create output image file %s: %s", outputName, err)
	}

	err = zbm.Encode(outputFile, img)
	if err != nil {
		return fmt.Errorf("couldn't pack output image %s: %s", outputName, err)
	}

	err = outputFile.Close()
	if err != nil {
		return fmt.Errorf("couldn't close output image file %s: %s", outputName, err)
	}

	return nil
}
