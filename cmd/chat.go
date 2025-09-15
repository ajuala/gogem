/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"regexp"
	"context"
	"net/http"
	"unicode/utf8"
	"encoding/json"
	"path/filepath"

	"github.com/ajuala/gogem/ai"
	"github.com/ergochat/readline"
	"google.golang.org/genai"
	"github.com/spf13/cobra"
	"github.com/google/shlex"
	"github.com/mattn/go-zglob"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat starts chat session with Gemini. Chat from the active session may be exported",
	Run: func(cmd *cobra.Command, args []string) {
		var histJSON string

		histfile, err := cmd.Flags().GetString("histfile")

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		if histfile == "" {
			d, err := cmd.Flags().GetString("histjson")

			if err != nil {
				eprint(err)
				os.Exit(1)
			}

			histJSON = d

		} else {
			data, err := os.ReadFile(histfile)
			if err != nil {
				eprint(err)
				os.Exit(1)
			}

			histJSON = string(data)
		}

		sysPrompt := strings.TrimSpace(sysPrompt)
		model := strings.TrimSpace(model)
		temp, topK, topP := getTempTopKP()
		client, chat, err := ai.CreateChat(sysPrompt, histJSON, model, apiKey, temp, topK, topP)

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		err = chatREPL(client, chat)
		if err != nil {
			eprint(err)
			os.Exit(1)
		}
	},
}

