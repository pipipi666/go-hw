package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, errInfo := os.Stat(fromPath)

	if errInfo != nil {
		return ErrUnsupportedFile
	}

	size := fileInfo.Size()

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	reader, errReader := os.Open(fromPath)

	if errReader != nil {
		return errReader
	}

	defer reader.Close()

	_, errSeek := reader.Seek(offset, io.SeekStart)

	if errSeek != nil {
		return errSeek
	}

	writer, errWriter := os.Create(toPath)

	if errWriter != nil {
		return errWriter
	}

	defer writer.Close()

	count := size

	if (size > limit) && (limit > 0) {
		count = limit
	}

	bar := pb.Full.Start64(count)
	defer bar.Finish()
	barReader := bar.NewProxyReader(reader)

	_, err := io.CopyN(writer, barReader, count)

	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}
