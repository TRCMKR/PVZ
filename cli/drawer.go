package cli

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

var (
	ErrWrongKey = errors.New("wrong key")
)

type key byte

const (
	leftArrow key = iota
	rightArrow
	upArrow
	downArrow

	exitKey  = 254
	wrongKey = 255
)

type keyHandler func([]byte) (key, error)

type drawer func(strings []string)

const (
	pageSize   = 7
	windowSize = 5
)

func (app *App) draw(strings []string, md mode) {
	var drawers = map[mode]drawer{
		raw:      app.rawDrawer,
		paged:    app.pagedDrawer,
		scrolled: app.scrolledDrawer,
	}

	drawers[md](strings)
}

func (app *App) rawDrawer(strings []string) {
	if strings[0] == clearScreen {
		fmt.Print(clearScreen)
		return
	}
	for _, str := range strings {
		app.stringBuilder.WriteString(str)
		app.stringBuilder.WriteByte('\n')
	}

	fmt.Print(app.stringBuilder.String())
	app.stringBuilder.Reset()
}

func (app *App) clearNLinesUp(n int) {
	for range n {
		fmt.Print("\033[A")  // Move cursor up one line
		fmt.Print("\033[2K") // Clear entire line
	}
	fmt.Print("\033[G")
}

func (app *App) makePages(strings []string) []string {
	var pages []string
	for i, str := range strings {
		app.stringBuilder.WriteString(str)
		app.stringBuilder.WriteString("\r\n")
		if (i+1)%pageSize == 0 {
			pages = append(pages, app.stringBuilder.String())
			app.stringBuilder.Reset()
		}
	}
	if len(strings)%pageSize != 0 {
		for i := len(strings) % pageSize; i < pageSize; i++ {
			app.stringBuilder.WriteString("--\r\n")
		}
		pages = append(pages, app.stringBuilder.String())
		app.stringBuilder.Reset()
	}

	return pages
}

func (app *App) printCurrentPage(pages []string, currentPage int, pageCount int) {
	fmt.Print(pages[currentPage], "\r")
	fmt.Printf("← %d/%d →\r\n", currentPage+1, pageCount)
}

func (app *App) input(handler keyHandler, exitChar byte) (key, error) {
	buf := make([]byte, 1)

	_, err := os.Stdin.Read(buf)
	if err != nil {
		fmt.Printf("error reading input: %v\n", err)
		return wrongKey, err
	}

	if buf[0] == exitChar {
		return exitKey, nil
	}

	pressedKey, err := handler(buf)
	if err != nil {
		return wrongKey, err
	}

	return pressedKey, nil
}

func (app *App) getArrowKeys(buf []byte) (key, error) {
	if buf[0] != '\x1b' {
		return wrongKey, ErrWrongKey
	}
	seq := make([]byte, 2)
	_, err := os.Stdin.Read(seq)
	if err != nil || seq[0] != '[' {
		return wrongKey, ErrWrongKey
	}

	switch seq[1] {
	case 'A': // Arrow up
		return upArrow, nil
	case 'B': // Arrow down
		return downArrow, nil
	case 'C': // Arrow right
		return rightArrow, nil
	case 'D': // Arrow left
		return leftArrow, nil
	default:
		return wrongKey, ErrWrongKey
	}
}

func (app *App) pagedDrawer(strings []string) {
	if len(strings) == 0 {
		return
	}
	pages := app.makePages(strings)

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("error opening terminal: %v\n", err)
		return
	}
	defer term.Restore(fd, oldState)

	currentPage := 0
	pageCount := (len(strings) + pageSize - 1) / pageSize
	fmt.Println("--Viewing return table. Press q to quit--\r")
	app.printCurrentPage(pages, currentPage, pageCount)
	var pressedKey key
	for {
		pressedKey, err = app.input(app.getArrowKeys, 'q')
		if err != nil {
			continue
		}

		if pressedKey == exitKey {
			break
		}

		switch pressedKey {
		case leftArrow:
			if currentPage == 0 {
				continue
			}

			currentPage--
		case rightArrow:
			if currentPage == pageCount-1 {
				continue
			}

			currentPage++
		default:
			continue
		}

		app.clearNLinesUp(pageSize + 1)
		app.printCurrentPage(pages, currentPage, pageCount)
	}

	app.clearNLinesUp(pageSize + 3)
	fmt.Print("Success: table viewed\r\n")
}

func (app *App) printWindow(strings []string, currentPosition int) {
	lastPosition := len(strings)
	for i := currentPosition; i < currentPosition+windowSize; i++ {
		if i >= lastPosition {
			app.stringBuilder.WriteString("--\r\n")
			continue
		}
		app.stringBuilder.WriteString(strings[i])
		app.stringBuilder.WriteString("\r\n")
	}
	fmt.Printf("↑ %d elements\r\n", currentPosition)
	fmt.Print(app.stringBuilder.String())
	fmt.Printf("↓ %d elements\r\n", max(0, lastPosition-currentPosition-windowSize))
	app.stringBuilder.Reset()
}

func (app *App) scrolledDrawer(strings []string) {
	if len(strings) == 0 {
		return
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("error opening terminal: %v\n", err)
		return
	}
	defer term.Restore(fd, oldState)

	currentPosition := 0
	lastPosition := len(strings) - 1
	fmt.Println("--Viewing return table. Press q to quit--\r")
	app.printWindow(strings, currentPosition)
	var pressedKey key
	for {
		pressedKey, err = app.input(app.getArrowKeys, 'q')
		if err != nil {
			continue
		}

		if pressedKey == exitKey {
			break
		}

		switch pressedKey {
		case upArrow:
			if currentPosition == 0 {
				continue
			}

			currentPosition--
		case downArrow:
			if currentPosition+windowSize-1 >= lastPosition {
				continue
			}

			currentPosition++
		default:
			continue
		}

		app.clearNLinesUp(windowSize + 2)
		app.printWindow(strings, currentPosition)
	}

	app.clearNLinesUp(windowSize + 4)
	fmt.Print("Success: table viewed\r\n")
}
