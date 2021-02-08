/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
package plugins

import (
	"math/rand"
	"time"
)

func BytesEqual(array1, array2 []byte) bool {
	len1 := len(array1)
	len2 := len(array2)
	if len1 != len2 {
		return false
	}
	for index := 0; index < len1; index++ {
		if array1[index] != array2[index] {
			return false
		}
	}
	return true
}

func BytesCopy(src []byte, srcPos uint, dest []byte, destPos uint, length uint) {
	var index uint
	for index = 0; index < length; index++ {
		dest[destPos + index] = src[srcPos + index]
	}
}

func RandomBytes(size uint) []byte {
	now := time.Now().UnixNano()
	random := rand.New(rand.NewSource(now))
	array := make([]byte, size)
	var index uint
	for index = 0; index < size; index++ {
		array[index] = uint8(random.Intn(256))
	}
	return array
}
