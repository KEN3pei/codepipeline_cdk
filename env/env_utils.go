package env

import (
	"os"
	"strings"

	"github.com/aws/jsii-runtime-go"
)

// ENV値が未設定ならnilを返す
func GetNilOrStrEnv(s string) *string {
	return nilGetEnv(os.Getenv(s))
}

// 空文字ならnilを返す
func nilGetEnv(s string) *string {
	env := jsii.String(os.Getenv(s))
	if *env == "" {
		return nil
	}
	return env
}

// []*stringのスライスをENVから作成する
func GetStringsEnv(s string) *[]*string {
	return toPtrSlice(os.Getenv(s))
}

func toPtrSlice(s string) *[]*string {
	arr := strings.Split(s, ",")
	var ptrSlice []*string

	// 各要素をポインタに変換してptrSliceに追加する
	for _, s := range arr {
		// 関数を経由することで異なるポインタ値を返す
		ptrSlice = append(ptrSlice, StrToPtr(s))
	}

	return &ptrSlice
}

// Goでは関数の引数の引き渡しの際に値渡しが行われるため、関数経由でポインタ値を返すことで異なる値を返すことが可能
func StrToPtr(v string) *string {
	return &v
}
