package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const mask = 256 - 1

type color int64

func (c color) String() string {
	// return strconv.FormatInt(int64(c), 16)
	return fmt.Sprintf("%06x", int(c))
}

var palette = []color{}

func main() {
	s, err := ioutil.ReadFile("palette.txt")
	if err != nil {
		panic(err)
	}
	ls := bytes.Split(s, []byte{'\n'})
	for _, x := range ls {
		i, err := strconv.ParseInt(string(x[1:]), 16, 64)
		if err != nil {
			continue
		}
		palette = append(palette, color(i))
	}
	s, err = ioutil.ReadFile("themes/source.json")
	if err != nil {
		panic(err)
	}
	ls = bytes.Split(s, []byte{'\n'})
	for _, x := range ls {
		fmt.Println(replaceColor(string(x)))
	}
}

func replaceColor(line string) string {
	re := regexp.MustCompile(`"#(.*)"`)
	findings := re.FindAllString(line, 1)
	if len(findings) == 0 {
		return line
	}
	finding := findings[0]
	d := len(finding) - 8
	i, err := strconv.ParseInt(finding[2:len(finding)-2-d], 16, 64)
	if err != nil {
		panic(err)
	}
	c := color(i)
	return strings.Replace(line, finding, fmt.Sprintf(`"#%s"`, findClosestBaseColor(c)), 1)
}

func findClosestBaseColor(c color) color {
	minDist := uint64(math.MaxUint64)
	R, G, B := decomp(c)
	var res color
	for _, v := range palette {
		R_, G_, B_ := decomp(v)
		dist := sqDiff(R, R_) + sqDiff(G, G_) + sqDiff(B, B_)
		if dist < minDist {
			minDist = dist
			res = v
		}
	}
	return res
}

func decomp(color color) (a, b, c uint16) {
	a = uint16(color & mask)
	b = uint16((color & (mask << 8)) >> 8)
	c = uint16((color & (mask << 16)) >> 16)
	return
}

func sqDiff(a, b uint16) uint64 {
	if a < b {
		a, b = b, a
	}
	diff := uint64(a - b)
	return diff * diff
}
