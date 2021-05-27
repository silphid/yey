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

func ContainerName(contextsPath string, context Context) string {
	return fmt.Sprintf(
		"yey-%s-%s-%s-%s",
		filepath.Base(filepath.Dir(contextsPath)),
		context.Name,
		hash(contextsPath),
		hash(context.String()),
	)
}
