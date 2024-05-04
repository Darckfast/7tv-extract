package utils

import (
	"fmt"
	"os"
	"sync/atomic"

	"7tv-extract/pkg/types"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	ResolutionsAttempt                 = []uint{128, 96, 64}
	totalEmotesConverted atomic.Uint32 = atomic.Uint32{}
	lastEmoteConverted   string        = ""
	mw                   *imagick.MagickWand
)

func InitMagick() {
}

func DoConversion(shortEmote *types.ShortEmoteList) {
	totalEmotesConverted.Add(1)
	lastEmoteConverted = shortEmote.EmoteName

	if shortEmote.Size > 256*1024*5 {
		ConvertFileV2(shortEmote, 64)
		return
	}

	for _, resolution := range ResolutionsAttempt {
		ConvertFileV2(shortEmote, resolution)

		fs, _ := os.Stat(shortEmote.OutputPath)

		if fs.Size() <= 256*1024 {
			break
		}
	}
}

func ConvertFileV2(
	shortEmote *types.ShortEmoteList,
	resolution uint,
) {
	imagick.Initialize()
	mw := imagick.NewMagickWand()

	defer imagick.Terminate()
	defer mw.Destroy()
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

	PrintLine(
		fmt.Sprintf("\r[%d/%d] Converting emotes [ %s ]",
			totalEmotesConverted.Load(),
			TotalEmotes,
			lastEmoteConverted,
		),
	)
}
