package scanner

import (
	"strings"
	"unsafe"

	"github.com/t14raptor/go-fast/ast"
)

func isLineTerminator(chr rune) bool {
	switch chr {
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}

// Irregular line breaks - '\u{2028}' (LS) and '\u{2029}' (PS)
// These are 3-byte UTF-8 sequences starting with 0xE2.
const lsOrPsFirst byte = 0xE2

var lsBytes2And3 = [2]byte{0x80, 0xA8}
var psBytes2And3 = [2]byte{0x80, 0xA9}

// Matches: '\r', '\n', 0xE2 (first byte of LS/PS).
var lineBreakTable [256]bool

// Matches: '*', '\r', '\n', 0xE2 (first byte of LS/PS).
var multiLineCommentTable [256]bool

func init() {
	for _, b := range []byte{'\r', '\n', lsOrPsFirst} {
		lineBreakTable[b] = true
	}
	for _, b := range []byte{'*', '\r', '\n', lsOrPsFirst} {
		multiLineCommentTable[b] = true
	}
}

// skipSingleLineComment skips a single-line comment (// already consumed).
// Does NOT consume the line terminator.
func (s *Scanner) skipSingleLineComment() {
	pos := s.src.pos
	base := s.src.base
	end := s.src.len

	for pos < end {
		b := *(*byte)(unsafe.Add(base, pos))
		if !lineBreakTable[b] {
			pos++
			continue
		}

		if b != lsOrPsFirst {
			s.src.pos = pos
			return
		}

		s.src.pos = pos + 1
		twoMore, ok := s.src.PeekTwoBytes()
		if !ok {
			pos = s.src.pos
			continue
		}
		if twoMore == lsBytes2And3 || twoMore == psBytes2And3 {
			s.src.pos = pos
			return
		}
		s.src.pos += 2
		pos = s.src.pos
	}
	s.src.pos = pos
}

// skipMultiLineComment skips a multi-line comment (/* already consumed).
// Sets s.Token.OnNewLine if the comment contains a line terminator.
// After finding a line break, switches to a faster path that only looks for `*/`.
func (s *Scanner) skipMultiLineComment() {
	pos := s.src.pos
	base := s.src.base
	end := s.src.len

	for pos < end {
		b := *(*byte)(unsafe.Add(base, pos))
		if !multiLineCommentTable[b] {
			pos++
			continue
		}

		switch b {
		case '*':
			pos++
			if pos < end && *(*byte)(unsafe.Add(base, pos)) == '/' {
				s.src.pos = pos + 1
				return
			}
		case '\r', '\n':
			s.Token.OnNewLine = true
			s.src.pos = pos + 1
			s.skipMultiLineCommentAfterLineBreak()
			return
		default:
			pos++
			s.src.pos = pos
			twoMore, ok := s.src.PeekTwoBytes()
			if !ok {
				pos = s.src.pos
				continue
			}
			if twoMore == lsBytes2And3 || twoMore == psBytes2And3 {
				s.Token.OnNewLine = true
				s.src.pos += 2
			} else {
				s.src.pos += 2
			}
			pos = s.src.pos
		}
	}

	s.src.pos = pos
	s.error(unterminatedMultiLineComment(s.unterminatedRange()))
}

// skipMultiLineCommentAfterLineBreak is the fast path for multi-line comment scanning
// after a line break has been found. Only needs to search for `*/`.
func (s *Scanner) skipMultiLineCommentAfterLineBreak() {
	remaining := s.src.Slice(s.src.Offset(), s.src.EndOffset())
	idx := strings.Index(remaining, "*/")
	if idx >= 0 {
		// Found `*/`, advance past it
		s.src.SetPosition(s.src.Offset() + ast.Idx(idx) + 2)
	} else {
		// Unterminated comment - advance to end
		s.src.SetPosition(s.src.EndOffset())
	}
}
