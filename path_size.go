package code

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func humanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return strconv.FormatInt(size, 10) + "B"
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	value := float64(size) / unit
	i := 0
	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}

	return strconv.FormatFloat(value, 'f', 1, 64) + units[i]
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}

func dirSize(path string, recursive, all bool) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64

	for _, entry := range entries {
		if !all && isHidden(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			if !recursive {
				continue
			}
			nested, err := dirSize(filepath.Join(path, entry.Name()), recursive, all)
			if err != nil {
				return 0, err
			}
			total += nested
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		total += info.Size()
	}

	return total, nil
}

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	var size int64

	if !all && isHidden(fileInfo.Name()) {
		return "0B", nil
	}

	if fileInfo.IsDir() {
		size, err = dirSize(path, recursive, all)
		if err != nil {
			return "", err
		}
	} else {
		size = fileInfo.Size()
	}

	if human {
		return humanSize(size), nil
	}

	return strconv.FormatInt(size, 10) + "B", nil
}
