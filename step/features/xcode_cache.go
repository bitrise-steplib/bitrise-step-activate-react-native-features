package features


const XcodeCacheEnabledMsg = "Xcode cache feature is enabled"

type XcodeCacheInput struct {
	XcodeCacheEnabled bool `env:"xcode_cache_enabled,required"`
}

type XcodeCache struct{}

func XcodeCacheFeature(inputParser InputParser, logger Logger) *XcodeCache {
	var input XcodeCacheInput
	if err := inputParser.Parse(&input); err != nil {
		return nil
	}

	if !input.XcodeCacheEnabled {
		return nil
	}

	logger.Debugf(XcodeCacheEnabledMsg)
	return &XcodeCache{}
}

func (x *XcodeCache) ActivateArgs() []string {
	return []string{"activate", "xcode"}
}

func (x *XcodeCache) ServiceArgs() []string {
	return []string{"xcelerate", "start-proxy"}
}
