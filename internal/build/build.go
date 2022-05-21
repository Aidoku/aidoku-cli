package build

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Aidoku/aidoku-cli/internal/common"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/fatih/color"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/valyala/fastjson"
)

type source struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	File       string `json:"file"`
	Icon       string `json:"icon"`
	Lang       string `json:"lang"`
	Version    int    `json:"version"`
	NSFW       int    `json:"nsfw"`
	MinVersion string `json:"minVersion,omitempty"`
	MaxVersion string `json:"maxVersion,omitempty"`
}

func BuildWrapper(zipPatterns []string, output string) error {
	os.RemoveAll(output)
	var fileList []string
	for _, arg := range zipPatterns {
		base, pattern := doublestar.SplitPattern(filepath.ToSlash(arg))
		fsys := os.DirFS(base)
		files, err := doublestar.Glob(fsys, pattern)
		if err != nil {
			color.Red("error: invalid glob pattern %s", arg)
			continue
		}
		for _, file := range files {
			fileList = append(fileList, base+"/"+file)
		}
	}
	if len(fileList) == 0 {
		return errors.New("no files given")
	}
	err := os.MkdirAll(output, os.FileMode(0777))
	if err != nil {
		color.Red("fatal: could not create output folder")
		return err
	}
	os.MkdirAll(output+"/icons", os.FileMode(0777))
	os.MkdirAll(output+"/sources", os.FileMode(0777))
	return BuildSource(fileList, output)
}

func BuildSource(zipFiles []string, output string) error {
	var wg sync.WaitGroup
	sourceList := struct {
		sync.Mutex
		data []source
	}{}
	sourceIds := struct {
		sync.Mutex
		data map[string]string
	}{}
	sourceIds.data = make(map[string]string)
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
			zipFileHash := fmt.Sprintf("%x", fnv1a.HashString64(zipFile))
			tempImageFile := fmt.Sprintf("%s/icons/%s.png", output, zipFileHash)
			for _, f := range r.File {
				if f.Name == "Payload/source.json" {
					rc, err := f.Open()
					if err != nil {
						color.Red("error: couldn't read source info for %s: %s", zipFile, err)
						os.Remove(tempImageFile)
						return
					}
					buf := new(strings.Builder)
					io.Copy(buf, rc)

					raw, err := parser.Parse(buf.String())
					if err != nil {
						color.Red("error: source.json is malformed for %s: %s", zipFile, err)
						os.Remove(tempImageFile)
						return
					}

					info := raw.Get("info")
					sourceInfo.Id = string(info.GetStringBytes("id"))
					if val, ok := sourceIds.data[sourceInfo.Id]; ok {
						color.Red("error: duplicate source identifier %s in %s, first found in %s", sourceInfo.Id, zipFile, val)
						os.Remove(tempImageFile)
						return
					}
					sourceIds.Lock()
					sourceIds.data[sourceInfo.Id] = zipFile
					sourceIds.Unlock()

					sourceInfo.Lang = string(info.GetStringBytes("lang"))
					sourceInfo.Name = string(info.GetStringBytes("name"))
					sourceInfo.Version = info.GetInt("version")
					sourceInfo.NSFW = info.GetInt("nsfw")
					sourceInfo.File = fmt.Sprintf("%s-v%d.aix", sourceInfo.Id, sourceInfo.Version)
					sourceInfo.Icon = fmt.Sprintf("%s-v%d.png", sourceInfo.Id, sourceInfo.Version)
					if minVersion := info.GetStringBytes("minVersion"); minVersion != nil {
						sourceInfo.MinVersion = string(minVersion)
					}
					if maxVersion := info.GetStringBytes("maxVersion"); maxVersion != nil {
						sourceInfo.MaxVersion = string(maxVersion)
					}

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
					img, err := os.Create(tempImageFile)
					if err != nil {
						color.Red("error: Couldn't create temporary icon file %s/icons/%s.png: %s", output, zipFileHash, err)
						hasIcon = false
						return
					}
					io.Copy(img, rc)
					img.Sync()
					img.Close()
				}
			}
			imageFile := fmt.Sprintf("%s/icons/%s", output, sourceInfo.Icon)
			if !hasIcon {
				color.Yellow("warning: %s doesn't have an icon, generating placeholder", zipFile)
				err = common.GeneratePng(imageFile)
				if err != nil {
					return
				}

			} else {
				os.Rename(tempImageFile, imageFile)
			}
		}(file)
	}
	wg.Wait()

	b, err := json.Marshal(sourceList.data)
	if err != nil {
		color.Red("fatal: couldn't serialize source list: %s", err.Error())
		return err
	}

	fm, err := os.Create(output + "/index.min.json")
	if err != nil {
		return err
	}
	fm.Write(b)
	fm.Sync()
	fm.Close()

	b, err = json.MarshalIndent(sourceList.data, "", "  ")
	if err != nil {
		color.Red("fatal: couldn't serialize source list: %s", err.Error())
		return err
	}

	f, err := os.Create(output + "/index.json")
	if err != nil {
		return err
	}
	f.Write(b)
	f.Sync()
	f.Close()
	return nil
}
