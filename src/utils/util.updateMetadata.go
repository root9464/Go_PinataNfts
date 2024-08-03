package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"root/src/structs"
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
	imagesMutex := new(sync.Mutex)

	jsonDataMap := make(map[string]interface{})
	jsonMutex := new(sync.Mutex)

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
				imagesMutex.Lock()
				*images = append(*images, structs.ImageData{
					FileName: filename,
					Data:     nil,
				})
				imagesMutex.Unlock()
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

	imagesMutex.Lock()
	for i, image := range *images {
		filenameWithoutExt := strings.TrimSuffix(image.FileName, ".png")
		jsonMutex.Lock()
		data, ok := jsonDataMap[filenameWithoutExt]
		jsonMutex.Unlock()
		if ok {
			(*images)[i].Data = data.(*structs.JSONData)

		}
	}
	imagesMutex.Unlock()

	jsonData, err := json.Marshal(&structs.Response{Images: *images})
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(jsonData)
}

func PrintData() {
	r, err := zip.OpenReader("build.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	images := []structs.ImageData{}
	jsonDataMap := make(map[string]interface{})

	for _, file := range r.File {
		path := filepath.ToSlash(file.Name)

		if strings.HasPrefix(path, jsonDir) && filepath.Ext(path) == ".json" {
			data, err := processJSON(file)
			if err != nil {
				log.Println("error:", err)
				continue
			}
			jsonDataMap[strings.TrimSuffix(filepath.Base(file.Name), ".json")] = data
		}

		if strings.HasPrefix(path, imageDir) && filepath.Ext(path) == ".png" {
			images = append(images, structs.ImageData{
				FileName: processImage(file),
				Data:     nil,
			})
		}
	}

	for i, image := range images {
		if data, ok := jsonDataMap[strings.TrimSuffix(image.FileName, ".png")]; ok {
			images[i].Data = data.(*structs.JSONData)
		}
	}

	fmt.Println("Images:")
	for _, image := range images {
		fmt.Printf("  - %s\n", image.FileName)
		if image.Data != nil {
			fmt.Println("    Data:")
			fmt.Printf("      Name: %s\n", image.Data.Name)
			fmt.Printf("      Description: %s\n", image.Data.Description)
			fmt.Printf("      Image: %s\n", image.Data.Image)
			fmt.Printf("      DNA: %s\n", image.Data.DNA)
			fmt.Printf("      Edition: %d\n", image.Data.Edition)
			fmt.Printf("      Date: %d\n", image.Data.Date)
			fmt.Println("      Attributes:")
			for _, attr := range image.Data.Attributes {
				fmt.Printf("        - %s: %s\n", attr.TraitType, attr.Value)
			}
			fmt.Printf("      Compiler: %s\n", image.Data.Compiler)
		}
	}
}
