/*
 * Minio Client (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"log"

	cli "github.com/minio/cli"
)

// Help template for xfile.
var xfileHelpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
  {{.Description}}

USAGE:
  xfile {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
  {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}{{if .Flags}}
FLAGS:
  {{range .Flags}}{{.}}
  {{end}}{{end}}

VERSION:
  ` + "0.0.1" +
	`{{ "\n"}}`

var (
	// global flags for xfile.
	globalFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "Show help.",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "Print result in json format.",
		},
	}
)

func Main() {
	// Set up app.
	app := cli.NewApp()
	app.Name = "xfile"
	app.Usage = "Extract meta information from input data"
	app.Description = "Pass the file content to intelligent engine to identify relevant information related to the content."
	app.Flags = globalFlags
	app.Commands = []cli.Command{xCmd}
	app.CustomAppHelpTemplate = xfileHelpTemplate
	app.CommandNotFound = func(ctx *cli.Context, command string) {
		log.Fatalln("Command `" + command + "` not found.")
	}

	app.RunAndExitOnError()
}
