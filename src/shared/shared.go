package shared

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mediabuyerbot/go-crx3/pb"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const extensionUrl = "https://fileportal.org/AQAHCZ-dpgFKC91APG-VdtmTShKbAovBspv67PR_Ot2lWg"

type Manifest struct {
	Version string `json:"version"`
}

func DownloadExtension() (path string, error error) {
	resp, err := retryablehttp.Get(extensionUrl)

	dir, err := ioutil.TempDir("", "lume_")
	if err != nil {
		log.Fatal(err)
	}

	extensionPath := filepath.Join(dir, "extension.crx")

	file, e := os.Create(extensionPath)
	if e != nil {
		panic(e)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return extensionPath, nil

}

func GetExtensionVersion(path string) (version string, error error) {
	crx, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var (
		headerSize = binary.LittleEndian.Uint32(crx[8:12])
		metaSize   = uint32(12)
		v          = crx[metaSize : headerSize+metaSize]
		header     pb.CrxFileHeader
		signedData pb.SignedData
	)

	if err := proto.Unmarshal(v, &header); err != nil {
		return "", err
	}
	if err := proto.Unmarshal(header.SignedHeaderData, &signedData); err != nil {
		return "", err
	}

	data := crx[len(v)+int(metaSize):]
	bytesReader := bytes.NewReader(data)
	size := int64(len(data))

	reader, err := zip.NewReader(bytesReader, size)
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
		if file.Name == "manifest.json" {
			var manifest Manifest

			manifestData, _ := file.Open()
			defer func(manifestData io.ReadCloser) {
				_ = manifestData.Close()
			}(manifestData)

			byteValue, _ := ioutil.ReadAll(manifestData)
			err := json.Unmarshal(byteValue, &manifest)
			if err != nil {
				return "", err
			}

			return manifest.Version, nil
		}
	}

	return "", nil
}
