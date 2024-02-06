package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	ctools "github.com/gookit/color"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"golang.org/x/term"

	"github.com/tomasruud/dog/ansi"
)

func main() {
	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	bg := flag.String("bg", "light", "Background color, availble values: [white, gray, black]")
	colors := flag.String("color", "auto", "Color level, available values: [auto, 16, 256, rgb]")

	flag.Usage = func() {
		fmt.Fprintln(flag.Output(), "Dump out graphics (dog) 🐶")
		fmt.Fprintln(flag.Output(), "A friendly, cat-like image previewer for your terminal.")
		fmt.Fprintln(flag.Output(), "")
		fmt.Fprintln(flag.Output(), "Syntax: dog [OPTIONS] [FILE]")
		fmt.Fprintln(flag.Output(), "")
		fmt.Fprintln(flag.Output(), "Available options:")
		flag.PrintDefaults()
	}

	flag.Parse(os.Args[1:])

	var in io.Reader
	if term.IsTerminal(int(os.Stdin.Fd())) {
		if flag.NArg() < 1 {
			flag.Usage()
			os.Exit(1)
		}

		raw, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			exit("Unable to read file", err)
		}

		in = bytes.NewBuffer(raw)
	} else {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			exit("Unable to read pipe", err)
		}
		in = bytes.NewBuffer(raw)
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		exit("Unable to determine terminal dimensions", err)
	}

	img, _, err := image.Decode(in)
	if err != nil {
		exit("Unable to decode image data", err)
	}

	matte := color.White
	switch *bg {
	case "black":
		matte = color.Black
	case "gray":
		matte = color.Gray16{0xAAAA}
	}

	level := ctools.TermColorLevel()

	switch *colors {
	case "16":
		level = ctools.Level16
	case "256":
		level = ctools.Level256
	case "rgb":
		level = ctools.LevelRgb
	}

	if level == ctools.LevelNo {
		exit("Unable to determine color level", errors.New("level was none"))
	}

	enc := ansi.Encoder{
		MaxW:       w,
		MaxH:       h,
		Matte:      matte,
		ColorLevel: level,
	}
	if err := enc.Encode(os.Stdout, img); err != nil {
		exit("Unable to encode image data", err)
	}
}

func exit(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}
