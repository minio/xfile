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
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

// getStreamFromURL - get data stream from URI
func getStreamFromURL(fileURI string) (io.ReadCloser, error) {
	url, err := url.Parse(fileURI)
	if err != nil {
		return nil, err
	}
	switch url.Scheme {
	case "http", "https":
		resp, err := http.Get(fileURI)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	case "file", "":
		f, err := os.Open(fileURI)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("URL Not supported")
	}
}
