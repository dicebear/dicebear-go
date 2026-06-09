package dicebear

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Cross-language avatar parity. Renders each shared fixture case and asserts the
// SVG matches byte-for-byte the output committed under
// <repo>/tests/fixtures/parity/avatars/, the same fixtures the JS, PHP, Python
// and Rust suites render against.

func TestAvatarParity(t *testing.T) {
	styleNames := []string{"initials", "thumbs", "glass", "shape-grid", "notionists"}

	for _, name := range styleNames {
		stylePath := fixturePath(t, filepath.Join("styles", name+".json"))
		styleJSON, err := os.ReadFile(stylePath)
		if err != nil {
			t.Fatalf("read style %s: %v", name, err)
		}

		style, err := NewStyle(styleJSON)
		if err != nil {
			t.Fatalf("parse style %s: %v", name, err)
		}

		var cases []struct {
			ID              string          `json:"id"`
			Options         map[string]any  `json:"options"`
			SVG             string          `json:"svg"`
			ResolvedOptions json.RawMessage `json:"resolvedOptions"`
		}
		readFixture(t, filepath.Join("avatars", name+".json"), &cases)

		for _, c := range cases {
			avatar, err := NewAvatar(style, c.Options)
			if err != nil {
				t.Fatalf("render %s/%s: %v", name, c.ID, err)
				continue
			}

			if got := avatar.SVG(); got != c.SVG {
				t.Errorf("%s/%s SVG mismatch at byte %d", name, c.ID, firstDiff(got, c.SVG))
			}

			// Deep-equal (order-independent), like the JS/PHP/Python/Rust
			// suites. Decode with UseNumber so a JSON integer (1) and float
			// (1.0) compare unequal — this pins whole-number options as JSON
			// integers, the same guarantee the Rust suite makes.
			gotJSON, err := json.Marshal(avatar.ResolvedOptions())
			if err != nil {
				t.Fatalf("marshal resolved options %s/%s: %v", name, c.ID, err)
			}
			got := decodeNumberAware(t, gotJSON)
			want := decodeNumberAware(t, c.ResolvedOptions)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s/%s resolved options mismatch\n got: %s\nwant: %s", name, c.ID, gotJSON, c.ResolvedOptions)
			}
		}
	}
}

// decodeNumberAware unmarshals JSON into `any` with UseNumber, so numbers stay
// as json.Number (distinguishing 1 from 1.0) instead of collapsing to float64.
func decodeNumberAware(t *testing.T, b []byte) any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return v
}

// firstDiff returns the index of the first byte where a and b differ, or the
// length of the shorter string if one is a prefix of the other.
func firstDiff(a, b string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}
