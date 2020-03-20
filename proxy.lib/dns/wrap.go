/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-06 22:18:28
 * @LastEditTime: 2019-12-17 23:03:57
 */

package dns

// New 构造新的 DNS 解析器
func New(resolv ResolvFunc) Resolver {
	return newDNSResolver(resolv)
}

var defaultResolver = New(defaultResolv)

// Resolv 解析 DNS
func Resolv(domain string) (ip string, err error) {
	return defaultResolver.Resolv(domain)
}
