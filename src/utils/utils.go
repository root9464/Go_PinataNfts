package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"root/src/structs"

	"github.com/joho/godotenv"
)

func TestAuthentication() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	JWT := os.Getenv("USER_TOKEN_JWT")
	HOST := os.Getenv("PINATA_HOST")

	req, err := http.NewRequest("GET", HOST+"/testAuthentication", nil)
	if err != nil {
		log.Fatal("error creating request", err)
	}

	req.Header.Set("Authorization", "Bearer "+JWT)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error sending request", err)
	}

	defer resp.Body.Close()
	status := resp.StatusCode
	if status == 200 {
		fmt.Println("testAuthentication: ✅")
	} else {
		fmt.Println("testAuthentication: ❌", status)
	}
}

func GetPinataResponseFuncs() (func(), func() (*structs.PinataResponsePinsList, string, error)) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	JWT := os.Getenv("USER_TOKEN_JWT")
	HOST := os.Getenv("PINATA_HOST")

	getPinataResponse := func() (*structs.PinataResponsePinsList, string, error) {
		client := new(http.Client)
		req, err := http.NewRequest("GET", HOST+"/pinList", nil)
		if err != nil {
			log.Fatal("error creating request", err)
		}

		req.Header.Set("Authorization", "Bearer "+JWT)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("error sending request", err)
		}
		defer resp.Body.Close()

		var pinataResponse structs.PinataResponsePinsList
		err = json.NewDecoder(resp.Body).Decode(&pinataResponse)
		if err != nil {
			log.Fatal("error decoding response", err)
		}

		for i := range pinataResponse.Rows {
			pinataResponse.Rows[i].URL = "https://gateway.pinata.cloud/ipfs/" + pinataResponse.Rows[i].IpfsPinHash
		}

		jsonResponse, err := json.Marshal(pinataResponse)
		if err != nil {
			fmt.Println("error convert to json", err)
		}

		return &pinataResponse, string(jsonResponse), nil

	}

	printPinataResponse := func() {
		pinataResponse, _, err := getPinataResponse()
		if err != nil {
			log.Fatal("error getting pinata response", err)
		}

		fmt.Printf("Count: %d\n", pinataResponse.Count)
		for _, row := range pinataResponse.Rows {
			fmt.Printf("ID: %s\n", row.ID)
			fmt.Printf("IpfsPinHash: %s\n", row.IpfsPinHash)
			fmt.Printf("URL: https://gateway.pinata.cloud/ipfs/%s\n", row.IpfsPinHash)
			fmt.Printf("Size: %d\n", row.Size)
			fmt.Printf("UserID: %s\n", row.UserID)
			fmt.Printf("DatePinned: %s\n", row.DatePinned)
			fmt.Printf("DateUnpinned: %s\n", row.DateUnpinned)
			fmt.Printf("Metadata:\n")
			fmt.Printf("  Name: %s\n", row.Metadata.Name)
			fmt.Printf("  Keyvalues: %s\n", row.Metadata.Keyvalues)
			fmt.Printf("MimeType: %s\n", row.MimeType)
			fmt.Printf("NumberOfFiles: %d\n", row.NumberOfFiles)
			fmt.Println()
		}
	}

	return printPinataResponse, getPinataResponse
}
