package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"../document"
	termbox "github.com/nsf/termbox-go"
)

var doc = document.NewDocument(1)
var currentPos = document.EndPos

var connections []net.Conn

func render() {
	width, _ := termbox.Size() //@TODO should deal with height
	termbox.Clear(termbox.ColorBlack, termbox.ColorWhite)

	currentIdx, _ := doc.GetLineIndex(currentPos)
	currentIdx-- //Keeps the rendered cursor in the expected location

	content, err := doc.GetContent()
	if err != nil {
		panic(err)
	}

	runes := []rune(content)
	currX := 0
	currY := 0

	if currentIdx == 0 {
		termbox.SetCursor(0, 0)
	}

	for idx := range runes {
		ch := runes[idx]

		if ch == '\n' {
			ch = '\x00'
			currY++
			currX = 0
		} else {
			termbox.SetCell(currX, currY, ch, termbox.ColorBlack, termbox.ColorWhite)
			currX++
		}

		if currX >= width {
			currY++
			currX = 0
		}

		//When currentIdx is at the end it is beyond the range of displayed runes because of the empty string at EndPos
		if idx == currentIdx-1 {
			termbox.SetCursor(currX, currY)
		}
	}

	termbox.Flush()
}

func runUi() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for {

		currentIdx, _ := doc.GetLineIndex(currentPos)
		currentIdx-- //To mirror what happens graphically

		render()

		contentSize, _ := doc.GetLength()

		e := termbox.PollEvent()
		switch {
		// Quit
		case e.Key == termbox.KeyCtrlC:
			os.Exit(0)
		case e.Key == termbox.KeyBackspace || e.Key == termbox.KeyBackspace2:
			err = doc.DeleteLeft(currentPos)

			deletedIdx, _ := doc.GetLineIndex(currentPos)
			deletedLine := doc.GetLine(deletedIdx)
			toSendBytes := append([]byte{2}, deletedLine.ToBytes()...) //1 means delete

			for connIdx := range connections {
				connections[connIdx].Write(toSendBytes)
			}
		case e.Key == termbox.KeyArrowLeft:
			if currentIdx > 0 {
				currentPos, err = doc.Move(currentPos, -1)
			}
		case e.Key == termbox.KeyArrowRight:
			if currentIdx < contentSize-2 { //-2 for the Start and End Pos
				currentPos, err = doc.Move(currentPos, 1)
			}
		case e.Key == termbox.KeyEnter:
			//To do this we need to change the delimeter for the messages
			insertChar('\n')
		case e.Key == termbox.KeySpace:
			insertChar(' ')
		case e.Ch == '.' || e.Ch == '<' || e.Ch == '>': //Forbidden chars

		case e.Ch != 0: //Characters
			insertChar(e.Ch)
		default:
		}

		if err != nil {
			panic(err)
		}
	}
}

func insertChar(ch rune) {
	newPos, err := doc.InsertLeft(currentPos, string(ch))
	if err != nil {
		panic(err)
	}

	newIdx, _ := doc.GetLineIndex(newPos)
	newLine := doc.GetLine(newIdx)
	toSendBytes := append([]byte{1}, newLine.ToBytes()...) //1 means insert

	for connIdx := range connections {
		connections[connIdx].Write(toSendBytes)
	}
}

func listen(listenPort string) {
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		panic(err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		connections = append(connections, connection)

		numChars, _ := doc.GetLength()
		numChars -= 2

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {

	// Begin reading the data
	reader := bufio.NewReader(connection)
	for {
		bytes, err := reader.ReadBytes(0)
		if err != nil {
			panic(err)
		}

		operation := bytes[0]

		line := document.FromBytes(bytes[1 : len(bytes)-1]) //Remove the operation and delimiter
		switch operation {
		case 1:
			doc.Insert(line.Pos, line.Content)
		case 2:
			doc.DeleteLeft(line.Pos)
		}

		render()

	}
}

func main() {
	listenPort := os.Args[1]

	go listen(listenPort)

	for idx := 2; idx < len(os.Args); idx++ {
		connection, err := net.Dial("tcp", os.Args[idx])
		if err != nil {
			panic(err)
		}

		connections = append(connections, connection)

		go handleConnection(connection)
	}

	var f, _ = os.OpenFile("testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)

	runUi()
}
