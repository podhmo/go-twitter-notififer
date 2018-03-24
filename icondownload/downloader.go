package icondownload

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ChimeraCoder/anaconda"
)

// Downloader :
type Downloader struct {
	images *sync.Map // userid -> filepath
	Base   string

	qCh    chan *anaconda.User
	errCh  chan error
	doneCh chan struct{}
}

// New :
func New(base string) *Downloader {
	var m sync.Map
	return &Downloader{
		Base:   base,
		images: &m,
		qCh:    make(chan *anaconda.User),
		errCh:  make(chan error),
		doneCh: make(chan struct{}),
	}
}

// Data :
func (d *Downloader) Data() map[int64]string {
	data := map[int64]string{}
	d.images.Range(func(key interface{}, value interface{}) bool {
		data[key.(int64)] = value.(string)
		return true
	})
	return data
}

// Load :
func (d *Downloader) Load(data map[int64]string) {
	for k, v := range data {
		d.images.Store(k, v)
	}
}

// Error :
func (d *Downloader) Error() <-chan error {
	return d.errCh
}

// Start :
func (d *Downloader) Start(concurrent int) {
	var wg sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		i := i
		wg.Add(1)
		go func() {
			log.Printf("start downloader %d", i)
			for {
				select {
				case user, ok := <-d.qCh:
					if !ok {
						log.Printf("end downloader %d", i)
						wg.Done()
						return
					}
					k := user.Id
					if _, ok := d.images.Load(k); !ok {
						err := d.download(user)
						if err != nil {
							go func() {
								d.images.Delete(k)
								d.errCh <- err
							}()
						}
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		go func() {
			d.doneCh <- struct{}{}
			close(d.doneCh)
		}()
	}()
}

// Stop :
func (d *Downloader) Stop() <-chan struct{} {
	go func() {
		close(d.errCh)
		close(d.qCh) // todo: graceful shutdown
	}()
	return d.doneCh
}

// Download :
func (d *Downloader) Download(user *anaconda.User) (filepath string, found bool) {
	k := user.Id
	path, found := d.images.Load(k)
	if found {
		return path.(string), found
	}
	go func() { d.qCh <- user }()
	return "", false
}

// download :
func (d *Downloader) download(user *anaconda.User) error {
	log.Printf("for %s, start downloading %s\n", user.Name, user.ProfileImageURL)
	defer log.Printf("for %s, end downloading\n", user.Name)

	k := user.Id
	path := filepath.Join(d.Base, strings.Replace(user.ProfileImageURL, "/", "~1", -1))
	d.images.Store(k, path)

	if err := os.MkdirAll(filepath.Dir(path), 0744); err != nil {
		return err
	}

	img, err := os.Create(path)
	if err != nil {
		return err
	}
	defer img.Close()

	resp, err := http.Get(user.ProfileImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(img, resp.Body); err != nil {
		return err
	}
	return nil
}
