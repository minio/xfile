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
	"io"
	"io/ioutil"

	"github.com/rakyll/magicmime"
)

func GuessFileTypeFromName(filename string) (string, error) {
	return "", nil
}

func GuessFileTypeFromStream(r io.Reader) (string, error) {
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		return "", err
	}
	defer magicmime.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}

	mimetype, err := magicmime.TypeByBuffer(buf)
	if err != nil {
		return "", err
	}

	return mimetype, nil
}
