package liner

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// print prompt
// interrupt
// read symbol
// move cursor left/right
// read enter and print new empty line
// backspace
// insert symbol in line

const (
	csiMoveCursorLeft     string = "\x1b[D"
	csiMoveCursorRight    string = "\x1b[C"
	csiCleanLine          string = "\x1b[2K"
	csiSetCursorToColumnD string = "\x1b[%dG"
)

// TODO: add gracefull shutdown

// Liner - struct for reading terminal input
type Liner struct {
	prompt      string
	currentLine []rune
	cursorPos   int

	readchan chan string
	waiter   chan struct{}
}

func NewLiner(prompt string) *Liner {
	return &Liner{
		prompt:      prompt,
		currentLine: []rune{},
		cursorPos:   0,
		readchan:    make(chan string),
		waiter:      make(chan struct{}),
	}
}

func (l *Liner) Close() {
	close(l.readchan)
}

func (l *Liner) ReadLiner() <-chan string {
	return l.readchan
}

func (l *Liner) Waiter() chan<- struct{} {
	return l.waiter
}

// Run - start terminal
func (l *Liner) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\r%s", l.prompt)

	readline:
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				log.Panic(err)
			}

			switch r {
			case 10, 13:
				l.enter()
				break readline
			case 127:
				l.backspace()
			case 27:
				nextRunes, _ := reader.Peek(2)
				if len(nextRunes) == 2 && nextRunes[0] == '[' {
					reader.Discard(2)

					switch nextRunes[1] {
					case 'A':
						// fmt.Println("its up arrow!")
					case 'C':
						l.moveRight()
					case 'B':
						// fmt.Println("its down arrow!")
					case 'D':
						l.moveLeft()
					}
				}
			case 4:
				// ctrl+D
				return
			default:
				l.insertRune(r)
			}
		}
	}
}

func (l *Liner) reDrawLine() {
	fmt.Printf("%s\r", csiCleanLine)                               // clean line
	fmt.Printf("%s%s", l.prompt, string(l.currentLine))            // print current string
	fmt.Printf(csiSetCursorToColumnD, 1+len(l.prompt)+l.cursorPos) // move cursor to pos
}

func (l *Liner) insertRune(r rune) {
	if l.cursorPos == len(l.currentLine) {
		l.currentLine = append(l.currentLine, r)
		l.cursorPos++
		fmt.Print(string(r))
	} else {
		l.currentLine = append(l.currentLine[:l.cursorPos], append([]rune{r}, l.currentLine[l.cursorPos:]...)...)
		l.currentLine[l.cursorPos] = r
		l.cursorPos++

		l.reDrawLine()
	}
}

func (l *Liner) enter() {
	fmt.Printf("\n\r")

	l.readchan <- string(l.currentLine)

	<-l.waiter // wait for (ATTENTION POTENTIAL DEADLOCK)

	l.currentLine = []rune{}
	l.cursorPos = 0
}

func (l *Liner) backspace() {
	if l.cursorPos > 0 {
		l.currentLine = append(l.currentLine[:l.cursorPos-1], l.currentLine[l.cursorPos:]...)
		l.cursorPos--

		l.reDrawLine()
	}
}

func (l *Liner) moveLeft() {
	if l.cursorPos > 0 {
		l.cursorPos--
		fmt.Print(csiMoveCursorLeft)
	}
}

func (l *Liner) moveRight() {
	if l.cursorPos < len(l.currentLine) {
		l.cursorPos++
		fmt.Print(csiMoveCursorRight)
	}
}
