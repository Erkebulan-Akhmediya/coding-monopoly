// Package locale holds shared multi-language text helpers for content fields.
package locale

// Text is a triple of English, Russian, and Kazakh content strings.
type Text struct {
	En string `json:"en"`
	Ru string `json:"ru"`
	Kz string `json:"kz"`
}

// New returns a Text with the given language values.
func New(en, ru, kz string) Text {
	return Text{En: en, Ru: ru, Kz: kz}
}

// FromEn builds a Text that uses the same string for all languages.
func FromEn(en string) Text {
	return Text{En: en, Ru: en, Kz: en}
}

// All returns non-empty unique variants for grading / matching.
func (t Text) All() []string {
	seen := make(map[string]struct{}, 3)
	out := make([]string, 0, 3)
	for _, value := range []string{t.En, t.Ru, t.Kz} {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
