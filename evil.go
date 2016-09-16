package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
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

func changeFile(fileName *string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var l string
		var err error

		_, cy := v.Cursor()
		if l, err = v.Line(cy); err != nil {
			l = ""
		}

		*fileName = l

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
		loadFile(*fileName, g, mv)

		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	currentFileName := os.Args[1]
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
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, changeFile(&currentFileName)); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain(currentFileName)); err != nil {
		return err
	}
	return nil
}

func saveMain(fileName string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		f, err := os.OpenFile(fileName, os.O_WRONLY, 0777)
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
}

func loadFile(fileName string, g *gocui.Gui, v *gocui.View) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(v, "%s", b)
	v.Editable = true
	v.Wrap = true
	if err := g.SetCurrentView("main"); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		files, _ := ioutil.ReadDir("./")
		for _, f := range files {
			fmt.Fprintln(v, f.Name())
		}
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		v.Title = os.Args[1]
		if err != gocui.ErrUnknownView {
			return err
		}
		loadFile(os.Args[1], g, v)
	}
	return nil
}

func main() {
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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
