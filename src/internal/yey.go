package yey

import (
	"encoding/hex"
	"fmt"
	"hash/crc64"
	"io"
	"path/filepath"
	"regexp"
)

var (
	spaces            = regexp.MustCompile(`\s+`)
	special           = regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
	dashes            = regexp.MustCompile(`-+`)
	alphaNumericStart = regexp.MustCompile(`^[a-zA-Z0-9]`)
	trailingDashes    = regexp.MustCompile(`-+$`)
)

func hash(value string) string {
	hasher := crc64.New(crc64.MakeTable(crc64.ECMA))
	io.WriteString(hasher, value)
	return hex.EncodeToString(hasher.Sum(nil))
}

// ContainerName returns the container name to use for given yey rc path
// and context
func ContainerName(path string, context Context) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		ContainerPathPrefix(path),
		sanitizeContextName(context.Name),
		hash(context.String()),
	)
}

func sanitizeContextName(value string) string {
	return special.ReplaceAllString(value, "-")
}

func sanitizePathName(value string) string {
	value = spaces.ReplaceAllString(value, "_")
	value = special.ReplaceAllString(value, "")
	value = dashes.ReplaceAllString(value, "-")
	value = trailingDashes.ReplaceAllString(value, "")

	if value == "" {
		return ""
	}

	if !alphaNumericStart.MatchString(value) {
		value = "0" + value
	}

	return value
}

func ContainerPathPrefix(path string) string {
	pathBase := sanitizePathName(filepath.Base(filepath.Dir(path)))
	if pathBase == "" {
		return fmt.Sprintf("yey-%s", hash(path))
	}
	return fmt.Sprintf("yey-%s-%s", pathBase, hash(path))
}

func ImageName(dockerfile []byte) string {
	return fmt.Sprintf("yey-%s", hash(string(dockerfile)))
}
