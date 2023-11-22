//go:build windows

package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/rivo/tview"
)

var winConfig string

type firstboot struct {
	firstboot bool
}

var firstText = `Welcome to the HSVR Bot firstboot installation utility.
If you need to reinstall the bot, just set firstboot to true in the firstboot configuration file.
If you only want to move your installation, set modifyinstall to true in the firstboot configuration file.
If you erroneously set firstboot or modifyinstall to true, exit the application and set either/both flags to false.


Do you wish to continue installing HSVR ELO Statistics Bot?`

func init() {
	var boot firstboot
	Db, err := scribble.New("config", nil)
	if err != nil {
		return
	}
	Db.Read("config", "firstboot", &boot)
	if boot.firstboot {
		app := tview.NewApplication()

		text := tview.NewTextView()
		text.SetBorder(true)
		text.SetText(firstText)
		text.SetBackgroundColor(tcell.ColorBlack)
		text.SetTextColor(tcell.ColorGreen)
		text.SetBorderColor(tcell.ColorGreen)

		continueButton := tview.NewButton("Continue").SetSelectedFunc(func() {
		})
		cancelButton := tview.NewButton("Cancel").SetSelectedFunc(func() {
			os.Exit(1)
		})

		continueButton.SetBorder(true).SetRect(0, 0, 22, 3)

		cancelButton.SetBorder(true).SetRect(0, 0, 22, 3)

		cancelButton.SetBorderColor(tcell.ColorBlack)
		cancelButton.SetBackgroundColor(tcell.ColorBlack)
		cancelButton.SetBackgroundColorActivated(tcell.ColorWhite)
		cancelButton.SetLabelColorActivated(tcell.ColorGreen)
		cancelButton.SetLabelColor(tcell.ColorGreen)

		cancelButton.SetActivatedStyle(tcell.StyleDefault.Attributes(tcell.AttrBlink + tcell.AttrBold + tcell.AttrUnderline))
		cancelButton.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

		continueButton.SetBorderColor(tcell.ColorBlack)
		continueButton.SetBackgroundColor(tcell.ColorBlack)
		continueButton.SetBackgroundColorActivated(tcell.ColorWhite)
		continueButton.SetLabelColorActivated(tcell.ColorGreen)
		continueButton.SetLabelColor(tcell.ColorGreen)

		continueButton.SetActivatedStyle(tcell.StyleDefault.Attributes(tcell.AttrBlink + tcell.AttrBold + tcell.AttrUnderline))
		continueButton.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrDim))

		buttonGrid := tview.NewGrid().SetColumns(2)

		buttonGrid.AddItem(continueButton, 0, 0, 1, 1, 0, 0, true)
		buttonGrid.AddItem(cancelButton, 0, 1, 1, 1, 0, 0, false)

		mainArea := tview.NewGrid().
			SetRows(2).AddItem(text, 0, 0, 10, 1, 0, 0, false)
		mainArea.AddItem(buttonGrid, 1, 0, 1, 1, 0, 0, true)
		mainArea.SetBorder(true).SetTitle("HSVR ELO Statistics Bot firstboot Installation")

		if err := app.SetRoot(mainArea,
			true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

	}
}
