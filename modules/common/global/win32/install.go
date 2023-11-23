//go:build windows

package win32

import (
	"fmt"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func installService(name, desc string) error {
	exepath := "C:\\Program Files\\247bvr\\247bvr.exe"

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
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	TextLog <- "Service does not exist... Continuing."
	TextLog <- "Creating service..."
	s, err = m.CreateService(name, exepath, mgr.Config{StartType: mgr.StartAutomatic, Description: desc, ServiceStartName: "247bvr"}, "is", "auto-started")
	if err != nil {
		return err
	}

	TextLog <- "Service Created"
	defer s.Close()

	TextLog <- "Setting up Event Log Source..."
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	TextLog <- "Event Log Source set up"

	TextLog <- "Install done"
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
