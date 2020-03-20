/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 16:46:48
 * @LastEditTime: 2019-12-07 16:49:19
 */

package renewconn

import "net"

// Renew 翻新连接
func Renew(ori net.Conn, head []byte) (renew net.Conn) {
	return newRenewConn(ori, head)
}
