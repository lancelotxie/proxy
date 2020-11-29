package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/pkg/errors"

	"github.com/lancelotXie/proxy/proxy.lib/grpc.client/configuration"
)

// Method 方法
type Method string

const (
	// Set 设置 key/value
	Set Method = "set"
	// Get 获取 value
	Get Method = "get"
	// Save 保存
	Save Method = "save"
)

var (
	// ErrNoMethod 未找到方法
	ErrNoMethod = errors.New("no method found")
	// ErrInvalidMethod 非法的方法
	ErrInvalidMethod = errors.New("invalid method")
	// ErrMissingArgs 缺少参数
	ErrMissingArgs = errors.New("missing arg(s)")
)

func main() {
	port := flag.Int("ctrl-port", 8085, "port for controller server")
	flag.Parse()

	err := configuration.Init("127.0.0.1", *port)
	if err != nil {
		log.Println(err)
		return
	}

	args := flag.Args()

	m, err := getMethod(args)
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	err = handleMethod(ctx, m, args[1:])
	if err != nil {
		log.Println(err)
	}
}

func getMethod(args []string) (m Method, err error) {
	if len(args) < 1 {
		err = errors.WithStack(ErrNoMethod)
		return
	}

	m = Method(args[0])
	return
}

func handleMethod(ctx context.Context, m Method, args []string) (err error) {
	switch m {
	case Get:
		err = handleGet(ctx, args)
	case Set:
		err = handleSet(ctx, args)
	case Save:
		err = handleSave(ctx)
	default:
		err = errors.WithStack(ErrInvalidMethod)
		err = errors.WithMessage(err, string(m))
	}
	return
}

func handleGet(ctx context.Context, args []string) (err error) {
	if len(args) < 1 {
		err = errors.WithStack(ErrMissingArgs)
		return
	}

	key := args[0]
	value, err := get(ctx, key)
	if err != nil {
		return
	}

	fmt.Println(key, ": ", value)
	return
}

func get(ctx context.Context, key string) (value interface{}, err error) {
	return configuration.Get(ctx, key)
}

func handleSet(ctx context.Context, args []string) (err error) {
	if len(args) < 2 {
		err = errors.WithStack(ErrMissingArgs)
		return
	}

	key := args[0]
	_value := args[1]
	var value interface{}

	value, errConvert := toFloat(_value)
	if errConvert != nil {
		value, errConvert = toBool(_value)
	}
	if errConvert != nil {
		value = _value
	}

	err = set(ctx, key, value)
	return
}

func set(ctx context.Context, key string, value interface{}) (err error) {
	return configuration.Set(ctx, key, value)
}

func toBool(src string) (dst bool, err error) {
	return strconv.ParseBool(src)
}

func toFloat(src string) (dst float64, err error) {
	return strconv.ParseFloat(src, 64)
}

func handleSave(ctx context.Context) error {
	return save(ctx)
}

func save(ctx context.Context) error {
	return configuration.Save(ctx)
}
