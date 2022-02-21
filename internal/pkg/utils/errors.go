package utils

import (
	"fmt"
)

//ErrBadAccent indicate bad accent error
type ErrBadAccent struct {
	BadAccents []string
}

//NewErrBadAccent creates new error
func NewErrBadAccent(badAccents []string) *ErrBadAccent {
	return &ErrBadAccent{BadAccents: badAccents}
}

func (r *ErrBadAccent) Error() string {
	return fmt.Sprintf("Wrong accents: %v", r.BadAccents)
}

//ErrWordTooLong indicates too long word
type ErrWordTooLong struct {
	Word string
}

//NewErrWordTooLong creates new too long error
func NewErrWordTooLong(word string) *ErrWordTooLong {
	return &ErrWordTooLong{Word: word}
}

func (r *ErrWordTooLong) Error() string {
	return fmt.Sprintf("Wrong accent, too long word: '%s'", r.Word)
}
