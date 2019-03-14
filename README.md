


https://dev.azure.com/marcellanz/file-read-challenge-go/

go get github.com/davecheney/gcvis


func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func BytesToString2(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	runtime.KeepAlive(&bytes) // this line is essential.
	return s
}

func BytesToString3(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func BytesToString4(bytes []byte) (s string) {
	if len(bytes) == 0 {
		return ""
	}
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = uintptr(unsafe.Pointer(&bytes[0]))
	str.Len = len(bytes)
	return s
}

		//s0 := string(scanner.Bytes())
		//if len(s0) == 0 {
		//}

		//willScan := true
		//s, err := bufio.NewReader(file).ReadString('\n')
		//fmt.Printf("%v\n", s)
		//if(err != nil) {
		//	willScan = false
		//}
		//lines = append(lines, s)

		//b := scanner.Bytes()
		//bytes := make([]byte, len(b))
		//copy(bytes, b)
		//lines = append(lines, BytesToString4(bytes))
