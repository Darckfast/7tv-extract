package utils_test

import (
	"os"
	"sync"
	"testing"

	"7tv-extract/pkg/types"
	"7tv-extract/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestMagickReadImage(t *testing.T) {
	utils.InitMagick()

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
	utils.InitMagick()
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

func TestMagickReadImageMultiThread(t *testing.T) {
	utils.InitMagick()
	emotes := []types.ShortEmoteList{}

	emotes = append(emotes, types.ShortEmoteList{
		EmoteName:  "yo",
		FullPath:   "../../.github/images/yo.webp",
		OutputPath: "../../.github/images/yo.gif",
		IsAnimated: true,
	}, types.ShortEmoteList{
		EmoteName:  "wtf",
		FullPath:   "../../.github/images/RIPGE.webp",
		OutputPath: "../../.github/images/RIPGE.png",
		IsAnimated: false,
	})

	defer os.Remove(emotes[0].OutputPath)
	defer os.Remove(emotes[1].OutputPath)

	os.Remove(emotes[0].OutputPath)
	os.Remove(emotes[1].OutputPath)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		utils.DoConversion(&emotes[0])
		wg.Done()
		fs, err := os.Stat(emotes[0].OutputPath)

		assert.False(t, err != nil)
		assert.Greater(t, fs.Size(), int64(0))
	}()

	wg.Add(1)
	go func() {
		utils.DoConversion(&emotes[1])
		wg.Done()
		fs, err := os.Stat(emotes[1].OutputPath)

		assert.False(t, err != nil)
		assert.Greater(t, fs.Size(), int64(0))
	}()
	wg.Wait()
}
