package utils

import (
	"bytes"
	"fmt"
	"image/gif"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"7tv-extract/pkg/types"
)

var (
	HasMagick                          = false
	HasGifsicle                        = false
	totalEmotesConverted atomic.Uint32 = atomic.Uint32{}
	lastEmoteConverted   string        = ""
)

func CheckForMagick() bool {
	_, err := exec.LookPath("magick")
	exec.Command("set", "MAGICK_OCL_DEVICE=true")

	HasMagick = err == nil
	return HasMagick
}

func CheckForGifsicle() bool {
	_, err := exec.LookPath("gifsicle")
	HasGifsicle = err == nil

	return HasGifsicle
}

func ConvertFile(
	shortEmote *types.ShortEmoteList,
	username string,
	limiter chan int,
	wg *sync.WaitGroup,
) {
	defer func() {
		<-limiter
		wg.Done()
	}()

	if !HasMagick {
		fmt.Println("skipping auto conversion, the file will remain as .webp")
		return
	}

	path, _ := os.Getwd()

	extension := "png"

	if shortEmote.IsAnimated {
		extension = "gif"
	}

	var outb, errb bytes.Buffer

	fileName := filepath.Join(
		username,
		shortEmote.EmoteName+"."+shortEmote.Extension)

	fileName = strings.Replace(fileName, ":", "Colon", 1)

	defer os.Remove(fileName)

	outputFileName := strings.Replace(fileName, "webp", extension, 1)
	cmd := exec.Command("magick",
		filepath.Join(path, fileName),
		"-coalesce",
		"-layers",
		"optimize-transparency",
		"-scale",
		"64x64",
		"-fuzz",
		"7%",
		"+dither",
		outputFileName,
	)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while converting file %s, %s\n", fileName, err.Error())
		fmt.Println("Error", errb.String(), outb.String())
	}

	totalEmotesConverted.Add(1)
	lastEmoteConverted = shortEmote.EmoteName

	if extension == "gif" && HasGifsicle {
		fileInfo, err := os.Stat(outputFileName)
		if err != nil {
			fmt.Println("Error stat file", outputFileName, err.Error())
			return
		}

		if fileInfo.Size() >= 256*1024 {
			fileGif, _ := os.Open(outputFileName)

			defer fileGif.Close()

			gifDecoded, err := gif.DecodeAll(fileGif)
			if err != nil {
				fmt.Println("Error decoding GIF", fileName, err.Error())
				return
			}

			totalFrams := len(gifDecoded.Image)

			for i := 0; i < totalFrams; i++ {
				if i%2 == 0 {
					continue
				}

				deleteArg := fmt.Sprintf("#%d", i)
				cmdGifslice := exec.Command("gifsicle",
					"-U",
					outputFileName,
					"--colors",
					"255",
					"--delete",
					deleteArg,
					"-o",
					outputFileName,
				)

				cmdGifslice.Stdout = &outb
				cmdGifslice.Stderr = &errb

				if err := cmdGifslice.Run(); err != nil {
					fmt.Println("Error while converting file", err.Error())
				}
			}
		}
	}

	PrintLine(
		fmt.Sprintf("\r[%d/%d] Converting emotes [ %s ]",
			totalEmotesConverted.Load(),
			TotalEmotes,
			lastEmoteConverted,
		),
	)
}
