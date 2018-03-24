package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"os/signal"
	"os/user"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gen2brain/beeep"
	"github.com/podhmo/egoist/icondownload"
	"gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	path     string
	showIcon bool
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
	app.Flag("config", "config file path").Short('c').Required().ExistingFileVar(&opt.path)
	app.Flag("show-icon", "show profile icon").BoolVar(&opt.showIcon)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(err.Error())
	}

	var downloader *icondownload.Downloader
	if opt.showIcon {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Join(u.HomeDir, ".config/egoist")
		downloader = icondownload.New(filepath.Join(dir, "img"))
		if f, _ := os.Open(filepath.Join(dir, "mapping.json")); f != nil {
			decoder := json.NewDecoder(f)
			var data map[int64]string
			if err := decoder.Decode(&data); err != nil {
				log.Println(err)
			}
			f.Close()
			log.Println("load", len(data), "items")
			downloader.Load(data)
		}
		downloader.Start(4)

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			select {
			case err := <-downloader.Error():
				log.Println(err)
			case <-ch:
				<-downloader.Stop()
				data := downloader.Data()
				log.Println("save", len(data), "items")
				f, err := os.Create(filepath.Join(dir, "mapping.json"))
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				encoder := json.NewEncoder(f)
				if err := encoder.Encode(&data); err != nil {
					log.Println(err)
				}
				f.Close()
				os.Exit(0)
			}
		}()
	}

	if err := run(opt.path, downloader); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(path string, downloader *icondownload.Downloader) error {
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
			imgpath := "assets/information.png"
			if downloader != nil {
				if path, ok := downloader.Download(&tweet.User); ok {
					imgpath = path
				}
			}

			if err := beeep.Notify(tweet.User.Name, tweet.FullText, imgpath); err != nil {
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
	return nil
}
