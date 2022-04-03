package platform

import (
	"fmt"
	"github.com/lumeweb/extension-installer/src/shared"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func StartInstall() {
	extension, err := shared.DownloadExtension()
	if err != nil {
		return
	}

	installDir := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Lume Web")
	extensionDestFile := filepath.Join(installDir, "extension.crx")

	_ = os.MkdirAll(installDir, os.ModePerm)
	extData, _ := ioutil.ReadFile(extension)
	_ = ioutil.WriteFile(extensionDestFile, extData, os.ModePerm)

	manifest, err := shared.GetExtensionInfo(extensionDestFile)

	if err != nil {
		return
	}

	installExtensionForBrowser("Google\\Chrome", extensionDestFile, manifest)
	installExtensionForBrowser("BraveSoftware\\Brave", extensionDestFile, manifest)

}

func installExtensionForBrowser(registryPrefix string, file string, manifest *shared.Manifest) {
	allowList := fmt.Sprintf("Software\\Policies\\%s\\ExtensionInstallAllowlist", registryPrefix)

	arch := ""

	if runtime.GOARCH == "amd64" {
		arch = "Wow6432Node\\"
	}

	extensionKey := fmt.Sprintf("Software\\%s%s\\Extensions\\%s", registryPrefix, arch, manifest.Id)

	err := ensureRegistryPathExists(allowList)
	if err != nil {
		return
	}
	err = ensureRegistryPathExists(extensionKey)
	if err != nil {
		return
	}

	key, _ := maybeCreateKey(extensionKey, true)
	_ = key.SetStringValue("path", file)
	_ = key.SetStringValue("version", manifest.Version)

	_ = key.Close()

}

func ensureRegistryPathExists(path string) error {
	parts := strings.Split(path, "\\")
	for index, _ := range parts {
		path := strings.Join(parts[0:index], "\\")
		_, err := maybeCreateKey(path, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func maybeCreateKey(path string, ret bool) (retkey *registry.Key, error error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
		if err != nil {
			return nil, err
		}
	}
	if ret {
		return &key, nil
	}

	key.Close()
	return nil, nil
}
