package macros

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// ApplyMacros applies macros to the input string based on the provided query and plugin context.
// It replaces the following macros in the input string:
//   - ${__from}: Replaced with the start time of the query time range.
//   - ${__timeFrom}: Replaced with the start time of the query time range. ( alias for from macro. also respect local timeShift defined in the panels )
//   - ${__to}: Replaced with the end time of the query time range.
//   - ${__timeTo}: Replaced with the end time of the query time range. ( alias for to macro. also respect local timeShift defined in the panels )
//   - ${__user.name}: Replaced with the name of the plugin context user.
//   - ${__user.email}: Replaced with the email of the plugin context user.
//   - ${__user.login}: Replaced with the login of the plugin context user.
//
// Parameters:
//   - input: The input string to apply macros to.
//   - query: The data query containing the time range.
//   - pluginCtx: The plugin context containing the user information.
//
// Returns:
//   - The input string with macros replaced.
//   - An error if there was an error applying the macros.
func ApplyMacros(input string, query backend.DataQuery, pluginCtx backend.PluginContext) (string, error) {
	timeRange := query.TimeRange
	var err error
	timeMacros := []func(input string, timeRange backend.TimeRange) (string, error){
		from,     // ${__from}
		timeFrom, // ${__timeFrom}
		to,       // ${__to}
		timeTo,   // ${__timeTo}
	}
	for _, f := range timeMacros {
		input, err = f(input, timeRange)
		if err != nil {
			return input, err
		}
	}
	if pluginCtx.User != nil {
		input = strings.ReplaceAll(input, "${__user.name}", pluginCtx.User.Name)
		input = strings.ReplaceAll(input, "${__user.email}", pluginCtx.User.Email)
		input = strings.ReplaceAll(input, "${__user.login}", pluginCtx.User.Login)
	}
	return input, nil
}

type macroFunc func(string, []string) (string, error)

func getMatches(macroName, input string) ([][]string, error) {
	macroRegex := fmt.Sprintf("\\$__%s\\b(?:\\((.*?)\\))?", macroName)
	if strings.HasPrefix(macroName, "$$") { // prefix $$ is used to denote macro from frontend or grafana global variable
		macroRegex = fmt.Sprintf("\\${__%s:?(.*?)}", strings.TrimPrefix(macroName, "$$"))
	}
	rgx, err := regexp.Compile(macroRegex)
	if err != nil {
		return nil, err
	}
	return rgx.FindAllStringSubmatch(input, -1), nil
}

func applyMacro(macroKey string, queryString string, macro macroFunc) (string, error) {
	matches, err := getMatches(macroKey, queryString)
	if err != nil {
		return queryString, err
	}
	for _, match := range matches {
		if len(match) == 0 {
			continue
		}
		args := []string{}
		if len(match) > 1 {
			args = strings.Split(match[1], ",")
		}
		res, err := macro(queryString, args)
		if err != nil {
			return queryString, err
		}
		queryString = strings.Replace(queryString, match[0], res, -1)
	}
	return strings.TrimSpace(queryString), nil
}
