/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 14:00:27
 * @LastEditTime: 2019-12-07 14:02:31
 */

package parsers

import (
	"context"
	"net"

	"github.com/pkg/errors"
)

type mockParser struct{}

func (p mockParser) Match(b []byte) error {
	if string(b) != "valid" {
		return errors.New("invalid")
	}
	return nil
}

func (p mockParser) Parse(ctx context.Context, conn net.Conn) (renew net.Conn, err error) {
	return
}

func (p mockParser) String() string {
	return "mock"
}
