package parser

import (
	"bytes"
 	"fmt"
)


type _RegExp_parser struct {
	str    string
	length int

	chr       rune // The current character
	chrOffset int  // The offset of current character
	offset    int  // The offset after current character (may be greater than 1)

	errors  []error
	invalid bool // The input is an invalid JavaScript RegExp

	goRegexp *bytes.Buffer
}

// TODO Better error reporting, use the offset, etc.
func (self *_RegExp_parser) error(offset int, msg string, msgValues ...interface{}) error {
	err := fmt.Errorf(msg, msgValues...)
	self.errors = append(self.errors, err)
	return err
}
