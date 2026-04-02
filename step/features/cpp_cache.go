package features


const CPPCacheEnabledMsg = "C++ cache feature is enabled"

type CPPCacheInput struct {
	CPPCacheEnabled bool `env:"cpp_cache_enabled,required"`
}

type CPPCache struct{}

func CPPCacheFeature(inputParser InputParser, logger Logger) *CPPCache {
	var input CPPCacheInput
	if err := inputParser.Parse(&input); err != nil {
		return nil
	}

	if !input.CPPCacheEnabled {
		return nil
	}

	logger.Debugf(CPPCacheEnabledMsg)
	return &CPPCache{}
}

func (c *CPPCache) ActivateArgs() []string {
	return []string{"activate", "c++"}
}

func (c *CPPCache) ServiceArgs() []string {
	return []string{"ccache", "storage-helper", "start"}
}
