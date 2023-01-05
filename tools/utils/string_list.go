package utils

import (
	"fmt"
)

func ConvertIntListToStringList(intList []int) []string {
	var stringList []string
	for _, string_ := range intList {
		stringList = append(stringList, fmt.Sprintf("%d", string_))
	}
	return stringList
}

func DifferTwoStringList(olds []string, news []string) (added []string, deleted []string) {
	for _, old := range olds {
		find := false
		for _, new_ := range news {
			if old == new_ {
				find = true
				break
			}
		}
		if !find {
			deleted = append(deleted, old)
		}
	}
	for _, new_ := range news {
		find := false
		for _, old := range olds {
			if old == new_ {
				find = true
				break
			}
		}
		if !find {
			added = append(added, new_)
		}
	}
	return
}

func IsStrInList(str string, list []string) bool {
	for _, str_in := range list {
		if str_in == str {
			return true
		}
	}
	return false
}
