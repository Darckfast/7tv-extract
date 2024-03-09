package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//go:generate goversioninfo

var hasMagick = true

func main() {
	tv7UserId := os.Args[1:]

	if len(tv7UserId) == 0 {
		reader := bufio.NewReader(os.Stdin)
		log.Print("Enter user id: ")
		userId, _ := reader.ReadString('\n')

		userId = strings.TrimSpace(userId)

		tv7UserId = append(tv7UserId, userId)
	}

	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + tv7UserId[0])
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	var emotes Emotes

	json.NewDecoder(resp.Body).Decode(&emotes)

	shortEmoteList := []ShortEmoteList{}

	for _, emote := range emotes.EmoteSet.Emotes {
		baseUrl := emote.Data.Host.URL

		emoteFileName := emote.Data.Host.Files[len(emote.Data.Host.Files)-1]
		emoteFileName.Name = "4x.webp"

		shortEmoteList = append(shortEmoteList, ShortEmoteList{
			FullUrl:    "https:" + baseUrl + "/" + emoteFileName.Name,
			Extension:  "webp",
			EmoteName:  emote.Data.Name,
			IsAnimated: emote.Data.Animated,
		})
	}

	stringContent, _ := json.MarshalIndent(shortEmoteList, "", " ")

	os.MkdirAll(filepath.Join(emotes.Username, "emotes"), os.ModePerm)
	os.WriteFile(filepath.Join(emotes.Username, "emotes_list.json"), stringContent, os.ModePerm)

	limiter := make(chan int, runtime.NumCPU())

	var wg sync.WaitGroup

	CheckForMagick()

	for index, shortEmote := range shortEmoteList {
		limiter <- 1
		wg.Add(1)

		go func(shortEmote ShortEmoteList, index int) {
			log.Printf("%d/%d Downloading emote %s\n", index+1, len(shortEmoteList), shortEmote.EmoteName)
			resp, err := http.Get(shortEmote.FullUrl)
			if err != nil {
				log.Println("Error making the request", err.Error())
				<-limiter
				return
			}

			defer resp.Body.Close()

			fileName := filepath.Join(
				emotes.Username,
				"emotes",
				shortEmote.EmoteName+"."+shortEmote.Extension)

			out, err := os.Create(fileName)
			if err != nil {
				log.Println("Error creating file", err.Error())
				<-limiter
				return
			}

			defer func() {
				out.Close()
				log.Printf("%d/%d Converting emote %s", index+1, len(shortEmoteList), shortEmote.EmoteName)
				ConvertFile(fileName, shortEmote.IsAnimated)
				os.Remove(fileName)

				wg.Done()
			}()

			io.Copy(out, resp.Body)

			<-limiter
		}(shortEmote, index)
	}
	wg.Wait()

	log.Println("Completed", emotes.User.Username, tv7UserId)

	if !hasMagick {
		log.Println("For auto conversion first install ImageMagick https://imagemagick.org/script/download.php")
	}
}

func CheckForMagick() {
	_, err := exec.LookPath("magick")
	hasMagick = err == nil
}

func ConvertFile(fileName string, isAnimated bool) {
	if !hasMagick {
		log.Println("skipping")
		return
	}

	path, _ := os.Getwd()

	extension := "png"

	if isAnimated {
		extension = "gif"
	}

	cmd := exec.Command("magick", "mogrify", "-format", extension, filepath.Join(path, fileName))

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
