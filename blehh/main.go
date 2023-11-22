package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bradfitz/iter"
	"github.com/gdamore/tcell/v2"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/rivo/tview"
)

var firstText = "Welcome to the HSVR Bot firstboot installation utility." + "\n" +
	"If you need to reinstall the bot, just set installed to false in the firstboot configuration file." + "\n" +
	"If you only want to move your installation, set `modifyinstall` to true in the firstboot configuration file." + "\n" +
	"If you erroneously set firstboot or modifyinstall to true, exit the application and set either/both flags to false" + "\n" +
	"Do you wish to continue installing HSVR ELO Statistics Bot?"

var scroll = []string{
	"If you need to reinstall the bot, just set installed to false in the firstboot configuration file." + "\n" +
		"If you only want to move your installation, set `modifyinstall` to true in the firstboot configuration file." + "\n" +
		"If you erroneously set firstboot or modifyinstall to true, exit the application and set either/both flags to false" + "\n" +
		"Do you wish to continue installing HSVR ELO Statistics Bot?",

	"If you only want to move your installation, set `modifyinstall` to true in the firstboot configuration file." + "\n" +
		"If you erroneously set firstboot or modifyinstall to true, exit the application and set either/both flags to false" + "\n" +
		"Do you wish to continue installing HSVR ELO Statistics Bot?",

	"If you erroneously set firstboot or modifyinstall to true, exit the application and set either/both flags to false" + "\n" +
		"Do you wish to continue installing HSVR ELO Statistics Bot?",

	"Do you wish to continue installing HSVR ELO Statistics Bot?",

	"",
}

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

			form.AddButton("Accept", func() {})
			form.AddButton("Deny", func() {
				go politeGoodbye(text, form)
			})
			app.Draw()
			break

		}
	}
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

																	log.Fatal(fmt.Errorf("They're out"))
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

var winConfig string

type firstboot struct {
	installed bool
}

var EULA = []string{
	"               GLWT(Good Luck With That) Public License",
	"Copyright (c) Everyone, except Author",
	"",
	"Everyone is permitted to copy, distribute, modify, merge, sell, publish,",
	"sublicense or whatever they want with this software but at their OWN RISK.",
	"",
	"		   Preamble",
	"",
	"The author has absolutely no clue what the code in this project does.",
	"It might just work or not, there is no third option.",
	"",
	"",
	"GOOD LUCK WITH THAT PUBLIC LICENSE",
	"TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION, AND MODIFICATION",
	"",
	"0. You just DO WHATEVER YOU WANT TO as long as you NEVER LEAVE A",
	"TRACE TO TRACK THE AUTHOR of the original product to blame for or hold",
	"responsible.",
	"",
	"IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER",
	"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING",
	"FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER",
	"DEALINGS IN THE SOFTWARE.",
	"",
	"Good luck and Godspeed.",
}

var politeBye = []string{
	"	                                                                                                                                                  ",
	"                                                                                                                                                     ",
	"FFFFFFFFFFFFFFFFFFFFFF                                      kkkkkkkk                     OOOOOOOOO        ffffffffffffffff    ffffffffffffffff       ",
	"F::::::::::::::::::::F                                      k::::::k                   OO:::::::::OO     f::::::::::::::::f  f::::::::::::::::f      ",
	"F::::::::::::::::::::F                                      k::::::k                 OO:::::::::::::OO  f::::::::::::::::::ff::::::::::::::::::f     ",
	"FF::::::FFFFFFFFF::::F                                      k::::::k                O:::::::OOO:::::::O f::::::fffffff:::::ff::::::fffffff:::::f     ",
	" F:::::F       FFFFFFuuuuuu    uuuuuu      cccccccccccccccc k:::::k    kkkkkkk     O::::::O   O::::::O f:::::f       fffffff:::::f       ffffff      ",
	" F:::::F             u::::u    u::::u    cc:::::::::::::::c k:::::k   k:::::k      O:::::O     O:::::O f:::::f             f:::::f                   ",
	" F::::::FFFFFFFFFF   u::::u    u::::u   c:::::::::::::::::c k:::::k  k:::::k       O:::::O     O:::::Of:::::::ffffff      f:::::::ffffff             ",
	" F:::::::::::::::F   u::::u    u::::u  c:::::::cccccc:::::c k:::::k k:::::k        O:::::O     O:::::Of::::::::::::f      f::::::::::::f             ",
	" F:::::::::::::::F   u::::u    u::::u  c::::::c     ccccccc k::::::k:::::k         O:::::O     O:::::Of::::::::::::f      f::::::::::::f             ",
	" F::::::FFFFFFFFFF   u::::u    u::::u  c:::::c              k:::::::::::k          O:::::O     O:::::Of:::::::ffffff      f:::::::ffffff             ",
	" F:::::F             u::::u    u::::u  c:::::c              k:::::::::::k          O:::::O     O:::::O f:::::f             f:::::f                   ",
	" F:::::F             u:::::uuuu:::::u  c::::::c     ccccccc k::::::k:::::k         O::::::O   O::::::O f:::::f             f:::::f                   ",
	"FF:::::::FF           u:::::::::::::::uuc:::::::cccccc:::::ck::::::k k:::::k        O:::::::OOO:::::::Of:::::::f           f:::::::f                 ",
	"F::::::::FF            u:::::::::::::::u c:::::::::::::::::ck::::::k  k:::::k        OO:::::::::::::OO f:::::::f           f:::::::f                 ",
	"F::::::::FF             uu::::::::uu:::u  cc:::::::::::::::ck::::::k   k:::::k         OO:::::::::OO   f:::::::f           f:::::::f                 ",
	"FFFFFFFFFFF               uuuuuuuu  uuuu    cccccccccccccccckkkkkkkk    kkkkkkk          OOOOOOOOO     fffffffff           fffffffff                 ",
}

