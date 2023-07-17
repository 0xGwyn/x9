package main

import (
	"bufio"
	"fmt"
	"math"
	"net/url"
	"os"
	"strings"

	"github.com/projectdiscovery/goflags"
	fileUtil "github.com/projectdiscovery/utils/file"
	sliceUtil "github.com/projectdiscovery/utils/slice"
)

type Options struct {
	list             string
	parameters       string
	chunk            int
	values           goflags.StringSlice
	generateStrategy goflags.StringSlice
	valueStrategy    string
	output           string
}

var options *Options

func main() {
	options = ParseOptions()
	urls := getUrls()
	params := getParams()

	var output []string
	if sliceUtil.Contains(options.generateStrategy, "normal") {
		output = append(output, normalStrat(urls, params)...)
	}
	if sliceUtil.Contains(options.generateStrategy, "combine") {
		output = append(output, combineStrat(urls)...)
	}
	if sliceUtil.Contains(options.generateStrategy, "ignore") {
		output = append(output, ignoreStrat(urls, params)...)
	}

	writeOutput(output)

}

func ParseOptions() *Options {
	options := &Options{}

	flags := goflags.NewFlagSet()
	flags.SetDescription("A tool designed for URL modification with specific modes to manipulate parameters and their values")

	flags.StringVarP(&options.list, "list", "l", "", "List of URLS to edit")
	flags.StringVarP(&options.parameters, "parameters", "p", "", "Paramter wordlist")
	flags.IntVarP(&options.chunk, "chunk", "c", 15, "Number of parameters in each URL")
	flags.StringSliceVarP(&options.values, "value", "v", nil, "Value for the parameters", goflags.StringSliceOptions)

	generateStrategyHelp := `
	Select the mode strategy from the available choices:
					normal:  Remove all parameters and put the wordlist
					combine: Pitchfork combine on the existing parameters
					ignore:  Don't touch the URL and append the paramters to the URL
				`
	flags.StringSliceVarP(&options.generateStrategy, "generate-strategy", "gs", nil, generateStrategyHelp, goflags.CommaSeparatedStringSliceOptions)

	valueStrategyHelp := `Select the strategy from the available choices:
					replace: Replace the current URL values with the given values
					suffix:  Append the value to the end of the parameters
				`
	flags.StringVarP(&options.valueStrategy, "value-strategy", "vs", "suffix", valueStrategyHelp)

	flags.StringVarP(&options.output, "output", "o", "", "File to write output results")

	if err := flags.Parse(); err != nil {
		panic(err)
	}

	return options
}

func (options *Options) validateOptions() error {
	// to be implemented
	return nil
}

