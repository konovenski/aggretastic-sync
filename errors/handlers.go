//Package errors contain high-level handlers for errors.
package errors

//Panic with msg error if err is not nill
func PanicOnError(err error, msg error) {
	if err != nil {
		if msg != nil {
			panic(msg.Error() + err.Error())
		}
		panic(err)
	}
}
