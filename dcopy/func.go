/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: func.go
 * @time: 2020/5/23 4:14 下午
 * @project: deepcopy
 */

package dcopy

import (
	"strings"
)


func littleCamelCase(str string) string {
	if len(str) < 1 {
		return str
	}
	first := string(str[0])
	tail := str[1:]
	return strings.ToLower(first) + tail
}