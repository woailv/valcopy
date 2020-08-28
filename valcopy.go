package valcopy

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//不通对象自动赋值
//src,dist 须为结构体指针
//数组字段须为值类型数组，不能为指针类型数组
//自动类型转换string->int64,uint64,int64,uint64->string,time.Time->string, string->time.Time，默认值转换:(i)=>i
//主要用于实体对象转pb文件对象，其余类型转换可能出现BUG
func ValMap(src, dist interface{}, mapFunc map[string]func(srcFieldVal interface{}) interface{}, pre ...string) {
	srv := reflect.ValueOf(src)
	for ; srv.Kind() == reflect.Ptr; {
		srv = srv.Elem()
	}
	srt := srv.Type()

	drv := reflect.ValueOf(dist)
	for ; drv.Kind() == reflect.Ptr; {
		drv = drv.Elem()
	}
	drt := drv.Type()

	preRes := "."
	if len(pre) == 1 {
		preRes = pre[0]
	}
	if mapFunc == nil {
		mapFunc = map[string]func(srcFieldVal interface{}) interface{}{}
	}

	switch drt.Kind() {
	case reflect.Int:
		drv.Set(reflect.ValueOf(i2Int(srv.Interface())))
	case reflect.Slice:
		distType := drt.Elem()
		n := 0
		for ; distType.Kind() == reflect.Ptr; {
			distType = distType.Elem()
			n++
		}
		switch distType.Kind() {
		case reflect.Struct:
			for k := 0; k < srv.Len(); k++ {
				distElem := reflect.New(distType)
				ValMap(srv.Index(k).Interface(), distElem.Interface(), mapFunc, preRes+drt.Name()+".") //drt.Name() 待测试
				distElem = distElem.Elem()
				for a := n; a > 0; a-- {
					distElem = distElem.Addr()
				}
				drv.Set(reflect.Append(drv, distElem))
			}
		case reflect.String:
			for k := 0; k < srv.Len(); k++ {
				drv.Set(reflect.Append(drv, reflect.ValueOf(i2String(srv.Index(k).Interface()))))
			}
		case reflect.Uint64:
			for k := 0; k < srv.Len(); k++ {
				drv.Set(reflect.Append(drv, reflect.ValueOf(i2Uint64(srv.Index(k).Interface()))))
			}
		case reflect.Int64:
			for k := 0; k < srv.Len(); k++ {
				drv.Set(reflect.Append(drv, reflect.ValueOf(i2Int64(srv.Index(k).Interface()))))
			}
		default:
			panic(fmt.Sprintf("未能处理的数组元素类型:%s", distType.Kind()))
		}
	case reflect.Struct:
		for i := 0; i < drv.NumField(); i++ {
			drtfi := drt.Field(i)
			drvfi := drv.Field(i)
			_, ok := srt.FieldByName(drtfi.Name)
			if !ok {
				continue
			}
			srvfi := srv.FieldByName(drtfi.Name)
			switch drtfi.Type.Kind() {
			case reflect.Ptr:
				if srvfi.IsNil() {
					continue
				}
				if drvfi.IsNil() {
					drvfi.Set(reflect.New(drtfi.Type.Elem()))
				}
				ValMap(srvfi.Interface(), drvfi.Interface(), mapFunc, preRes+drtfi.Name+".")
			case reflect.Slice:
				distType := drtfi.Type.Elem()
				n := 0
				for ; distType.Kind() == reflect.Ptr; {
					distType = distType.Elem()
					n++
				}
				switch distType.Kind() {
				case reflect.Struct:
					for k := 0; k < srvfi.Len(); k++ {
						distElem := reflect.New(distType)
						ValMap(srvfi.Index(k).Interface(), distElem.Interface(), mapFunc, preRes+drtfi.Name+".")
						distElem = distElem.Elem()
						for a := n; a > 0; a-- {
							distElem = distElem.Addr()
						}
						drvfi.Set(reflect.Append(drvfi, distElem))
					}
				case reflect.String:
					for k := 0; k < srvfi.Len(); k++ {
						drvfi.Set(reflect.Append(drvfi, reflect.ValueOf(i2String(srvfi.Index(k).Interface()))))
					}
				case reflect.Uint64:
					for k := 0; k < srvfi.Len(); k++ {
						drvfi.Set(reflect.Append(drvfi, reflect.ValueOf(i2Uint64(srvfi.Index(k).Interface()))))
					}
				case reflect.Int64:
					for k := 0; k < srvfi.Len(); k++ {
						drvfi.Set(reflect.Append(drvfi, reflect.ValueOf(i2Int64(srvfi.Index(k).Interface()))))
					}
				default:
					panic(fmt.Sprintf("未能处理的数组元素类型:%s", distType.Kind()))
				}
			default:
				key := preRes + drtfi.Name
				f, ok := mapFunc[key]
				if !ok {
					switch drvfi.Interface().(type) {
					case bool:
						f = i2Bool
					case int:
						f = i2Int
					case int8:
						f = i2Int8
					case int32:
						f = i2Int32
					case int64:
						f = i2Int64
					case uint:
						f = i2Uint
					case uint8:
						f = i2Uint8
					case uint32:
						f = i2Uint32
					case uint64:
						f = i2Uint64
					case float32:
						f = i2Float32
					case float64:
						f = i2Float64
					case string:
						f = i2String
					case time.Time:
						f = i2Time
					}
					if f == nil {
						f = func(srcFieldVal interface{}) interface{} {
							return srcFieldVal
						}
					}
				}
				res := f(srvfi.Interface())
				drvfi.Set(reflect.ValueOf(res))
			}
		}
	default:
		panic(fmt.Sprintf("未处理的列:%s类型:%s", drt.Name(), drt.Kind()))
	}
}

