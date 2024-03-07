package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	tv7UserId := os.Args[1:]

	if len(tv7UserId) == 0 {
		fmt.Println("No 7tv id")
		return
	}

	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + tv7UserId[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	var emotes Emotes

	json.NewDecoder(resp.Body).Decode(&emotes)

	shortEmoteList := []ShortEmoteList{}

	for _, emote := range emotes.EmoteSet.Emotes {
		baseUrl := emote.Data.Host.URL
		emoteFileName := emote.Data.Host.Files[len(emote.Data.Host.Files)-1]

		shortEmoteList = append(shortEmoteList, ShortEmoteList{
			FullUrl:   "https:" + baseUrl + "/" + emoteFileName.Name,
			Extension: strings.ToLower(emoteFileName.Format),
			EmoteName: emote.Data.Name,
		})
	}

	stringContent, _ := json.MarshalIndent(shortEmoteList, "", " ")

	os.MkdirAll(filepath.Join(emotes.Username, "emotes"), os.ModePerm)
	os.WriteFile(filepath.Join(emotes.Username, "emotes_list.json"), stringContent, os.ModePerm)

	for index, shortEmote := range shortEmoteList {
		fmt.Printf("%d/%d Downloading emote %s\n", index+1, len(shortEmoteList), shortEmote.EmoteName)

		resp, err := http.Get(shortEmote.FullUrl)
		if err != nil {
			fmt.Println("Error making the request", err.Error())
			continue
		}

		defer resp.Body.Close()

		out, err := os.Create(filepath.Join(
			emotes.Username,
			"emotes",
			shortEmote.EmoteName+"."+shortEmote.Extension))
		if err != nil {
			fmt.Println("Error creating file", err.Error())
			continue
		}

		defer out.Close()

		io.Copy(out, resp.Body)
	}

	fmt.Println("Download completed", emotes.User.Username, tv7UserId)
}
