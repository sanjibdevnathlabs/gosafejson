//go:build !jsoniter_sloppy
// +build !jsoniter_sloppy

package gosafejson

import (
	"fmt"
	"io"
)

// Reverted skipNumber
func (iter *Iterator) skipNumber() {
	if !iter.trySkipNumber() {
		iter.unreadByte()
		if iter.Error != nil && iter.Error != io.EOF {
			return
		}
		// Use ReadFloat64 for strict validation (handles leading zeros)
		iter.ReadFloat64()
		// If ReadFloat64 reported a format error (or any non-EOF error),
		// that error is the result for Skip(). Return immediately.
		if iter.Error != nil && iter.Error != io.EOF {
			return // Keep error from ReadFloat64, do not attempt ReadBigFloat
		}
		// If ReadFloat64 succeeded or only hit EOF, the skip is successful for validation purposes.
		// We don't need to read the potentially larger number with ReadBigFloat for Skip().
	}
}

// Reverted trySkipNumber
func (iter *Iterator) trySkipNumber() bool {
	dotFound := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if dotFound {
				iter.ReportError("validateNumber", `more than one dot found in number`)
				return true // already failed
			}
			if i+1 == iter.tail {
				return false
			}
			c = iter.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				iter.ReportError("validateNumber", `missing digit after dot`)
				return true // already failed
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if iter.head == i {
					return false // if - without following digits
				}
				iter.head = i
				return true // must be valid
			}
			return false // may be invalid
		}
	}
	return false
}

// Reverted skipString
func (iter *Iterator) skipString() {
	if !iter.trySkipString() {
		iter.unreadByte()
		iter.ReadString()
	}
}

// Reverted trySkipString
func (iter *Iterator) trySkipString() bool {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			iter.head = i + 1
			return true // valid
		} else if c == '\\' {
			return false
		} else if c < ' ' {
			iter.ReportError("trySkipString",
				fmt.Sprintf(`invalid control character found: %d`, c))
			return true // already failed
		}
	}
	return false
}

// Reverted skipObject
func (iter *Iterator) skipObject() {
	iter.unreadByte()
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		iter.Skip()
		return true
	})
}

// Reverted skipArray
func (iter *Iterator) skipArray() {
	iter.unreadByte()
	iter.ReadArrayCB(func(iter *Iterator) bool {
		iter.Skip()
		return true
	})
}
