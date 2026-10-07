package main

type suggestCase struct {
	name   string
	prefix string
	k      int
	want   []string
}

var suggestCases = []suggestCase{
	// basics
	{"basic prefix", "ber", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
	{"uppercase prefix", "BER", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
	{"mixed case prefix", "bErLiN", 10, []string{"Berlin"}},
	{"leading whitespace trim", "  ber", 10, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
	{"trailing space is kept", "new ", 10, []string{"New York", "New Delhi"}},
	{"full name as prefix", "berlin", 10, []string{"Berlin"}},
	{"empty prefix returns nothing", "", 10, []string{}},

	// no results
	{"no match", "zzz", 10, []string{}},
	{"prefix longer than any name", "berlinnn", 10, []string{}},
	{"substring is not a prefix", "erg", 10, []string{}},

	// k handling
	{"k zero", "ber", 0, []string{}},
	{"k negative", "ber", -5, []string{}},
	{"k one", "ber", 1, []string{"Berlin"}},
	{"k smaller than matches", "a", 3, []string{"Amsterdam", "Athens", "Augsburg"}},
	{"k equals matches", "ber", 4, []string{"Berlin", "Bergen", "Bern", "Bergamo"}},
	{"k larger than matches", "a", 100, []string{"Amsterdam", "Athens", "Augsburg", "Aachen", "Amberg"}},

	// ordering and ties
	{"tie broken by name", "ba", 10, []string{"Bamberg", "Bayreuth"}},
	{"several names sharing a prefix", "new", 10, []string{"New York", "Newark", "New Delhi"}},
	{"zero population still returned", "tiny", 10, []string{"Tinyville"}},

	// non-ASCII
	{"umlaut prefix", "mü", 10, []string{"München", "Münster", "Mülheim"}},
	{"uppercase umlaut prefix", "MÜ", 10, []string{"München", "Münster", "Mülheim"}},

	// space handling
	{"prefix with inner space", "new y", 10, []string{"New York"}},
}
