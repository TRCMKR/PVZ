package cli

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type drawer func(strings []string)

var drawers = map[mode]drawer{
	raw:      rawDrawer,
	paged:    pagedDrawer,
	scrolled: scrolledDrawer,
}

const (
	pageSize   = 7
	windowSize = 5
)

func draw(strings []string, md mode) {
	drawers[md](strings)
}

func rawDrawer(strings []string) {
	if strings[0] == clearScreen {
		fmt.Print(clearScreen)
		return
	}
	for _, str := range strings {
		stringBuilder.WriteString(str)
		stringBuilder.WriteByte('\n')
	}

	fmt.Print(stringBuilder.String())
	stringBuilder.Reset()
}

func clearNLinesUp(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[A")  // Move cursor up one line
		fmt.Print("\033[2K") // Clear entire line
	}
	fmt.Print("\033[G")
}

func makePages(strings []string) []string {
	var pages []string
	for i, str := range strings {
		stringBuilder.WriteString(str)
		stringBuilder.WriteString("\r\n")
		if (i+1)%pageSize == 0 {
			pages = append(pages, stringBuilder.String())
			stringBuilder.Reset()
		}
	}
	if len(strings)%pageSize != 0 {
		for i := len(strings) % pageSize; i < pageSize; i++ {
			stringBuilder.WriteString("--\r\n")
		}
		pages = append(pages, stringBuilder.String())
		stringBuilder.Reset()
	}

	return pages
}

func pagedDrawer(strings []string) {
	if len(strings) == 0 {
		return
	}
	pages := makePages(strings)

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
	fmt.Print(pages[currentPage], "\r")
	fmt.Printf("%d/%d\r\n", currentPage+1, pageCount)
	buf := make([]byte, 1)
	for {
		_, err = os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading input: %v\n", err)
			return
		}

		// Catch arrow left / arrow right
		if buf[0] == '\x1b' {
			seq := make([]byte, 2)
			_, err = os.Stdin.Read(seq)
			if err == nil && seq[0] == '[' {
				switch seq[1] {
				case 'C':
					if currentPage == pageCount-1 {
						continue
					}
					clearNLinesUp(pageSize + 1)
					currentPage++
					fmt.Print(pages[currentPage], "\r")
					fmt.Printf("%d/%d\r\n", currentPage+1, pageCount)
				case 'D':
					if currentPage == 0 {
						continue
					}
					clearNLinesUp(pageSize + 1)
					currentPage--
					fmt.Print(pages[currentPage], "\r")
					fmt.Printf("%d/%d\r\n", currentPage+1, pageCount)
				}
			}
			continue
		}
		if buf[0] == 'q' {
			clearNLinesUp(pageSize + 2)
			break
		}
	}
	fmt.Print("Success: table viewed\r\n")
}

func printWindow(strings []string, currentPosition int) {
	lastPosition := len(strings)
	for i := currentPosition; i < currentPosition+windowSize; i++ {
		if i >= lastPosition {
			stringBuilder.WriteString("--\r\n")
			continue
		}
		stringBuilder.WriteString(strings[i])
		stringBuilder.WriteString("\r\n")
	}
	fmt.Printf("↑ %d elements\r\n", currentPosition)
	fmt.Print(stringBuilder.String())
	fmt.Printf("↓ %d elements\r\n", max(0, lastPosition-currentPosition-windowSize))
	stringBuilder.Reset()
}

func scrolledDrawer(strings []string) {
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
	printWindow(strings, currentPosition)
	buf := make([]byte, 1)
	for {
		_, err = os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading input: %v\n", err)
			return
		}

		// Catch arrow up / arrow down
		if buf[0] == '\x1b' {
			seq := make([]byte, 2)
			_, err = os.Stdin.Read(seq)
			if err == nil && seq[0] == '[' {
				switch seq[1] {
				case 'A':
					if currentPosition == 0 {
						continue
					}

					clearNLinesUp(windowSize + 2)
					currentPosition--
					printWindow(strings, currentPosition)
				case 'B':
					if currentPosition+windowSize-1 >= lastPosition {
						continue
					}
					clearNLinesUp(windowSize + 2)
					currentPosition++
					printWindow(strings, currentPosition)
				}
			}
			continue
		}
		if buf[0] == 'q' {
			clearNLinesUp(windowSize + 3)
			break
		}
	}
	fmt.Print("Success: table viewed\r\n")
}