func i2Bool(i interface{}) interface{} {
	switch x := i.(type) {
	case bool:
		return x
	case int:
		return x > 0
	case int8:
		return x > 0
	case int32:
		return x > 0
	case int64:
		return x > 0
	case uint:
		return x > 0
	case uint8:
		return x > 0
	case uint32:
		return x > 0
	case uint64:
		return x > 0
	case float32:
		return x > 0
	case float64:
		return x > 0
	case string:
		return x != ""
	}
	return false
}

func i2Int(i interface{}) interface{} {
	var data int
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = x
	case int8:
		data = int(x)
	case int32:
		data = int(x)
	case int64:
		data = int(x)
	case uint:
		data = int(x)
	case uint8:
		data = int(x)
	case uint32:
		data = int(x)
	case uint64:
		data = int(x)
	case float32:
		data = int(x)
	case float64:
		data = int(x)
	case string:
		data, _ = strconv.Atoi(x)
	}
	return data
}

func i2Int8(i interface{}) interface{} {
	var data int8
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = int8(x)
	case int8:
		data = x
	case int32:
		data = int8(x)
	case int64:
		data = int8(x)
	case uint:
		data = int8(x)
	case uint8:
		data = int8(x)
	case uint32:
		data = int8(x)
	case uint64:
		data = int8(x)
	case float32:
		data = int8(x)
	case float64:
		data = int8(x)
	case string:
		t, _ := strconv.Atoi(x)
		data = int8(t)
	}
	return data
}

func i2Int32(i interface{}) interface{} {
	var data int32
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = int32(x)
	case int8:
		data = int32(x)
	case int32:
		data = x
	case int64:
		data = int32(x)
	case uint:
		data = int32(x)
	case uint8:
		data = int32(x)
	case uint32:
		data = int32(x)
	case uint64:
		data = int32(x)
	case float32:
		data = int32(x)
	case float64:
		data = int32(x)
	case string:
		t, _ := strconv.Atoi(x)
		data = int32(t)
	case time.Time:
		data = int32(x.Unix())
	}
	return data
}

func i2Int64(i interface{}) interface{} {
	var data int64
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = int64(x)
	case int8:
		data = int64(x)
	case int32:
		data = int64(x)
	case int64:
		data = x
	case uint:
		data = int64(x)
	case uint8:
		data = int64(x)
	case uint32:
		data = int64(x)
	case uint64:
		data = int64(x)
	case float32:
		data = int64(x)
	case float64:
		data = int64(x)
	case string:
		data = str2Int64(x)
	case time.Time:
		data = x.Unix()
	}
	return data
}

func i2Uint(i interface{}) interface{} {
	var data uint
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = uint(x)
	case int8:
		data = uint(x)
	case int32:
		data = uint(x)
	case int64:
		data = uint(x)
	case uint:
		data = (x)
	case uint8:
		data = uint(x)
	case uint32:
		data = uint(x)
	case uint64:
		data = uint(x)
	case float32:
		data = uint(x)
	case float64:
		data = uint(x)
	case string:
		data = uint(str2Uint64(x))
	case time.Time:
		data = uint(x.Unix())
	}
	return data
}

func i2Uint8(i interface{}) interface{} {
	var data uint8
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = uint8(x)
	case int8:
		data = uint8(x)
	case int32:
		data = uint8(x)
	case int64:
		data = uint8(x)
	case uint:
		data = uint8(x)
	case uint8:
		data = x
	case uint32:
		data = uint8(x)
	case uint64:
		data = uint8(x)
	case float32:
		data = uint8(x)
	case float64:
		data = uint8(x)
	case string:
		t, _ := strconv.Atoi(x)
		data = uint8(t)
	}
	return data
}

