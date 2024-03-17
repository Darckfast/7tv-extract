package main

import (
	"bufio"
	"bytes"
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

	shortEmoteList, emotes := GetEmoteList(tv7UserId[0])
	if shortEmoteList == nil {
		return
	}

	stringContent, _ := json.MarshalIndent(shortEmoteList, "", " ")

	os.MkdirAll(filepath.Join(emotes.Username, "emotes"), os.ModePerm)
	os.WriteFile(filepath.Join(emotes.Username, "emotes_list.json"), stringContent, os.ModePerm)

	threads := runtime.NumCPU()

	log.Printf("Using %d threads\n", threads)
	limiter := make(chan int, threads)

	CheckForMagick()

	var wg sync.WaitGroup

	for index, shortEmote := range *shortEmoteList {
		limiter <- 1
		wg.Add(1)

		DownloadEmote(
			&shortEmote,
			index,
			len(*shortEmoteList),
			emotes.Username,
			limiter,
			&wg,
		)
	}
	wg.Wait()

	log.Println("Completed", emotes.User.Username, tv7UserId)

	if !hasMagick {
		log.Println("For auto conversion first install ImageMagick https://imagemagick.org/script/download.php")
	}
}

func GetEmoteList(userId string) (*[]ShortEmoteList, *Emotes) {
	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + userId)
	if err != nil {
		log.Println(err.Error())
		return nil, nil
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

	return &shortEmoteList, &emotes
}

func DownloadEmote(
	shortEmote *ShortEmoteList,
	index int,
	total int,
	username string,
	limiter chan int,
	wg *sync.WaitGroup,
) {
	log.Printf("%d/%d Downloading emote %s\n", index+1, total, shortEmote.EmoteName)
	resp, err := http.Get(shortEmote.FullUrl)
	if err != nil {
		log.Println("Error making the request", err.Error())
		<-limiter
		wg.Done()
		return
	}

	defer resp.Body.Close()

	fileName := filepath.Join(
		username,
		"emotes",
		shortEmote.EmoteName+"."+shortEmote.Extension)

	out, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file", err.Error())
		<-limiter
		wg.Done()
		return
	}

	defer func() {
		out.Close()
		log.Printf("%d/%d Converting emote %s", index+1, total, shortEmote.EmoteName)
		ConvertFile(fileName, shortEmote.IsAnimated, limiter)
		os.Remove(fileName)

		wg.Done()
	}()

	io.Copy(out, resp.Body)
}

func CheckForMagick() bool {
	_, err := exec.LookPath("magick")
	hasMagick = err == nil
	exec.Command("set", "MAGICK_OCL_DEVICE=true")

	return hasMagick
}

func ConvertFile(fileName string, isAnimated bool, limiter chan int) {
	if !hasMagick {
		log.Println("skipping")
		<-limiter
		return
	}

	path, _ := os.Getwd()

	extension := "png"

	if isAnimated {
		extension = "gif"
	}

	var outb, errb bytes.Buffer

	outputFileName := strings.Replace(fileName, "webp", extension, 1)
	cmd := exec.Command("magick", filepath.Join(path, fileName), "-coalesce", "-layers", "optimize-transparency", outputFileName)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		log.Println("Error while converting file", err.Error())
		log.Println("Error", errb.String(), outb.String())
	}

	<-limiter
}
