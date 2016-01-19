package main

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var ErrSteamNotFound = errors.New("steam not found")

// FindSteam returns a valid path to the Steam install, if found.
func FindSteam() (string, error) {
	var steamPath string
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// TODO: Try several options despite GOOS/GOARCH?
	switch runtime.GOOS {
	case "darwin":
		steamPath = filepath.Join(usr.HomeDir, "Library", "Application Support", "Steam")
	case "linux":
		steamPath = filepath.Join(usr.HomeDir, ".local", "share", "Steam")
	case "windows":
		if runtime.GOARCH == "amd64" {
			steamPath = filepath.Join("C:", "Program Files (x86)", "Steam")
		} else if runtime.GOARCH == "386" {
			steamPath = filepath.Join("C:", "Program Files", "Steam")
		}
	}
	if steamPath == "" {
		return "", ErrSteamNotFound
	}

	_, err = os.Stat(steamPath)
	if os.IsNotExist(err) {
		return "", ErrSteamNotFound
	}
	// Return the error just in case it's something like permission denied.
	return steamPath, err
}
