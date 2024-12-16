/*
zwf_unpack converts Gamewave .zwf sounds to one of the more popular formats.

This program can output wav files.
*/
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/namgo/GameWaveFans/pkg/zwf"
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
	fmt.Println("Usage: zwf_unpack [-o output_dir] <input_file/input_dir>")
	fmt.Println("Unpacks sounds from .zwf audio format used by the Gamewave console")
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
				outputName = strings.TrimSuffix(inputName, filepath.Ext(inputName)) + ".wav"
			}
			err := unpackSound(inputName, outputName)
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
	return func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			if strings.ToLower(filepath.Ext(path)) == ".zwf" {
				outputName = strings.TrimSuffix(path, filepath.Ext(path)) + ".wav"
				err := unpackSound(path, outputName)
				if err != nil {
					fail := true
					failed = &fail
				}
			}
		}
		return nil
	}
}

func unpackSound(inputName, outputName string) error {
	// file deepcode ignore PT: This is CLI tool, this is intended to be traversable
	file, err := os.Open(inputName)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %s", inputName, err)
	}

	buf, err := zwf.Decode(file)
	if err != nil {
		return fmt.Errorf("couldn't parse audio file %s: %s", inputName, err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("couldn't close audio file %s: %s", inputName, err)
	}

	outputFile, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("couldn't create output image file %s: %s", outputName, err)
	}

	ext := filepath.Ext(strings.ToLower(outputName))
	switch ext {
	case ".wav":
		fallthrough
	case ".wave":
		err = writeWave(outputFile, buf)
	default:
		err = fmt.Errorf("unknown output format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("couldn't pack output audio %s: %s", outputName, err)
	}

	err = outputFile.Close()
	if err != nil {
		return fmt.Errorf("couldn't close audio file %s: %s", outputName, err)
	}

	return nil
}

func writeWave(outputFile io.WriteSeeker, buffer *audio.IntBuffer) error {
	var err error
	enc := wav.NewEncoder(outputFile, buffer.Format.SampleRate, buffer.SourceBitDepth, buffer.Format.NumChannels, 1)
	if errTmp := enc.Write(buffer); err != nil {
		err = fmt.Errorf("could not write %s: %s", outputName, errTmp)
	}
	if errTmp := enc.Close(); err != nil {
		err = fmt.Errorf("could not close wav encoder: %s", errTmp)
	}
	return err
}
