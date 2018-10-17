package models

type Question struct {
	QId     int                      `json:"qid"`
	Title   string                   `json:"title"`
	Type    int8                     `json:"type"`
	Options []map[string]interface{} `json:"options"`
}
