package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

const config_dir = ".config/playlistmaker"
const config_filename = "plm_conf.json"

type AppConfig struct {
	config_dir           string
	config_file          string
	music_library_config MusicLibraryConfig
}

func (ac *AppConfig) load_config() error {
	fmt.Println("Loading configurations...")

	home_dir := os.Getenv("HOME")
	if home_dir == "" {
		return fmt.Errorf("no HOME environment variable found")
	}

	config_full_path := home_dir + "/" + config_dir
	_, err := os.Stat(config_full_path)
	if os.IsNotExist(err) {
		fmt.Println("Config dir does not exist, will create it...")
		err := os.MkdirAll(config_full_path, 0750)
		if err != nil {
			return fmt.Errorf("failed to create configuration dir: %v", err)
		}
		fmt.Println("OK")
	} else if err != nil {
		return fmt.Errorf("config dir validation error: %v", err)
	}

	config_file_full_path := home_dir + "/" + config_dir + "/" + config_filename
	_, err = os.Stat(config_file_full_path)
	if os.IsNotExist(err) {
		fmt.Println("Config file does not exist, creating new file...")
		var media_full_dir string
		fmt.Print("Type media full path: ")
		_, err := fmt.Scan(&media_full_dir)
		if err != nil {
			return err
		}
		ac.music_library_config.new_config(media_full_dir)

		byte_data, _ := json.Marshal(ac.music_library_config)
		err = os.WriteFile(config_file_full_path, byte_data, 0600)
		if err != nil {
			return fmt.Errorf("failed to create new configuration file: %v", err)
		}
		fmt.Println("OK")

	} else if err != nil {
		return fmt.Errorf("config file validation error: %v", err)
	} else {
		fmt.Println("Loading libray configurations from file...")
		byte_data, err := os.ReadFile(filepath.Clean(config_file_full_path))
		if err != nil {
			return fmt.Errorf("failed to read configuration file: %v", err)
		}
		err = json.Unmarshal(byte_data, &ac.music_library_config)
		if err != nil {
			return fmt.Errorf("failed to parse configuration file json format: %v", err)
		}
		fmt.Println("OK")
	}

	fmt.Println("Validating media root directories existence...")
	err = ac.music_library_config.validate_dirs()
	if err != nil {
		return err
	}
	fmt.Println("OK")

	ac.config_dir = config_full_path
	ac.config_file = config_filename

	fmt.Println("All configurations loaded successfully!")

	return nil
}

const songs_dir = "all_songs"
const songs_playlist_dir = "playlists"
const radio_playlist_dir = "internet_radios"

type MusicLibraryConfig struct {
	Media_dir          string `json:"media_dir"`
	Songs_dir          string `json:"songs_dir"`
	Songs_playlist_dir string `json:"songs_playlist_dir"`
	Radio_playlist_dir string `json:"radio_playlist_dir"`
}

func (mlc *MusicLibraryConfig) new_config(media_dir string) {
	mlc.Media_dir = media_dir
	mlc.Songs_dir = media_dir + "/" + songs_dir
	mlc.Songs_playlist_dir = media_dir + "/" + songs_playlist_dir
	mlc.Radio_playlist_dir = media_dir + "/" + radio_playlist_dir
}

func (mlc *MusicLibraryConfig) validate_dirs() error {
	v := reflect.ValueOf(mlc).Elem()
	for i := 0; i < v.NumField(); i++ {
		_, err := os.Stat(v.Field(i).String())
		if os.IsNotExist(err) {
			return fmt.Errorf("media directory validation failure: %v", err)
		}
	}
	return nil
}
