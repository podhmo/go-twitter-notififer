package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gen2brain/beeep"
	"gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	path string
}

// config :
type config struct {
	ConsumerKey       string `json:"ConsumerKey"`
	ConsumerSecret    string `json:"ConsumerSecret"`
	AccessToken       string `json:"AccessToken"`
	AccessTokenSecret string `json:"AccessTokenSecret"`
}

func main() {
	var opt opt
	app := kingpin.New("egoist", "twitter client")
	app.Flag("--config", "config file path").Short('c').Required().ExistingFileVar(&opt.path)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(err.Error())
	}

	if err := run(opt.path); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(path string) error {
	var c config
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return err
	}

	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api := anaconda.NewTwitterApiWithCredentials(
		c.AccessToken,
		c.AccessTokenSecret,
		c.ConsumerKey,
		c.ConsumerSecret,
	)

	for x := range api.UserStream(nil).C {
		switch tweet := x.(type) {
		case anaconda.Tweet:
			if err := beeep.Notify(tweet.User.Name, tweet.FullText, "assets/information.png"); err != nil {
				log.Printf("%+v", err)
			}
		case anaconda.StatusDeletionNotice:
			log.Println("deleted")
		default:
			msg := fmt.Sprintf("unknown type(%T) : %v", x, x)
			if len(msg) > 140 {
				msg = msg[:140] + "..."
			}
			log.Println(msg)
		}
	}
}
