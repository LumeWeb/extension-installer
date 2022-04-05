package platform

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
	"golang.org/x/sys/windows/registry"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	sysruntime "runtime"
	"strconv"
	"strings"
)

func StartInstall(app *App) {
	currentApp = app
	updateStatus("Downloading extension..")

	extension, err := DownloadExtension()
	if err != nil {
		return
	}

	installDir := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Lume Web")
	extensionDestFile := filepath.Join(installDir, "extension.crx")

	_ = os.MkdirAll(installDir, os.ModePerm)
	extData, _ := ioutil.ReadFile(extension)
	_ = ioutil.WriteFile(extensionDestFile, extData, os.ModePerm)

	_ = os.Remove(extension)

	manifest, err := GetExtensionInfo(extensionDestFile)

	if err != nil {
		return
	}

	updateStatus(fmt.Sprintf("Installing extension version %s..\n", manifest.Version))

	installExtensionForBrowser("Google\\Chrome", extensionDestFile, manifest)
	installExtensionForBrowser("BraveSoftware\\Brave", extensionDestFile, manifest)

	deleteProfileUninstallSetting(manifest.Id, "Google", "Chrome")
	deleteProfileUninstallSetting(manifest.Id, "BraveSoftware", "Brave-Browser")

	setInstallState(2)

	updateStatus("Done")
}

func installExtensionForBrowser(registryPrefix string, file string, manifest *Manifest) {
	arch := ""

	if sysruntime.GOARCH == "amd64" {
		arch = "WOW6432Node\\"
	}

	allowList := fmt.Sprintf("Software\\%sPolicies\\%s\\ExtensionInstallAllowlist", arch, registryPrefix)
	extensionKey := fmt.Sprintf("Software\\%s%s\\Extensions\\%s", arch, registryPrefix, manifest.Id)

	err := ensureRegistryPathExists(allowList)
	if err != nil {
		return
	}
	err = ensureRegistryPathExists(extensionKey)
	if err != nil {
		return
	}

	key, _ := maybeCreateKey(allowList, true)
	extList, _ := key.ReadValueNames(-1)

	extFound := false

	for _, extIndex := range extList {
		extId, _, _ := key.GetStringValue(extIndex)
		if extId == manifest.Id {
			extFound = true
		}
	}

	if !extFound {
		_ = key.SetStringValue(strconv.Itoa(len(extList)), manifest.Id)
	}

	_ = key.Close()

	key, _ = maybeCreateKey(extensionKey, true)
	err = key.SetStringValue("path", file)
	err = key.SetStringValue("version", manifest.Version)

	_ = key.Close()

}

func ensureRegistryPathExists(path string) error {
	parts := strings.Split(path, "\\")
	for index, _ := range parts {
		partsPath := strings.Join(parts[:index+1], "\\")
		_, err := maybeCreateKey(partsPath, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func maybeCreateKey(path string, ret bool) (retkey *registry.Key, error error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
		if err != nil {
			return nil, err
		}
	}
	if ret {
		return &key, nil
	}

	_ = key.Close()
	return nil, nil
}

func deleteProfileUninstallSetting(extensionId string, vendorName string, vendorBrowserName string) {
	profilePrefLocation := filepath.Join(os.Getenv("LOCALAPPDATA"), vendorName, vendorBrowserName, "User Data", "Default", "Preferences")

	exist, _ := fileExists(profilePrefLocation)
	if !exist {
		return
	}

	prefFile, err := ioutil.ReadFile(profilePrefLocation)

	if err != nil {
		return
	}

	uninstallPath := "extensions.external_uninstalls"

	uninstalls := gjson.GetBytes(prefFile, uninstallPath)

	if !uninstalls.Exists() {
		return
	}

	newUninstalls := make([]string, 0)

	uninstalls.ForEach(func(key, value gjson.Result) bool {
		newUninstalls = append(newUninstalls, value.String())
		return true
	})

	if !slices.Contains(newUninstalls, extensionId) {
		return
	}

	foundExtIndex := slices.Index(newUninstalls, extensionId)

	newUninstalls[foundExtIndex] = newUninstalls[len(newUninstalls)-1]
	newUninstalls[len(newUninstalls)-1] = ""
	newUninstalls = newUninstalls[:len(newUninstalls)-1]

	prefFile, _ = sjson.SetBytes(prefFile, uninstallPath, newUninstalls)
	_ = ioutil.WriteFile(profilePrefLocation, prefFile, fs.ModePerm)
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
