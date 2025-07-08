package tool

func PanicIfErr(err ...error) {
	for _, e := range err {
		if e != nil {
			print(e.Error())
			panic(e)
		}
	}
}
func HandleErr(err error, f func()) {
	if err != nil {
		print(err.Error())
		f()
	}
}
