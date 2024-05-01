package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func main() {
	keepGoing()

}

func keepGoing() {
	fmt.Println("Pick a song number. 1-10")
	var songNumber string
	fmt.Scanln(&songNumber)
	if songNumber == "q" {
		fmt.Println("You just had to follow the damn train")
		return
	}

	songToPlay := fmt.Sprintf("./audio/BEATS/Track_00%s.mp3", songNumber)
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

	<- done

	keepGoing()
	// defer streamer.Close()
}