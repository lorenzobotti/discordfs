package discordfs

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
)

type CompressionLevel int16
type CompressionType int16
type Compression struct {
	Level CompressionLevel `json:"level"`
	Type  CompressionType  `json:"type"`
}

const (
	NoCompression      CompressionLevel = flate.NoCompression
	BestSpeed          CompressionLevel = flate.BestSpeed
	BestCompression    CompressionLevel = flate.BestCompression
	DefaultCompression CompressionLevel = flate.DefaultCompression
	HuffmanOnly        CompressionLevel = flate.HuffmanOnly
)

const (
	EachChunk CompressionType = iota << 8
	WholeFile
)

// these next functions use a switch case instead of a map
// because i suspect the map is slower with such few entries
// todo: research if it's actually slower

func (c CompressionLevel) String() string {
	switch c {
	case NoCompression:
		return "none"
	case BestSpeed:
		return "best_speed"
	case BestCompression:
		return "best_compression"
	case DefaultCompression:
		return "default"
	case HuffmanOnly:
		return "huffman"
	}

	return ""
}

func (c CompressionType) String() string {
	switch c {
	case EachChunk:
		return "each_chunk"
	case WholeFile:
		return "whole_file"
	}

	return ""
}

func (c CompressionLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.String() + `"`), nil
}

func (c CompressionType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.String() + `"`), nil
}

func (c *CompressionLevel) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"none"`:
		*c = NoCompression
	case `"best_speed"`:
		*c = BestSpeed
	case `"best_compression"`:
		*c = BestCompression
	case `"default"`:
		*c = DefaultCompression
	case `"huffman"`:
		*c = HuffmanOnly
	default:
		return fmt.Errorf("unknown value: %s", string(input))
	}
	return nil
}

func (c *CompressionType) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"each_chunk"`:
		*c = EachChunk
	case `"whole_file"`:
		*c = WholeFile
	default:
		return fmt.Errorf("unknown value: %s", string(input))
	}
	return nil
}

// readAndCompress reads all the data from `src`, compresses it,
// and returns a bytes.Buffer with the gzipped data
//
// CAREFUL: this copies the WHOLE READER into memory, without caring
// too much about its size. it won't complain if you give it a
// 5TB file, but your pc sure will
func readAndCompressInto(dst io.Writer, src io.Reader, level CompressionLevel) error {
	compressor, err := gzip.NewWriterLevel(dst, int(level))
	if err != nil {
		return fmt.Errorf("invalid compression level: %w", err)
	}

	_, err = io.Copy(compressor, src)
	if err != nil {
		return fmt.Errorf("can't copy from compressor: %w", err)
	}

	err = compressor.Close()
	if err != nil {
		return fmt.Errorf("can't close compressor: %w", err)
	}

	return nil
}

func readAndDecompressInto(dst io.Writer, src io.Reader) error {
	decompressor, err := gzip.NewReader(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, decompressor)
	if err != nil {
		return fmt.Errorf("can't copy from decompressor: %w", err)
	}

	err = decompressor.Close()
	if err != nil {
		return fmt.Errorf("can't close gzip reader: %w", err)
	}

	_, err = io.Copy(dst, decompressor)
	if err != nil {
		return fmt.Errorf("can't copy from decompressor: %w", err)
	}

	return nil
}
