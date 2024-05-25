package utils

import (
	"os"
	"sync/atomic"

	"7tv-extract/pkg/types"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	ResolutionsAttempt                 = []uint{128, 96, 64}
	totalEmotesConverted atomic.Uint32 = atomic.Uint32{}
	lastEmoteConverted   string        = ""
	MAX_SIZE_LIMIT        = 256 * 1024 * 5
	HARD_SIZE_LIMIT int64 = 256 * 1024
)

func DoConversion(shortEmote *types.ShortEmoteList) {
	totalEmotesConverted.Add(1)
	lastEmoteConverted = shortEmote.EmoteName

	defer os.Remove(shortEmote.FullPath)

	if shortEmote.Size > MAX_SIZE_LIMIT {
		ConvertFileV2(shortEmote, 64)
		return
	}

	for _, resolution := range ResolutionsAttempt {
		ConvertFileV2(shortEmote, resolution)

		fs, _ := os.Stat(shortEmote.OutputPath)

		if fs.Size() <= HARD_SIZE_LIMIT {
			break
		}
	}
}

var mw = imagick.NewMagickWand()

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

	Progress(int(totalEmotesConverted.Load()), TotalEmotes)
}
