package internal_test

import (
	"testing"

	"github.com/alvissraghnall/stewfish/internal"
)

func TestWeight_CreationAndRetrieval(t *testing.T) {
	w := internal.W(100, 200)

	if got := w.Mg(); got != 100 {
		t.Errorf("expected Mg() to be 100, got %d", got)
	}

	if got := w.Eg(); got != 200 {
		t.Errorf("expected Eg() to be 200, got %d", got)
	}
}

func TestWeight_Add(t *testing.T) {
	w1 := internal.W(100, 200)
	w2 := internal.W(50, 25)

	w1.Add(w2)
	if gotMg, gotEg := w1.Mg(), w1.Eg(); gotMg != 150 || gotEg != 225 {
		t.Errorf("expected Add to result in W(150, 225), got W(%d, %d)", gotMg, gotEg)
	}
}

func TestWeight_Sub(t *testing.T) {
	w1 := internal.W(100, 200)
	w2 := internal.W(50, 25)

	w1.Sub(w2)
	if gotMg, gotEg := w1.Mg(), w1.Eg(); gotMg != 50 || gotEg != 175 {
		t.Errorf("expected Sub to result in W(50, 175), got W(%d, %d)", gotMg, gotEg)
	}
}

func TestWeight_NegativeValues(t *testing.T) {
	w1 := internal.W(-30, -20)
	w2 := internal.W(10, -5)

	w1.Add(w2)
	if gotMg, gotEg := w1.Mg(), w1.Eg(); gotMg != -20 || gotEg != -25 {
		t.Errorf("expected Add to result in W(-20, -25), got W(%d, %d)", gotMg, gotEg)
	}

	w1.Sub(w2)
	if gotMg, gotEg := w1.Mg(), w1.Eg(); gotMg != -30 || gotEg != -20 {
		t.Errorf("expected Sub to result in W(-30, -20), got W(%d, %d)", gotMg, gotEg)
	}
}
