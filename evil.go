package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"

	"github.com/jroimartin/gocui"
	"github.com/mgutz/ansi"
)

var currentFileName string

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" || v.Name() == "info" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("side")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func changeFile(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	currentFileName = l

	mv, err := g.View("main")
	if err != nil {
		panic(err)
	}

	mv.Clear()
	if err := mv.SetOrigin(0, 0); err != nil {
		return err
	}
	if err := mv.SetCursor(0, 0); err != nil {
		return err
	}
	loadFile(currentFileName, g, mv)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, changeFile); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	f, err := os.OpenFile(currentFileName, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFile(fileName string, g *gocui.Gui, v *gocui.View) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(v, "%s", b)
	v.Editable = true
	v.Wrap = true
	v.Title = fileName
	if err := g.SetCurrentView("main"); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, 0, 27, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		files, _ := ioutil.ReadDir("./")
		for _, f := range files {
			fmt.Fprintln(v, f.Name())
		}
	}
	if v, err := g.SetView("main", 27, 0, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = os.Args[1]
		loadFile(os.Args[1], g, v)
		g.Editor = gocui.EditorFunc(simpleEditor(g))
	}
	if v, err := g.SetView("info", 27, maxY-2, maxX, maxY); err != nil {
		v.Editable = true
		if err != gocui.ErrUnknownView {
			return err
		}
		mv, err := g.View("main")
		if err != nil {
			return err
		}
		cx, cy := mv.Cursor()
		var mode string
		if g.Mouse {
			mode = "INSERT"
		} else {
			mode = "NORMAL"
		}
		_, oy := v.Origin()
		v.Title = fmt.Sprintf("%s─────────────────────────────────L%d/C%d", mode, cy+oy, cx)
	}
	return nil
}

func moveCursor(g *gocui.Gui, v *gocui.View, fun func(x, y int) (int, int)) {
	x, y := v.Cursor()
	nx, ny := fun(x, y)
	var l string
	var err error
	if l, err = v.Line(ny); err != nil {
		l = ""
	}
	if nx > utf8.RuneCountInString(l) {
		nx, ny = 0, ny
	} else if nx < 0 {
		nx, ny = 0, ny
	}
	if err := v.SetCursor(nx, ny); err != nil {
		ox, oy := v.Origin()
		dy := 0
		if ny > y {
			dy = 1
		} else if ny < y && oy-1 >= 0 {
			dy = -1
		}
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			panic(err)
		}
	}
}

func switchToCommand(g *gocui.Gui, v *gocui.View) error {
	err := g.SetCurrentView("info")
	if err != nil {
		return err
	}
	iv, err := g.View("info")
	if err != nil {
		return err
	}
	iv.Clear()
	if err := v.SetCursor(4, 0); err != nil {
		panic(err)
	}
	return nil
}

func runCommand(g *gocui.Gui, v *gocui.View) {
	v.Clear()
	switch v.ViewBuffer() {
	case "w":
		saveMain(g, v)
	default:
		v.Clear()
		if err := v.SetCursor(0, 0); err != nil {
			panic(err)
		}
		fmt.Fprintf(v, ansi.Color("Bad command: %s", "cyan+b+h"), v.ViewBuffer())
	}
}

func simpleEditor(g *gocui.Gui) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		if v.Name() == "main" {
			switch {
			case ch == 'i' && !g.Mouse:
				g.Mouse = true
			case ch == ':' && !g.Mouse:
				switchToCommand(g, v)
			case (ch == 'j' && !g.Mouse) || key == gocui.KeyArrowDown:
				moveCursor(g, v, func(x, y int) (int, int) {
					return x, y + 1
				})
			case (ch == 'k' && !g.Mouse) || key == gocui.KeyArrowUp:
				moveCursor(g, v, func(x, y int) (int, int) {
					return x, y - 1
				})
			case (ch == 'l' && !g.Mouse) || key == gocui.KeyArrowRight:
				moveCursor(g, v, func(x, y int) (int, int) {
					return x + 1, y
				})
			case (ch == 'h' && !g.Mouse) || key == gocui.KeyArrowLeft:
				moveCursor(g, v, func(x, y int) (int, int) {
					return x - 1, y
				})
			case ch != 0 && mod == 0 && g.Mouse:
				v.EditWrite(ch)
			case key == gocui.KeySpace && g.Mouse:
				v.EditWrite(' ')
			case (key == gocui.KeyBackspace || key == gocui.KeyBackspace2) && g.Mouse:
				v.EditDelete(true)
			case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
				moveCursor(g, v, func(x, y int) (int, int) {
					return x - 1, y
				})
			case key == gocui.KeyEsc:
				g.Mouse = false
				moveCursor(g, v, func(x, y int) (int, int) {
					return x - 1, y
				})
			}
			iv, err := g.View("info")
			if err != nil {
				panic(err)
			}
			mv, err := g.View("main")
			if err != nil {
				panic(err)
			}
			cx, cy := mv.Cursor()
			_, oy := v.Origin()
			var mode string
			if g.Mouse {
				mode = "INSERT"
				iv.Title = fmt.Sprintf("%s─────────────────────────────────L%d/C%d", mode, oy+cy, cx)
			} else {
				mode = "NORMAL"
				iv.Title = fmt.Sprintf("%s─────────────────────────────────L%d/C%d", mode, oy+cy, cx)
			}
		} else if v.Name() == "info" {
			switch {
			case ch != 0 && mod == 0:
				v.EditWrite(ch)
			case key == gocui.KeySpace:
				v.EditWrite(' ')
			case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
				v.EditDelete(true)
			case key == gocui.KeyEnter:
				runCommand(g, v)
				nextView(g, v)
			}
		}
	}
}

func main() {
	currentFileName = os.Args[1]
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true
	g.InputEsc = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
