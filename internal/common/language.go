package common

var language = "pt"

// SetLanguage define o idioma atual do CLI
func SetLanguage(lang string) {
	if lang != "" {
		language = lang
	}
}

// Lang retorna o idioma atual
func Lang() string { return language }
