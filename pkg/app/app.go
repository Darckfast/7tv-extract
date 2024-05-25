package app

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"7tv-extract/pkg/utils"
)

func Run(tv7UserId string) (int, string) {
	if len(tv7UserId) == 0 && len(os.Args) > 1 {
		tv7UserId = os.Args[1:][0]
	}

	if len(tv7UserId) == 0 {
		reader := bufio.NewReader(os.Stdin)
		utils.Info("Enter user id:")
		userId, _ := reader.ReadString('\n')

		userId = strings.TrimSpace(userId)

		tv7UserId = userId
	}

	shortEmoteList, emotes := utils.GetEmoteList(tv7UserId)
	if shortEmoteList == nil {
		return 0, ""
	}

	threads := runtime.NumCPU() + 1

	fmt.Printf("Using %s threads\n", utils.PINK_STYLE.Render(fmt.Sprintf("%d", threads)))
	limiter := make(chan int, threads)

	wg := sync.WaitGroup{}

	utils.Info("Downloading raw emotes")
	for _, shortEmote := range *shortEmoteList {
		wg.Add(1)
		limiter <- 1

		go func() {
			utils.DownloadEmote(
				&shortEmote,
			)

			<-limiter
			wg.Done()
		}()

	}

	wg.Wait()

	utils.Info("Converting emotes")

	for _, shortEmote := range *shortEmoteList {
		utils.DoConversion(&shortEmote)
	}

	totalEmotes := len(*shortEmoteList)
	fmt.Printf("Completed - downloaded %d emotes from %s profile\n",
		totalEmotes,
		utils.INFO_HG_STYLE.Render(emotes.User.Username),
	)

	dirPath := ""

	if totalEmotes > 0 {
		dirPath = (*shortEmoteList)[0].DirPath
	}

	utils.Info("Emotes available in ")
	utils.Highlight(dirPath)

	return totalEmotes, dirPath
}
