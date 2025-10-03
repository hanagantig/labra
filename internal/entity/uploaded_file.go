package entity

import (
	"crypto/sha256"
	"encoding/hex"
)

type UploadedFileStatus string

const (
	UploadedFileStatusNew         UploadedFileStatus = "new"
	UploadedFileStatusRecognizing UploadedFileStatus = "recognizing"
	UploadedFileStatusRecognized  UploadedFileStatus = "recognized"
	UploadedFileStatusDuplicated  UploadedFileStatus = "duplicated"
)

type UploadedFile struct {
	ID          int
	bytes       []byte
	Fingerprint string
	FileID      string
	PipelineID  string
	UserID      int
	ProfileID   int
	FileType    string
	Source      string
	Status      UploadedFileStatus
	Details     string
}

func NewUploadedFile(bytes []byte) UploadedFile {
	f := UploadedFile{
		bytes:  bytes,
		Status: UploadedFileStatusNew,
	}

	f.Fingerprint = f.GenerateFingerprint()
	return f
}

func (u UploadedFile) GenerateFingerprint() string {
	if len(u.bytes) == 0 {
		return ""
	}

	hash := sha256.New()
	hash.Write(u.bytes)

	return hex.EncodeToString(hash.Sum(nil))
}

func (u UploadedFile) Bytes() []byte {
	return u.bytes
}
