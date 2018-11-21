package main

import (
	"fmt"
	"kroetnet/match"
	"net/http"
	"strconv"
)

func httpServer(c chan int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count, ok := r.URL.Query()["playerCount"]
		if !ok || len(count[0]) < 1 {
			fmt.Println("Url Param 'playerCount' is missing")
			return
		}
		fmt.Fprint(w, "OK")
		pCount, _ := strconv.Atoi(count[0])
		c <- pCount
	})
	http.ListenAndServe(":8888", nil)
}

func main() {
	// c := make(ch1n int, 1)
	playerCount := 1
	// go httpServer(c)
	// for v := range c {
	//   playerCount = v
	//   break
	// }
	// main goroutine
	fmt.Println("Start match with playerCount:", playerCount)
	port := ":2448"
	match := match.NewMatch(playerCount, 15, port)
	match.StartServer()

}
