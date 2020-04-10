package main

import (
	"encoding/json"
	"os"
)

const (
	ConfigDir  = ".config/gomics" // relative to user's home
	ConfigFile = "config"         // relative to config dir
	ImageDir   = "images"         // relative to config dir
)

type Config struct {
	ZoomMode            string
	Enlarge             bool
	Shrink              bool
	LastDirectory       string
	Fullscreen          bool
	WindowWidth         int
	WindowHeight        int
	NSkip               int
	Random              bool
	Seamless            bool
	HFlip               bool
	VFlip               bool
	DoublePage          bool
	MangaMode           bool
	OneWide             bool
	EmbeddedOrientation bool
	Interpolation       int
	ImageDiffThres      float32
	SceneScanSkip       int
	SmartScroll         bool
	Bookmarks           []Bookmark
	HideIdleCursor      bool
}

func (c *Config) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	if err = d.Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) Save(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (c *Config) Defaults() {
	c.ZoomMode = "BestFit"
	c.Shrink = true
	c.Enlarge = false
	c.WindowWidth = 640
	c.WindowHeight = 480
	c.NSkip = 10
	c.Seamless = true
	c.Interpolation = 2
	c.EmbeddedOrientation = true
	c.ImageDiffThres = 0.4
	c.SceneScanSkip = 5
	c.SmartScroll = true
	c.HideIdleCursor = true
}
