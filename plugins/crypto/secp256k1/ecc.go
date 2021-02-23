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
package secp256k1

/*
#include "ecc.h"
#include "micro-ecc/uECC.c"
*/
import "C"

import (
	"unsafe"
)

/**
 *  Generate key pair
 *
 * @return public key & private key
 */
func Generate() (pub, pri []byte) {
	pub = make([]byte, 64)
	pri = make([]byte, 32)
	pubPtr := (*C.uchar)(unsafe.Pointer(&pub[0]))
	priPtr := (*C.uchar)(unsafe.Pointer(&pri[0]))
	res := C.uECC_make_key(pubPtr, priPtr, C.uECC_secp256k1())
	if res != 1 {
		panic("failed to generate ECC private key")
	}
	return pub, pri
}

/**
 *  Get public key from private key
 *
 * @param pri - private key data (32 bytes)
 * @return public key data (64 bytes)
 */
func GetPublicKey(pri []byte) []byte {
	pub := make([]byte, 64)
	pubPtr := (*C.uchar)(unsafe.Pointer(&pub[0]))
	priPtr := (*C.uchar)(unsafe.Pointer(&pri[0]))
	res := C.uECC_compute_public_key(priPtr, pubPtr, C.uECC_secp256k1())
	if res == 1 {
		return pub
	} else {
		return nil
	}
}

/**
 *  Sign data with private key
 *
 * @param pri - private key data (32 bytes)
 * @param digest - message digest
 * @return signature (64 bytes)
 */
func Sign(pri, digest []byte) []byte {
	sig := make([]byte, 64)
	keyPtr := (*C.uchar)(unsafe.Pointer(&pri[0]))
	digPtr := (*C.uchar)(unsafe.Pointer(&digest[0]))
	sigPtr := (*C.uchar)(unsafe.Pointer(&sig[0]))
	C.uECC_sign(keyPtr, digPtr, C.unsigned(len(digest)), sigPtr, C.uECC_secp256k1())
	return sig
}

/**
 *  Verify data and signature with public key
 *
 * @param pub - public key data (64 bytes)
 * @param digest - message digest
 * @param signature - signature (64 bytes)
 * @return true on matched
 */
func Verify(pub, digest, signature []byte) bool {
	keyPtr := (*C.uchar)(unsafe.Pointer(&pub[0]))
	digPtr := (*C.uchar)(unsafe.Pointer(&digest[0]))
	sigPtr := (*C.uchar)(unsafe.Pointer(&signature[0]))
	return C.uECC_verify(keyPtr, digPtr, C.unsigned(len(digest)), sigPtr, C.uECC_secp256k1()) == 1
}

func SignatureToDER(signature []byte) []byte {
	der := make([]byte, 72)
	sigPtr := (*C.uchar)(unsafe.Pointer(&signature[0]))
	derPtr := (*C.uchar)(unsafe.Pointer(&der[0]))
	cnt := C.ecc_sig_to_der(sigPtr, derPtr)
	return der[:cnt]
}
func SignatureFromDER(der []byte) []byte {
	sig := make([]byte, 64)
	sigPtr := (*C.uchar)(unsafe.Pointer(&sig[0]))
	derPtr := (*C.uchar)(unsafe.Pointer(&der[0]))
	cnt := len(der)
	C.ecc_der_to_sig(derPtr, C.int(cnt), sigPtr)
	return sig
}
