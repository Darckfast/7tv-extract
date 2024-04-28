package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"7tv-extract/pkg/types"

	"github.com/Nykakin/quantize"
	"github.com/gen2brain/avif"
	"github.com/nfnt/resize"
)

var (
	totalEmotesConverted atomic.Uint32 = atomic.Uint32{}
	lastEmoteConverted   string        = ""
)

func ConvertFileNative(
	shortEmote *types.ShortEmoteList,
	username string,
	limiter chan int,
	wg *sync.WaitGroup,
) {
	defer func() {
		<-limiter
		wg.Done()
	}()

	defer os.Remove(shortEmote.FilePath)
	extension := "png"

	if shortEmote.IsAnimated {
		extension = "gif"
	}

	outputFileName := strings.Replace(shortEmote.FilePath, shortEmote.Extension, extension, 1)

	file, err := os.Open(shortEmote.FilePath)
	if err != nil {
		fmt.Println("Error opening file", shortEmote, err.Error())
		return
	}

	defer file.Close()

	if !shortEmote.IsAnimated {
		img, err := avif.Decode(file)
		if err != nil {
			fmt.Println(shortEmote, err.Error())
			return
		}
		ResizeAndEncode(*shortEmote, outputFileName, 128, &img, nil)
		finalFile, _ := os.Stat(outputFileName)

		if finalFile.Size() > 256*1024 {
			ResizeAndEncode(*shortEmote, outputFileName, 96, &img, nil)
		}

		finalFile, _ = os.Stat(outputFileName)

		if finalFile.Size() > 256*1024 {
			ResizeAndEncode(*shortEmote, outputFileName, 64, &img, nil)
		}
	} else {
		imgs, err := avif.DecodeAll(file)
		if err != nil {
			fmt.Println(shortEmote, err.Error())
			panic(err)
		}
		ResizeAndEncode(*shortEmote, outputFileName, 128, nil, imgs)
		finalFile, _ := os.Stat(outputFileName)

		if finalFile.Size() > 256*1024 {
			ResizeAndEncode(*shortEmote, outputFileName, 96, nil, imgs)
		}

		finalFile, _ = os.Stat(outputFileName)

		if finalFile.Size() > 256*1024 {
			ResizeAndEncode(*shortEmote, outputFileName, 64, nil, imgs)
		}
	}

	totalEmotesConverted.Add(1)
	lastEmoteConverted = shortEmote.EmoteName

	PrintLine(
		fmt.Sprintf("\r[%d/%d] Converting emotes [ %s ]",
			totalEmotesConverted.Load(),
			TotalEmotes,
			lastEmoteConverted,
		),
	)
}

func ResizeAndEncode(
	shortEmote types.ShortEmoteList,
	outputFileName string,
	size uint,
	img *image.Image,
	imgs *avif.AVIF,
) {
	outFile, _ := os.Create(outputFileName)

	defer outFile.Close()

	if img != nil {
		newImg := resize.Resize(0, size, *img, resize.Lanczos3)
		png.Encode(outFile, newImg)
	} else if imgs != nil {
		var imgGif gif.GIF

		customPalette := color.Palette{
			image.Transparent,
		}

		quantizer := quantize.NewHierarhicalQuantizer()

		for _, img := range imgs.Image {
			resizedImg := resize.Resize(0, size, img, resize.Lanczos3)
			colors, _ := quantizer.Quantize(resizedImg, 5)

			palette := make([]color.Color, len(colors))
			for index, clr := range colors {
				palette[index] = clr
			}

			customPalette = append(customPalette, palette...)
			palettedImg := image.NewPaletted(resizedImg.Bounds(), customPalette)

			draw.Draw(palettedImg, palettedImg.Rect, resizedImg, resizedImg.Bounds().Min, draw.Src)
			imgGif.Image = append(imgGif.Image, palettedImg)
			imgGif.Delay = append(imgGif.Delay, 1)
			imgGif.Disposal = append(imgGif.Disposal, gif.DisposalPrevious)
		}

		gif.EncodeAll(outFile, &imgGif)
	}
}