func writeOutput(urls []string) {
	output := strings.Join(urls, "\n")
	// save to output file otherwise write to stdin
	if options.output != "" {
		outputFile, err := os.OpenFile(options.output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		fmt.Fprint(outputFile, output+"\n")

	} else {
		fmt.Fprint(os.Stdout, output+"\n")
	}
}

func getParams() []string {
	params := []string{}

	ch, err := fileUtil.ReadFile(options.parameters)
	if err != nil {
		panic(err)
	}
	for param := range ch {
		params = append(params, param)
	}

	return params
}

func getUrls() []string {
	urls := []string{}

	// read input from a file otherwise read from stdin
	if options.list != "" {
		ch, err := fileUtil.ReadFile(options.list)
		if err != nil {
			panic(err)
		}
		for url := range ch {
			urls = append(urls, url)
		}
	} else if fileUtil.HasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urls = append(urls, strings.TrimSpace(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	return urls
}

func combineStrat(urls []string) []string {
	modifiedUrls := []string{}

	for _, singleUrl := range urls {
		// parse each url
		parsedUrl, err := url.Parse(singleUrl)
		if err != nil {
			panic(err)
		}
		queryParams := parsedUrl.Query()
		numOfOldParams := len(queryParams)

		// only get new parameters so we don't accidentally override the current params
		oldKeys := []string{}
		for keys := range queryParams {
			oldKeys = append(oldKeys, keys)
		}

		for _, singeValue := range options.values {
			// each iteration contains a url with the number of parameters provided by the chunk size flag
			for iteration := 0; iteration < numOfOldParams; iteration++ {
				newQueryParams := url.Values{}

				// first add all parameters
				for _, key := range oldKeys {
					newQueryParams.Set(key, queryParams.Get(key))
				}
				// modify one parameter in each iteration
				if options.valueStrategy == "replace" {
					newQueryParams.Set(oldKeys[iteration], singeValue)
				} else {
					newQueryParams.Set(oldKeys[iteration], queryParams.Get(oldKeys[iteration])+singeValue)
				}

				parsedUrl.RawQuery = newQueryParams.Encode()
				modifiedUrls = append(modifiedUrls, parsedUrl.String())
			}
		}
	}

	return modifiedUrls
}

func ignoreStrat(urls []string, params []string) []string {
	modifiedUrls := []string{}

	for _, singleUrl := range urls {
		// parse each url
		parsedUrl, err := url.Parse(singleUrl)
		if err != nil {
			panic(err)
		}
		queryParams := parsedUrl.Query()
		numOfOldParams := len(queryParams)
		numOfIterations := int(math.Ceil(float64(len(params)) / float64(options.chunk-numOfOldParams)))

		// only get new parameters so we don't accidentally override the current params
		oldKeys := []string{}
		for keys := range queryParams {
			oldKeys = append(oldKeys, keys)
		}

		for _, singeValue := range options.values {
			newKeys, _ := sliceUtil.Diff(params, oldKeys)
			// each iteration contains a url with the number of parameters provided by the chunk size flag
			for iteration := 0; iteration < numOfIterations; iteration++ {
				newQueryParams := url.Values{}

				// add old parameters
				for key, value := range queryParams {
					newQueryParams.Set(key, value[len(value)-1])
				}

				// add new paramters
				for paramNum := 0; paramNum < options.chunk-numOfOldParams && len(newKeys) > 0; paramNum++ {
					newQueryParams.Set(pop(&newKeys), singeValue)
				}

				parsedUrl.RawQuery = newQueryParams.Encode()
				modifiedUrls = append(modifiedUrls, parsedUrl.String())
			}
		}
	}

	return modifiedUrls
}

func normalStrat(urls []string, params []string) []string {
	modifiedUrls := []string{}

	for _, singleUrl := range urls {
		// parse each url
		parsedUrl, err := url.Parse(singleUrl)
		if err != nil {
			panic(err)
		}
		queryParams := parsedUrl.Query()
		numOfOldParams := len(queryParams)
		numOfIterations := int(math.Ceil(float64(len(params)) / float64(options.chunk-numOfOldParams)))

		// only get new parameters so we don't accidentally override the current params
		oldKeys := []string{}
		for keys := range queryParams {
			oldKeys = append(oldKeys, keys)
		}

		for _, singeValue := range options.values {
			newKeys, _ := sliceUtil.Diff(params, oldKeys)
			// each iteration contains a url with the number of parameters provided by the chunk size flag
			for iteration := 0; iteration < numOfIterations; iteration++ {
				newQueryParams := url.Values{}

				// add old parameters
				for key := range queryParams {
					newQueryParams.Set(key, singeValue)
				}

				// add new paramters
				for paramNum := 0; paramNum < options.chunk-numOfOldParams && len(newKeys) > 0; paramNum++ {
					newQueryParams.Set(pop(&newKeys), singeValue)
				}

				parsedUrl.RawQuery = newQueryParams.Encode()
				modifiedUrls = append(modifiedUrls, parsedUrl.String())
			}
		}
	}

	return modifiedUrls
}

func newParamsOnlyStrat(urls []string, params []string) []string {
	modifiedUrls := []string{}
	numOfIterations := int(math.Ceil(float64(len(params)) / float64(options.chunk)))

	for _, singleUrl := range urls {
		// parse each url
		parsedUrl, err := url.Parse(singleUrl)
		if err != nil {
			panic(err)
		}

		// get the base url (without the params and values)
		baseUrl := parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path

		// parse the base url
		parsedUrl, err = url.Parse(baseUrl)
		if err != nil {
			panic(err)
		}
		for _, singeValue := range options.values {
			newKeys := params
			// each iteration contains a url with the number of parameters provided by the chunk size flag
			for iteration := 0; iteration < numOfIterations; iteration++ {
				newQueryParams := url.Values{}

				// set new parameters with the given values
				for paramNum := 0; paramNum < options.chunk && len(newKeys) > 0; paramNum++ {
					newQueryParams.Set(pop(&newKeys), singeValue)
				}

				// add parameters to a copy of the base url
				parsedUrl.RawQuery = newQueryParams.Encode()
				modifiedUrls = append(modifiedUrls, parsedUrl.String())
			}
		}
	}

	return modifiedUrls
}

// pops an item from the slice then removes it from the slice
func pop(aSlice *[]string) string {
	f := len(*aSlice)
	rv := (*aSlice)[f-1]
	*aSlice = (*aSlice)[:f-1]
	return rv
}
