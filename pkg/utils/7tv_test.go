package utils_test

import (
	"os"
	"testing"

	"7tv-extract/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestEmotesList(t *testing.T) {
	userId := "65187389"
	shortEmotes, emotes := utils.GetEmoteList(userId)

	assert.NotNil(t, shortEmotes)
	assert.NotNil(t, emotes)

	assert.Equal(t, emotes.Username, "darckfast")
	assert.Equal(t, emotes.ID, userId)
	assert.Equal(t, emotes.EmoteSet.EmoteCount, len(*shortEmotes))
}

func TestDownloadEmoteList(t *testing.T) {
	shortEmotes, _ := utils.GetEmoteList("65187389")

	for _, shortEmote := range *shortEmotes {
		utils.DownloadEmote(&shortEmote)

		fs, err := os.Stat(shortEmote.FullPath)

		assert.Nil(t, err)
		assert.NotZero(t, fs.Size())
	}

    defer os.Remove("darckfast")
}
