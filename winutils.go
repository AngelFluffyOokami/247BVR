//go:build windows

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global/win32"
	"github.com/bradfitz/iter"
	"github.com/gdamore/tcell/v2"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/rivo/tview"
	"golang.org/x/sys/windows"
)

var app *tview.Application

func init() {
	var boot firstboot
	Db, err := scribble.New("config", nil)
	if err != nil {
		return
	}
	Db.Read("config", "firstboot", &boot)
	if !boot.installed {
		app = tview.NewApplication()

		text := tview.NewTextView()
		text.SetBorder(true)
		text.SetText(firstText)
		text.SetBackgroundColor(tcell.ColorBlack)
		text.SetTextColor(tcell.ColorGreen)
		text.SetBorderColor(tcell.ColorGreen).SetTitle("HSVR USB 2.0 First Run Installation Utility")
		text.SetDisabled(true)
		text.SetWordWrap(true).SetRegions(true)

		text.SetBorder(true)
		text.SetDisabled(true)

		focusedstyle := tcell.StyleDefault
		focusedstyle = focusedstyle.Background(tcell.ColorGreen)

		unfocusedstyle := tcell.StyleDefault
		unfocusedstyle = unfocusedstyle.Background(tcell.ColorBlack)

		form := tview.NewForm()
		form.AddButton("Continue", func() {
			go continueFunc(text, form)
		}).SetButtonActivatedStyle(focusedstyle).SetButtonStyle(unfocusedstyle)
		form.AddButton("Cancel", func() { os.Exit(1) }).SetButtonActivatedStyle(focusedstyle).SetButtonStyle(unfocusedstyle)
		form.SetBorder(false)
		form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyLeft {
				app.SetFocus(form.SetFocus(0))
			} else if event.Key() == tcell.KeyRight {
				app.SetFocus(form.SetFocus(1))
			}
			return event
		})

		mainArea := tview.NewGrid().
			SetRows(0, 3).
			AddItem(text, 0, 0, 1, 2, 0, 0, false).
			AddItem(form, 1, 0, 1, 1, 0, 0, true)
		mainArea.SetFocusFunc(func() {})

		if err := app.SetRoot(mainArea,
			true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

	} else {
		continueInit()
	}
}

// Make sure user actually knows the installer is doing something,
// by pretending the installer takes longer than it actually does.
func longRandLoadTime() int {
	min := 1000
	max := 3000
	return rand.Intn(max-min) + min
}

func continueFunc(text *tview.TextView, form *tview.Form) {
	_, _, _, height := text.GetRect()

	text.SetTitle("EULA")

	form.ClearButtons()

	whitespaceheight := height - 4
	cursor := 0
	for {
		var newStr string
		if cursor < 5 {
			newStr = scroll[cursor]
		}

		for range iter.N(whitespaceheight - cursor) {
			newStr += "\n"
		}

		var eulastring string
		for j := 0; j < cursor; j++ {
			if cursor >= 24 {
				for b, v := range EULA {
					if b == 0 {
						eulastring += v
					} else {
						eulastring += "\n" + v
					}

				}

				break
			} else {
				if j == 0 {
					eulastring += EULA[j]
				} else {
					eulastring += "\n" + EULA[j]
				}

			}

		}

		// Move the cursor down one line
		cursor++

		text.SetText(newStr + eulastring)

		go app.Draw()

		time.Sleep(100 * time.Millisecond)
		// If the cursor has reached the bottom of the screen, reset it to the top
		if cursor == (height - 3) {

			form.AddButton("Accept", func() {
				go installMethod(text, form)
			})
			form.AddButton("Deny", func() {
				go politeGoodbye(text, form)
			})
			app.SetFocus(form.SetFocus(0))
			form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyLeft {
					app.SetFocus(form.SetFocus(0))
				} else if event.Key() == tcell.KeyRight {
					app.SetFocus(form.SetFocus(1))
				}
				return event
			})
			go app.Draw()
			break

		}
	}
}

func becomeAdmin() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func installMethod(text *tview.TextView, form *tview.Form) {
	text.SetText(installText)
	text.SetTitle("HSVR Installation")

	form.ClearButtons()
	form.AddButton("Windows Service (Recommended)", func() {

		for {
			if !admin() {
				text.SetTitle("Administrator Privileges required")
				text.SetText("It appears the installer lacks administrator privileges.\nPlease allow administrator privileges in the following popup.")
				form.ClearButtons()
				form.AddButton("Okay", func() {
					becomeAdmin()
					os.Exit(0)
				})
				app.SetFocus(form.SetFocus(0))
				go app.Draw()
				break

			} else {
				text.SetTitle("Installing...")
				app.SetRoot(text, true)
				go app.Draw()
				go win32.Service("install")
				go handleLog(text, form)
				return
			}
		}

	})
	form.AddButton("Current User", func() {})
	app.SetFocus(form.SetFocus(0))
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(form.SetFocus(0))
		} else if event.Key() == tcell.KeyRight {
			app.SetFocus(form.SetFocus(1))
		}
		return event
	})
	go app.Draw()
}

func shortRandLoadTime() int {
	min := 100
	max := 400
	return rand.Intn(max-min) + min
}

