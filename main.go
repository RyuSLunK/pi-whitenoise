package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	gohook "github.com/robotn/gohook"
)

func run() error {

	f, err := os.Open("whitenoise.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func main() {
	eventHook := robotgo.Start()

	f, err := os.Open("whitenoise.mp3")
	defer f.Close()

	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	isOn := false
	var e gohook.Event
	fmt.Println("Ready to Play whitenoise")
	for e = range eventHook {
		if e.Kind == gohook.MouseDown {
			if isOn {
				isOn = false
				player.Close()
				player = context.NewPlayer()
				fmt.Println("Stopping whitenoise")
			} else {
				isOn = true
				fmt.Println("Playing whitenoise")
				go func() { play(player, decoder) }()
			}
		}
	}

}

func play(player *oto.Player, decoder *mp3.Decoder) {
	io.Copy(player, decoder)
}
