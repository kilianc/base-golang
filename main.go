package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	fmt.Println("start")
	var wg sync.WaitGroup

	{
		seconds := 1
		fmt.Printf("waiting for %d seconds...\n", seconds)
		wait(3, &wg)
		wait(3, &wg)
		wg.Wait()
	}

	{
		seconds := 2
		fmt.Printf("waiting for %d seconds...\n", seconds)
		<-wait(seconds, nil)
	}

	url := "https://jsonplaceholder.typicode.com/users"

	// just testing unwrap patterns: unnecessarily async since there is nothing else going on
	// this is the normal way to do it: people, err := FetchPeople(url)
	people, err := (<-AsyncFetchPeople(url)).Unwrap()

	if err != nil {
		log.Fatal("error fetching people: ", err)
	}

	for _, person := range people {
		fmt.Println(person.String())
	}

	fmt.Println("end")
}
