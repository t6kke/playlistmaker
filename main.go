package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello from Play List Maker")

	app_config := AppConfig{
		music_library_config: MusicLibraryConfig{},
	}

	err := app_config.load_config()
	if err != nil {
		fmt.Printf("Failed to load configurations: %v\n", err)
	}

	fmt.Println(app_config)

	app_config.music_library_config.scanner()
}
