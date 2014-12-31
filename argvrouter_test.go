package argvrouter

import (
	"testing"
)

type RouteTestCase struct {
	route       *Route
	args        []string
	expectation bool
}

func TestIsPatternParam(t *testing.T) {
	if IsPatternParam("") {
		t.Errorf("IsPatternParam(\"\") should have been false!")
	}

	if IsPatternParam("foo") {
		t.Errorf("IsPatternParam(\"foo\") should have been false!")
	}

	if !IsPatternParam(":foo") {
		t.Errorf("IsPatternParam(\":foo\") should have been true!")
	}
}

func TestRouteMatches(t *testing.T) {
	cases := []RouteTestCase{
		{&Route{Pattern: []string{"ls"}}, []string{}, false},
		{&Route{Pattern: []string{"ls"}}, []string{""}, false},
		{&Route{Pattern: []string{"ls"}}, []string{"cd"}, false},
		{&Route{Pattern: []string{"ls"}}, []string{"ls"}, true},
	}

	for _, test_case := range cases {
		_, matched := RouteMatches(test_case.route, test_case.args)
		if matched != test_case.expectation {
			t.Errorf("RouteMatches(%q,%q) == %q expected %q", test_case.route, test_case.args, matched, test_case.expectation)
		}
	}

	matched_route, matched := RouteMatches(
		&Route{Pattern: []string{"ls", ":fname"}},
		[]string{"ls", "foo.txt"},
	)

	if !matched {
		t.Errorf("Expected (ls :fname), to match (ls foo.txt)!")
	}

	if matched_route.Params["fname"] != "foo.txt" {
		t.Errorf("Expected :fname to be bound to \"foo.txt\", it was: %q", matched_route.Params["fname"])
	}

	args := []string{"ls", "foo.txt"}
	route := &Route{Pattern: []string{"ls"}}
	matched_route, matched = RouteMatches(route, args)

	if matched {
		t.Errorf("Error: %q should not have been a match for route %q", args, route)
	}

}

func TestSplatRouteMatching(t *testing.T) {
	args := []string{"ls", "this.txt", "that.txt", "other.txt"}
	route := &Route{Pattern: []string{"ls", "*"}}
	_, matched := RouteMatches(route, args)
	if !matched {
		t.Errorf("Error: %q should have been a (splat) match for route %q", args, route)
	}
}

func TestRoutingTable(t *testing.T) {
	ClearRoutingTable()
	AddRoute(&Route{Pattern: []string{"ls"}})
	if len(RoutingTable) != 1 {
		t.Errorf("Expected the routing table to have 1 entry, it had %d", len(RoutingTable))
	}

	if RoutingTable[0].Pattern[0] != "ls" {
		t.Errorf("Expected the first routing table to be 'ls', it was %q", RoutingTable[0].Pattern)
	}
}

// TODO: FindMatchingRoute(args []string) *Route

func TestFindMatchingRoute(t *testing.T) {
	ClearRoutingTable()
	AddRoute(&Route{Pattern: []string{"ls"}})
	AddRoute(&Route{Pattern: []string{"ls", ":fname"}})

	if nil == FindMatchingRoute([]string{"ls"}) {
		t.Errorf("Expected to find a route matching \"ls\"")
	}

	found := FindMatchingRoute([]string{"ls", "foo.txt", "extra-arg"})
	if nil != found {
		t.Errorf("Expected to not find a route matching rotue for [ls,foo.txt,extra-arg] got: %q", found)
	}
}
