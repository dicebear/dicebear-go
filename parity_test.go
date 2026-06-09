package dicebear

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dicebear/dicebear-go/v10/internal/initials"
	"github.com/dicebear/dicebear-go/v10/internal/num"
	"github.com/dicebear/dicebear-go/v10/internal/prng"
)

// Cross-language parity: assert the primitives produce exactly the values in the
// shared fixtures under <repo>/tests/fixtures/parity/, the same the JS, PHP,
// Python and Rust suites run against. The tests skip gracefully when the
// fixtures are absent (e.g. in the split dicebear-go repo).

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join("..", "..", "..", "tests", "fixtures", "parity", name)
	if _, err := os.Stat(p); err != nil {
		t.Skipf("parity fixtures not available (%s)", p)
	}
	return p
}

func readFixture(t *testing.T, name string, v any) {
	t.Helper()
	b, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
}

func strID(s string) string { return s }

func TestFnv1aParity(t *testing.T) {
	var cases []struct {
		Input string `json:"input"`
		Hash  uint32 `json:"hash"`
		Hex   string `json:"hex"`
	}
	readFixture(t, "fnv1a.json", &cases)

	for _, c := range cases {
		if got := prng.Hash(c.Input); got != c.Hash {
			t.Errorf("Hash(%q) = %d, want %d", c.Input, got, c.Hash)
		}
		if got := prng.Hex(c.Input); got != c.Hex {
			t.Errorf("Hex(%q) = %s, want %s", c.Input, got, c.Hex)
		}
	}
}

func TestMulberry32Parity(t *testing.T) {
	var cases []struct {
		Seed     uint32 `json:"seed"`
		Sequence []struct {
			Float float64 `json:"float"`
			State int32   `json:"state"`
		} `json:"sequence"`
	}
	readFixture(t, "mulberry32.json", &cases)

	for _, c := range cases {
		m := prng.NewMulberry32(c.Seed)
		for i, step := range c.Sequence {
			f := m.NextFloat()
			s := m.SignedState()
			if f != step.Float {
				t.Errorf("seed %d step %d: float = %v, want %v", c.Seed, i, f, step.Float)
			}
			if s != step.State {
				t.Errorf("seed %d step %d: state = %d, want %d", c.Seed, i, s, step.State)
			}
		}
	}
}

type rangeFixture struct {
	Min  float64  `json:"min"`
	Max  float64  `json:"max"`
	Step *float64 `json:"step"`
}

func (r rangeFixture) toRange() *prng.Range {
	return &prng.Range{Min: r.Min, Max: r.Max, Step: r.Step}
}

func TestPrngParity(t *testing.T) {
	var f struct {
		GetValue []struct {
			Seed, Key string
			Result    float64
		} `json:"getValue"`
		Pick []struct {
			Seed, Key string
			Items     []string
			Result    *string
		} `json:"pick"`
		WeightedPick []struct {
			Seed, Key string
			Weights   map[string]float64
			Result    *string
		} `json:"weightedPick"`
		Bool []struct {
			Seed, Key  string
			Likelihood float64
			Result     bool
		} `json:"bool"`
		Float []struct {
			Seed, Key string
			Range     rangeFixture
			Result    float64
		} `json:"float"`
		Integer []struct {
			Seed, Key string
			Range     rangeFixture
			Result    int
		} `json:"integer"`
		Shuffle []struct {
			Seed, Key string
			Items     []string
			Result    []string
		} `json:"shuffle"`
	}
	readFixture(t, "prng.json", &f)

	for _, c := range f.GetValue {
		if got := prng.New(c.Seed).GetValue(c.Key); got != c.Result {
			t.Errorf("GetValue(%q, %q) = %v, want %v", c.Seed, c.Key, got, c.Result)
		}
	}

	for _, c := range f.Pick {
		got, ok := prng.Pick(prng.New(c.Seed), c.Key, c.Items, strID)
		if !matchOptional(got, ok, c.Result) {
			t.Errorf("Pick(%q, %q, %v) = (%q, %v), want %v", c.Seed, c.Key, c.Items, got, ok, c.Result)
		}
	}

	for _, c := range f.WeightedPick {
		got, ok := prng.New(c.Seed).WeightedPick(c.Key, c.Weights)
		if !matchOptional(got, ok, c.Result) {
			t.Errorf("WeightedPick(%q, %q, %v) = (%q, %v), want %v", c.Seed, c.Key, c.Weights, got, ok, c.Result)
		}
	}

	for _, c := range f.Bool {
		if got := prng.New(c.Seed).Bool(c.Key, c.Likelihood); got != c.Result {
			t.Errorf("Bool(%q, %q, %v) = %v, want %v", c.Seed, c.Key, c.Likelihood, got, c.Result)
		}
	}

	for _, c := range f.Float {
		if got := prng.New(c.Seed).Float(c.Key, c.Range.toRange()); got != c.Result {
			t.Errorf("Float(%q, %q) = %v, want %v", c.Seed, c.Key, got, c.Result)
		}
	}

	for _, c := range f.Integer {
		if got := prng.New(c.Seed).Integer(c.Key, c.Range.toRange()); got != c.Result {
			t.Errorf("Integer(%q, %q) = %d, want %d", c.Seed, c.Key, got, c.Result)
		}
	}

	for _, c := range f.Shuffle {
		got := prng.New(c.Seed).Shuffle(c.Key, c.Items)
		if !equalStrings(got, c.Result) {
			t.Errorf("Shuffle(%q, %q, %v) = %v, want %v", c.Seed, c.Key, c.Items, got, c.Result)
		}
	}
}

func TestNumberParity(t *testing.T) {
	var cases []struct {
		Input  float64 `json:"input"`
		Output string  `json:"output"`
	}
	readFixture(t, "numbers.json", &cases)

	for _, c := range cases {
		if got := num.Format(c.Input); got != c.Output {
			t.Errorf("num.Format(%v) = %q, want %q", c.Input, got, c.Output)
		}
	}
}

func TestInitialsParity(t *testing.T) {
	var cases []struct {
		Seed   string `json:"seed"`
		Result string `json:"result"`
	}
	readFixture(t, "initials.json", &cases)

	for _, c := range cases {
		if got := initials.FromSeed(c.Seed); got != c.Result {
			t.Errorf("initials.FromSeed(%q) = %q, want %q", c.Seed, got, c.Result)
		}
	}
}

func matchOptional(got string, ok bool, want *string) bool {
	if want == nil {
		return !ok
	}
	return ok && got == *want
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
