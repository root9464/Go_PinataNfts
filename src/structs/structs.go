package structs

type PinataResponsePinsList struct {
	Count int          `json:"count"`
	Rows  []PinataPins `json:"rows"`
}

type PinataPins struct {
	ID            string             `json:"id"`
	IpfsPinHash   string             `json:"ipfs_pin_hash"`
	Size          int                `json:"size"`
	URL           string             `json:"url"`
	UserID        string             `json:"user_id"`
	DatePinned    string             `json:"date_pinned"`
	DateUnpinned  string             `json:"date_unpinned"`
	Metadata      PinataPinsMetadata `json:"metadata"`
	MimeType      string             `json:"mime_type"`
	NumberOfFiles int                `json:"number_of_files"`
}

type PinataPinsMetadata struct {
	Name      string                 `json:"name"`
	Keyvalues map[string]interface{} `json:"keyvalues"`
}

type JSONData struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	DNA         string      `json:"dna"`
	Edition     int         `json:"edition"`
	Date        int64       `json:"date"`
	Attributes  []Attribute `json:"attributes"`
	Compiler    string      `json:"compiler"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type ImageData struct {
	FileName string    `json:"file_name"`
	Data     *JSONData `json:"data"`
}

type Response struct {
	Images []ImageData `json:"images"`
}
