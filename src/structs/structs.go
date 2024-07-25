package structs

type PinataResponsePinsList struct {
	Count int          `json:"count"`
	Rows  []PinataPins `json:"rows"`
}

type PinataPins struct {
	ID           string `json:"id"`
	IpfsPinHash  string `json:"ipfs_pin_hash"`
	URLPin       string `json:"-"`
	Size         int    `json:"size"`
	UserID       string `json:"user_id"`
	DatePinned   string `json:"date_pinned"`
	DateUnpinned string `json:"date_unpinned"`
	Metadata     struct {
		Name      string `json:"name"`
		Keyvalues string `json:"keyvalues"`
	} `json:"metadata"`
	MimeType      string `json:"mime_type"`
	NumberOfFiles int    `json:"number_of_files"`
}
