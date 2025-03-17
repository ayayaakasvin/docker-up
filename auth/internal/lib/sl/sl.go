package sl

import "log/slog"

func Err (err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Any (key string, value interface{}) slog.Attr {
	return slog.Attr{
		Key: key,
		Value: slog.AnyValue(value),
	}
}