//go:build windows

package win32

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"

	"github.com/bradfitz/iter"
	"github.com/tjarratt/babble"
)

func AddUser() bool {

	TextLog <- "Ensuring user not exists already..."
	_, err := user.Lookup("HSVRUSB")
	if err == nil {
		TextLog <- "User Already exists, skipping step"
		return true
	}
	TextLog <- "Loading Windows Powershell..."

	TextLog <- "Creating powershell runspace..."

	TextLog <- "Runspace created"
	pwd := newPassword()

	censor := ""
	for range iter.N(len(pwd)) {
		censor += "*"
	}
	TextLog <- "Attempting user create..."

	cmd := `New-LocalUser -AccountNeverExpires -Description "HSVR USB 2.0 Service User" -Disabled -FullName "HSVR UselessStatisticsBot 2.0" -Name "HSVRUSB" -Password (ConvertTo-SecureString "` + pwd + `" -AsPlainText -Force) -PasswordNeverExpires`

	exitcmd := `exit`

	TextLog <- "Summoning child powershell process..."
	cmdps := exec.Command("powershell", "-nologo", "-noprofile")
	TextLog <- "Process Summoned"

	TextLog <- "Adding User..."
	stdin, err := cmdps.StdinPipe()
	TextLog <- "Opening Stdin pipe...."
	if err != nil {
		TextLog <- err.Error()
		return false
	}
	cmdcensor := `New-LocalUser -AccountNeverExpires -Description "HSVR ELO statistics bot User" -Disabled -FullName "HSVR ELO Statistics Service User" -Name "HSVRUSB" -Password (ConvertTo-SecureString "` + censor + `" -AsPlainText -Force) -PasswordNeverExpires`
	TextLog <- cmdcensor
	TextLog <- "Collecting results..."
	fmt.Fprintln(stdin, cmd+"\n"+exitcmd)
	_, err = cmdps.CombinedOutput()
	if err != nil {
		TextLog <- err.Error()
		return false
	} else {
		TextLog <- "Created user HSVRUSB with password " + censor + " as ADMINISTRATOR"
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
	TextLog <- "mkdir C:\\Program Files\\HSVRUSB\\"
	err := os.MkdirAll("C:\\Program Files\\HSVRUSB\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists, skipping step"
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
	TextLog <- "mkdir C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\"
	err := os.MkdirAll("C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists, skipping step"
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
	TextLog <- "mkdir C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\config\\"
	err := os.MkdirAll("C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\config\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists, skipping step"
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
	TextLog <- "mkdir C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\vtol.vr\\"
	err := os.MkdirAll("C:\\Users\\HSVRUSB\\AppData\\Roaming\\HSVRUSB\\vtol.vr\\", 0755)

	if os.IsExist(err) {
		TextLog <- "Folder already exists, skipping step"
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

	dest, err := os.Create("C:\\Program Files\\HSVRUSB\\HSVRUSB.exe")
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

	TextLog <- "Finished copying " + fmt.Sprint(nbytes) + " bytes to C:\\Program Files\\HSVRUSB\\HSVRUSB.exe"
	return true
}
