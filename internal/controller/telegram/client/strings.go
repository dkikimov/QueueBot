package client

import "strings"

// cutStringByLines cuts string so 'linesToHave' lines are left in the end.
// It returns string with "...." before the lines.
func cutStringByLines(s string, linesToHave int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= linesToHave {
		return s
	}

	b := strings.Builder{}
	b.WriteString("....\n")
	b.WriteString(strings.Join(lines[len(lines)-linesToHave:], "\n"))

	return b.String()
}

// cutStringByLinesWithCurrent cuts string by lines where current index is in the middle and halfOfLinesAroundCenter lines are around it.
// It returns string with "...." before and after the lines.
func cutStringByLinesWithCurrent(s string, halfOfLinesAroundCenter int, currentIdx int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= 2*halfOfLinesAroundCenter {
		return s
	}

	b := strings.Builder{}

	b.WriteString("....\n")
	b.WriteString(strings.Join(lines[currentIdx-halfOfLinesAroundCenter:currentIdx+halfOfLinesAroundCenter+1], "\n"))
	b.WriteString("\n....")

	return b.String()
}
