package platform

import (
	"github.com/lumeweb/extension-installer/src/shared"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"path/filepath"
)

func StartInstall() {
	extension, err := shared.DownloadExtension()
	if err != nil {
		return
	}

	installDir := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Lume Web")
	extenstionDestFile := filepath.Join(installDir, "extension.crx")

	_ = os.MkdirAll(installDir, os.ModePerm)
	extData, _ := ioutil.ReadFile(extension)
	_ = ioutil.WriteFile(extenstionDestFile, extData, os.ModePerm)

	version := shared.GetExtensionVersion(extension)

}

func installExtensionForBrowser(pathBase string, file string, version string) {

}
func createRegistryKey(path string) (registry.Key, func()) {
	var access uint32 = registry.ALL_ACCESS
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, access)
	if err != nil {
		panic(err)
	}

	return key, func() {
		var err error
		if err = key.Close(); err != nil {
			panic(err)
		}
	}
}
