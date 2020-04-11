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
	HideIdleCursor      bool
	Bookmarks           []Bookmark
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
	*c = Config{
		// cannot use zero-values
		ZoomMode: "BestFit",
		Shrink: true,
		WindowWidth: 640,
		WindowHeight: 480,
		NSkip: 2,
		Seamless: true,
		Interpolation: 2,
		EmbeddedOrientation: true,
		ImageDiffThres: 0.4,
		SceneScanSkip: 5,
		SmartScroll: true,
		HideIdleCursor: true,
		// always retain these variables
		Bookmarks: c.Bookmarks,
		LastDirectory: c.LastDirectory,
	}
}
