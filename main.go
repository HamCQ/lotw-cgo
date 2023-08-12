package main

// import "C"
import (
	"fmt"
	"path"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

//安全加载
var D = syscall.NewLazyDLL("tqsllib2.dll")

//匹配错误
var (
	errString    = ""
	errNoStation = "The selected station location could not be found"
)

func main() {

	ret, _, err := D.NewProc("tqsl_init").Call()
	if err != nil {
		e := err.(syscall.Errno)
		println(err.Error(), "errno =", e)
	}
	if ret != 0 {
		fmt.Println("tqsl_init err")
		return
	}
	fmt.Println("tqsl_init done")

	var major int
	var minor int
	ret, _, _ = D.NewProc("tqsl_getVersion").Call(uintptr(unsafe.Pointer(&major)), uintptr(unsafe.Pointer(&minor)))
	fmt.Printf("tqsl_getVersion ret=%d major=%d minor=%d \n", ret, major, minor)
	errToString()

	ret, _, _ = D.NewProc("tqsl_getConfigVersion").Call(uintptr(unsafe.Pointer(&major)), uintptr(unsafe.Pointer(&minor)))
	fmt.Printf("tqsl_getConfigVersion ret=%d major=%d minor=%d \n", ret, major, minor)
	errToString()

	// var serial int32
	// ret, _, _ = D.NewProc("tqsl_getSerialFromTQSLFile").
	// 	Call(strPtr("BG5UWQ.tq6"), uintptr(unsafe.Pointer(&serial)), 0, 0, 0, 0, 0)
	// fmt.Printf("tqsl_getSerialFromTQSLFile ret=%d serial=%v \n", ret, serial)
	// errToString()

	var userdata unsafe.Pointer
	//暂时不处理回调信息
	callback := syscall.NewCallback(func(r int, t uintptr, v uintptr) (ret uintptr) {
		// fmt.Printf("from callback: %v\n", r)
		// runtime.Gosched()
		return 0
	})
	ret, _, _ = D.NewProc("tqsl_importTQSLFile").
		Call(strPtr("BG5UWQ.tq6"), callback, uintptr(unsafe.Pointer(&userdata)))
	fmt.Printf("tqsl_importTQSLFile ret=%d userinfo=%v \n", ret, userdata)
	errToString()

	var locp unsafe.Pointer
	var callsign = "BG5UWQ"
	var locationName = "emin"
	// var location_name = make([]byte, 256)
	ret, _, _ = D.NewProc("tqsl_getStationLocation").Call(uintptr(unsafe.Pointer(&locp)), strPtr(locationName))
	fmt.Printf("tqsl_getStationLocation ret=%d %v \n", ret, locp)
	errString = errToString()
	if errString == errNoStation {
		ret, _, _ = D.NewProc("tqsl_initStationLocationCapture").Call(uintptr(unsafe.Pointer(&locp)))
		fmt.Printf("tqsl_initStationLocationCapture ret=%d \n", ret)
		errToString()

		ret, _, _ = D.NewProc("tqsl_setStationLocationCaptureName").Call(uintptr(unsafe.Pointer(locp)), strPtr(locationName))
		fmt.Printf("tqsl_setStationLocationCaptureName ret=%d \n", ret)
		errToString()

		ret, _, _ = D.NewProc("tqsl_saveStationLocationCapture").Call(uintptr(unsafe.Pointer(locp)), 0)
		fmt.Printf("tqsl_saveStationLocationCapture ret=%d \n", ret)
		errToString()

		// var locationNum int
		// ret, _, _ = D.NewProc("tqsl_getNumLocationField").Call(uintptr(unsafe.Pointer(locp)), uintptr(unsafe.Pointer(&locationNum)))
		// fmt.Printf("tqsl_getNumLocationField ret=%d numf=%v \n", ret, locationNum)
		// errToString()
		// if locationNum != 0 {
		// 	var i int = 0
		// 	for {
		// 		if i == locationNum {
		// 			break
		// 		}
		// 		var locationLable = make([]byte, 256)
		// 		ret, _, _ = D.NewProc("tqsl_getLocationFieldDataLabel").
		// 			Call(uintptr(unsafe.Pointer(locp)), intPtr(i-1), uintptr(unsafe.Pointer(&locationLable[0])), 256)
		// 		temp := string(locationLable[:256])
		// 		//Call Sign / DXCC Entity / Grid Square / ITU Zone / CQ Zone / IOTA ID
		// 		fmt.Printf("tqsl_getLocationFieldDataLabel ret=%d temp=%v \n", ret, temp)
		// 		errToString()
		// 		i++
		// 	}
		// }
	}

	ret, _, _ = D.NewProc("tqsl_getLocationCallSign").
		Call(uintptr(unsafe.Pointer(locp)), strPtr(callsign), intPtr(len(callsign)+1))
	fmt.Printf("tqsl_getLocationCallSign ret=%d \n", ret)
	errToString()

	var dxcc int
	ret, _, _ = D.NewProc("tqsl_getLocationDXCCEntity").
		Call(uintptr(unsafe.Pointer(locp)), uintptr(unsafe.Pointer(&dxcc)))
	fmt.Printf("tqsl_getLocationDXCCEntity ret=%d %v \n", ret, dxcc)
	errToString()

	var tqslCertList **unsafe.Pointer
	var nCert unsafe.Pointer
	ret, _, _ = D.NewProc("tqsl_selectCertificates").
		Call(uintptr(unsafe.Pointer(&tqslCertList)), uintptr(unsafe.Pointer(&nCert)), strPtr(callsign), 0, 0, 0, 0)
	fmt.Printf("tqsl_selectCertificates ret=%d  \n", ret)
	errToString()

	var tqslConvp unsafe.Pointer
	ret, _, _ = D.NewProc("tqsl_beginADIFConverter").
		Call(uintptr(unsafe.Pointer(&tqslConvp)),
			strPtr("test.adi"), uintptr(unsafe.Pointer(tqslCertList)),
			uintptr(unsafe.Pointer(nCert)), uintptr(unsafe.Pointer(locp)))
	fmt.Printf("tqsl_beginADIFConverter ret=%d %v \n", ret, tqslConvp)
	errToString()

	ret, _, _ = D.NewProc("tqsl_setConverterAllowDuplicates").
		Call(uintptr(tqslConvp), 0)
	fmt.Printf("tqsl_setConverterAllowDuplicates ret=%d \n", ret)
	errToString()

	// var app = make([]byte, 256)
	// var p1 uintptr
	// ret, _, _ = D.NewProc("tqsl_setConverterAppName").
	// 	Call(uintptr(tqslConvp), strPtr("h"))
	// fmt.Printf("tqsl_setConverterAppName ret=%d app\n", ret)
	// errToString()

	for {
		ret, _, _ = D.NewProc("tqsl_getConverterGABBI").
			Call(uintptr(tqslConvp))
		fmt.Printf("tqsl_getConverterGABBI ret=%d \n", ret)
		errToString()
		if ret == 0 {
			break
		} else {
			p := (*byte)(unsafe.Pointer(ret))
			data := make([]byte, 0)

			for *p != 0 {
				data = append(data, *p)
				ret += unsafe.Sizeof(byte(0))
				p = (*byte)(unsafe.Pointer(ret))
			}
			fmt.Println(string(data))
		}
	}

	ret, _, _ = D.NewProc("tqsl_converterCommit").
		Call(uintptr(tqslConvp))
	fmt.Printf("tqsl_converterCommit ret=%d \n", ret)
	errToString()

	ret, _, _ = D.NewProc("tqsl_endConverter").
		Call(uintptr(unsafe.Pointer(&tqslConvp)))
	fmt.Printf("tqsl_endConverter ret=%d \n", ret)
	errToString()

	ret, _, _ = D.NewProc("tqsl_endStationLocationCapture").
		Call(uintptr(unsafe.Pointer(&locp)))
	fmt.Printf("tqsl_endStationLocationCapture ret=%d \n", ret)
	errToString()
}

func setDir() {
	ret, _, _ := D.NewProc("tqsl_setDirectory").
		Call(strPtr(getCurrentAbPathByCaller() + "\\data"))
	fmt.Printf("tqsl_setDirectory ret=%d \n", ret)
	errToString()
}

func UintPtrToString(r uintptr) string {
	p := (*uint16)(unsafe.Pointer(r))
	if p == nil {
		return ""
	}

	n, end, add := 0, unsafe.Pointer(p), unsafe.Sizeof(*p)
	for *(*uint16)(end) != 0 {
		end = unsafe.Add(end, add)
		n++
	}
	return string(utf16.Decode(unsafe.Slice(p, n)))
}

func errToString() string {
	ret, _, _ := D.NewProc("tqsl_getErrorString").Call()
	p := (*byte)(unsafe.Pointer(ret))
	data := make([]byte, 0)

	for *p != 0 {
		data = append(data, *p)
		ret += unsafe.Sizeof(byte(0))
		p = (*byte)(unsafe.Pointer(ret))
	}
	fmt.Println(string(data))
	return string(data)
}

// 获取字符串的指针
func strPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringBytePtr(s)))
}

// 获取数字的指针
func intPtr(n int) uintptr {
	return uintptr(n)
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