var fymnuh = `                                                                                                                                                                           
                                                                                                                                                                           
ffffffffffffffff                                                                                            hhhhhhh                                   hhhhhhh             
f::::::::::::::::f                                                                                          h:::::h                                   h:::::h             
f::::::::::::::::::f                                                                                        h:::::h                                   h:::::h             
f::::::fffffff:::::f                                                                                        h:::::h                                   h:::::h             
f:::::f       ffffffyyyyyyy           yyyyyyy mmmmmmm    mmmmmmm         nnnn  nnnnnnnn    uuuuuu    uuuuuu  h::::h hhhhh            uuuuuu    uuuuuu  h::::h hhhhh       
f:::::f               y:::::y         y:::::ymm:::::::m  m:::::::mm      n:::nn::::::::nn  u::::u    u::::u  h::::hh:::::hhh         u::::u    u::::u  h::::hh:::::hhh    
f:::::::ffffff         y:::::y       y:::::ym::::::::::mm::::::::::m     n::::::::::::::nn u::::u    u::::u  h::::::::::::::hh       u::::u    u::::u  h::::::::::::::hh  
f::::::::::::f          y:::::y     y:::::y m::::::::::::::::::::::m     nn:::::::::::::::nu::::u    u::::u  h:::::::hhh::::::h      u::::u    u::::u  h:::::::hhh::::::h 
f::::::::::::f           y:::::y   y:::::y  m:::::mmm::::::mmm:::::m       n:::::nnnn:::::nu::::u    u::::u  h::::::h   h::::::h     u::::u    u::::u  h::::::h   h::::::h
f:::::::ffffff            y:::::y y:::::y   m::::m   m::::m   m::::m       n::::n    n::::nu::::u    u::::u  h:::::h     h:::::h     u::::u    u::::u  h:::::h     h:::::h
f:::::f                    y:::::y:::::y    m::::m   m::::m   m::::m       n::::n    n::::nu::::u    u::::u  h:::::h     h:::::h     u::::u    u::::u  h:::::h     h:::::h
f:::::f                     y:::::::::y     m::::m   m::::m   m::::m       n::::n    n::::nu:::::uuuu:::::u  h:::::h     h:::::h     u:::::uuuu:::::u  h:::::h     h:::::h
f:::::::f                    y:::::::y      m::::m   m::::m   m::::m       n::::n    n::::nu:::::::::::::::uuh:::::h     h:::::h     u:::::::::::::::uuh:::::h     h:::::h
f:::::::f                     y:::::y       m::::m   m::::m   m::::m       n::::n    n::::n u:::::::::::::::uh:::::h     h:::::h      u:::::::::::::::uh:::::h     h:::::h
f:::::::f                    y:::::y        m::::m   m::::m   m::::m       n::::n    n::::n  uu::::::::uu:::uh:::::h     h:::::h       uu::::::::uu:::uh:::::h     h:::::h
fffffffff                   y:::::y         mmmmmm   mmmmmm   mmmmmm       nnnnnn    nnnnnn    uuuuuuuu  uuuuhhhhhhh     hhhhhhh         uuuuuuuu  uuuuhhhhhhh     hhhhhhh
						   y:::::y                                                                                                                                        
					      y:::::y                                                                                                                                         
					     y:::::y                                                                                                                                          
					    y:::::y                                                                                                                                           
					   yyyyyyy                                                                                                                                            
																																									   
																																									   `
var flower = `
█████████████▀▀▀▀▀███████▀▀▀▀▀█████████████
█████████▀░░▀▀█▄▄▄▄▄▄██▄▄▄▄▄▄█▀░░▀█████████
████████▄░░▄▄████▀▀▀▀▀▀▀▀▀████▄▄░░▄████████
████▀▀▀▀█████▀░░░░░░░░░░░░░░░▀█████▀▀▀▀████
██▀░░░░░░██▀░░░░░░██░░░██░░░░░░▀██░░░░░░▀██
█░░░▀▀▀▀███░░░░░░░██░░░██░░░░░░░███▀▀▀▀░░░█
█▄▄░░░░░░██░░░░▄░░▀▀░░░▀▀░░▄░░░░██░░░░░░▄▄█
████▄░░░░▀██░░░░███████████░░░░██▀░░░░▄████
██████████▀██▄░░░▀███████▀░░░▄██▀██████████
███████▀░░░████▄▄░░░░░░░░░▄▄████░░░▀███████
██████░░░▄▀░░▀▀▀███████████▀▀▀░░▀▄░░░██████
██████░░░▀░░░░░░░░▄▄▄█▄▄▄░░░░░░░░▀░░░██████
████████▄▄▄▄▄▄███████████████▄▄▄▄▄▄████████
██████████████████▀░░▀█████████████████████
█████████████████▀░░░▄█████████████████████
█████████████████░░░███████████████████████
██████████████████░░░▀█████████████████████
███████████████████▄░░░████████████████████
█████████████████████░░░███████████████████


`
