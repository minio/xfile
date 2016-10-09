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

package xfile

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

func GetVersion() string {
	return "0.0.1"
}

// XFile - xfile module
type XFile struct {
	mimesExperts map[string][]Expert
}

// New() - initialize and return a ready XFile instance
func New() XFile {
	xf := XFile{}
	// Register experts per supported file type
	xf.mimesExperts = make(map[string][]Expert)
	for _, e := range registeredExperts {
		for _, m := range e.GetSupportedMIMEs() {
			xf.mimesExperts[m] = append(xf.mimesExperts[m], e)

		}
	}
	return xf
}

// ExtractMetadataFromStream - given the file type, call for experts to fetch meta information.
func (xf XFile) ExtractMetadataFromStream(input io.Reader, fileType string) (*xfileResponse, error) {

	// Check if there is at least one expert which supports the passed file type
	experts, ok := xf.mimesExperts[fileType]
	if !ok || len(experts) == 0 {
		return nil, errors.New("`" + fileType + "` mime not supported")
	}

	// Setup pipe to transfer file data to all experts
	pipesRd := make([]io.ReadCloser, len(experts))
	pipesWr := make([]io.WriteCloser, len(experts))
	for i := range experts {
		r, w := io.Pipe()
		pipesRd[i] = r
		pipesWr[i] = w
	}

	var metas []interface{}
	var descs, keys []string

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Run experts inspection job in parallel
	for i, e := range experts {
		wg.Add(1)
		go func(idx int, expert Expert) {
			defer wg.Done()
			defer pipesRd[idx].Close()

			// Call for expert inspection
			d, k, m, err := expert.Inspect(fileType, pipesRd[idx])
			if err != nil {
				return
			}

			// Append inspection results
			mu.Lock()
			descs = append(descs, d)
			keys = append(keys, k...)
			if m != nil {
				metas = append(metas, struct {
					ExpertName string
					Meta       interface{}
				}{ExpertName: expert.GetName(), Meta: m})
			}
			mu.Unlock()

		}(i, e)
	}

	errPipes := make([]bool, len(experts))
	eof := false

	// Duplicate incoming file data to all experts. MultiWriter is avoided because it quits on the first encountered error
	for {
		// Read some data from input file
		buf := make([]byte, 4*1024)
		_, err := input.Read(buf)
		switch err {
		case nil:
		case io.EOF:
			eof = true
		default:
			break
		}
		// Distribute input data to all experts
		for i := range experts {
			if errPipes[i] {
				continue
			}
			_, err = pipesWr[i].Write(buf)
			if err != nil {
				errPipes[i] = true
			}
		}
		// Terminate when we don't have anymore data from input file
		if eof {
			break
		}
	}

	// Close all write pipes
	for i := range experts {
		pipesWr[i].Close()
	}

	// Wait until all experts terminate their jobs
	wg.Wait()

	return &xfileResponse{TextDescs: descs, Keywords: keys, Metas: metas}, nil
}

// ExtractMetadata - get meta information from the given file URI
func (xf XFile) ExtractMetadata(fileURI string) (string, *xfileResponse, error) {
	// Get data stream from file URI
	stream, err := getStreamFromURL(fileURI)
	if err != nil {
		return "", nil, err
	}
	defer stream.Close()

	// Setup new reader to save read data in saveBuf
	var saveBuf bytes.Buffer
	teeStream := io.TeeReader(stream, &saveBuf)

	// Guess file type
	fileType, err := GuessFileTypeFromStream(teeStream)
	if err != nil {
		return "", nil, err
	}

	// Reconstruct the original reader
	replay := io.MultiReader(&saveBuf, stream)

	// Run the real extraction work
	meta, err := xf.ExtractMetadataFromStream(replay, fileType)

	return fileType, meta, err
}
