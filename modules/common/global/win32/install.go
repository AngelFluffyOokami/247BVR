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
	s, err := m.OpenService(name)
	if err == nil {
		TextLog <- "Service Exists, ensuring event log sources are set up..."

		elog, err = eventlog.Open(name)

		if err != nil {

			TextLog <- err.Error()
			TextLog <- "Event log source probably not exists, attempting to add a new one"

			err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)

			if err != nil {
				s.Delete()
				TextLog <- err.Error()
				return fmt.Errorf("SetupEventLogSource() failed: %s", err)
			}

			TextLog <- "Verifying event log source registered succesfully"
			elog, err = eventlog.Open(name)
			if err != nil {
				s.Delete()

				TextLog <- err.Error()

				TextLog <- "Error setting up event log source, will now clean up and remove service."

				return fmt.Errorf("err setting up: %s", err)

			}
			TextLog <- "Event log works"
			elog.Info(1, "HSVR USB 2.0 Installation succesfully completed.")
			TextLog <- "HSVR USB 2.0 Installation succesfully completed"

			return nil
		}

		s.Close()
		TextLog <- "Service and Event log source exists"
		elog.Info(1, "HSVR USB 2.0 Installation succesfully completed.")
		TextLog <- "HSVR USB 2.0 Installation succesfully completed."

		return nil
	}

	TextLog <- "Service does not exist... Continuing."
	TextLog <- "Creating service..."

	usr, _ := user.Lookup("HSVRUSB")
	s, err = m.CreateService(name, exepath, mgr.Config{StartType: mgr.StartAutomatic, Description: desc, ServiceStartName: usr.Username}, "is", "auto-started")
	if err != nil {
		TextLog <- err.Error()
		return err
	}

	TextLog <- "Service Created"
	defer s.Close()

	TextLog <- "Setting up Event Log Source..."
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)

	if err != nil {
		s.Delete()
		TextLog <- err.Error()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
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
