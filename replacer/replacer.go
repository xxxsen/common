package replacer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type valueSearcherFunc func(key string) (interface{}, bool)

func innerReplace(str string, searcher valueSearcherFunc) string {
	sb := strings.Builder{}
	sb.Grow(len(str))
	for i := 0; i < len(str); i++ {
		leftBr := -1
		foundRightBr := false
		var j int
		for j = i; j < len(str); j++ {
			switch str[j] {
			case '{':
				leftBr = j
				continue
			case '}':
				foundRightBr = true
				break
			default:
				continue
			}
			//左括号未找到, 以起始点i开始复制到j(包含j)的全部数据
			if leftBr < 0 {
				sb.WriteString(str[i : j+1])
				break
			}
			key := str[leftBr+1 : j]
			value, ok := searcher(key)
			//左右括号均找到, 但是括号中的key在map中未找到, 此时不进行替换
			if !ok {
				sb.WriteString(str[i : j+1])
				break
			}
			//左右括号找到且key找到, 将内容进行替换
			sb.WriteString(str[i:leftBr])
			sb.WriteString(asStrValue(value))
			break
		}
		//当前循环直接遍历到结尾都没有找到右括号, 那么直接复制全部数据
		if !foundRightBr {
			sb.WriteString(str[i:j])
		}
		i = j
	}
	return sb.String()
}

func asStrValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case bool:
		return strconv.FormatBool(v)
	case complex64:
		return strconv.FormatComplex(complex128(v), 'g', -1, 64)
	case complex128:
		return strconv.FormatComplex(v, 'g', -1, 128)
	case nil:
		return "<nil>"
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		res, _ := json.Marshal(v)
		return string(res)
	}
}

func ReplaceByMap(str string, m map[string]interface{}) string {
	if len(str) == 0 || len(m) == 0 {
		return str
	}
	return innerReplace(str, func(key string) (interface{}, bool) {
		v, ok := m[key]
		return v, ok
	})
}

func ReplaceByList[T any](str string, in ...T) string {
	if len(str) == 0 || len(in) == 0 {
		return str
	}
	return innerReplace(str, func(key string) (interface{}, bool) {
		if len(key) == 0 || (len(key) > 1 && key[0] == '0') {
			return nil, false
		}
		index, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return nil, false
		}
		if int(index) < len(in) {
			return in[index], true
		}
		return nil, false
	})
}
