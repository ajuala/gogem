# A CLI program for Making API Calls to Google Gemini

This is a fun project which I used while studying the Google Gemini API documentation pages. The objective is to try create simple interfaces for some of Gemini’s API client calls using Gemini’s Go SDK and the Cobra CLI framework. The benefit of this is that users can craft their prompts and plug them in using more flexible scripting languages such as Bash and Python. Over time I hope to share usage examples on this repo.

## Installation

To build and use this program, you need to have Go installed on your machine. If you already have Go intalled, run the following command in your shell program:

```sh
go install github.com/ajuala/gogem@latest
```


**NOTE:** You need a Gemini API key to use this program. You can get one from Google Gemini and save it in your `GEMINI_API_KEY` environment variable.

## Available Commands

The `gogem` program has the following commands at the moment:
1. `gentext`: This command takes a text prompt and returns a text response.
2. `genimage`: This command takes a text prompt and returns an image generated based on the prompt. By default the command encodes the output as Base64 when printing to the standard output, but saves directly as `.png` if `--output` option is specified.
3. `genimage`: This command takes a text prompt and an image prompt, the text prompt is used by Gemini to edit the image and outputs the edited image. Like `genimage`, printing directly to the standard output encodes the data to Base64 by default, otherwise, it is saved as `.png` to specified output file.
4. `genspeech`: This command takes a text prompt and one of Gemini’s supported voices and generates a speech from the prompt. It returns a `.wav` data which can be saved to file out output as Base64 encoded data to the standard output by default.

To see available commands, run:
```sh
gogem --help
```

To see available options for each command run `./gemini help <command>` for any of the commands. For example:
```sh
gogem help gentext
```
