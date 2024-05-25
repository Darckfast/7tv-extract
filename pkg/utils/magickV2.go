package utils

import (
	"fmt"
	"os"

	"7tv-extract/pkg/types"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	ResolutionsAttempt        = []uint{128, 96, 64}
	emotesConverted           = 0
	emotesSkipped             = 0
	lastEmoteConverted string = ""
	MAX_SIZE_LIMIT            = 1024 * 1024 * 5
	HARD_SIZE_LIMIT    int64  = 256 * 1024
    mw = imagick.NewMagickWand()
)

func DoConversion(shortEmote *types.ShortEmoteList) {
	emotesConverted = emotesConverted + 1
	lastEmoteConverted = shortEmote.EmoteName

	defer os.Remove(shortEmote.FullPath)

	if shortEmote.Size > MAX_SIZE_LIMIT {
		emotesSkipped = emotesSkipped + 1
		return
	}

	for _, resolution := range ResolutionsAttempt {
		ConvertFileV2(shortEmote, resolution)

		fs, err := os.Stat(shortEmote.OutputPath)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		if fs.Size() <= HARD_SIZE_LIMIT {
			break
		}
	}
}

func ConvertFileV2(
	shortEmote *types.ShortEmoteList,
	resolution uint,
) {
	defer mw.Clear()

	mw.ReadImage(shortEmote.FullPath)
	mw.CoalesceImages()
	mw.OptimizeImageTransparency()
	mw.SetImageFuzz(0.07)
	mw.ScaleImage(resolution, 0)

	if shortEmote.IsAnimated {
		mw.WriteImages(shortEmote.OutputPath, true)
	} else {
		mw.WriteImage(shortEmote.OutputPath)
	}

	Progress(emotesConverted, TotalEmotes)
}
