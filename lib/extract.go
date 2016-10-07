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

// ExtractMetadataFromStream - given the file type, call for experts to fetch meta informatin
func (xf XFile) ExtractMetadataFromStream(reader io.Reader, fileType string) (*xfileResponse, error) {
	// Instantiate export module according to the filetype
	experts, ok := xf.mimesExperts[fileType]
	if !ok || len(experts) == 0 {
		return nil, errors.New("`" + fileType + "` mime unsupported")
	}

	var metas []interface{}
	var descs, keys []string

	// Loop over all experts and gather their results
	for _, e := range experts {
		d, k, m, err := e.Inspect(fileType, reader)
		if err != nil {
			return nil, err
		}
		descs = append(descs, d)
		keys = append(keys, k...)
		metas = append(metas, struct {
			ExpertName string
			Meta       interface{}
		}{ExpertName: e.GetName(), Meta: m})
	}

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
	meta, err := xf.ExtractMetadataFromStream(replay, fileType)

	return fileType, meta, err
}
