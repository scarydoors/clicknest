package errorutil

type joinedErr interface { Unwrap() [] error }

func IntoSlice(err error) []error {
	if jerr, ok := err.(joinedErr); ok {
		return jerr.Unwrap()
	} else {
		return []error{err}
	}
}
