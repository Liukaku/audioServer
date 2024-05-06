package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var song = ""

// I should really set this up to work with channels but i tried for like an hour and i really cba

// var song = make(chan string)

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1)
	// run the music handler on a separate thread/goroutine so the server can intercept requests
	go keepGoing()
	// <- song

	
	// wg.Wait()
	// http server for getting up to date requests 
	runServer()
}

func runServer(){
	// defer wg.Done()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received get request")
		io.WriteString(w, fmt.Sprintf("Received request: %s \n song number: %s", r.Method, song))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}


func songPrompt() string {
	fmt.Println("Pick a song number. 1-10")
	var songNumber string
	fmt.Scanln(&songNumber)
	if songNumber == "q" {
		fmt.Println("You just had to follow the damn train")
		os.Exit(3)
	}

	return songNumber
}

func keepGoing() {

	// defer wg.Done()
	// wg.Wait()

	k := songPrompt()
	song = k

	fmt.Println("received song")
	// val, ok := <- song

	// fmt.Println(ok)
	// fmt.Println(val)
	// song <- k
	
	songToPlay := fmt.Sprintf("./audio/BEATS/Track_00%s.mp3", song)
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
	
	done := make(chan bool)
	
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// <- song 

	<- done

	// close(song)

	keepGoing()
	// defer streamer.Close()
}