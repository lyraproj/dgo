package util

// Catch calls the given function and recovers and returns any panic that is thrown with an error. The function
// returns nil if f doesn't panic.
func Catch(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if er, ok := e.(error); ok {
				err = er
				return
			}
			// Just catch errors.
			panic(e)
		}
	}()
	f()
	return
}
