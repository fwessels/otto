package otto

import (
	"github.com/dlclark/regexp2"
	"unicode/utf8"
)

type _regExp2Object struct {
	regularExpression *regexp2.Regexp
	global            bool
	ignoreCase        bool
	multiline         bool
	source            string
	flags             string
}

func (runtime *_runtime) newRegExp2Object(pattern string, flags string) *_object {
	self := runtime.newObject()
	self.class = "RegExp"

	global := false
	ignoreCase := false
	multiline := false

	// TODO Pass in flags properly

	regularExpression, err := regexp2.Compile(pattern, 0)
	if err != nil {
		panic(runtime.panicSyntaxError("Invalid regular expression 1: %s", err.Error()[22:]))
	}

	self.value = _regExp2Object{
		regularExpression: regularExpression,
		global:            global,
		ignoreCase:        ignoreCase,
		multiline:         multiline,
		source:            pattern,
		flags:             flags,
	}
	self.defineProperty("global", toValue_bool(global), 0, false)
	self.defineProperty("ignoreCase", toValue_bool(ignoreCase), 0, false)
	self.defineProperty("multiline", toValue_bool(multiline), 0, false)
	self.defineProperty("lastIndex", toValue_int(0), 0100, false)
	self.defineProperty("source", toValue_string(pattern), 0, false)
	return self
}

func (self *_object) regExpValue() _regExp2Object {
	value, _ := self.value.(_regExp2Object)
	return value
}

func execRegExp(this *_object, target string) (match bool, result []int) {
	if this.class != "RegExp" {
		panic(this.runtime.panicTypeError("Calling RegExp.exec on a non-RegExp object"))
	}
	lastIndex := this.get("lastIndex").number().int64
	index := lastIndex
	global := this.get("global").bool()
	if !global {
		index = 0
	}
	if 0 > index || index > int64(len(target)) {
	} else {
		result = FindStringSubmatchIndex(this.regExpValue().regularExpression, target[index:])
	}
	if result == nil {
		//this.defineProperty("lastIndex", toValue_(0), 0111, true)
		this.put("lastIndex", toValue_int(0), true)
		return // !match
	}
	match = true
	startIndex := index
	endIndex := int(lastIndex) + result[1]
	// We do this shift here because the .FindStringSubmatchIndex above
	// was done on a local subordinate slice of the string, not the whole string
	for index, _ := range result {
		result[index] += int(startIndex)
	}
	if global {
		//this.defineProperty("lastIndex", toValue_(endIndex), 0111, true)
		this.put("lastIndex", toValue_int(endIndex), true)
	}
	return // match
}

func execResultToArray(runtime *_runtime, target string, result []int) *_object {
	captureCount := len(result) / 2
	valueArray := make([]Value, captureCount)
	for index := 0; index < captureCount; index++ {
		offset := 2 * index
		if result[offset] != -1 {
			valueArray[index] = toValue_string(target[result[offset]:result[offset+1]])
		} else {
			valueArray[index] = Value{}
		}
	}
	matchIndex := result[0]
	if matchIndex != 0 {
		matchIndex = 0
		// Find the rune index in the string, not the byte index
		for index := 0; index < result[0]; {
			_, size := utf8.DecodeRuneInString(target[index:])
			matchIndex += 1
			index += size
		}
	}
	match := runtime.newArrayOf(valueArray)
	match.defineProperty("input", toValue_string(target), 0111, false)
	match.defineProperty("index", toValue_int(matchIndex), 0111, false)
	return match
}
