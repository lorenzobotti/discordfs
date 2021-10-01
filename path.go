package discordfs

import (
	"path"
	"strings"
)

func CleanPath(p string) string {
	cleaned := path.Clean(p)
	if cleaned == "." {
		return "/"
	}

	if !strings.HasPrefix(cleaned, "/") {
		return "/" + cleaned
	}

	if cleaned == "/" {
		return cleaned
	} else {
		return strings.TrimSuffix(cleaned, "/")

	}
}

func pathElements(p string) []string {
	p = path.Clean(p)
	out := []string{}

	for {
		out = append(out, path.Base(p))
		p = path.Dir(p)

		if p == "." || p == "/" {
			break
		}
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out
}

func topFolder(p string) string {
	elements := pathElements(p)
	if len(elements) <= 1 {
		return "/"
	} else {
		return elements[0]
	}
}

func isInFolder(p string) bool {
	cleaned := path.Clean(p)

	return len(pathElements(cleaned)) != 1
}
