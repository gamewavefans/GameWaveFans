/*
cheese_unpack unpacks files from the end of .bin binary files.
*/
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/namgo/GameWaveFans/pkg/common"
	"github.com/spf13/pflag"
)

// flags
var (
	outputDir string
)

func parseFlags() {
	pflag.StringVarP(&outputDir, "output", "o", "", "name of the output folder")
	pflag.Parse()
}

func usage() {
	fmt.Println("Unpacks data found at the end of .bin files.")
	fmt.Println("This tool is not supported, as it's not meant for the end users.")
	fmt.Println("Flags:")
	pflag.PrintDefaults()
}
func main() {
	parseFlags()
	args := pflag.Args()
	if len(args) != 1 {
		usage()
		os.Exit(1)
	}
	inputName := args[0]
	f, err := os.Open(inputName)
	defer closeFile(f)
	if err != nil {
		fmt.Printf("Failed to get info about %s: %s", inputName, err)
		os.Exit(1)
	}

	cheeseAddress, err := findCheese(f)
	if err != nil {
		fmt.Printf("Failed to find cheese in %s: %s", inputName, err)
		os.Exit(1)
	}
	if cheeseAddress < 0 {
		fmt.Printf("Failed to find built-in files. is this the correct file?")
		os.Exit(1)
	}
	fmt.Printf("Found cheese at 0x%x, digging in\n", cheeseAddress)
	cheeseBlock, err := parseCheese(f, cheeseAddress)
	if err != nil {
		fmt.Printf("Failed to parse cheese in %s: %s", inputName, err)
		os.Exit(1)
	}
	fmt.Printf("Found %d pieces of cheese:\n", len(cheeseBlock))

	if outputDir != "" {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to creat output dir: %s", err)
			os.Exit(1)
		}
	}
	for _, piece := range cheeseBlock {
		fmt.Printf("* %s %d\n", piece.Name, len(piece.Data))
		outName := path.Join(outputDir, piece.Name)
		err = os.WriteFile(outName, piece.Data, 0660)
		if err != nil {
			fmt.Printf("Failed to save piece of cheese: %s", err)
			os.Exit(1)
		}
	}
}

func closeFile(f *os.File) func() {
	return func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Could not close file: %s", err)
			os.Exit(1)
		}
	}
}

// CheeseFile is one parsed file from the block
type CheeseFile struct {
	Name string
	Data []byte
}

// Cheese is a generic file container at the end of the binary files
type Cheese = []CheeseFile

func findCheese(f *os.File) (int, error) {
	buf, err := io.ReadAll(f)
	if err != nil {
		return -1, err
	}
	cheesseLoc := bytes.Index(buf, []byte{0x12, 0x34, 0x56, 0x78, 0x87, 0x65, 0x43, 0x21})
	return cheesseLoc, nil
}

func parseCheese(f *os.File, cheeseIndex int) (Cheese, error) {
	fileCount, err := common.ReadUint32Big(f, int64(cheeseIndex)+8)
	if err != nil {
		return nil, err
	}
	cheese := make(Cheese, 0)
	for i := 0; i < int(fileCount); i++ {
		nameBuffer := make([]byte, 40)
		_, err = f.Read(nameBuffer)
		if err != nil {
			return nil, err
		}
		n := bytes.IndexByte(nameBuffer[:], 0)

		intBuffer := make([]byte, 4)
		_, err = f.Read(intBuffer)
		if err != nil {
			return nil, err
		}
		address := binary.BigEndian.Uint32(intBuffer)

		_, err = f.Read(intBuffer)
		if err != nil {
			return nil, err
		}
		size := binary.BigEndian.Uint32(intBuffer)

		currentPos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		_, err = f.Seek(int64(cheeseIndex+int(address)), io.SeekStart)
		if err != nil {
			return nil, err
		}

		dataBuffer := make([]byte, size)
		_, err = f.Read(dataBuffer)
		if err != nil {
			return nil, err
		}

		_, err = f.Seek(currentPos, io.SeekStart)
		if err != nil {
			return nil, err
		}

		piece := CheeseFile{Name: string(nameBuffer[:n]), Data: dataBuffer}
		cheese = append(cheese, piece)
	}
	return cheese, nil
}
