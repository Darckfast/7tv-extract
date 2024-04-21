package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/gif"
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
var hasMagick = false
var hasGifsicle = false

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

	os.MkdirAll(filepath.Join(emotes.Username), os.ModePerm)

	threads := runtime.NumCPU()

	log.Printf("Using %d threads\n", threads)
	limiter := make(chan int, threads)

	CheckForMagick()
	CheckForGifsicle()

	if !hasMagick {
		log.Println("For auto conversion first install ImageMagick https://imagemagick.org/script/download.php")
	}

	if !hasGifsicle {
		log.Println("For higher gif compression first install https://www.lcdf.org/gifsicle/")
	}

	wg := sync.WaitGroup{}

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
		shortEmote.EmoteName+"."+shortEmote.Extension)

	fileName = strings.Replace(fileName, ":", "Colon", 1)

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

func CheckForGifsicle() bool {
	_, err := exec.LookPath("gifsicle")
	hasGifsicle = err == nil

	return hasGifsicle
}

func CheckForMagick() bool {
	_, err := exec.LookPath("")
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
	cmd := exec.Command("magick",
		filepath.Join(path, fileName),
		"-coalesce",
		"-layers",
		"optimize-transparency",
		"-scale",
		"64x64",
		"-fuzz",
		"7%",
		"+dither",
		outputFileName,
	)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		log.Println("Error while converting file", err.Error())
		log.Println("Error", errb.String(), outb.String())
	}

	if extension == "gif" && hasGifsicle {
		fileInfo, _ := os.Stat(outputFileName)

		if fileInfo.Size() >= 256*1024 {
			fileGif, _ := os.Open(outputFileName)

			defer fileGif.Close()

			gifDecoded, _ := gif.DecodeAll(fileGif)

			totalFrams := len(gifDecoded.Image)

			for i := 0; i < totalFrams; i++ {
				if i%2 == 0 {
					continue
				}

				deleteArg := fmt.Sprintf("#%d", i)
				cmdGifslice := exec.Command("gifsicle",
					"-U",
					outputFileName,
					"--colors",
					"255",
					"--delete",
					deleteArg,
					"-o",
					outputFileName,
				)

				cmdGifslice.Stdout = &outb
				cmdGifslice.Stderr = &errb

				if err := cmdGifslice.Run(); err != nil {
					log.Println("Error while converting file", err.Error())
				}
			}

			newFileInfo, _ := os.Stat(outputFileName)
			sizeDiff := float32(newFileInfo.Size() - fileInfo.Size())
			sizeDiff = sizeDiff / float32(fileInfo.Size())
			sizeDiff = sizeDiff * 100

			log.Printf("Reduced from %d to %d (-%.2f%%)", fileInfo.Size(), newFileInfo.Size(), sizeDiff)
		}
	}

	<-limiter
}
