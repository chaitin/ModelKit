package diff

import (
	"strings"
)

type ConflictBlock struct {
	OursContent   []string
	TheirsContent []string
	StartLine     int
	EndLine       int
}

type ConflictParser struct {
	lines       []string
	currentLine int
	conflicts   []ConflictBlock
}

func NewConflictParser(text string) *ConflictParser {
	return &ConflictParser{
		lines: strings.Split(text, "\n"),
	}
}

func (p *ConflictParser) ParseConflicts() []ConflictBlock {
	for p.currentLine < len(p.lines) {
		line := p.lines[p.currentLine]

		if p.isConflictStart(line) {
			conflict := p.parseConflictBlock()
			if conflict != nil {
				p.conflicts = append(p.conflicts, *conflict)
			}
		} else {
			p.currentLine++
		}
	}
	return p.conflicts
}

func (p *ConflictParser) isConflictStart(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 7 {
		return false
	}

	for i := range 7 {
		if trimmed[i] != '<' {
			return false
		}
	}
	return true
}

func (p *ConflictParser) isConflictSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) != 7 {
		return false
	}

	for i := range 7 {
		if trimmed[i] != '=' {
			return false
		}
	}
	return true
}

func (p *ConflictParser) isConflictEnd(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 7 {
		return false
	}

	for i := range 7 {
		if trimmed[i] != '>' {
			return false
		}
	}
	return true
}

func (p *ConflictParser) parseConflictBlock() *ConflictBlock {
	startLine := p.currentLine
	p.currentLine++

	conflict := &ConflictBlock{
		StartLine: startLine,
	}

	for p.currentLine < len(p.lines) {
		line := p.lines[p.currentLine]

		if p.isConflictSeparator(line) {
			p.currentLine++
			break
		}

		conflict.OursContent = append(conflict.OursContent, line)
		p.currentLine++
	}

	for p.currentLine < len(p.lines) {
		line := p.lines[p.currentLine]

		if p.isConflictEnd(line) {
			conflict.EndLine = p.currentLine
			p.currentLine++
			return conflict
		}

		conflict.TheirsContent = append(conflict.TheirsContent, line)
		p.currentLine++
	}

	return nil
}

func CountAddedLines(text string) int {
	lines := strings.Split(text, "\n")
	addedLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "/*") {
			addedLines++
		}
	}

	return addedLines
}

func (cb *ConflictBlock) CountAddedLinesInConflict() int {
	theirsText := strings.Join(cb.TheirsContent, "\n")
	return CountAddedLines(theirsText)
}

func (cb *ConflictBlock) CountNetAddedLines() int {
	oursText := strings.Join(cb.OursContent, "\n")
	theirsText := strings.Join(cb.TheirsContent, "\n")
	oursLines := CountAddedLines(oursText)
	theirsLines := CountAddedLines(theirsText)
	return theirsLines - oursLines
}

func ParseConflictsAndCountLines(text string) []int {
	parser := NewConflictParser(text)
	conflicts := parser.ParseConflicts()

	var addedLines []int
	for _, conflict := range conflicts {
		addedLines = append(addedLines, conflict.CountAddedLinesInConflict())
	}

	return addedLines
}
