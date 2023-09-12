package config

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Profile struct {
	Token string `toml:"token"`
	Team  string `toml:"team"`
}

func LoadProfile() Profile {
	profile := Profile{}

	usr, _ := user.Current()
	dir := usr.HomeDir
	profileBytes, err := os.ReadFile(filepath.Join(dir, "/.sailhouse/profile.toml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			profile.SaveProfile()
		} else {
			panic(err)
		}
	}

	err = toml.Unmarshal(profileBytes, &profile)
	if err != nil {
		panic(err)
	}

	return profile
}

func (p *Profile) SaveProfile() {
	usr, _ := user.Current()
	dir := usr.HomeDir

	profileBytes, err := toml.Marshal(p)
	if err != nil {
		panic(err)
	}

	// ensure the `~/.sailhouse` directory exists
	err = os.MkdirAll(filepath.Join(dir, "/.sailhouse"), 0700)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(dir, "/.sailhouse/profile.toml"), profileBytes, 0600)
	if err != nil {
		panic(err)
	}
}
