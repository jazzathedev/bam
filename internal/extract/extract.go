package extract

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jazzathedev/bam/internal/plugin"
)

func Extractor(pluginConfig plugin.PluginConfig, sourcePath, destDir string) error {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("Unable to make destination directory %s: %w", destDir, err)
	}

	var archiveType string

	if strings.HasSuffix(sourcePath, ".tar.gz") {
		archiveType = ".tgz"
	} else {
		archiveType = filepath.Ext(sourcePath)
	}

	if archiveType == "" {
		return fmt.Errorf("Unable to determine archive file type")
	}

	switch archiveType {
	case ".zip":
		return ExtractZip(sourcePath, destDir, pluginConfig.Install.StripComponents)
	}

	return fmt.Errorf("Archive type not supported, not extracted.")
}

func ExtractZip(sourcePath, destDir string, strip bool) error {
	reader, err := zip.OpenReader(sourcePath)
	if err != nil {
		return fmt.Errorf("Unable to read zip archive: %w", err)
	}

	defer reader.Close()

	// https://github.com/snyk/zip-slip-vulnerability
	// https://stackoverflow.com/a/24792688
	// I still love you StackOverflow even if your userbase is gone </3
	for _, file := range reader.File {
		fileName := file.Name

		// Strip the archive-name-as-folder `toml:strip_components`
		if strip {
			fileName = strings.SplitN(fileName, "/", 2)[1]

			if fileName == "" {
				continue
			}
		}

		fpath := filepath.Join(destDir, fileName)
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("Illegal file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}

		outFile.Close()
		rc.Close()
	}

	return nil
}
