package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"7tv-extract/pkg/utils"
)

func Run() {
	tv7UserId := os.Args[1:]

	if len(tv7UserId) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter user id: ")
		userId, _ := reader.ReadString('\n')

		userId = strings.TrimSpace(userId)

		tv7UserId = append(tv7UserId, userId)
	}

	shortEmoteList, emotes := utils.GetEmoteList(tv7UserId[0])
	if shortEmoteList == nil {
		return
	}

	os.MkdirAll(filepath.Join(emotes.Username), os.ModePerm)

	threads := runtime.NumCPU()

	fmt.Printf("Using %d threads\n", threads)
	limiter := make(chan int, threads)

	utils.CheckForMagick()
	utils.CheckForGifsicle()

	if !utils.HasMagick {
		fmt.Println("For auto conversion first install ImageMagick https://imagemagick.org/script/download.php")
	}

	if !utils.HasGifsicle {
		fmt.Println("For higher GIF compression first install https://www.lcdf.org/gifsicle/")
	}

	wg := sync.WaitGroup{}

	for _, shortEmote := range *shortEmoteList {
		wg.Add(1)
		limiter <- 1

		go utils.DownloadEmote(
			&shortEmote,
			emotes.Username,
			limiter,
			&wg,
		)

		break
	}
	wg.Wait()

	for _, shortEmote := range *shortEmoteList {
		limiter <- 1
		wg.Add(1)

		go utils.ConvertFileNative(
			&shortEmote,
			emotes.Username,
			limiter,
			&wg,
		)
		break
	}
	wg.Wait()

	fmt.Println("Completed", emotes.User.Username, tv7UserId)
}
