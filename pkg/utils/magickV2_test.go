package utils_test

import (
	"os"
	"testing"

	"7tv-extract/pkg/types"
	"7tv-extract/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestMagickReadImage(t *testing.T) {
	shortEmote := types.ShortEmoteList{
		EmoteName:  "wtf",
		FullPath:   "../../.github/images/RIPGE.webp",
		OutputPath: "../../.github/images/RIPGE.png",
		IsAnimated: false,
	}
	defer os.Remove(shortEmote.OutputPath)

	utils.DoConversion(&shortEmote)

	fs, err := os.Stat(shortEmote.OutputPath)

	assert.False(t, err != nil)
	assert.Greater(t, fs.Size(), int64(0))
}

func TestMagickReadAnimatedImage(t *testing.T) {
	shortEmote := types.ShortEmoteList{
		EmoteName:  "yo",
		FullPath:   "../../.github/images/yo.webp",
		OutputPath: "../../.github/images/yo.gif",
		IsAnimated: true,
	}

	defer os.Remove(shortEmote.OutputPath)
	utils.DoConversion(&shortEmote)

	fs, err := os.Stat(shortEmote.OutputPath)

	assert.False(t, err != nil)
	assert.Greater(t, fs.Size(), int64(0))
}
