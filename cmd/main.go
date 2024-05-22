package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var k = ""

// I should really set this up to work with channels but i tried for like an hour and i really cba


func main() {
	var enableServer string
	fmt.Println("Enable Server? y/N")
	fmt.Scanln(&enableServer)
	ch := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	// run the music handler on a separate thread/goroutine so the server can intercept requests
	go keepGoing(ch, &wg)
	// <- song

	if enableServer == "y" {
		go runServer()
	}
	
	wg.Wait()
	// http server for getting up to date requests 
}

func keepGoing(ch chan string, wg *sync.WaitGroup) {
	
	// defer wg.Done()
	// wg.Wait()
	
	songPrompt(ch)
	k = <- ch
	
	fmt.Println("received song")
	// val, ok := <- song
	
	// fmt.Println(ok)
	// fmt.Println(val)
	// song <- k
	songStr := fmt.Sprintf("song %s", k)
	fmt.Println(songStr)
	songToPlay := fmt.Sprintf("./audio/BEATS/Track_00%s.mp3", k)
	f, err := os.Open(songToPlay)
	
	if err != nil{
		fmt.Println("uh oh")
		panic(err)
	}
	
	streamer, format, err := mp3.Decode(f)
	
	if err != nil {
		panic(err)
	}
	
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	
	// done := make(chan bool)
	
	// fmt.Println("Skip current song? Y/n")
	// var skipStr string
	// for !<- done {
		// 	fmt.Scanln(&skipStr)
		// 	if skipStr != "n" {
			// 		done <- true
			// 	}
			// }
			
			loop := beep.Loop(1, streamer)
			ctrl := &beep.Ctrl{Streamer: loop, Paused: false }
			
			speaker.Play(beep.Seq(ctrl, beep.Callback(func(){
				ctrl.Paused = true
				// done <- true
				fmt.Println("fuck me")
			})))
			
			
	for !ctrl.Paused {
		if ctrl.Paused {
			fmt.Println("exiting early")
			break
		}
		// var skipStr string
		fmt.Println("Skip current song? Y/n")
		// fmt.Scanln(&skipStr)
		fmt.Scanln()
		
		// if skipStr != "n" {
			fmt.Println("skipping")
			speaker.Lock()
			ctrl.Paused = true
			speaker.Unlock()
			break
		// }

	}
	fmt.Println("alo early")

	fmt.Println("alo")
	wg.Add(1)
	keepGoing(ch, wg)

}

// func runServer(curr <- chan string, wg *sync.WaitGroup){
func runServer(){
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// wg.Add(1)
		fmt.Println("received get request")
		io.WriteString(w, fmt.Sprintf("Received request: %s \n song number: %s", r.Method, k))

	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func songPrompt(ch chan string) {
	fmt.Println("Pick a song number. 1-10")
	var songNumber string
	fmt.Scanln(&songNumber)
	if songNumber == "q" {
		fmt.Println("You just had to follow the damn train")
		os.Exit(3)
	}
	ch <- songNumber
}