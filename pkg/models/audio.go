package models

type AudioInfo struct {
	ID        int
	FileName  string
	AudioName string
	Artist    string
}

type AudioFile struct {
	Info AudioInfo
	File []byte
}
