package main

import (
	console_tool "deepdive/utils/console-tool"
	local_llm "deepdive/utils/local-llm"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
)

var engineURL = flag.String("engine-url", "http://localhost:7999/v1/chat/completions", "LLM API endpoint")
var defaultModel = flag.String("model", "qwen2:72b-instruct-q6_K", "which model to use")
var endpointToken = flag.String("token", "token-abc123", "token to use for endpoint")
var taskDescription = flag.String("task-description", "Bot will be provided with two variables - `Context` and `Name`. The context is a list of names of different persons in various forms and variations. The name is a name of a person. The goal is to provide a JSON list of all forms and variations of the person's name mentioned in context.", "task description")

func main() {
	lg := console_tool.ConsoleInit("mk-prompt")
	compilerPrompt := "Here is a task description for which I would like you to create a high-quality prompt template for:\n<task_description>\n%s\n</task_description>\nBased on task description, please create a well-structured prompt template that another AI could use to consistently complete the task. The prompt template should include:\n- Do not inlcude <input> or <output> section and variables in the prompt, assume user will add them at their own will. \n- Clear instructions for the AI that will be using this prompt, demarcated with <instructions> tags. The instructions should provide step-by-step directions on how to complete the task using the input variables. Also Specifies in the instructions that the output should not contain any xml tag. \n- Relevant examples if needed to clarify the task further, demarcated with <example> tags. Do not include variables in the prompt. Give three pairs of input and output examples.   \n- Include other relevant sections demarcated with appropriate XML tags like <examples>, <instructions>.\n- Use the same language as task description. \n- Output in ``` xml ``` and start with <instruction>\nPlease generate the full prompt template with at least 300 words and output only the prompt template."
	engine := &local_llm.LLMEngine{
		Endpoint: *engineURL,
		Token:    *endpointToken,
		Model:    *defaultModel,
	}
	lg.Info().Msgf("starting up with url: %s", *engineURL)

	prompt := local_llm.NewThread().
		AddUserMessage(fmt.Sprintf(compilerPrompt, *taskDescription))

	result, err := engine.Run(prompt, 0.4)
	if err != nil {
		lg.Fatal().Err(err).Msgf("error running llm engine")
	}

	fmt.Printf("Prompt:\n%s\n", aurora.White(result[0].Content))
}
