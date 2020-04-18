package document

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Position struct {
	posId uint8
	site  int
}

//These are the pre-existing lines at the start and end of the document. All following inserts will be somewhere between these
var (
	StartPos = []Position{{0, 0}}
	EndPos   = []Position{{math.MaxUint8, 0}}
)

//Random number between x and y non inclusive on both ends (x,y)
func randBetween(x int, y int) int {
	return int(rand.Intn(int(y-x-1)) + int(x) + 1)
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

//Returns 1 if x > y  -1 if x < y, and 0 if x = y
func Compare(x []Position, y []Position) int {
	// for idx := range x {
	// 	if len(y) == idx {
	// 		return 1
	// 	}

	// 	xPos := x[idx].posId
	// 	yPos := y[idx].posId

	// 	if xPos > yPos {
	// 		return 1
	// 	} else if xPos < yPos {
	// 		return -1
	// 	}
	// }

	// if len(x) < len(y) {
	// 	return -1
	// }

	// fmt.Printf("%s %s\n", ToString(x), ToString(y))

	for idx := range x {
		if idx == len(y) {
			return 1
		}

		xPos := x[idx].posId
		yPos := y[idx].posId

		xSite := x[idx].site
		ySite := y[idx].site

		if xPos == yPos && xSite == ySite {
			continue
		}

		if xPos < yPos {
			return -1
		} else if xPos == yPos && xSite < ySite {
			return -1
		} else {
			return 1
		}

	}

	if len(x) < len(y) {
		return -1
	}

	return 0
}

//According the definitions given the paper
//@TODO make sure I am assigning the sites correctly, I need to read that part of paper more
func GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	//@TODO verify l < r

	diffenceLen := len(r) - len(l)
	smallestLen := min(len(r), len(l))
	addFinalDigit := false
	var lastIdx int

	for idx := 0; idx < smallestLen; idx++ {
		lastIdx = idx

		lPos := l[idx]
		rPos := r[idx]

		var difference = rPos.posId - lPos.posId

		if difference == 0 {
			//Add the digit, same so it doesnt matter
			pos = append(pos, Position{rPos.posId, rPos.site})
			addFinalDigit = true
		} else if difference == 1 {
			if idx < len(l)-1 {
				pos = append(pos, Position{lPos.posId, lPos.site})
				addFinalDigit = true
				break
			} else {
				pos = append(pos, Position{lPos.posId, lPos.site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), lPos.site})
				addFinalDigit = false
			}
		} else if difference > 1 {
			if idx < len(l)-1 {
				pos = append(pos, Position{lPos.posId, lPos.site})
				addFinalDigit = true
				break
			} else {
				pos = append(pos, Position{uint8(randBetween(int(lPos.posId), int(rPos.posId))), lPos.site})
				addFinalDigit = false
			}
		} else {
			panic("Difference GeneratePositionBetween() is less than 0! This should never happen")
		}
	}

	if addFinalDigit {
		if diffenceLen < 0 || (lastIdx < len(r)-1 && lastIdx < len(l)-1) {
			lastIdx++
			nextLeftPos := l[lastIdx].posId
			for nextLeftPos == math.MaxUint8 {
				pos = append(pos, Position{math.MaxUint8, site})
				lastIdx++
				nextLeftPos = l[lastIdx].posId
			}

			if nextLeftPos == math.MaxUint8 || nextLeftPos == math.MaxUint8-1 {
				pos = append(pos, Position{math.MaxUint8, site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
			} else {
				pos = append(pos, Position{uint8(randBetween(int(nextLeftPos), math.MaxUint8)), site})
			}
		} else if diffenceLen > 0 {
			lastIdx++
			nextRightPos := r[lastIdx].posId
			for nextRightPos == 0 {
				pos = append(pos, Position{0, site})
				lastIdx++
				nextRightPos = r[lastIdx].posId
			}

			if nextRightPos == 0 || nextRightPos == 1 {
				pos = append(pos, Position{0, site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
			} else {
				pos = append(pos, Position{uint8(randBetween(0, int(nextRightPos))), site})
			}
		}
	}

	return
}

func (this *Position) ToString() string {
	return fmt.Sprintf("<%d,%d>", this.posId, this.site)
}

func ToString(pos []Position) string {
	var posBytes bytes.Buffer
	for posIdx := range pos {
		if posIdx != 0 {
			posBytes.WriteString(".")
		}
		posBytes.WriteString(pos[posIdx].ToString())
	}

	return posBytes.String()
}

func FromString(str string) (pos Position) {
	str = str[1 : len(str)-2] //Remove first and last indices
	sep := strings.Split(str, ",")

	posIdParsed, _ := strconv.ParseUint(sep[0], 10, 8)
	siteParsed, _ := strconv.ParseUint(sep[1], 10, 32)

	pos.posId = uint8(posIdParsed)
	pos.site = int(siteParsed)

	return
}
