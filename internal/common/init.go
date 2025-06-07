package common

func init() {
	cfg := LoadConfig()
	SetLanguage(cfg.Language)
}