func decideRandTime() int {
	min := 1
	max := 10
	chances := rand.Intn(max-min) + min
	if chances <= 4 {
		return shortRandLoadTime()
	} else {
		return longRandLoadTime()
	}
}

var logvar []string

func handleLog(text *tview.TextView, form *tview.Form) {
	for {
		// Pretend it takes longer than it does for :sparkles:(✨ reference omg???) user experience :sparkles:(✨ reference omg???) or some shit
		time.Sleep(time.Duration(decideRandTime()) * time.Millisecond)
		select {
		case currentLog := <-win32.TextLog:
			logvar = append(logvar, currentLog)
			var logs string
			for b, v := range logvar {
				if b == 0 {
					logs += v
				} else {
					logs += "\n" + v
				}

			}
			text.SetText(logs)
			go app.Draw()
		case installed := <-win32.Installed:
			if installed {

				text.SetTitle("You like installing services, don't you?")

				logvar = append(logvar, "Press esc key to close installation utility")
				var logs string
				for b, v := range logvar {
					if b == 0 {
						logs += v
					} else {
						logs += "\n" + v
					}

				}

				app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Key() == tcell.KeyEsc {
						os.Exit(0)
					}
					return event
				})

				text.SetText(logs)
				go app.Draw()
				foobar := keepAlive()
				fmt.Println(foobar)
			} else {
				log.Fatal("error installing")
				return
			}
		}
	}
}

func keepAlive() int {
	foo := 1
	bar := 2
	for {
		foo = foo + bar
	}

}
func admin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	if windows.GetCurrentProcessToken().IsElevated() {
		fmt.Println("elevated too")
	}
	return true
}
func politeGoodbye(text *tview.TextView, form *tview.Form) {
	_, _, _, height := text.GetRect()

	text.SetTitle("Sod Off, won't ya then?")
	whitespaceheight := height - 24
	form.ClearButtons()
	cursor := 0

	if height > 24 {

		for {
			var newStr string
			if cursor < 24 {
				eulaarr := EULA
				for range iter.N(cursor + 1) {
					eulaarr = RemoveIndex(eulaarr, 0)
				}
				for b, v := range eulaarr {
					if b == 0 {
						newStr += v

					} else {
						newStr += "\n" + v
					}
				}
			} else if cursor == 24 {
				newStr = ""
			}
			if whitespaceheight-cursor >= 1 {
				for range iter.N(whitespaceheight - cursor) {
					newStr += "\n"
				}

			}

			var goobyestr string
			for j := 0; j < cursor; j++ {
				if cursor >= 18 {
					for b, v := range politeBye {
						if b == 0 {
							goobyestr += v
						} else {
							goobyestr += "\n" + v
						}

					}

					break
				} else {
					if j == 0 {
						goobyestr += politeBye[j]
					} else {
						goobyestr += "\n" + politeBye[j]
					}

				}

			}

			// Move the cursor down one line
			cursor++

			text.SetText(newStr + goobyestr)

			app.Draw()
			time.Sleep(100 * time.Millisecond)
			// If the cursor has reached the bottom of the screen, reset it to the top
			if cursor == (height-whitespaceheight)+1 {

				form.AddButton("Exit", func() {
					os.Exit(0)
				})
				form.AddButton("Nuh uh", func() {
					text.SetText(fymnuh)
					form.ClearButtons()
					form.AddButton("Shut up", func() {
						form.ClearButtons()
						form.AddButton("They're coming", func() {
							form.ClearButtons()
							form.AddButton("They're near", func() {
								form.ClearButtons()
								form.AddButton("They're in the walls", func() {
									form.ClearButtons()
									form.AddButton("Rip them out", func() {
										form.ClearButtons()
										form.AddButton("Too late", func() {
											form.AddButton("They're in your veins", func() {
												form.ClearButtons()
												form.AddButton("Veins make you crazy", func() {
													form.ClearButtons()
													form.AddButton("Crazy?", func() {
														form.ClearButtons()
														form.AddButton("You were crazy once", func() {
															form.ClearButtons()
															form.AddButton("They put you in a rubber room", func() {
																form.ClearButtons()
																form.AddButton("A rubber room with veins.", func() {

																	log.Fatal(fmt.Errorf("they're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out\nthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're outthey're out"))
																})
																text.SetText(flower)
																text.SetTextColor(tcell.ColorRed)
																text.SetBackgroundColor(tcell.ColorDarkRed)
																form.SetButtonBackgroundColor(tcell.ColorDarkRed)
																form.SetButtonActivatedStyle(tcell.StyleDefault.Background(tcell.ColorRed))

																go app.Draw()
															})
															go app.Draw()
														})
														go app.Draw()
													})
													go app.Draw()
												})
												go app.Draw()
											})
											go app.Draw()
										})
										go app.Draw()
									})
									go app.Draw()
								})
								go app.Draw()
							})
							go app.Draw()
						})
						go app.Draw()
					})
					go app.Draw()
				})
				go app.Draw()
				break

			}
		}
	}

}

func RemoveIndex(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

type firstboot struct {
	installed bool
}