func chatREPL(client *genai.Client, chat *genai.Chat) error {
	// Chat session
	type ChatItem struct{
		question string
		response string
	}

	var sess []ChatItem

	preamble := "Chat session started\n"
	helpMsg := "Available commands:\n" +
	":q,:quit  Quits chat session\n" +
	":sh,:savehist <filepath> Saves chat session history to <filepath>\n" +
	":ls,listdir [path] Lists files and directories in [path]\n" +
	":cd <directory> Changes working directory to <directory>\n" +
	":ex,export <filepath> Exports chat in current session to <filepath>\n" +
	":vim Opens a Vim session for multiline prompt editing\n" +
	":p,:prompt <prompt text> Single line prompt command\n" +
	":up,:upload <filepath> File to upload\n" +
	":files Shows loaded file list\n" +
	":cf,:clearfiles Clears all files from list\n" +
	"h,:help Prints this help message\n"

	upload := func(fpath, mime string) (*genai.File, error) {

		if mime == "" {
			buf := make([]byte, 512)

			f, err := os.Open(fpath)
			if err != nil {
				return nil, err
			}

			n, err := f.Read(buf)

			f.Close()

			if err != nil {
				return nil, err
			}


			mime = http.DetectContentType(buf[:n])
		}

		ctx := context.Background()

		return client.Files.UploadFromPath(ctx, fpath, &genai.UploadFileConfig{
			MIMEType: mime,
		})
	}

	resolvePath := func(path string) string {
		if expandedUser, err := expandUser(path); err != nil {
			path = expandedUser
		}

		return filepath.Clean(os.ExpandEnv(path))
	}

	listDir := func(pth string) {
		pth = resolvePath(pth)


		match, err := zglob.Glob(pth)

						if err != nil {
							eprint(err)
							return
						}

						if len(match) == 1 {
							pthIsDir, err := isDir(match[0])
							if err != nil {
								eprint(err)
							} else if pthIsDir {
								match, err := zglob.Glob(filepath.Join(match[0], "*"))
								if err != nil {
									eprint(err)
								} else {
									for _, item := range match {
										if isdir, _ := isDir(item); isdir {
											fmt.Printf("%s/\n", item)
										} else {
										fmt.Println(item)
										}
									}
								}
							}
						} else {
							for _, f := range match {
								if isdir, _ := isDir(f); isdir {
									fmt.Printf("%s/\n", f)
								} else {
								fmt.Println(f)
								}
							}
						}
	}


	rl, err := readline.New(">> ")

	if err != nil {
		return err
	}

	defer rl.Close()

	editorInitialContent := "Edit text"

	fmt.Println(preamble)
	fmt.Println(helpMsg)

	// Line begins with a command prefix patterb
		cmdRefRe := regexp.MustCompile("^[ \\t]*(:[a-z]+|/\\d+(?:[ \\t]*/\\d+))[ \\t]*(.*)")
		refsRe := regexp.MustCompile("/(\\d+)[ ]*")
		// cmdRe := regexp.MustCompile("^:([a-z]+)[ ]+(.*)")
		promptRe := regexp.MustCompile("^:p(?:rompt)?\\b\\s*(.+)")

		var parts []genai.Part
		var stagedFiles []string

	loop:
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if ! strings.HasPrefix(line, ":") {
			eprint("error: please prefix your input with a command")
			continue
		}


		switch line {
		case ":quit", ":q":
			break loop
			case ":help", ":h":
				fmt.Println(helpMsg)
				continue
		}

		prompt := ""

		if m := promptRe.FindStringSubmatch(line); m != nil {
			prompt = m[1]
		} else {

		switch line {
			case ":vim", ":nvim":
				useNvim := false
				if line == ":nvim" {
					useNvim = true
				}

				text, err := vimEditor(editorInitialContent, ".md", useNvim)
				if err != nil {
					eprint(err)
					continue
				}

				if text == "" || text == editorInitialContent {
					continue
				}

				prompt = strings.TrimSpace(text)


			default:
				m := cmdRefRe.FindStringSubmatch(line)

				if strings.HasPrefix(m[1], ":") {
					cmdList, err := shlex.Split(line[1:])
					if err != nil {
						eprint(err)
						continue
					}

					if len(cmdList) == 0 {
						continue
					}

					// Commands switch
					switch cmdList[0] {
					case "upload", "up":
						// upload to remote and include in parts variable
						if !(len(cmdList) != 2 || len(cmdList) != 3) {
							eprint(fmt.Sprintf("command :upload requires 1 argument, %d given", len(cmdList) -  1))
							continue
						}

						mime := ""
						
						if len(cmdList) == 3 {
							mime = cmdList[2]
						}

						f, err := upload(cmdList[1], mime)
						if err != nil {
							eprint(err)
							continue
						}

						parts = append(parts, *genai.NewPartFromURI(f.URI, f.MIMEType))
						stagedFiles = append(stagedFiles, cmdList[1])
						fmt.Println("file loaded")
						continue
					case "load":
						// load file into parts variable
						if !(len(cmdList) != 2 || len(cmdList) != 3) {
							eprint(fmt.Sprintf("command :load requires a single argument, %d given", len(cmdList) - 1))
							continue
						}


						data, err := os.ReadFile(cmdList[1])
						if err != nil {
							eprint(err)
							continue
						}

						mime := ""
						if len(cmdList) == 3 {
							mime = cmdList[2]
						} else {
						mime = http.DetectContentType(data)
						}
						parts = append(parts, genai.Part{
							InlineData: &genai.Blob{
								MIMEType: mime,
								Data: data,
							},
						})

						stagedFiles = append(stagedFiles, cmdList[1])

						continue

					case "cf", "clearfiles":
						parts = nil
						stagedFiles = nil
						fmt.Println("staged files cleared")
						continue

					case "show":
						if len(sess)  == 0 {
							eprint("error: current chat session is empty")
						} else {
							switch len(cmdList[1:]) {
						case 0:
							// print last response
							fmt.Println(sess[len(sess) - 1].response)
						case 1:
							// print the previous n response
							if i, err := strconv.Atoi(cmdList[1]); err != nil {
								eprint(err)
							} else if i < 0 {
								eprint("error: index cannot be negative")
							} else if len(sess) < i {
								eprint(fmt.Sprintf("error: index is out of range, current session has %d only items", len(sess)))
							} else {
								fmt.Println(sess[len(sess) - i].response)
							}

						default:
							eprint(fmt.Sprintf("error: print requires one argument, %d given", len(cmdList)))
						}
						}

						continue

					case "export": //export session
					switch len(cmdList) {
					case 1:
						eprint("error: export requires a filename argument, none given")
					case 2:
						if len(sess) == 0 {
							eprint("error: current session is empty")
						} else {
							var outList []string
							for _, item := range sess {
								input := item.question
								response := fmt.Sprintf("\n----- response start -----\n\n%s\n\n----- response end -----\n", item.response)
								outList = append(outList, strings.Join([]string{input, response}, "\n"))
							}

							outString := strings.Join(outList, "\n\n")
							err := os.WriteFile(cmdList[1], []byte(outString), 0644)
							if err != nil {
								eprint(err)
							}
						}

					default:
						eprint(fmt.Sprintf("command export requires one argument, %d given", len(cmdList) - 1))

					continue
					}

				case "listdir", "ls":
					switch len(cmdList[1:]) {
					case 0:
						matches, err := zglob.Glob("./*")
						if err != nil {
							eprint(err)
						} else {
							for _, f := range matches {
								if isdir, _ := isDir(f); isdir {
									fmt.Printf("%s/\n", f)
								} else {
								fmt.Println(f)
								}
							}
						}

					case 1:
						listDir(cmdList[1])
					default:
						for _, ptn := range cmdList[1:len(cmdList) - 1] {
							listDir(ptn)
							fmt.Println("\n-----\n")
						}

						listDir(cmdList[len(cmdList) - 1])
					}
					continue


				case "pwd":
					if len(cmdList) != 1 {
						eprint("error: command pwd takes no argument")
					} else {
						if dir, err := os.Getwd(); err != nil {
							eprint(err)
						} else {
							fmt.Println(dir)
						}
					}

					continue

				case "cd":
					switch len(cmdList[1:]) {
					case 0:
						continue
					case 1:
						path := cmdList[1]
						if expandedUser, err := expandUser(path); err == nil {
							path = expandedUser
						}

						path = os.ExpandEnv(path)
						if err := os.Chdir(path); err != nil {
							eprint(err)
						} else {
							fmt.Println("changed working directory")
						}
					default:
						eprint(fmt.Sprintf("error: command cd takes one argument, %d given", len(cmdList) - 1))
					}
					continue

				case "savehist", "sh":
					if len(cmdList[1:]) != 1 {
						eprint(fmt.Sprintf("error: command savehist takes one argument, %d given", len(cmdList[1:])))
					} else {
						path := resolvePath(cmdList[1])
						history := chat.History(false)
						histJSON, err := json.Marshal(history)

						if err != nil {
							eprint(err)
						} else {
							if err := os.WriteFile(path, histJSON, 0644); err != nil {
								eprint(err)
							}


						}
						fmt.Println(path)
					}
					continue

				case "read", "r":
					if len(cmdList[1:]) != 2 {
						eprint(fmt.Sprintf("error: command read takes one argument <filepath>, %d given", len(cmdList[1:])))
						continue
					}

					b, err := os.ReadFile(cmdList[1])
					if err != nil {
						eprint(err)
						continue
					}

					if !utf8.Valid(b) {
						eprint("error: file contains non-valid utf-8 data")
						continue
					}

					prompt = string(b)

				case "files":
					if len(stagedFiles) == 0 {
						fmt.Println("file list is empty")
					} else {

						for _, f := range stagedFiles {
							fmt.Println(f)
						}
					}
					continue

				default:
					eprint("error: command not implented")
					continue
					}

					continue

				} else {
					// WIP
					n := refsRe.FindAllStringSubmatch(m[1], -1)
					fmt.Println(n)
				}

				continue





		}
		}



					ctx := context.Background()
					parts_ := append(parts, *genai.NewPartFromText(prompt))
					stream := chat.SendMessageStream(ctx, parts_...)

					outText := ""


					for chunk, err := range stream {
						if err != nil {
							eprint(err)

				// Reset variables
				prompt = ""

				if parts != nil {
				parts = nil
				stagedFiles = nil
				fmt.Println("staged files cleared")
				}

							continue loop
						}

						part := chunk.Candidates[0].Content.Parts[0]
						fmt.Printf("%s", part.Text)
						outText += part.Text
					}

					fmt.Printf("\n")

						sess = append(sess, ChatItem{prompt, outText})




				// Reset variables
				prompt = ""

				if parts != nil {
				parts = nil
				stagedFiles = nil
				fmt.Println("staged files cleared")
				}

	}

	return nil
}



func expandUser(path string) (string, error) {
    if strings.HasPrefix(path, "~") {
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return strings.Replace(path, "~", home, 1), nil
    }
    return path, nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		// Handle the case where the path does not exist
		if os.IsNotExist(err) {
			return false, err // Or return false, nil if you consider "doesn't exist" as "not a directory" without error
		}
		// Handle other potential errors (e.g., permission denied)
		return false, err
	}

	return fileInfo.IsDir(), nil
}



func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	chatCmd.Flags().String("histfile", "", "Loads chat history from file")
	chatCmd.Flags().String("histjson", "", "Loads chat history from JSON data")
	chatCmd.MarkFlagsMutuallyExclusive("histfile", "histjson")
}
