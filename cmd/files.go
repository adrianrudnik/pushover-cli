package cmd

import (
	"errors"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"os"
)

func parseAttachment(path string) (ftype types.Type, err error) {
	file, err := os.Open(path)

	if err != nil {
		return types.Type{}, err
	}

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		file.Close()
		return types.Type{}, err
	}

	if !filetype.IsImage(head) {
		file.Close()
		return types.Type{}, errors.New("file is not an image")
	}

	// Extract type
	ftype, err = filetype.Get(head)
	if err != nil {
		file.Close()
		return types.Type{}, err
	}

	return ftype, file.Close()
}
