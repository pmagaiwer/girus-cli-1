package common

// T retorna a string traduzida de acordo com o idioma atual
func T(pt, es string) string {
	if Lang() == "es" {
		return es
	}
	return pt
}
