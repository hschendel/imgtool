package main

import (
	"image"
	"flag"
	"fmt"
	"os"
	"image/png"
	"image/jpeg"
	"image/gif"
)

type SpriteMap interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

func imageEmpty(img image.Image) bool {
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0 {
				return false
			}
		}
	}
	return true
}

type command interface {
	Execute(image.Image) image.Image
	Explanation() string
}

type args struct {
	InputFilename    string
	OutputFilename   string
	Command string
}

var commands = map[string]command {
	"flip-x": &flipX{},
	"flip-y": &flipY{},
}

func (a *args) commandValid() bool {
	_, found := commands[a.Command]
	return found
}

func (a *args) parse() bool {
	flag.StringVar(&a.InputFilename, "i", "", "Input filename, stdin if left empty")
	flag.StringVar(&a.OutputFilename, "o", "", "Output filename, stdout if left empty")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s executes one of the following commands on the input and writes the result to the output\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr,"  Commands are:")
		for commandName, command := range commands {
			fmt.Fprintf(os.Stderr,"    %s\t%s\n", commandName, command.Explanation())
		}
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return false
	}
	a.Command = flag.Arg(0)
	if !a.commandValid() {
		fmt.Fprintln(os.Stderr, "Invalid command", a.Command)
		return false
	}

	return true
}

func saveImage(img image.Image, format string, file *os.File) {
	var encodeErr error
	switch format {
	case "jpeg":
		encodeErr = jpeg.Encode(file, img, nil)
	case "gif":
		encodeErr = gif.Encode(file, img, nil)
	default:
		encodeErr = png.Encode(file, img)
	}
	file.Close()
	if encodeErr != nil {
		fmt.Fprintln(os.Stderr, "Cannot encode image:", encodeErr)
		os.Exit(5)
	}
}

func main() {
	var args args
	if !args.parse() {
		os.Exit(1)
	}

	inFile := os.Stdin
	if args.InputFilename != "" {
		var openErr error
		inFile, openErr = os.Open(args.InputFilename)
		if openErr != nil {
			fmt.Fprintln(os.Stderr, "Cannot open", args.InputFilename + ":", openErr)
			os.Exit(2)
		}
		defer inFile.Close()
	}


	img, imageFormat, decodeErr := image.Decode(inFile)
	if decodeErr != nil {
		fmt.Fprintln(os.Stderr, "Cannot decode", args.InputFilename + ":", decodeErr)
		os.Exit(3)
	}

	inFile.Close()

	command := commands[args.Command]
	resultImg := command.Execute(img)

	outFile := os.Stdout
	if args.OutputFilename != "" {
		var createErr error
		outFile, createErr = os.Create(args.OutputFilename)
		if createErr != nil {
			fmt.Fprintln(os.Stderr, "Cannot write to", args.OutputFilename + ":", createErr)
			os.Exit(4)
		}
	}

	saveImage(resultImg, imageFormat, outFile)
}
