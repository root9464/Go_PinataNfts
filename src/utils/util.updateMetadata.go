package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"root/src/structs"
	"sort"
	"strings"
	"sync"
)

const (
	jsonDir  = "build/json/"
	imageDir = "build/images/"
)

func ResponseData() *bytes.Buffer {
	r, err := zip.OpenReader("build.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	images := new([]structs.ImageData)

	for _, file := range r.File {
		path := filepath.ToSlash(file.Name)

		if strings.HasPrefix(path, imageDir) && filepath.Ext(path) == ".png" {
			*images = append(*images, structs.ImageData{
				FileName: processImage(file),
				Data:     nil,
			})
		} else if strings.HasPrefix(path, jsonDir) && filepath.Ext(path) == ".json" {
			data, err := processJSON(file)
			if err != nil {
				log.Println(err)
				continue
			}

			filename := strings.TrimSuffix(filepath.Base(file.Name), ".json") + ".png"
			for i, image := range *images {
				if image.FileName == filename {
					(*images)[i].Data = data
					break
				}
			}
			*images = append(*images, structs.ImageData{
				FileName: filename,
				Data:     data,
			})
		}
	}

	jsonData, err := json.Marshal(&structs.Response{Images: *images})
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(jsonData)
}

func PrintZipData() {
	r, err := zip.OpenReader("build.zip")
	if err != nil {
		log.Fatal("error opening zip file", err)
	}
	defer r.Close()

	jsonDataMap := make(map[string]interface{})
	jsonMutex := new(sync.Mutex)

	imageFiles := make(map[string]bool)
	imageMutex := new(sync.Mutex)

	wg := new(sync.WaitGroup)
	var errChan = make(chan error)

	for _, file := range r.File {
		path := filepath.ToSlash(file.Name)

		if strings.HasPrefix(path, jsonDir) && filepath.Ext(path) == ".json" {
			wg.Add(1)
			go func(file *zip.File) {
				defer wg.Done()
				data, err := processJSON(file)
				if err != nil {
					errChan <- err
					return
				}

				filename := strings.TrimSuffix(filepath.Base(file.Name), ".json")
				jsonMutex.Lock()
				jsonDataMap[filename] = data
				jsonMutex.Unlock()
			}(file)
		}

		if strings.HasPrefix(path, imageDir) && filepath.Ext(path) == ".png" {
			wg.Add(1)
			go func(file *zip.File) {
				defer wg.Done()
				filename := processImage(file)
				imageMutex.Lock()
				imageFiles[filename] = true
				imageMutex.Unlock()
			}(file)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		log.Println("error:", err)
	}

	filenames := make([]string, 0, len(imageFiles))
	for filename := range imageFiles {
		filenames = append(filenames, filename)
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		filenameWithoutExt := strings.TrimSuffix(filename, ".png")
		jsonMutex.Lock()
		data, ok := jsonDataMap[filenameWithoutExt]
		jsonMutex.Unlock()
		if ok {
			fmt.Printf("Имя файла: %s.png\n", filenameWithoutExt)
			fmt.Println("Json:")
			for key, value := range data.(map[string]interface{}) {
				fmt.Printf("     %s = %v\n", key, value)
			}
			fmt.Println()
		} else {
			fmt.Printf("Имя файла: %s.png\nJSON объект не найден\n\n", filenameWithoutExt)
		}
	}
}
