package utils

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strings"
	"sync"

	"7tv-extract/pkg/types"

	"github.com/gen2brain/avif"
	"github.com/nfnt/resize"
)

var mu sync.Mutex

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

	outputFileName := strings.Replace(shortEmote.FilePath, "avif", extension, 1)

	file, err := os.Open(shortEmote.FilePath)
	if err != nil {
		fmt.Println("Error opening file", shortEmote, err.Error())
		return
	}

	defer file.Close()

	if outputFileName[:3] == "png" {
		mu.Lock()
		img, err := avif.Decode(file)
		if err != nil {
			fmt.Println(shortEmote, err.Error())
			panic(err)
		}
		mu.Unlock()
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
		mu.Lock()
		imgs, err := avif.DecodeAll(file)
		if err != nil {
			fmt.Println(shortEmote, err.Error())
			panic(err)
		}
		mu.Unlock()
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

	if outputFileName[:3] == "png" {
		newImg := resize.Resize(size, 0, *img, resize.Lanczos3)
		png.Encode(outFile, newImg)
	} else {
		var imgGif gif.GIF
		customPalette := append(palette.WebSafe, image.Transparent)

		for _, img := range imgs.Image {
			resizedImg := resize.Resize(size, 0, img, resize.Lanczos3)
			palettedImg := image.NewPaletted(resizedImg.Bounds(), customPalette)

			draw.Draw(palettedImg, palettedImg.Rect, resizedImg, resizedImg.Bounds().Min, draw.Over)
			imgGif.Image = append(imgGif.Image, palettedImg)
			imgGif.Delay = append(imgGif.Delay, 0)
		}

		gif.EncodeAll(outFile, &imgGif)
	}
}
