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
