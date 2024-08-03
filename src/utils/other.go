package utils

import (
	"archive/zip"
	"encoding/json"
	"path/filepath"
	"root/src/structs"
)

func processImage(file *zip.File) string {
	return filepath.Base(file.Name)
}

func processJSON(file *zip.File) (*structs.JSONData, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data := new(structs.JSONData)
	err = json.NewDecoder(reader).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
