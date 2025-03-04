package cli

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

var (
	errWrongKey = errors.New("wrong key")
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

func (a *App) draw(strings []string, md mode) {
	var drawers = map[mode]drawer{
		raw:      a.rawDrawer,
		paged:    a.pagedDrawer,
		scrolled: a.scrolledDrawer,
	}

	drawers[md](strings)
}

func (a *App) rawDrawer(strings []string) {
	if strings[0] == clearScreen {
		fmt.Print(clearScreen)

		return
	}
	for _, str := range strings {
		a.stringBuilder.WriteString(str)
		a.stringBuilder.WriteByte('\n')
	}

	fmt.Print(a.stringBuilder.String())
	a.stringBuilder.Reset()
}

func (a *App) clearNLinesUp(n int) {
	for range n {
		fmt.Print("\033[A")  // Move cursor up one line
		fmt.Print("\033[2K") // Clear entire line
	}
	fmt.Print("\033[G")
}

func (a *App) makePages(lines []string) []string {
	var pages []string
	for i, str := range lines {
		a.stringBuilder.WriteString(str)
		a.stringBuilder.WriteString("\r\n")
		if (i+1)%pageSize == 0 {
			pages = append(pages, a.stringBuilder.String())
			a.stringBuilder.Reset()
		}
	}
	if len(lines)%pageSize == 0 {
		return pages
	}

	for i := len(lines) % pageSize; i < pageSize; i++ {
		a.stringBuilder.WriteString("--\r\n")
	}
	pages = append(pages, a.stringBuilder.String())
	a.stringBuilder.Reset()

	return pages
}

func (a *App) printCurrentPage(pages []string, currentPage int, pageCount int) {
	fmt.Print(pages[currentPage], "\r")
	fmt.Printf("← %d/%d →\r\n", currentPage+1, pageCount)
}

func (a *App) input(handler keyHandler, exitChar byte) (key, error) {
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

func getKeyByLastByte(last byte) key {
	switch last {
	case 'A':
		return upArrow
	case 'B':
		return downArrow
	case 'C':
		return rightArrow
	case 'D':
		return leftArrow
	default:
		return wrongKey
	}
}

func (a *App) getArrowKeys(buf []byte) (key, error) {
	if buf[0] != '\x1b' {
		return wrongKey, errWrongKey
	}
	seq := make([]byte, 2)
	_, err := os.Stdin.Read(seq)
	if err != nil || seq[0] != '[' {
		return wrongKey, errWrongKey
	}

	keyPressed := getKeyByLastByte(seq[1])
	if keyPressed == wrongKey {
		return wrongKey, errWrongKey
	}

	return keyPressed, nil
}

func (a *App) setupTerminal(fd int) (*term.State, error) {
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("error opening terminal: %v\n", err)
	}

	return oldState, err
}

func (a *App) restoreTerminal(fd int, oldState *term.State) {
	err := term.Restore(fd, oldState)
	if err != nil {
		fmt.Printf("error restoring terminal: %v\n", err)
	}
}

func (a *App) changePage(pressedKey key, currentPage int, pageCount int) (int, bool) {
	switch pressedKey {
	case leftArrow:
		if currentPage == 0 {
			return currentPage, false
		}

		currentPage--
	case rightArrow:
		if currentPage == pageCount-1 {
			return currentPage, false
		}

		currentPage++
	default:
		return currentPage, false
	}

	return currentPage, true
}

func (a *App) pagedDrawer(strings []string) {
	if len(strings) == 0 {
		return
	}
	pages := a.makePages(strings)

	fd := int(os.Stdin.Fd())
	oldState, err := a.setupTerminal(fd)
	if err != nil {
		return
	}
	defer a.restoreTerminal(fd, oldState)

	currentPage := 0
	var hasChanged bool
	pageCount := (len(strings) + pageSize - 1) / pageSize
	fmt.Println("--Viewing return table. Press q to quit--\r")
	a.printCurrentPage(pages, currentPage, pageCount)
	var pressedKey key
	for {
		pressedKey, err = a.input(a.getArrowKeys, 'q')
		if err != nil {
			continue
		}

		if pressedKey == exitKey {
			break
		}

		currentPage, hasChanged = a.changePage(pressedKey, currentPage, pageCount)
		if !hasChanged {
			continue
		}

		a.clearNLinesUp(pageSize + 1)
		a.printCurrentPage(pages, currentPage, pageCount)
	}

	a.clearNLinesUp(pageSize + 2)
	fmt.Print("Success: table viewed\r\n")
}

func (a *App) printWindow(strings []string, currentPosition int) {
	lastPosition := len(strings)
	for i := currentPosition; i < currentPosition+windowSize; i++ {
		if i >= lastPosition {
			a.stringBuilder.WriteString("--\r\n")

			continue
		}
		a.stringBuilder.WriteString(strings[i])
		a.stringBuilder.WriteString("\r\n")
	}
	fmt.Printf("↑ %d elements\r\n", currentPosition)
	fmt.Print(a.stringBuilder.String())
	fmt.Printf("↓ %d elements\r\n", max(0, lastPosition-currentPosition-windowSize))
	a.stringBuilder.Reset()
}

func (a *App) changeWindow(pressedKey key, currentPosition int, lastPosition int) (int, bool) {
	switch pressedKey {
	case upArrow:
		if currentPosition == 0 {
			return currentPosition, false
		}

		currentPosition--
	case downArrow:
		if currentPosition+windowSize-1 >= lastPosition {
			return currentPosition, false
		}

		currentPosition++
	default:
		return currentPosition, false
	}

	return currentPosition, true
}

func (a *App) scrolledDrawer(strings []string) {
	if len(strings) == 0 {
		return
	}

	fd := int(os.Stdin.Fd())
	oldState, err := a.setupTerminal(fd)
	if err != nil {
		return
	}
	defer a.restoreTerminal(fd, oldState)

	currentPosition := 0
	lastPosition := len(strings) - 1
	var hasChanged bool
	fmt.Println("--Viewing return table. Press q to quit--\r")
	a.printWindow(strings, currentPosition)
	var pressedKey key
	for {
		pressedKey, err = a.input(a.getArrowKeys, 'q')
		if err != nil {
			continue
		}

		if pressedKey == exitKey {
			break
		}

		currentPosition, hasChanged = a.changeWindow(pressedKey, currentPosition, lastPosition)
		if !hasChanged {
			continue
		}

		a.clearNLinesUp(windowSize + 2)
		a.printWindow(strings, currentPosition)
	}

	a.clearNLinesUp(windowSize + 3)
	fmt.Print("Success: table viewed\r\n")
}
