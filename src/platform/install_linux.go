package platform

import (
	"github.com/sqweek/dialog"
	"os/user"
	"strconv"
)

func StartInstall() {
	updateStatus(STATUS_DOWNLOADING)

	extension, err := DownloadExtension()
	if err != nil {
		return
	}

}

func IsAdmin() bool {
	u, _ := user.Current()
	uid, _ := strconv.Atoi(u.Uid)
	return uid == 0
}

func ReLaunchAsAdmin() {
	dialog.Message("Please run as root").Title("Error").Info()
}
