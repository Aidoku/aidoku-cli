package build

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/beerpiss/aidoku-cli/internal/common"
	"github.com/fatih/color"
	"github.com/valyala/fastjson"
)

type source struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	File    string `json:"file"`
	Icon    string `json:"icon"`
	Lang    string `json:"lang"`
	Version int    `json:"version"`
	NSFW    int    `json:"nsfw"`
}

func BuildWrapper(zipPatterns []string, output string) error {
	os.RemoveAll(output)
	var fileList []string
	for _, arg := range zipPatterns {
		files, err := filepath.Glob(arg)
		if err != nil {
			color.Red("error: invalid glob pattern %s", arg)
			continue
		}
		fileList = append(fileList, files...)
	}
	if len(fileList) == 0 {
		return errors.New("no files given")
	}
	os.MkdirAll(output, os.FileMode(0644))
	os.MkdirAll(output+"/icons", os.FileMode(0644))
	os.MkdirAll(output+"/sources", os.FileMode(0644))
	return BuildSource(fileList, output)
}

func BuildSource(zipFiles []string, output string) error {
	var wg sync.WaitGroup
	sourceList := struct {
		sync.Mutex
		data []source
	}{}
	for _, file := range zipFiles {
		wg.Add(1)
		go func(zipFile string) {
			defer wg.Done()
			r, err := zip.OpenReader(zipFile)
			if err != nil {
				color.Red("error: %s is not a valid package", zipFile)
				return
			}
			defer r.Close()

			var sourceInfo source
			var parser fastjson.Parser
			hasIcon := false
			for _, f := range r.File {
				if f.Name == "Payload/source.json" {
					rc, err := f.Open()
					if err != nil {
						color.Red("error: couldn't read source info for %s", zipFile)
						return
					}
					buf := new(strings.Builder)
					io.Copy(buf, rc)

					raw, err := parser.Parse(buf.String())
					if err != nil {
						color.Red("error: source.json is malformed for %s", zipFile)
						return
					}

					info := raw.Get("info")
					sourceInfo.ID = string(info.GetStringBytes("id"))
					sourceInfo.Lang = string(info.GetStringBytes("lang"))
					sourceInfo.Name = string(info.GetStringBytes("name"))
					sourceInfo.Version = info.GetInt("version")
					sourceInfo.NSFW = info.GetInt("nsfw")
					sourceInfo.File = fmt.Sprintf("%s-v%d.aix", sourceInfo.ID, sourceInfo.Version)
					sourceInfo.Icon = fmt.Sprintf("%s-v%d.png", sourceInfo.ID, sourceInfo.Version)

					common.CopyFileContents(zipFile, output+"/sources/"+sourceInfo.File)
					sourceList.Lock()
					sourceList.data = append(sourceList.data, sourceInfo)
					sourceList.Unlock()
				} else if f.Name == "Payload/Icon.png" {
					hasIcon = true
					rc, err := f.Open()
					if err != nil {
						color.Red("error: couldn't read icon for %s", zipFile)
						return
					}
					img, err := os.Create(fmt.Sprintf("%s/icons/%s.png", output, filepath.Base(zipFile)))
					if err != nil {
						color.Red("error: Couldn't create temporary icon file %s/icons/%s.png", output, filepath.Base(zipFile))
						hasIcon = false
						return
					}
					io.Copy(img, rc)
					img.Sync()
					img.Close()
				}
			}

			if !hasIcon {
				color.Yellow("warning: %s doesn't have an icon, generating placeholder", zipFile)
				img, err := os.Create(fmt.Sprintf("%s/icons/%s", output, sourceInfo.Icon))
				if err != nil {
					color.Red("error: Couldn't write icon file %s: %s", sourceInfo.Icon, err.Error())
					return
				}
				transparent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==")
				io.Copy(img, bytes.NewReader(transparent))
				img.Sync()
				img.Close()

			} else {
				os.Rename(fmt.Sprintf("%s/icons/%s.png", output, filepath.Base(zipFile)), fmt.Sprintf("%s/icons/%s", output, sourceInfo.Icon))
			}
		}(file)
	}
	wg.Wait()
	b, err := json.Marshal(sourceList.data)
	if err != nil {
		color.Red("Couldn't serialize source list: ", err.Error())
		return err
	}

	fm, err := os.Create(output + "/index.min.json")
	if err != nil {
		return err
	}
	defer fm.Close()
	fm.Write(b)
	fm.Sync()

	b, err = json.MarshalIndent(sourceList.data, "", "  ")
	if err != nil {
		color.Red("Couldn't serialize source list: ", err.Error())
		return err
	}

	f, err := os.Create(output + "/index.json")
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(b)
	f.Sync()
	return nil
}
