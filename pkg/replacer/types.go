package replacer

type Options struct {
	BasePath      string
	FilePattern   string
	ShowReport    bool
	ScanRecursive bool
	TruncateFiles bool
	StopOnErrors  bool
	StopOnMissing bool
	StopOnEmpty   bool
	Strict        bool
}

func WithPipedOptions() *Options {
	return &Options{
		"",
		"",
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}
}

type Stats struct {
}
