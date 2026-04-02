package features


const GradleCacheEnabledMsg = "Gradle cache feature is enabled"

type GradleCacheInput struct {
	GradleCacheEnabled bool `env:"gradle_cache_enabled,required"`
}

type GradleCache struct{}

func GradleCacheFeature(inputParser InputParser, logger Logger) *GradleCache {
	var input GradleCacheInput
	if err := inputParser.Parse(&input); err != nil {
		return nil
	}

	if !input.GradleCacheEnabled {
		return nil
	}

	logger.Debugf(GradleCacheEnabledMsg)
	return &GradleCache{}
}

func (g *GradleCache) ActivateArgs() []string {
	return []string{"activate", "gradle", "--analytics", "--cache"}
}

func (g *GradleCache) ServiceArgs() []string {
	return nil
}
