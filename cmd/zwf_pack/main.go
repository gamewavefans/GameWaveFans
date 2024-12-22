/*
zwf_pack converts audio to Gamewave .zwf files.
*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/transforms"
	"github.com/go-audio/wav"
	"github.com/iafan/cwalk"
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
	fmt.Println("Packs audio to .zwf format used by Gamewave console")
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
				failed = true
			}
		} else {
			if outputName == "" || len(args) > 1 {
				outputName = strings.TrimSuffix(inputName, filepath.Ext(inputName)) + ".zwf"
			}
			err := packSound(inputName, outputName)
			if err != nil {
				fmt.Printf("Failed to unpack %s: %s\n", inputName, err)
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
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".wav" {
				outputName = filepath.Join(basePath, strings.TrimSuffix(path, filepath.Ext(path))+".zwf")
				return packSound(filepath.Join(basePath, path), outputName)
			}
		}
		return nil
	}
}

func packSound(inputName, outputName string) error {
	// file deepcode ignore PT: This is CLI tool, this is intended to be traversable
	file, err := os.Open(inputName)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %s", inputName, err)
	}
	decoder := wav.NewDecoder(file)
	if err != nil {
		return fmt.Errorf("couldn't create wav decoder for %s: %s", inputName, err)
	}
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return fmt.Errorf("couldn't get audio buffer %s: %s", inputName, err)
	}

	if buffer.SourceBitDepth != 16 {
		bufferFloat := buffer.AsFloatBuffer()
		transforms.Quantize(bufferFloat, 16)
		// return fmt.Errorf("Expected 16bit sample, got %d", buffer.SourceBitDepth)
		buffer = bufferFloat.AsIntBuffer()
	}

	if buffer.Format.NumChannels != 2 {
		if buffer.Format.NumChannels == 1 {
			// return fmt.Errorf("Expected stereo sound; got %d channels", buffer.Format.NumChannels)
			bufferFloat := buffer.AsFloat32Buffer()
			err = transforms.MonoToStereoF32(bufferFloat)
			if err != nil {
				return fmt.Errorf("Got mono sound, but couldn't convert to stereo: %s", err)
			}
			buffer = bufferFloat.AsIntBuffer()
		} else {
			return fmt.Errorf("Expected mono or stereo sound; got %d channels", buffer.Format.NumChannels)
		}
	}

	if buffer.Format.SampleRate != 22050 {
		return fmt.Errorf("Expected 22050Hz; got %dHz", buffer.Format.SampleRate)
	}

	fmt.Printf("Packing %s\n", inputName)

	err = file.Close()
	if err != nil {
		return fmt.Errorf("couldn't close audio file %s: %s", inputName, err)
	}

	outputFile, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("couldn't create output zwf file %s: %s", outputName, err)
	}

	err = zwf.Encode(outputFile, buffer)
	if err != nil {
		return fmt.Errorf("couldn't pack output image %s: %s", outputName, err)
	}

	err = outputFile.Close()
	if err != nil {
		return fmt.Errorf("couldn't close output zwf file %s: %s", outputName, err)
	}

	return nil
}
