package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common/global/win32"
	"github.com/bradfitz/iter"
	"github.com/gdamore/tcell/v2"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/rivo/tview"
	"golang.org/x/sys/windows"
)

var app *tview.Application

func main() {
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
		text.SetBorderColor(tcell.ColorGreen).SetTitle("HSVR ELO Statistics Bot firstboot Installation")

		text.SetChangedFunc(func() {
			app.Draw()
		}).SetWordWrap(true).SetRegions(true)

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
		if err := app.SetRoot(mainArea,
			true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

	}
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

		app.Draw()
		time.Sleep(100 * time.Millisecond)
		// If the cursor has reached the bottom of the screen, reset it to the top
		if cursor == (height - 3) {

			form.AddButton("Accept", func() {
				go installMethod(text, form)
			})
			form.AddButton("Deny", func() {
				go politeGoodbye(text, form)
			})
			app.Draw()
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

	form.ClearButtons()
	form.AddButton("Windows Service (Recommended)", func() {

		for {
			if !admin() {
				becomeAdmin()
			} else {
				text.SetTitle("Installing...")
				app.SetRoot(text, true)
				go app.Draw()
				go win32.Service("install")
				go handleLog(text, form)
			}
		}

	})
	form.AddButton("Current User", func() {})
	go app.Draw()
}

var logvar []string

func handleLog(text *tview.TextView, form *tview.Form) {
	for {
		select {
		case currentLog := <-win32.TextLog:
			_, _, _, height := text.GetInnerRect()
			if len(logvar) <= height {
				logvar = append(logvar, currentLog)
			} else {
				RemoveIndex(logvar, 0)
				logvar = append(logvar, currentLog)
			}

			var logs string
			for _, v := range logvar {
				if logs == "" {
					logs += v
				} else {
					logs += "\n" + v
				}
			}

			text.SetText(logs)
			go app.Draw()
		case installed := <-win32.Installed:
			if installed {
				return
			} else {
				log.Fatal("error installing")
				return
			}
		}
	}
}

func admin() bool {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Fatalf("SID Error: %s", err)
		return false
	}
	defer windows.FreeSid(sid)

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	// Also note that an admin is _not_ necessarily considered
	// elevated.
	// For elevation see https://github.com/mozey/run-as-admin
	return token.IsElevated()
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
