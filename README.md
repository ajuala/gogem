# A CLI program for Making API Calls to Google Gemini

This is a fun project which I used while studying the Google Gemini API documentation pages. The objective is to try create simple interfaces for some of Gemini’s API client calls using Gemini’s Go SDK and the Cobra CLI framework. The benefit of this is that users can craft their prompts and plug them in using more flexible scripting languages such as Bash and Python. Over time I hope to share usage examples on this repo.

## Installation

To build and use this program, you need to have Go installed on your machine. If you already have Go intalled, run the following command in your shell program:

```sh
go install github.com/ajuala/gogem@latest
```


**NOTE:** You need a Gemini API key to use this program. You can get one from [Google Gemini](https://aistudio.google.com/apikey) and save it in your `GEMINI_API_KEY` environment variable.

## Available Commands

The `gogem` program has the following commands at the moment:
1. `gentext`: This command takes a text prompt and returns a text response.
2. `genimage`: This command takes a text prompt and returns an image generated based on the prompt. By default the command encodes the output as Base64 when printing to the standard output, but saves directly as `.png` if `--output` option is specified.
3. `genimage`: This command takes a text prompt and an image prompt, the text prompt is used by Gemini to edit the image and outputs the edited image. Like `genimage`, printing directly to the standard output encodes the data to Base64 by default, otherwise, it is saved as `.png` to specified output file.
4. `genspeech`: This command takes a text prompt and one of Gemini’s supported voices and generates a speech from the prompt. It returns a `.wav` data which can be saved to file out output as Base64 encoded data to the standard output by default.
5. `chat`: Starts a chat session with Gemini on your terminal.

To see available commands, run:
```sh
gogem --help
```

To see available options for each command run `./gemini help <command>` for any of the commands. For example:
```sh
gogem help gentext
```

## Examples

### Generating Texts

For simple text generation use the `gentext` command:

```sh
gogem gentext --prompt "How many lives does a cat have?"
```

If successful, this prints the response to your standard output. You can redirect output to a file or use the `--output` option to write directly to file.

```sh
gogem gentext --prompt "How many lives does a cat have?" > answer.md
```

```sh
gogem gentext --prompt "How many lives does a cat have?" --output answer.md
```

Prompts can also be read from the standard input stream if the `--prompt` option is omitted.

```sh
gogem gentext --output answer.md < prompt.txt
```

### Generating Images

To generate images use the `genimage` command:

```sh
gogem genimage --output banana-army.png --prompt "Generate image of an army of nano bananas marching to war"
```

> **Note**
> Always provide the `--output` option when working with `genimage`, `editimage` and `genspeech` comandss. If you don't your output would be printed to the standard output stream as Base64 encoded data.

As with the `gentext` command, prompt can be read from the standard input stream if `--prompt` option is omitted or explicitly set to a hyphen (`-`).

```sh
echo 'Generate image of an army of nano bananas marching to war' | gogem genimage --prompt - --output march-to-war.png
```
