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
	"encoding/json"
	"fmt"
	"log"

	"github.com/minio/cli"
	"github.com/minio/xfile/lib"
)

var xCmd = cli.Command{
	Name:   "x",
	Usage:  "Extract meta information from an object.",
	Action: mainX,
	CustomHelpTemplate: `NAME:
xfile {{.Name}} - {{.Usage}}

USAGE:
   xfile {{.Name}} [INPUT_URL]

FLAGS:
   {{range .Flags}}{{.}}
{{end}}
`,
}

// X command - eXtract meta information from an object
func mainX(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		cli.ShowCommandHelpAndExit(ctx, "x", 1)
	}

	// Init new xfile engine
	x := xfile.New()
	// Fetch meta information from the URL
	fileType, meta, err := x.ExtractMetadata(ctx.Args().First())
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}

	log.Println("File type found: " + fileType)
	log.Println("Extracted meta information: ")

	// Print extracted meta info
	if ctx.Bool("json") {
		// Print json output when json flag is enabled
		output, err := json.Marshal(meta)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(output))
	} else {
		// Print normal output
		log.Print("Desc:\n")
		for _, d := range meta.TextDescs {
			log.Print("- " + d + "\n")
		}
		log.Print("Keywords:\n")
		for _, k := range meta.Keywords {
			log.Print("- " + k + "\n")
		}
		log.Print("Meta:\n")
		for _, m := range meta.Metas {
			output, err := json.Marshal(m)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("- " + string(output))
		}
	}
}
