package underline

import (
	"errors"
	"fmt"
	"strings"

	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
)

// 用于解析在 Edgeless 插件包 v1 中 使用的 Underline 文件名格式。

var (
	ErrUnderlineParse = errors.New("error: underline/filename parse error")
)

type Underline interface {
	Name() string
	Version() u.Version
	Author() string
	Category() v.Option[string]
	String() string
}

type underline struct {
	name     string
	version  u.Version
	author   string
	category v.Option[string]
}

func (s *underline) Name() string {
	return s.name
}

func (s *underline) Version() u.Version {
	return s.version
}

func (s *underline) Author() string {
	return s.author
}

func (s *underline) Category() v.Option[string] {
	return s.category
}

func (s *underline) String() string {
	cate := s.Category().Or("[-]")
	return fmt.Sprintf("Underline{%s-%s-%s-%s}", s.name, s.version, s.author, cate)
}

func NewUnderline(id string, category v.Option[string]) v.Result[Underline, error] {
	s := strings.Split(id, "_")
	l := len(s)
	if l < 3 {
		return v.Err[Underline](ErrUnderlineParse)
	}

	if l == 3 {
		return v.Ok[Underline, error](&underline{s[0], u.NewVersion(s[1]), s[2], category})
	}

	return v.Ok[Underline, error](&underline{s[0], u.NewVersion(s[1]), s[2], v.Some(s[3])})
}
