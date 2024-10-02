package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/api/models"
)

func main() {

	new_user := models.NewUser("Big")
	var tmp int = 12

	fmt.Fprintf(os.Stdout, "\n\nSup, %s --%s age: %d \n\n", new_user.GetUsername(), new_user.GetStatus_String(), tmp)
	fmt.Println("Using goroutines, iterate the bytes of [new_user.id] and print them ..\n .. they are indexed in the order the goroutines were created .. ")

	var wg sync.WaitGroup
	wg.Add(len(new_user.ID))
	for i, id := range new_user.GetId() {
		go func(wg *sync.WaitGroup) {
			// by calling a function using [go nameOfFunction()] a goroutine is created
			// goroutines run concurrently
			time.Sleep(2500 * time.Millisecond)
			fmt.Printf("%d: %d \n", i, id)
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	// sleep to wait for the goroutines to finish
	// time.Sleep(2 * time.Second)
	fmt.Println("new_user_id: ", new_user.GetId())
	fmt.Println("\n\n\n\n\n!! >> See Dockerfile CMD \n\n ")
	fmt.Println("\n\n!! >>  go build /cmd/client/server.go \n\n ")
	fmt.Println("\n\n!! >> PLACEHOLDER (main.go) !! \n\n ")

}
