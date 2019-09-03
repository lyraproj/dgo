package internal

import "testing"

func Test_decodeEndToken(t *testing.T) {
	if `unexpected delimiter ']'` != endOfArray.Error() {
		t.Failed()
	}
}
