package tool

func PanicIfErr(err ...error) {
	for _, e := range err {
		if e != nil {
			panic(e)
		}
	}
}
func HandleErr(err error, f func()) {
	if err != nil {
		f()
	}
}
