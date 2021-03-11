package config

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/radovskyb/watcher"
)

type Model struct {
	Proxy map[string]string
	Tls map[string]string
}

var _cfg atomic.Value

func Get() *Model {
	if _cfg.Load() != nil {
		return _cfg.Load().(*Model)
	}

	return nil
}

func Init() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s config.toml\n", os.Args[0])
	}

	if err := load(os.Args[1]); err != nil {
		return err
	}

	watch()

	return nil
}

func load(file string) error {
	var cfg Model

	if _, err := toml.DecodeFile(file, &cfg); err != nil {
		return err
	}
	_cfg.Store(&cfg)
	return nil
}

func watch() {
	w := watcher.New()
	w.FilterOps(watcher.Write)
	_ = w.Add(os.Args[1])

	go func() {
		for {
			select {
			case event := <-w.Event:
				_ = load(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
}