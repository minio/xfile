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

import "C"

import (
	"bufio"
	"image"
	"io"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
	"github.com/jdeng/gomxnet"
)

type mxnetInfo struct{}

func (u mxnetInfo) GetName() string {
	return "mxnetInfo"
}

func (u mxnetInfo) GetSupportedMIMEs() []string {
	return []string{"image/jpeg", "image/png"}
}

func (u mxnetInfo) Inspect(contentType string, r io.Reader) (string, []string, interface{}, error) {

	var batch uint32 = 1
	img, _, _ := image.Decode(r)
	img = imaging.Fill(img, 224, 224, imaging.Center, imaging.Lanczos)

	symbol, err := ioutil.ReadFile("./Inception-symbol.json")
	if err != nil {
		return "", nil, nil, err
	}
	params, err := ioutil.ReadFile("./Inception-0009.params")
	if err != nil {
		return "", nil, nil, err
	}
	synset, err := os.Open("./synset.txt")
	if err != nil {
		return "", nil, nil, err
	}

	pred, err := gomxnet.NewPredictor(
		gomxnet.Model{symbol, params},
		gomxnet.Device{gomxnet.CPU, 0},
		[]gomxnet.InputNode{{"data", []uint32{batch, 3, 224, 224}}},
	)

	if err != nil {
		return "", nil, nil, err
	}

	input, _ := gomxnet.InputFrom([]image.Image{img}, gomxnet.ImageMean{117.0, 117.0, 117.0})
	pred.Forward("data", input)
	output, _ := pred.GetOutput(0)
	pred.Free()

	dict := []string{}
	scanner := bufio.NewScanner(synset)
	for scanner.Scan() {
		dict = append(dict, scanner.Text())
	}

	keywords := make([]string, 0)
	outputLen := uint32(len(output)) / batch
	var b uint32 = 0
	for ; b < batch; b++ {
		out := output[b*outputLen : (b+1)*outputLen]
		index := make([]int, len(out))
		gomxnet.Argsort(out, index)

		for i := 0; i < 3; i++ {
			// fmt.Printf("%d: %f, %d, %s\n", i, out[i], index[i], dict[index[i]])
			keywords = append(keywords, dict[index[i]])
		}

	}

	return contentType, keywords, nil, nil
}
