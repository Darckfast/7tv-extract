package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"

	"7tv-extract/pkg/types"
)

var (
	discordEmotes         *types.DiscordEmotes = nil
	totalEmotesDownloaded atomic.Uint32        = atomic.Uint32{}
	lastEmoteDownloaded   string               = ""

	TotalEmotes int = 0
)

func GetEmoteList(userId string) (*[]types.ShortEmoteList, *types.Emotes) {
	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + userId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	defer resp.Body.Close()
	var emotes types.Emotes

	json.NewDecoder(resp.Body).Decode(&emotes)

	shortEmoteList := []types.ShortEmoteList{}

	for _, emote := range emotes.EmoteSet.Emotes {
		TotalEmotes = TotalEmotes + 1
		baseUrl := emote.Data.Host.URL

		emoteFile := emote.Data.Host.Files[len(emote.Data.Host.Files)-1]
		emoteFile.Name = "4x.webp"

		outExtension := "png"

		if emote.Data.Animated {
			outExtension = "gif"
		}

		fileName := filepath.Join(
			emotes.Username,
			emote.Name+".webp")

		fileName = strings.Replace(fileName, ":", "Colon", 1)
		currentDir, _ := os.Getwd()
		dirPath := filepath.Join(currentDir, emotes.Username)
		outFileName := fileName[:len(fileName)-4] + outExtension

		shortEmoteList = append(shortEmoteList, types.ShortEmoteList{
			FullUrl:    "https:" + baseUrl + "/" + emoteFile.Name,
			FullPath:   fileName,
			DirPath:    dirPath,
			OutputPath: outFileName,
			EmoteName:  emote.Data.Name,
			IsAnimated: emote.Data.Animated,
			Size:       emoteFile.Size,
		})
	}

	sort.Slice(shortEmoteList, func(i, j int) bool {
		return shortEmoteList[i].Size < shortEmoteList[j].Size
	})

	fmt.Println("")

	return &shortEmoteList, &emotes
}

func DownloadEmote(
	shortEmote *types.ShortEmoteList,
) {
	totalEmotesDownloaded.Add(1)
	lastEmoteDownloaded = shortEmote.EmoteName

	os.MkdirAll(shortEmote.DirPath, os.ModePerm)

	resp, err := http.Get(shortEmote.FullUrl)
	if err != nil {
		fmt.Println("Error making the request", err.Error())
		return
	}

	defer resp.Body.Close()

	out, err := os.Create(shortEmote.FullPath)
	if err != nil {
		fmt.Println("Error creating file", err.Error())
		return
	}

	defer out.Close()

	io.Copy(out, resp.Body)

	Progress(int(totalEmotesDownloaded.Load()), TotalEmotes)
}
