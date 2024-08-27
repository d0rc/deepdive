package main

import (
	console_tool "deepdive/utils/console-tool"
	disk_cache "deepdive/utils/disk-cache"
	"deepdive/utils/graphs"
	local_llm "deepdive/utils/local-llm"
	"deepdive/utils/parser"
	retry_tool "deepdive/utils/retry-tool"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/logrusorgru/aurora"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var engineURL = flag.String("engine-url", "http://localhost:7999/v1/chat/completions", "LLM API endpoint")

// var engineURL = flag.String("engine-url", "https://api.together.xyz/v1/chat/completions", "LLM API endpoint")
// var defaultModel = flag.String("model", "meta-llama/Meta-Llama-3.1-405B-Instruct-Turbo", "which model to use")
var defaultModel = flag.String("model", "qwen2-72b-32k:latest", "which model to use")
var pageCacheDB = flag.String("page-cache", "page-cache.db", "a path to web cache file for this research")
var endpointToken = flag.String("token", "", "token to use for endpoint")
var userQuestion = flag.String("question", "What are the best seaside beaches near Milan, Italy?", //"What are the best seaside beaches near Milan, Italy?",
	"the questions you seek answer for")
var outputGraphPath = flag.String("output-graph-path", "tree.md", "a path to the knowledge graph")

func main() {
	maxThreads := make(chan struct{}, 4)
	lg := console_tool.ConsoleInit("mk-prompt")
	engine := &local_llm.LLMEngine{
		Endpoint: *engineURL,
		Token:    *endpointToken,
		Model:    *defaultModel,
	}
	lg.Info().Msgf("starting up with question: %s", aurora.BrightYellow(*userQuestion))

	ts := time.Now()
	moreQuestions, err := getMoreQuestions(engine, *userQuestion)
	if err != nil {
		lg.Fatal().Msgf("error getting more questions: %v", err)
	}
	fmt.Printf("Got these questions in %v\n", time.Since(ts))
	moreQuestions = append(moreQuestions, *userQuestion)
	for idx, question := range moreQuestions {
		if question == *userQuestion {
			fmt.Printf("%04d. Q: %s\n",
				aurora.BrightGreen(idx),
				aurora.BrightYellow(question))
		} else {
			fmt.Printf("%04d. Q: %s\n",
				aurora.BrightGreen(idx),
				aurora.BrightCyan(question))
		}
	}

	ts = time.Now()
	keywords := make([]string, 0)
	keywordsLock := sync.RWMutex{}
	wg := sync.WaitGroup{}
	for _, question := range moreQuestions {
		wg.Add(1)
		maxThreads <- struct{}{}
		go func(question string) {
			defer func() {
				wg.Done()
				<-maxThreads
			}()
			searchQueries, err := getSearchQueries(engine, question)
			if err != nil {
				lg.Fatal().Msgf("error getting search queries: %v", err)
			}

			keywordsLock.Lock()
			keywords = append(keywords, searchQueries...)
			keywordsLock.Unlock()
		}(question)
	}

	wg.Wait()

	fmt.Printf("Generated %d search queries in %v\n",
		aurora.BrightGreen(len(keywords)),
		aurora.White(time.Since(ts)))
	summaries := ""
	summariesLock := sync.RWMutex{}
	summariesTotal := 0
	summariesAccepted := 0
	wg = sync.WaitGroup{}

	graph := graphs.NewGraph()
	graphLock := sync.RWMutex{}

	for idx, keyword := range keywords {
		wg.Add(1)
		maxThreads <- struct{}{}
		go func(idx int, keyword string) {
			defer func() {
				wg.Done()
				<-maxThreads
			}()
			fmt.Printf("%04d. %s\n",
				aurora.BrightGreen(idx),
				aurora.White(keyword))

			ts = time.Now()
			results, err := search(keyword)
			if err != nil {
				lg.Fatal().Msgf("error running search: %v", err)
			}

			fmt.Printf("\tSearch took: %v, got %d results.\n", time.Since(ts), results.NumberOfResults)
			summariesSlice, err := getSearchSummary(engine, *userQuestion, keyword, *results)
			if err != nil {
				lg.Fatal().Msgf("error getting search summary: %v", err)
			}

			for _, summary := range summariesSlice {
				fmt.Printf("\tSummary:\n\t%s\n", aurora.White(summary))

				responseClass, err := getResponseClassification(engine, *userQuestion, summary)
				if err != nil {
					lg.Fatal().Msgf("error classifing summary: %v", err)
				}

				if responseClass == ResponseIsGood {
					// extract NEs
					namedEntities, err := retry_tool.RetryCallWithCount(func() ([]string, error) {
						return getListOfNEs(engine, summary)
					}, 40)
					if err != nil {
						lg.Fatal().Msgf("error extracting named entities: %v", err)
					}

					fmt.Printf("Extracted following named entities:\n%s\n", strings.Join(namedEntities, "\n"))

					for _, entity := range namedEntities {
						connections, err := retry_tool.RetryCallWithCount(func() (map[string][]string, error) {
							return getListOfEntityConnections(engine, namedEntities, entity, summary)
						}, 40)
						if err != nil {
							lg.Fatal().Msgf("error extracting entity connections...!")
						}
						fmt.Printf("Connections: %v", connections)
						graphLock.Lock()
						for k, v := range connections {
							for _, el := range v {
								graph.AddEdge(entity, el, k)
							}
						}
						_ = os.WriteFile(*outputGraphPath, []byte(graph.RenderMermaid()), 0666)
						graphLock.Unlock()
					}
				}

				summariesLock.Lock()
				summariesTotal++
				if responseClass == ResponseIsGood {
					summariesAccepted++
					fmt.Printf("Accepting summary: [%d of %d]!\n", aurora.BrightGreen(summariesAccepted), aurora.BrightYellow(summariesTotal))
					summaries += fmt.Sprintf("Keywords: %s\nSummary: %s\n", keyword, summary)
				} else {
					fmt.Printf("Skipping summary: [%d of %d]!\n", aurora.BrightGreen(summariesAccepted), aurora.BrightYellow(summariesTotal))
					if strings.Contains(summary, "Buondonno") {
						fmt.Printf("It should be a good one...!\n")
					}
				}
				summariesLock.Unlock()
			}
		}(idx, keyword)
	}

	wg.Wait()

	fmt.Printf("Final summaries:\n%s", summaries)
}

type SearchResults struct {
	Query           string `json:"query"`
	NumberOfResults int    `json:"number_of_results"`
	Results         []SearchResultEntry
}

type SearchResultEntry struct {
	Url       string   `json:"url"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Engine    string   `json:"engine"`
	Engines   []string `json:"engines"`
	Positions []int    `json:"positions"`
}

func search(query string) (*SearchResults, error) {
	// http://localhost:4000/search?q=matteo+trimarchi&format=json
	results := &SearchResults{}
	query = url.QueryEscape(query)

	res, err := http.Get(fmt.Sprintf("http://localhost:4000/search?q=%s&format=json", query))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

//go:embed prompts/get-more-questions.md
var getMoreQuestionsPrompt string

//go:embed prompts/get-more-targeted-questions.md
var getMoreTargetedQuestions string

func getMoreQuestions(engine *local_llm.LLMEngine, question string) ([]string, error) {
	generators := []string{getMoreQuestionsPrompt, getMoreTargetedQuestions}
	allResults := make([]string, 0)
	allResultsLock := sync.RWMutex{}
	wg := sync.WaitGroup{}
	for _, generator := range generators {
		wg.Add(1)
		go func(generator string) {
			defer wg.Done()
			prompt := local_llm.NewThread().
				AddSystemMessage(generator).
				AddUserMessage(fmt.Sprintf("Original Question: \"%s\"", question))

			results, err := engine.Run(prompt, 0.4)
			if err != nil {
				fmt.Printf("[%s] %v\n", aurora.BrightRed("ERROR-LLM"), err)
				return
			}

			result := make([]string, 0)
			message := results[0].Content

			err = json.Unmarshal([]byte(message), &result)
			if err != nil {
				fmt.Printf("[%s] %v\n", aurora.BrightRed("ERROR-JSON"), err)
				return
			}

			allResultsLock.Lock()
			allResults = append(allResults, result...)
			allResultsLock.Unlock()
		}(generator)
	}

	wg.Wait()

	return allResults, nil
}

//go:embed prompts/get-search-summary.md
var getSearchSummaryPrompt string

func splitStringIntoChunksUTF(inputStr string, chunkSize int) []string {
	var chunks []string
	start := 0
	for start < len(inputStr) {
		runeCount := utf8.RuneCountInString(inputStr[start:])
		if runeCount <= chunkSize {
			chunks = append(chunks, inputStr[start:start+runeCount])
			start += runeCount
		} else {
			end := start + chunkSize
			chunk := inputStr[start:end]
			chunks = append(chunks, chunk)
			start = end
		}
	}
	return chunks
}

type summaryChunk struct {
	Content  string
	Url      string
	Keywords string
	Question string
}

func getSearchSummary(engine *local_llm.LLMEngine, question, keywords string, results SearchResults) ([]string, error) {
	//searchResultsYaml := ""
	chunksToSummarise := make([]summaryChunk, 0)
	for _, res := range results.Results {
		chunksToSummarise = append(chunksToSummarise, summaryChunk{Content: fmt.Sprintf(`- url: %s
  title: %s
  content: %s
`, res.Url, res.Title, res.Content), Url: "", Question: question, Keywords: keywords})

		markdown, err := getWebPageAsMarkdown(res.Url, nil)
		if err != nil || markdown == "" || len(res.Content) > len(markdown) {
			fmt.Printf("Error loading page %s, due to error %v\n", res.Url, err)
			// split markdown by max_characters_per_chunk = 3000
			for _, chunk := range splitStringIntoChunksUTF(markdown, 3000) {
				chunksToSummarise = append(chunksToSummarise, summaryChunk{Content: fmt.Sprintf(`- url: %s
  title: %s
  content: %s
`, res.Url, res.Title, chunk), Question: question, Url: res.Url, Keywords: keywords})
			}
		}
	}

	summaries := make([]string, 0)
	for _, chunk := range chunksToSummarise {
		prompt := local_llm.NewThread().
			AddSystemMessage(getSearchSummaryPrompt).
			AddUserMessage(fmt.Sprintf(`Question: %s
Keywords: [%s]
Search Engine Output:
%s
`, question, keywords, chunk))

		res, err := engine.Run(prompt, 0.4)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, fmt.Sprintf("Summary for question: %s\nUrl: %s\n%s\n", chunk.Question, chunk.Url, res[0].Content))
	}

	return summaries, nil
}

//go:embed prompts/get-search-queries.md
var getSearchQueriesPrompt string

func getSearchQueries(engine *local_llm.LLMEngine, question string) ([]string, error) {
	prompt := local_llm.NewThread().
		AddSystemMessage(getSearchQueriesPrompt).
		AddUserMessage(fmt.Sprintf("Input Question: \"%s\"", question))

	results, err := engine.Run(prompt, 0.4)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	message := results[0].Content

	_, err = parser.TryParseData(func(subContext string) ([]string, error) {
		err = json.Unmarshal([]byte(subContext), &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}, message)

	if err != nil {
		return nil, err
	}

	return result, nil
}

//go:embed prompts/get-list-of-NEs.md
var getListOfNEsPrompt string

func getListOfNEs(engine *local_llm.LLMEngine, context string) ([]string, error) {
	prompt := local_llm.NewThread().
		AddSystemMessage(getListOfNEsPrompt).
		AddUserMessage(fmt.Sprintf("Context: \"%s\"", context))

	results, err := engine.Run(prompt, 0.4)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	message := results[0].Content

	message = strings.ReplaceAll(message, "\",\n]", "\"\n]")
	_, err = parser.TryParseData(func(subContext string) ([]string, error) {
		err := json.Unmarshal([]byte(subContext), &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}, message)

	if err != nil {
		return nil, err
	}

	return result, nil
}

//go:embed prompts/get-entity-connections.md
var getEntityConnectionsPrompt string

func getListOfEntityConnections(engine *local_llm.LLMEngine, entities []string, entity, context string) (map[string][]string, error) {
	prompt := local_llm.NewThread().
		AddSystemMessage(getEntityConnectionsPrompt).
		AddUserMessage(fmt.Sprintf("List of Entities: %s\nContext: %s\nThe Entity of Interest: '%s'\n",
			mkJson(entities), context, entity))

	results, err := engine.Run(prompt, 0.4)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	message := results[0].Content

	err = json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func mkJson(entities []string) string {
	result, err := json.Marshal(entities)
	if err != nil {
		panic("failed to create json for a list of strings...!")
	}

	return string(result)
}

//go:embed prompts/get-response-classification.md
var getResponseClassificationPrompt string

type ResponseClass int

const (
	ResponseIsGood ResponseClass = iota
	ResponseIsBad
	ResponseUndefined
)

func getResponseClassification(engine *local_llm.LLMEngine, question, context string) (ResponseClass, error) {
	prompt := local_llm.NewThread().
		AddSystemMessage(getResponseClassificationPrompt).
		AddUserMessage(fmt.Sprintf("Question: %s\nContext: %s", question, context))

	results, err := engine.Run(prompt, 0.4)
	if err != nil {
		return ResponseUndefined, err
	}

	if strings.Contains(results[0].Content, "PROBLEM-SOLVED") {
		return ResponseIsGood, nil
	}

	return ResponseIsBad, nil
}

var diskCache = func() *disk_cache.DiskCache {
	dc, err := disk_cache.NewDiskCache(*pageCacheDB)
	if err != nil {
		panic("error loading page-cache.db")
	}

	return dc
}()

func getWebPageAsMarkdown(urlString string, proxies []string) (string, error) {
	if tmpMd, ok := diskCache.Get(urlString); ok {
		return tmpMd, nil
	}
	// Set up the HTTP client with timeout and retries
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				if len(proxies) > 0 {
					return url.Parse(proxies[0])
				}
				return nil, nil
			},
		},
	}

	// Make the HTTP request
	resp, err := client.Get(urlString)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// fmt.Printf("Content-type: %v\n", resp.Header.Get("Content-Type"))

	if !(strings.Contains(resp.Header.Get("Content-Type"), "text/html")) {
		fmt.Printf("Content-type: %v\n", resp.Header.Get("Content-Type"))
		return "", err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert the HTML to Markdown
	converter := md.NewConverter("", true, nil)
	converter.Remove("img")
	markdown, err := converter.ConvertString(string(body))
	if err != nil {
		return "", err
	}

	_ = diskCache.Set(urlString, markdown)

	return markdown, nil
}
