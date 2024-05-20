package app_test

import (
	"os"
	"testing"

	"7tv-extract/pkg/app"

	"github.com/stretchr/testify/assert"
)

func TestAppRun(t *testing.T) {
	totalEmotes, emotesDir := app.Run("65187389")
	defer os.RemoveAll(emotesDir)

	fs, err := os.Stat(emotesDir)

	assert.Nil(t, err)
	assert.True(t, fs.IsDir())

	files, _ := os.ReadDir(emotesDir)

	assert.Equal(t, len(files), totalEmotes)
}

