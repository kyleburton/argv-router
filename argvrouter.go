// This package implements 'routing table' like functionality for command line programs.
package argvrouter

import (
	"strings"
)

// The handler function for a route.
type RouteHandler func(*Route)

// A function to return potential completions for a partial route match.
type RouteParameterCompletions func(route *Route, param, word string) []string

type Route struct {
	Pattern       []string
	Params        map[string]string
	Args          []string
	Handler       RouteHandler
	HelpText      *string
	CompletionsFn RouteParameterCompletions
}

var RoutingTable []*Route

func ClearRoutingTable() {
	RoutingTable = make([]*Route, 0)
}

func AddRoute(route *Route) {
	RoutingTable = append(RoutingTable, route)
}

func IsPatternParam(s string) bool {
	return strings.HasPrefix(s, ":")
}

func (self *Route) PatternEndsWithSplat() bool {
	return "*" == self.Pattern[len(self.Pattern)-1]
}

func RouteMatches(route *Route, args []string) (*Route, bool) {

	if !route.PatternEndsWithSplat() {
		if len(args) != len(route.Pattern) {
			return nil, false
		}
	}

	var res *Route = &Route{
		Pattern:       route.Pattern,
		Params:        make(map[string]string),
		Handler:       route.Handler,
		CompletionsFn: route.CompletionsFn,
	}

	for idx, part := range route.Pattern {
		var arg string
		if len(args) > idx {
			arg = args[idx]
		}
		res.Args = args[idx:]

		if "*" == part {
			return res, true
		}

		if IsPatternParam(part) {
			res.Params[part[1:]] = arg
			continue
		}

		if part == arg {
			continue
		}

		return nil, false
	}

	return res, true
}

func FindMatchingRoute(args []string) *Route {
	for _, route := range RoutingTable {
		res, matched := RouteMatches(route, args)
		if matched {
			return res
		}
	}

	return nil
}
