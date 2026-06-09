// Package color holds the color helpers used by the DiceBear renderer and the
// option resolver: hex normalization, WCAG relative luminance, contrast sorting,
// and exclusion filtering. They are exported so consumers can reproduce the
// engine's color math, mirroring the `Color` utility in the JS core and the
// `color` module in the Rust core.
package color

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// ToHex normalizes any hex format to 6- or 8-digit lowercase with a "#" prefix.
func ToHex(hex string) string {
	h := strings.ToLower(strings.TrimPrefix(hex, "#"))
	c := []rune(h)

	switch len(c) {
	case 3:
		return "#" + dup(c[0]) + dup(c[1]) + dup(c[2])
	case 4:
		return "#" + dup(c[0]) + dup(c[1]) + dup(c[2]) + dup(c[3])
	default:
		return "#" + h
	}
}

// ToRGBHex is like ToHex, but strips the alpha channel and always returns
// 6-digit hex.
func ToRGBHex(hex string) string {
	h := ToHex(hex)

	// JS slice(0, 7) counts UTF-16 code units and never splits a char; using
	// runes keeps a multi-byte char (reachable via non-hex input) from being
	// split mid-byte.
	r := []rune(h)
	if len(r) > 7 {
		return string(r[:7])
	}
	return h
}

// ParseHex parses a hex color into an (r, g, b) triple of 8-bit channel values.
// Malformed channels fall back to 0.
func ParseHex(hex string) (uint8, uint8, uint8) {
	h := []rune(ToHex(hex))[1:]
	digits := string(h)

	return channel(digits, 0), channel(digits, 2), channel(digits, 4)
}

// Luminance returns the WCAG 2.1 relative luminance with sRGB linearization.
//
// See https://www.w3.org/WAI/GL/wiki/Relative_luminance
func Luminance(hex string) float64 {
	r, g, b := ParseHex(hex)

	return 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
}

// SortByContrast returns a new slice sorted by descending contrast against
// refColor. The sort is stable, matching the JS `b.ratio - a.ratio` sort.
//
// See https://www.w3.org/WAI/GL/wiki/Contrast_ratio
func SortByContrast(candidates []string, refColor string) []string {
	refLum := Luminance(refColor)

	type entry struct {
		color string
		ratio float64
	}

	withRatio := make([]entry, len(candidates))
	for i, c := range candidates {
		lum := Luminance(c)
		ratio := (math.Max(lum, refLum) + 0.05) / (math.Min(lum, refLum) + 0.05)
		withRatio[i] = entry{color: c, ratio: ratio}
	}

	sort.SliceStable(withRatio, func(a, b int) bool {
		return withRatio[a].ratio > withRatio[b].ratio
	})

	out := make([]string, len(withRatio))
	for i, e := range withRatio {
		out[i] = e.color
	}
	return out
}

// FilterNotEqualTo returns a new slice with excluded colors removed, falling
// back to the original candidates when filtering would empty the list.
func FilterNotEqualTo(candidates, excluded []string) []string {
	normalized := make(map[string]struct{}, len(excluded))
	for _, c := range excluded {
		normalized[ToRGBHex(c)] = struct{}{}
	}

	filtered := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if _, ok := normalized[ToRGBHex(c)]; !ok {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		out := make([]string, len(candidates))
		copy(out, candidates)
		return out
	}

	return filtered
}

// dup returns the rune repeated twice as a string.
func dup(r rune) string {
	return string([]rune{r, r})
}

// channel parses two hex digits at start into a channel value, falling back to 0.
func channel(digits string, start int) uint8 {
	if start+2 > len(digits) {
		return 0
	}
	v, err := strconv.ParseUint(digits[start:start+2], 16, 8)
	if err != nil {
		return 0
	}
	return uint8(v)
}

// linearize converts an 8-bit sRGB channel value into linear-light space.
func linearize(c uint8) float64 {
	s := float64(c) / 255.0

	if s <= 0.04045 {
		return s / 12.92
	}

	return math.Pow((s+0.055)/1.055, 2.4)
}
