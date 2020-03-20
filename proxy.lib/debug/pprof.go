/*
 * @Author: EagleXiang
 * @LastEditors  : EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-19 20:55:32
 * @LastEditTime : 2019-12-19 21:02:55
 */

package debug

import (
	"log"
	"net/http"
	_ "net/http/pprof" // 开启 pprof
)

// TryRun 尝试运行 debug 相关程序
func TryRun(sig string) {
	if sig != "on" {
		return
	}

	log.Println("开启 pprof 服务")
	go http.ListenAndServe("0.0.0.0:6060", nil)
}
