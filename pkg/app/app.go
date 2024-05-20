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
	if len(tv7UserId) == 0 && len(os.Args) == 1 {
		fmt.Println("No user id")
		return 0, ""
	}

	if len(tv7UserId) == 0 {
		tv7UserId = os.Args[1:][0]
	}

	if len(tv7UserId) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter user id: ")
		userId, _ := reader.ReadString('\n')

		userId = strings.TrimSpace(userId)

		tv7UserId = userId
	}

	shortEmoteList, emotes := utils.GetEmoteList(tv7UserId)
	if shortEmoteList == nil {
		return 0, ""
	}

	threads := runtime.NumCPU() + 1

	fmt.Printf("Using %d threads\n", threads)
	limiter := make(chan int, threads)

	wg := sync.WaitGroup{}

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

	for _, shortEmote := range *shortEmoteList {
		utils.DoConversion(&shortEmote)
	}

	totalEmotes := len(*shortEmoteList)
	fmt.Println("Completed", emotes.User.Username, tv7UserId, totalEmotes)

	dirPath := ""

	if totalEmotes > 0 {
		dirPath = (*shortEmoteList)[0].DirPath
	}

	return totalEmotes, dirPath
}
