package env

func LoadEnv() error {
	if err := LoadAppEnv(); err != nil {
		return err
	}
	if err := LoadDiscordEnv(); err != nil {
		return err
	}
	if err := LoadDBEnv(); err != nil {
		return err
	}
	if err := LoadRedisEnv(); err != nil {
		return err
	}

	return nil
}