func i2Uint32(i interface{}) interface{} {
	var data uint32
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = uint32(x)
	case int8:
		data = uint32(x)
	case int32:
		data = uint32(x)
	case int64:
		data = uint32(x)
	case uint:
		data = uint32(x)
	case uint8:
		data = uint32(x)
	case uint32:
		data = x
	case uint64:
		data = uint32(x)
	case float32:
		data = uint32(x)
	case float64:
		data = uint32(x)
	case string:
		data = uint32(str2Uint64(x))
	case time.Time:
		data = uint32(x.Unix())
	}
	return data
}

func i2Uint64(i interface{}) interface{} {
	var data uint64
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = uint64(x)
	case int8:
		data = uint64(x)
	case int32:
		data = uint64(x)
	case int64:
		data = uint64(x)
	case uint:
		data = uint64(x)
	case uint8:
		data = uint64(x)
	case uint32:
		data = uint64(x)
	case uint64:
		data = x
	case float32:
		data = uint64(x)
	case float64:
		data = uint64(x)
	case string:
		data = str2Uint64(x)
	case time.Time:
		data = uint64(x.Unix())
	}
	return data
}

func i2Float32(i interface{}) interface{} {
	var data float32
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = float32(x)
	case int8:
		data = float32(x)
	case int32:
		data = float32(x)
	case int64:
		data = float32(x)
	case uint:
		data = float32(x)
	case uint8:
		data = float32(x)
	case uint32:
		data = float32(x)
	case uint64:
		data = float32(x)
	case float32:
		data = x
	case float64:
		data = float32(x)
	case string:
		t, _ := strconv.ParseFloat(x, 32)
		data = float32(t)
	}
	return data
}
func i2Float64(i interface{}) interface{} {
	var data float64
	switch x := i.(type) {
	case bool:
		if x {
			data = 1
		}
	case int:
		data = float64(x)
	case int8:
		data = float64(x)
	case int32:
		data = float64(x)
	case int64:
		data = float64(x)
	case uint:
		data = float64(x)
	case uint8:
		data = float64(x)
	case uint32:
		data = float64(x)
	case uint64:
		data = float64(x)
	case float32:
		data = float64(x)
	case float64:
		data = x
	case string:
		data, _ = strconv.ParseFloat(x, 32)
	}
	return data
}
func i2String(i interface{}) interface{} {
	var data string
	switch x := i.(type) {
	case bool:
		if x {
			data = "true"
		} else {
			data = "false"
		}
	case int:
		data = strconv.Itoa(x)
	case int8:
		data = strconv.Itoa(int(x))
	case int32:
		data = strconv.Itoa(int(x))
	case int64:
		data = strconv.FormatInt(x, 10)
	case uint:
		data = strconv.FormatUint(uint64(x), 10)
	case uint8:
		data = strconv.FormatUint(uint64(x), 10)
	case uint32:
		data = strconv.FormatUint(uint64(x), 10)
	case uint64:
		data = strconv.FormatUint(uint64(x), 10)
	case float32:
		data = strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		data = strconv.FormatFloat(x, 'f', -1, 64)
	case string:
		data = x
	case time.Time:
		data = time2Str(x)
	}
	return data
}

func i2Time(i interface{}) interface{} {
	var data time.Time
	switch x := i.(type) {
	case int:
		data = time.Unix(int64(x), 0)
	case int32:
		data = time.Unix(int64(x), 0)
	case int64:
		data = time.Unix(x, 0)
	case uint:
		data = time.Unix(int64(x), 0)
	case uint32:
		data = time.Unix(int64(x), 0)
	case uint64:
		data = time.Unix(int64(x), 0)
	case string:
		data = str2Time(x)
	case time.Time:
		data = x
	}
	return data
}

func time2Str(tm interface{}) string {
	t := tm.(time.Time)
	r := t.Format("2006-01-02 15:04:05")
	if r == "0001-01-01 00:00:00" || r == "1970-01-01 08:00:00" {
		return ""
	}
	return r
}

func str2Int64(str string) int64 {
	if str == "" {
		str = "0"
	}
	data, _ := strconv.ParseInt(str, 10, 64)
	return data
}

func str2Uint64(str string) uint64 {
	data, _ := strconv.ParseUint(str, 10, 64)
	return data
}

func str2Time(str string) time.Time {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	return tm
}
