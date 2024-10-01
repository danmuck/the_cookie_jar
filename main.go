package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/danmuck/the_cookie_jar/sandbox"
	"github.com/danmuck/the_cookie_jar/sandbox/db_types"
)

func main() {

	new_user := db_types.NewUser("Big")
	var tmp int = 12

	fmt.Fprintf(os.Stderr, "\n\nSup, %s --%s age: %d \n\n", new_user.GetUsername(), new_user.GetStatus_String(), tmp)
	fmt.Println("Using goroutines, iterate the bytes of [new_user.id] and print them ..\n .. they are indexed in the order the goroutines were created .. ")

	var wg sync.WaitGroup
	wg.Add(len(new_user.ID))
	for i, id := range new_user.GetId() {
		go func(wg *sync.WaitGroup) {
			// by calling a function using [go nameOfFunction()] a goroutine is created
			// goroutines run concurrently
			time.Sleep(250 * time.Millisecond)
			fmt.Printf("%d: %d \n", i, id)
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	// sleep to wait for the goroutines to finish
	// time.Sleep(2 * time.Second)
	fmt.Println("id: ", new_user.GetId())

	fmt.Println("\nServer connecting .. \n ")
	server, err := sandbox.RunServer("test_org_key")
	if err != nil {
		log.Fatal(err)
	}

	go server.Serve()

	test_users := []string{"Dan M.", "Michael R.", "Michael Y.", "Saqib M.", "Cordell H."}
	for {
		fmt.Println(">> waiting to insert random maker .. ctrl-c to quit")
		time.Sleep(5 * time.Second)

		i := int64(rand.Intn(len(test_users)))
		rs := test_users[i]
		server.DB_AddUser(*db_types.NewUser(rs))
	}

}
