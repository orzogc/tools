package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	x := flag.Int("x", 0, "裁剪范围顶点的 x 坐标，默认为 0")
	y := flag.Int("y", 0, "裁剪范围顶点的 y 坐标，默认为 0")
	width := flag.Int("width", 0, "裁剪宽度（必需），会限制在图片范围内")
	height := flag.Int("height", 0, "裁剪高度（必需），会限制在图片范围内")
	input := flag.String("input", "", "要裁剪的 PNG 图片文件（必需）")
	output := flag.String("output", "", "裁剪图片的输出文件（必需）")
	helpShort := flag.Bool("h", false, "打印帮助信息")
	helpLong := flag.Bool("help", false, "打印帮助信息")

	flag.Parse()

	if flag.NFlag() == 0 || *helpShort || *helpLong {
		fmt.Println("crop_png [参数]")
		flag.PrintDefaults()

		return
	}

	if *width == 0 || *height == 0 {
		log.Fatal("需要设置 width 和 height")
	}
	if *input == "" || *output == "" {
		log.Fatal("需要设置 input 和 output")
	}

	inputFile, err := os.Open(*input)
	if err != nil {
		log.Fatalf("open image error: %v", err)
	}
	defer inputFile.Close()

	originalImage, err := png.Decode(inputFile)
	if err != nil {
		log.Fatalf("decode png image error: %v", err)
	}

	bounds := originalImage.Bounds()
	x0 := max(min(*x, *x+*width), bounds.Min.X)
	y0 := max(min(*y, *y+*height), bounds.Min.Y)
	x1 := min(max(*x, *x+*width), bounds.Max.X)
	y1 := min(max(*y, *y+*height), bounds.Max.Y)
	cropRect := image.Rect(x0, y0, x1, y1)
	croppedImage := originalImage.(SubImager).SubImage(cropRect)

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatalf("create image error: %v", err)
	}
	defer outputFile.Close()
	if err := png.Encode(outputFile, croppedImage); err != nil {
		log.Fatalf("encode png image error: %v", err)
	}
}
