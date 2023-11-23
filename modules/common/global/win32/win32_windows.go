//go:build windows

package win32

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/tjarratt/babble"
)

func AddUser() bool {

	TextLog <- "Loading Windows Powershell..."

	TextLog <- "Creating powershell runspace..."

	TextLog <- "Runspace created"
	pwd := newPassword()

	TextLog <- "Attempting user create..."
	TextLog <- "New-LocalUser -AccountNeverExpires -Description \"HSVR ELO statistics bot User\" -Disabled -FullName \"HSVR ELO Statistics Service User\" -Name \"247bvr\" -Password \"" + pwd + "\" -PasswordNeverExpires -Confirm"

	//	env := `[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";$APP_PATH\bin", [System.EnvironmentVariableTarget]::User)`
	cmd := `New-LocalUser -AccountNeverExpires -Description "HSVR ELO statistics bot User" -Disabled -FullName "HSVR ELO Statistics Service User" -Name "247bvr" -Password (ConvertTo-SecureString "` + pwd + `" -AsPlainText -Force) -PasswordNeverExpires`
	exitcmd := `exit`
	cmdps := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmdps.StdinPipe()
	if err != nil {
		TextLog <- err.Error()
		return false
	}
	cmdps.StdinPipe()
	fmt.Fprintln(stdin, cmd+"\n"+exitcmd)
	out, err := cmdps.CombinedOutput()
	if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- string(out)
	}

	if err != nil {
		TextLog <- err.Error()
		return false
	}
	return true
}
func newPassword() string {
	TextLog <- "Generating password."
	babbler := babble.NewBabbler()
	babbler.Count = 4
	babbler.Separator = "-"
	pwd := babbler.Babble()
	return pwd
}

func populatePF() bool {
	TextLog <- "mkdir C:\\Program Files\\247bvr\\"
	err := os.MkdirAll("C:\\Program Files\\247bvr\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists"
		return true
	} else if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- "Permission bits set to 0755"
		return true
	}
}

func populateAD() bool {
	TextLog <- "mkdir C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\"
	err := os.MkdirAll("C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists"
		return true
	} else if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- "Permission bits set to 0755"
		return true
	}
}

func populateADConf() bool {
	TextLog <- "mkdir C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\config\\"
	err := os.MkdirAll("C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\config\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists"
		return true
	} else if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- "Permission bits set to 0755"
		return true
	}
}

func populateADDb() bool {
	TextLog <- "mkdir C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\vtol.vr\\"
	err := os.MkdirAll("C:\\Users\\247bvr\\AppData\\Roaming\\247bvr\\vtol.vr\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists"
		return true
	} else if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- "Permission bits set to 0755"
		return true
	}
}
func populatePaths() bool {

	populateAD()
	populateADConf()
	populateADDb()
	populatePF()

	TextLog <- "Attempting to ascert current process binary path..."
	execPath, err := os.Executable()
	TextLog <- execPath

	if err != nil {
		TextLog <- err.Error()
		return false
	}

	TextLog <- "Atempting to load current binary into memory..."
	sourceFileStat, err := os.Stat(execPath)
	if err != nil {
		TextLog <- err.Error()
		return false
	}

	if !sourceFileStat.Mode().IsRegular() {
		TextLog <- "Error encountered loading file: File is not regular."
		return false
	}
	source, err := os.Open(execPath)
	if err != nil {
		TextLog <- err.Error()
		return false
	}

	defer source.Close()
	TextLog <- "Binary loaded into memory"
	TextLog <- "Attempting to copy into place..."

	dest, err := os.Create("C:\\Program Files\\247bvr\\247bvr.exe")
	if err != nil {
		TextLog <- err.Error()
		return false
	}
	defer dest.Close()

	TextLog <- "File created, copying contents"

	nbytes, err := io.Copy(dest, source)
	if err != nil {
		TextLog <- err.Error()
		return false
	}

	TextLog <- "Finished copying " + fmt.Sprint(nbytes) + " to C:\\Program Files\\247bvr\\247bvr.exe"
	return true
}
