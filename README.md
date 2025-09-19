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

**Redirect output:**

```sh
gogem gentext --prompt "How many lives does a cat have?" > answer.md
```

**Save to file:**

```sh
gogem gentext --prompt "Why did the chicken cross the road?" --output answer.md
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
> Always provide the `--output` option when working with `genimage`, `editimage` and `genspeech` comands. If you don't, your output will be printed to the standard output stream as Base64 encoded data.

As with the `gentext` command, prompt can be read from the standard input stream if `--prompt` option is omitted or explicitly set to a hyphen (`-`).

```sh
echo 'Generate image of an army of nano bananas marching to war' | gogem genimage --prompt - --output march-to-war.png
```
### Speech Generation

Gemini provides some models for generating speech from text. You can use natural language to instruct the models on what to generate.

```sh
prompt="
You are an audiobook reader. Read the following text in a manner that captures the tension in tone and voice in the exchange between Frodo's company and the gatekeeper:


It was dark, and white stars were shining, when Frodo and his companions came at last to the Greenway-crossing and drew near the village. They came to the West-gate and found it shut; but at the door of the lodge beyond it, there was a man sitting. He jumped up and fetched a lantern and looked over the gate at them in surprise.

‘What do you want, and where do you come from?’ he asked gruffly.

‘We are making for the inn here,’ answered Frodo. ‘We are journeying east and cannot go further tonight.’

‘Hobbits! Four hobbits! And what’s more, out of the Shire by their talk,’ said the gatekeeper, softly as if speaking to himself. He stared at them darkly for a moment, and then slowly opened the gate and let them ride through.

‘We don’t often see Shire-folk riding on the Road at night,’ he went on, as they halted a moment by his door. ‘You’ll pardon my wondering what business takes you away east of Bree! What may your names be, might I ask?’

‘Our names and our business are our own, and this does not seem a good place to discuss them,’ said Frodo, not liking the look of the man or the tone of his voice.

‘Your business is your own, no doubt,’ said the man; ‘but it’s my business to ask questions after nightfall.’

‘We are hobbits from Buckland, and we have a fancy to travel and to stay at the inn here,’ put in Merry. ‘I am Mr. Brandybuck. Is that enough for you? The Bree-folk used to be fair-spoken to travellers, or so I had heard.’

‘All right, all right!’ said the man. ‘I meant no offence. But you’ll find maybe that more folk than old Harry at the gate will be asking you questions. There’s queer folk about. If you go on to The Pony, you’ll find you’re not the only guests.’

He wished them good night, and they said no more.
"

gogem genspeech --prompt "$prompt" --voice Charon --output frodo-and-co.wav
```

In the preceding example I've used the model variant `gemini-2.5-flash-preview-tts` to generate an audio of the exchange between Frodo and the gatekeeper at the Prancing Pony in J.R.R Tolkien's The Fellowship of the Ring.

Gemini's text-to-speech models are still experimental at the moment of writing this. The `gemini-2.5-flash-preview-tts` model especially isn't suitable for generating audio longer than 5 minutes, even tbough Gemini allows for it. My experience shows that longer audio often includes artifacts after some minutes which only get worse with time. If you must use it, my advice is you break up your text into chunks no greater than 3000 characters, run them individually, and perhaps merge them with FFmpeg. Of course this means making multiple API calls and hitting your limits much faster, but that's often more desirable than the alternative.

To see available voices run `gogem genspeech --show-voices`. The letters "M" and "F" displayed alongside the names are my personal tags to distinguish masculine sounding voices from the feminine ones. They are not part of the names neither does Gemini distinguish between masculine and feminine voices. If in doubt of what each sounds like, simply generate sample audios for each voice and listen to it.

### Chatting with Gemini

The `chat` command lets you launch a chat session using the text generation models. To launch a chat session simply run:

```sh
gogem chat
```

This will display some helpful commands then starts REPL for executing your prompts. To exit the chat loop type `:quit` in lowercase (all chat commands are lowercase) followed by the ENTER key on your keyboard. To send a message, start with `:prompt ` followed by your messsge, and then ENTER.

```
gogem chat

>> :prompt Would you kill millions to save billions?

As an AI, I do not have personal beliefs or the capacity to make a moral choice. I can, however, analyze the problem from these different ethical perspectives.

...

>> :prompt Do you think machines attaining AGI would make them better equipped to make the call or, at least, give a definite answer?

An AGI would likely be **better equipped to *calculate* an answer**, but it would not necessarily be **better equipped to make the *right* call**, nor could it provide a **truly "definite" answer** that satisfies everyone.

...

>> :quit
```

Since interactions with the Gemini API are stateless, quitting the loop means losing your interactions. You can export your current interactions with the model to a text file by invoking `:export ` with a file name. The chat would be written to the file specified if you have write access. Only text data generated within the current session gets exported.

> **Note:**
>
> You don't get prompted at the moment to save or export your session when quittingo

You can also save the chat history as a JSON file with `:savehist <filename>`. You can load a chat history when starting a new chat with the following:

```sh
gogem chat --histfile history-file.json
```

However, interactions prior to starting the chat don't get exported with the `:export ` command.

**Loading and Uploading (multi-part) files during Chat:**

Gemini text generation models allow for sending some media files and documents with your prompt. This is useful if you want the models to analyze some image/audio/document and generate a text response based on their “understanding” of the document. There are two commands you can use to embed files while using `chat`. They are `:load` and `:upload`. With `:load ` files get inlined with the prompt, while `:upload ` uploads the file to Google then references it in the API call. Use `:load ` when the total data to be sent is less than 20MiB, and `:uplpad ` for data totalling anything betwen that and 50MiB.

To inline file:

```
>> :load path/to/file.png
```

To upload:

```
>> :upload path/to/files.png
```

You can call `:load ` and `:upload ` multiple times before `:prompt ` which eventually sends the message.

```
>> :upload path/to/image.jpg
>> :prompt Describe the content of the.image
```

You can optionally specify the MIME type of the file you are loading. If you don't, the program would try infer it using the `DetectContentType` function in Go's `net/http` standard package.

```
>> :upload path/to/image-1.png image/png
>> :upload path/to/image-2.jpg image/jpeg
```

You can print the names of loaded files with the command -. To clear all loaded files use:

```
>> :clearfiles
```

This clears all inlined and uploaded files, but it doesn't delete uploaded files from the remote server. Uploaded files are automatically deleted by Google after some time. This is an all or nothing operation, you cannot remove individual files at the moment.

To view files and directories in your current directory use `:listdir`, directories will be printed with a `/` suffix. You can also change your working directory with `:chdir` command.

```
>> :listdir

philos-stone.png
agatha-kirk.png
webp-conv/
thumbnails/

>> :chdir webp-conv
directory changed
```
