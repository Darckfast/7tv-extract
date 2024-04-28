package utils

import (
	"fmt"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"7tv-extract/pkg/types"

	"github.com/gen2brain/avif"
	"github.com/nfnt/resize"
)

var (
// HasMagick                          = false
// HasGifsicle                        = false
// totalEmotesConverted atomic.Uint32 = atomic.Uint32{}
// lastEmoteConverted   string        = ""
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

	fileName := filepath.Join(
		username,
		shortEmote.EmoteName+"."+shortEmote.Extension)

	fileName = strings.Replace(fileName, ":", "Colon", 1)

	defer os.Remove(fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file", fileName, err.Error())
		return
	}

	defer file.Close()

	img, err := avif.Decode(file)
	// _, err2 := avif.DecodeAll(file)
	if err != nil {
		fmt.Println(fileName, err.Error())
		panic(err)
	}

	fmt.Println("Decoded")
	extension := "png"

	if shortEmote.IsAnimated {
		extension = "gif"
	}

	outputFileName := strings.Replace(fileName, "webp", extension, 1)
	outFile, _ := os.Create(outputFileName)

	defer outFile.Close()

	newImg := resize.Resize(128, 0, img, resize.Lanczos3)

	if extension == "png" {
		png.Encode(outFile, newImg)
	} else {
		gif.Encode(outFile, img, &gif.Options{NumColors: 255})
	}

	totalEmotesConverted.Add(1)
	lastEmoteConverted = shortEmote.EmoteName

	// if extension == "gif" && HasGifsicle {
	// 	fileInfo, err := os.Stat(outputFileName)
	// 	if err != nil {
	// 		fmt.Println("Error stat file", outputFileName, err.Error())
	// 		return
	// 	}
	//
	// 	if fileInfo.Size() >= 256*1024 {
	// 		fileGif, _ := os.Open(outputFileName)
	//
	// 		defer fileGif.Close()
	//
	// 		gifDecoded, err := gif.DecodeAll(fileGif)
	// 		if err != nil {
	// 			fmt.Println("Error decoding GIF", fileName, err.Error())
	// 			return
	// 		}
	//
	// 		totalFrams := len(gifDecoded.Image)
	//
	// 		for i := 0; i < totalFrams; i++ {
	// 			if i%2 == 0 {
	// 				continue
	// 			}
	//
	// 			deleteArg := fmt.Sprintf("#%d", i)
	// 			cmdGifslice := exec.Command("gifsicle",
	// 				"-U",
	// 				outputFileName,
	// 				"--colors",
	// 				"255",
	// 				"--delete",
	// 				deleteArg,
	// 				"-o",
	// 				outputFileName,
	// 			)
	//
	// 			cmdGifslice.Stdout = &outb
	// 			cmdGifslice.Stderr = &errb
	//
	// 			if err := cmdGifslice.Run(); err != nil {
	// 				fmt.Println("Error while converting file", err.Error())
	// 			}
	// 		}
	// 	}
	// }

	PrintLine(
		fmt.Sprintf("\r[%d/%d] Converting emotes [ %s ]",
			totalEmotesConverted.Load(),
			TotalEmotes,
			lastEmoteConverted,
		),
	)
}
