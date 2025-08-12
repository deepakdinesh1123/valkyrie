package common

import (
	"archive/tar"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/db/jsonschema"
)

func CreateTarArchive(files map[string]string, additional_files jsonschema.ExecReqFiles, dir string) (string, error) {
	tarFilePath := filepath.Join(dir, fmt.Sprintf("%d.tar", time.Now().UnixNano()))
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return "", err
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(content)),
			Mode: 0744,
		}); err != nil {
			return "", err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return "", err
		}
	}

	for _, file := range additional_files.Files {
		decodedContent, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 content for %s: %w", file.Name, err)
		}
		if err := tw.WriteHeader(&tar.Header{
			Name: file.Name,
			Size: int64(len(decodedContent)),
			Mode: 0744,
		}); err != nil {
			return "", err
		}
		if _, err := tw.Write(decodedContent); err != nil {
			return "", err
		}
	}

	return tarFilePath, nil
}
