package document

import (
	"bytes"
	"fmt"
	"strings"
)

//Not neccesarily an actual line with an EOL, more like a set of characters. The paper uses this term so I will too
type Line struct {
	Pos     []Position
	Content string
}

func (this *Line) ToString() string {
	var posBytes bytes.Buffer
	for posIdx := range this.Pos {
		if posIdx != 0 {
			posBytes.WriteString(".")
		}
		posBytes.WriteString(this.Pos[posIdx].ToString())
	}

	return fmt.Sprintf("%s%s", posBytes.String(), this.Content)
}

func (this *Line) ToBytes() []byte {
	bytes := []byte(this.ToString())
	bytes = append(bytes, 0) //Delimiter
	return bytes
}

func FromBytes(bytes []byte) (line Line) {
	sep := strings.Split(string(bytes), ".")
	lastSep := strings.SplitAfter(sep[len(sep)-1], ">") //Split the last term into the final position and the content
	sep[len(sep)-1] = lastSep[0]

	line.Content = lastSep[1]

	line.Pos = make([]Position, len(sep))

	for idx := range sep {
		line.Pos[idx] = FromString(sep[idx])
	}
	return
}
