package yey

import (
	"encoding/hex"
	"fmt"
	"hash/crc64"
	"io"
	"path/filepath"
)

func hash(value string) string {
	hasher := crc64.New(crc64.MakeTable(crc64.ECMA))
	io.WriteString(hasher, value)
	return hex.EncodeToString(hasher.Sum(nil))
}

func ContainerName(path string, context Context) string {
	return fmt.Sprintf(
		"yey-%s-%s-%s-%s",
		filepath.Base(filepath.Dir(path)),
		hash(path),
		context.Name,
		hash(context.String()),
	)
}
