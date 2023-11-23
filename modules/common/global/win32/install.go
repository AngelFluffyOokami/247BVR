//go:build windows

package win32

import (
	"fmt"
	"os/user"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func installService(name, desc string) error {
	exepath := "C:\\Program Files\\HSVRUSB\\HSVRUSB.exe"

	TextLog <- "Establishing connection to the Service Control Manager..."
	m, err := mgr.Connect()

	if err != nil {
		return err
	}

	TextLog <- "Connected to Service Control Manager"
	defer m.Disconnect()

	TextLog <- "Checking if service exists."
	_, err = m.OpenService(name)
	if err == nil {
		TextLog <- "Service Exists, ensuring event log sources are set up..."

		err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)

		if err != nil {

			TextLog <- err.Error()
			TextLog <- "Assuming event log source exists, testing..."
			elog, err = eventlog.Open(name)
			if err != nil {
				TextLog <- err.Error()
				return err
			}

			TextLog <- "Verifying event log source is working"
			if err = elog.Info(1, "Event Log health test, ignore."); err != nil {
				TextLog <- err.Error()
				TextLog <- "error setting up event log source"
				return err
			}

			TextLog <- "Event log works"
			elog.Info(1, "HSVR USB 2.0 Installation succesfully completed.")
			TextLog <- "HSVR USB 2.0 Installation succesfully completed"

			return nil
		} else {
			TextLog <- "Event log source does not exist, but has just been added, testing..."
			elog, err = eventlog.Open(name)
			if err != nil {

				TextLog <- err.Error()

				TextLog <- "error setting up event log source"

				return err
			}
			if err = elog.Info(1, "Event Log health test, ignore."); err != nil {
				TextLog <- err.Error()
				TextLog <- "error setting up event log source"
				return err
			}
			TextLog <- "Event log works"
			elog.Info(1, "HSVR USB 2.0 Installation succesfully completed.")
			TextLog <- "HSVR USB 2.0 Installation succesfully completed"
			return nil

		}

	}

	TextLog <- "Service does not exist... Continuing."
	TextLog <- "Creating service..."

	usr, _ := user.Lookup("HSVRUSB")
	s, err := m.CreateService(name, exepath, mgr.Config{StartType: mgr.StartAutomatic, Description: desc, ServiceStartName: usr.Username}, "is", "auto-started")
	if err != nil {
		TextLog <- err.Error()
		return err
	}

	TextLog <- "Service Created"
	defer s.Close()

	TextLog <- "Setting up Event Log Source..."
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)

	if err != nil {
		TextLog <- err.Error()
		TextLog <- "Error is /probably/ benign, leaving error unhandled"
	}
	TextLog <- "Event Log Source set up"

	TextLog <- "Verifying event log source registered succesfully"
	elog, err = eventlog.Open(name)

	if err != nil {
		TextLog <- err.Error()
		return err
	}
	TextLog <- "Event log works"
	elog.Info(1, "HSVR USB 2.0 Installation succesfully completed.")
	TextLog <- "HSVR USB 2.0 Installation succesfully completed"
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}
