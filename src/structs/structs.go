package structs

type PinataResponse struct {
	Count int `json:"count"`
	Rows  []struct {
		ID           string `json:"id"`
		IpfsPinHash  string `json:"ipfs_pin_hash"`
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
	} `json:"rows"`
}
