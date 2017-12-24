package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var j interface{}
		json.Unmarshal([]byte(line[:len(line)-1]), &j)
		m := j.(map[string]interface{})
		fmt.Println(sum(m))
	}
}

func sumList(l []interface{}) int {
	total := 0
	for i := range l {
		switch v := l[i].(type) {
		case int:
			total += v
		case int64:
			total += int(v)
		case float64:
			total += int(v)
		case map[string]interface{}:
			total += sum(v)
		case []interface{}:
			total += sumList(v)
		case string:
		default:
			fmt.Println(i, v, "is of a type I don't know how to handle")
		}
	}
	return total
}

func sum(m map[string]interface{}) int {
	total := 0
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			if vv == "red" {
				return 0
			}
		case int:
			total += vv
		case int64:
			total += int(vv)
		case float64:
			total += int(vv)
		case []interface{}:
			total += sumList(vv)
		case map[string]interface{}:
			total += sum(vv)
		default:
			fmt.Println(k, vv, "is of a type I don't know how to handle")
		}
	}
	return total
}
