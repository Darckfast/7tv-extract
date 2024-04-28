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
	"sync"
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
	LoadDiscordJSON()

	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + userId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	defer resp.Body.Close()
	var emotes types.Emotes

	json.NewDecoder(resp.Body).Decode(&emotes)

	shortEmoteList := []types.ShortEmoteList{}

	skipped := 0
	for _, emote := range emotes.EmoteSet.Emotes {
		if IsDuplicate(emote.Name) {
			skipped = skipped + 1
			PrintLine(fmt.Sprintf("\rSkipping duplicate [ %s ] %d", emote.Name, skipped))
			continue
		}

		TotalEmotes = TotalEmotes + 1
		baseUrl := emote.Data.Host.URL

		emoteFile := emote.Data.Host.Files[len(emote.Data.Host.Files)-1]
		emoteFile.Name = "4x.avif"

		shortEmoteList = append(shortEmoteList, types.ShortEmoteList{
			FullUrl:    "https:" + baseUrl + "/" + emoteFile.Name,
			Extension:  "avif",
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

func IsDuplicate(emoteName string) bool {
	if discordEmotes == nil {
		return false
	}

	for _, discEmote := range *discordEmotes {
		if discEmote.Name == emoteName {
			return true
		}
	}

	return false
}

func LoadDiscordJSON() {
	fmt.Println("Looking for discord.json")

	if _, err := os.Stat("discord.json"); err == nil {
		content, _ := os.ReadFile("discord.json")
		json.Unmarshal(content, &discordEmotes)

		fmt.Println("Found discord.json, ignoring duplicates")
	}
}

func DownloadEmote(
	shortEmote *types.ShortEmoteList,
	username string,
	limiter chan int,
	wg *sync.WaitGroup,
) {
	resp, err := http.Get(shortEmote.FullUrl)
	if err != nil {
		fmt.Println("Error making the request", err.Error())
		wg.Done()
		return
	}

	defer resp.Body.Close()

	fileName := filepath.Join(
		username,
		shortEmote.EmoteName+"."+shortEmote.Extension)

	fileName = strings.Replace(fileName, ":", "Colon", 1)

	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file", err.Error())
		wg.Done()
		return
	}

	defer func() {
		out.Close()
		wg.Done()
		<-limiter
	}()

	io.Copy(out, resp.Body)
	totalEmotesDownloaded.Add(1)
	lastEmoteDownloaded = shortEmote.EmoteName

	PrintLine(
		fmt.Sprintf("\r[%d/%d] Downloading emotes [ %s ]",
			totalEmotesDownloaded.Load(),
			TotalEmotes,
			lastEmoteDownloaded,
		),
	)
}
