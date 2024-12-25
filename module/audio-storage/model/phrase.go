package model

type Phrase struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	OriginalFileName string `json:"original_filename"`
	FilePath         string `json:"file_path"`
	Content          []byte `json:"content"`
}
