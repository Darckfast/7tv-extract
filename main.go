package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
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
		emoteFileName.Name = "4x.avif"

		shortEmoteList = append(shortEmoteList, ShortEmoteList{
			FullUrl:    "https:" + baseUrl + "/" + emoteFileName.Name,
			Extension:  strings.ToLower(emoteFileName.Format),
			EmoteName:  emote.Data.Name,
			IsAnimated: emote.Data.Animated,
		})
	}

	stringContent, _ := json.MarshalIndent(shortEmoteList, "", " ")

	os.MkdirAll(filepath.Join(emotes.Username, "emotes"), os.ModePerm)
	os.WriteFile(filepath.Join(emotes.Username, "emotes_list.json"), stringContent, os.ModePerm)

	limiter := make(chan int, runtime.NumCPU())

	for index, shortEmote := range shortEmoteList {
		limiter <- 1

		go func(shortEmote ShortEmoteList, index int) {
			fmt.Printf("%d/%d Downloading emote %s\n", index+1, len(shortEmoteList), shortEmote.EmoteName)
			resp, err := http.Get(shortEmote.FullUrl)
			if err != nil {
				fmt.Println("Error making the request", err.Error())
				<-limiter
				return
			}

			defer resp.Body.Close()

			fileName := filepath.Join(
				emotes.Username,
				"emotes",
				shortEmote.EmoteName+"."+shortEmote.Extension)

			fileNameConv := ""

			if shortEmote.IsAnimated {
				fileNameConv = filepath.Join(
					emotes.Username,
					"emotes",
					shortEmote.EmoteName+".gif")
			} else {
				fileNameConv = filepath.Join(
					emotes.Username,
					"emotes",
					shortEmote.EmoteName+".png")
			}

			out, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error creating file", err.Error())
				<-limiter
				return
			}

			defer func() {
				out.Close()
				os.Remove(fileName)
			}()

			io.Copy(out, resp.Body)

			err = ffmpeg_go.Input(fileName).
				Output(fileNameConv).
				OverWriteOutput().
				GlobalArgs("-hide_banner", "-loglevel", "panic", "-y").
				ErrorToStdOut().
				Silent(true).
				Run()

			if err != nil {
				fmt.Println("Error converting file", err.Error())
			}
			<-limiter
		}(shortEmote, index)
	}

	fmt.Println("Download completed", emotes.User.Username, tv7UserId)
}
