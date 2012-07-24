package decisiondriver

import (
	"github.com/fluffle/sp0rkle/lib/util"
	"math/rand"
	"reflect"
	"testing"
)

// deterministic randomizer
var mytestrand *rand.Rand = util.NewRand(42)

func TestRand(t *testing.T) {
	tests := []string{
		"no plugin value",
		"this has a <plugin=rand 100>% chance of working",
		"<plugin=rand 40-50> should be between 40 and 50",
		"there's a <plugin=rand 60-100 %.3f%%> chance of accurate random",
		"<plugin=rand dongs> shouldn't work, but <plugin=rand 20> should",
	}
	expected := []string{
		"no plugin value",
		"this has a 37% chance of working",
		"41 should be between 40 and 50",
		"there's a 84.164% chance of accurate random",
		"0 shouldn't work, but 1 should",
	}
	for i, s := range tests {
		ret := rand_replacer(s, mytestrand)
		if ret != expected[i] {
			t.Errorf("Test: %s\nExpected: %s\nGot: %s\n", s, expected[i], ret)
		}
	}
}

func TestDecide(t *testing.T) {
	tests := []string{
		"<plugin=decide singlevalue>",
		"<plugin=decide AAA BBB CCC>",
		"<plugin=decide DDD | EEE>",
		"<plugin=decide FFF | GGG | HHH>",
		"<plugin=decide spam | spam and sausage | eggs | ham | spam eggs and spam>",
		"<plugin=decide \"spam\" \"spam and sausage\" \"eggs\" \"ham\" \"spam spam spam spam eggs and spam\">",
		"<plugin=decide \"cheese\" \"ham>",
		"<plugin=decide 'cheese' 'carrots' 'sausage'>",
		"<plugin=decide \"foo bar\" \"foo's bar\" \"something with spaces in it\">",
		"<plugin=decide \"foobar\" \"bar\" \"cheese\">",
	}
	expected := []string{
		"singlevalue", //if their is only one option, accept that
		"AAA",
		"DDD",
		"HHH",
		"ham",
		"spam spam spam spam eggs and spam",
		"Unbalanced quotes",
		"cheese",
		"foo's bar",
		"cheese",
	}
	for i, s := range tests {
		ret := rand_decider(s, mytestrand)
		if ret != expected[i] {
			t.Errorf("Test: %s\nExpected: [%s]\nGot: [%s]\n\n", s, expected[i], ret)
		}
	}
}

func TestChoices(t *testing.T) {
	tests := []string{
		"singlevalue",
		"AAA BBB CCC",
		"DDD | EEE",
		"FFF | GGG | HHH",
		"spam | spam and sausage | eggs | ham | spam eggs and spam",
		`"spam" "spam and sausage" "eggs" "ham" "spam spam spam spam eggs and spam"`,
		`"cheese" "ham`,
		`'cheese' 'carrots' 'sausage'`,
		`"foo bar" "foo's bar" "something with spaces in it"`,
		`"foobar" "bar" "cheese"`,
	}
	expected := [][]string{
		[]string{"singlevalue"},
		[]string{"AAA", "BBB", "CCC"},
		[]string{"DDD ", " EEE"},
		[]string{"FFF ", " GGG ", " HHH"},
		[]string{"spam ", " spam and sausage ", " eggs ", " ham ", " spam eggs and spam"},
		[]string{"spam", "spam and sausage", "eggs", "ham", "spam spam spam spam eggs and spam"},
		[]string{"Unbalanced quotes"},
		[]string{"cheese", "carrots", "sausage"},
		[]string{"foo bar", "foo's bar", "something with spaces in it"},
		[]string{"foobar", "bar", "cheese"},
	}
	for i, s := range tests {
		ret := choices(s)
		// We don't trim space from the possible choices *in* choices()
		// as it's only really necessary to do so for the one chosen.
		if !reflect.DeepEqual(expected[i], ret) {
			t.Errorf("Test: %s\nExpected: %#v\nGot: %#v\n\n", s, expected[i], ret)
		}
	}
}
