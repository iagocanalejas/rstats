package assert

import (
	"log/slog"

	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
)

func Assert(condition bool, msg string, data ...any) {
	if !condition {
		prettylog.Error("Assert#condition not met")
		prettylog.Fatal(msg, data...)
	}
}

func Contains[T comparable](item T, list []T, msg string, data ...any) {
	for _, l := range list {
		if l == item {
			return
		}
	}

	prettylog.Error("Contains#item not found")
	prettylog.Fatal(msg, data...)
}

func Nil(item any, msg string, data ...any) {
	if item != nil {
		prettylog.Error("Nil#not nil encountered")
		prettylog.Fatal(msg, data...)
	}
}

func NotNil(item any, msg string, data ...any) {
	if item == nil {
		slog.Error("NotNil#nil encountered")
		prettylog.Error("NotNil#nil encountered")
		prettylog.Fatal(msg, data...)
	}
}

func NoError(err error, msg string, data ...any) {
	if err != nil {
		prettylog.Error("NoError#error encountered")
		prettylog.Fatal(msg, data...)
	}
}
