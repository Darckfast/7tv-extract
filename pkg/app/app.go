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

	threads := runtime.NumCPU()

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

	fmt.Println("Completed", emotes.User.Username, tv7UserId)
}
