# A Toy Project That May Well Succeed

This is a playful project I use while reading through Google Gemini’s API documentation pages. Much of the code is copied from the Gemini documentation and tweaked to make this program. The objective is to try create simple interfaces for some of Gemini’s API client calls using Gemini’s Go SDK and the Cobra CLI framework. The benefit of this is that users can craft their prompts and plug them in using more flexible scripting languages such as Bash and Python. Over time I hope to share usage examples on this repo.

## Compiling the Program

If you want to try out the this program you will need to compile it into an executable program. You need to already have Go installed on your machine to do so. Assumming you already have Go intalled, here are steps for building an executable file:

- Clone or this repository or download it as a zip file. Note: if you download, you would need to extract the content into a folder. To clone run the following in your shell:
  ```sh
  git clone github.com/ajuala/ggemini.git
  ```
- Make the project directory your working directory.
  ```sh
  cd gggemini
  ```
- Install project dependencies.
  ```sh
  go get
  ```
- Run Go build command
  ```sh
  go build
  ```

If everything go successfully, you should have a `ggemini` or `ggemini.exe` executable file inside the project directory. Feel free to move it to a more convenient directory.

**NOTE:** You need a Gemini API key to use this program. You can get one from Google Gemini and save it in your `GEMINI_API_KEY` environment variable.

## Available Commands

The `ggemini` program has the following commands at the moment:
1. `gentext`: This command takes a text prompt and returns a text response.
2. `genimage`: This command takes a text prompt and returns an image generated based on the prompt. By default the command encodes the output as Base64 when printing to the standard output, but saves directly as `.png` if `--output` option is specified.
3. `genimage`: This command takes a text prompt and an image prompt, the text prompt is used by Gemini to edit the image and outputs the edited image. Like `genimage`, printing directly to the standard output encodes the data to Base64 by default, otherwise, it is saved as `.png` to specified output file.
4. `genspeech`: This command takes a text prompt and one of Gemini’s supported voices and generates a speech from the prompt. It returns a `.wav` data which can be saved to file out output as Base64 encoded data to the standard output by default.

To see available commands, run:
```sh
./ggemini --help
```

To see available options for each command run `./gemini help <command>` for any of the commands. For example:
```sh
./ggemini help gentext
```
