package platform

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mediabuyerbot/go-crx3"
	"github.com/mediabuyerbot/go-crx3/pb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

var currentApp *App

const STATUS_DOWNLOADING string = "Downloading your personal gateway into web3..."
const STATUS_COMPLETE string = "All done, fully powered and ready to go!"

const extensionUrl = "https://fileportal.org/AQAHCZ-dpgFKC91APG-VdtmTShKbAovBspv67PR_Ot2lWg"

type Manifest struct {
	Version string `json:"version"`
	Id      string
}

func DownloadExtension() (path string, error error) {
	client := retryablehttp.NewClient()
	client.Logger = nil
	resp, err := client.Get(extensionUrl)
	dir, err := ioutil.TempDir("", "lume_")
	if err != nil {
		log.Fatal(err)
	}

	extensionPath := filepath.Join(dir, "extension.crx")

	file, e := os.Create(extensionPath)
	if e != nil {
		panic(e.(any))
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

func GetExtensionInfo(path string) (manifest *Manifest, err error) {
	version, err := getExtensionVersion(path)
	if err != nil {
		return nil, err
	}

	id, err := crx3.Extension(path).ID()
	if err != nil {
		return nil, err
	}

	return &Manifest{Id: id, Version: version}, nil
}

func getExtensionVersion(path string) (version string, error error) {
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

func InstructionsPrompt() {
	sendCommand("finished")
}

func NewApp() *App {
	return &App{}
}

func (app *App) Startup(ctx context.Context) {
	app.ctx = ctx
	currentApp = app
	runtime.EventsOn(app.ctx, "install", func(optionalData ...interface{}) {
		setInstallState(1)
		StartInstall()
	})
}

func updateStatus(status string) {
	sendCommand("status", status)
}

func setInstallState(state int) {
	sendCommand("setInstallState", state)
}
func sendCommand(status string, args ...interface{}) {
	runtime.EventsEmit(currentApp.ctx, status, args...)
}
