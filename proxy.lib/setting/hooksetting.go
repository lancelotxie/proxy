package setting

// HookFunc0 第一步钩子函数
type HookFunc0 func(value interface{}) error

// HookFunc1 第二步钩子函数，ok 参数表示第一步是否全部成功
type HookFunc1 func(value interface{}, ok bool)

// Hooker 钩子
type Hooker struct {
	H0 HookFunc0
	H1 HookFunc1
}

type hookSetting struct {
	kindBindSetting
	hooks map[string][]Hooker
}

func newhookSetting(base *kindBindSetting) (s *hookSetting) {
	s = new(hookSetting)
	s.kindBindSetting = *base
	s.hooks = make(map[string][]Hooker)
	return
}

// Hook 注册钩子
func (s *hookSetting) Hook(key string, h Hooker) {
	s.hooks[key] = append(s.hooks[key], h)
}

// Set 设置 key/value，设置前会调用 Hooker 的检查方法，设置成功后会调用钩子方法
func (s *hookSetting) Set(key string, value interface{}) (err error) {
	hs := s.hooks[key]

	// 执行第一轮钩子函数
	for _, h := range hs {
		err = h.H0(value)
		if err != nil {
			break
		}
	}

	// 判断第一轮钩子函数是否全部成功
	ok := err == nil
	if !ok {
		return
	}

	// 全部成功则执行第二轮钩子函数
	for _, h := range hs {
		if h.H1 != nil {
			h.H1(value, ok)
		}
	}

	err = s.kindBindSetting.Set(key, value)
	return
}

func (s *hookSetting) MarshalJSON() (b []byte, err error) {
	b, err = s.kindBindSetting.MarshalJSON()
	return
}

func (s *hookSetting) UnmarshalJSON(b []byte) (err error) {
	err = s.kindBindSetting.UnmarshalJSON(b)
	return
}
