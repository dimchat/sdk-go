/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package format

import (
	"encoding/json"
	. "github.com/dimchat/mkm-go/format"
)

type JSONCoder struct {}

func (coder JSONCoder) Init() ObjectCoder {
	return coder
}

//-------- IObjectCoder

func (coder JSONCoder) Encode(object interface{}) string {
	bytes, err := json.Marshal(object)
	if err == nil {
		return StringFromBytes(bytes)
	} else {
		//panic("failed to encode to JsON string")
		return ""
	}
}

func (coder JSONCoder) Decode(str string) interface{} {
	bytes := BytesFromString(str)
	for _, ch := range bytes {
		if ch == '{' {
			// decode to map
			var dict map[string]interface{}
			err := json.Unmarshal(bytes, &dict)
			if err == nil {
				return dict
			} else {
				return nil
			}
		} else if ch == '[' {
			// decode to array
			var array []interface{}
			err := json.Unmarshal(bytes, &array)
			if err == nil {
				return array
			} else {
				return nil
			}
		} else if ch != ' ' && ch != '\t' {
			// error
			break
		}
	}
	//panic(bytes)
	return nil
}
