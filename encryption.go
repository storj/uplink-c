// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"reflect"
	"unsafe"

	"storj.io/uplink"
)

// EncryptionKey represents a key for encrypting and decrypting data.
type EncryptionKey struct {
	*uplink.EncryptionKey
}

//export uplink_derive_encryption_key
// uplink_derive_encryption_key derives a salted encryption key for passphrase using the
// salt.
//
// This function is useful for deriving a salted encryption key for users when
// implementing multitenancy in a single app bucket.
func uplink_derive_encryption_key(passphrase *C.uplink_const_char, salt unsafe.Pointer, length C.size_t) C.UplinkEncryptionKeyResult {
	if passphrase == nil {
		return C.UplinkEncryptionKeyResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.UplinkEncryptionKeyResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var goSalt []byte
	hGoSalt := (*reflect.SliceHeader)(unsafe.Pointer(&goSalt))
	hGoSalt.Data = uintptr(salt)
	hGoSalt.Len = ilength
	hGoSalt.Cap = ilength

	encKey, err := uplink.DeriveEncryptionKey(C.GoString(passphrase), goSalt)
	if err != nil {
		return C.UplinkEncryptionKeyResult{
			error: mallocError(err),
		}
	}

	return C.UplinkEncryptionKeyResult{
		encryption_key: (*C.UplinkEncryptionKey)(mallocHandle(universe.Add(&EncryptionKey{encKey}))),
	}
}

//export uplink_free_encryption_key_result
// uplink_free_encryption_key_result frees the resources associated with encryption key.
func uplink_free_encryption_key_result(result C.UplinkEncryptionKeyResult) {
	uplink_free_error(result.error)
	freeEncryptionKey(result.encryption_key)
}

func freeEncryptionKey(encryptionKey *C.UplinkEncryptionKey) {
	if encryptionKey == nil {
		return
	}
	defer C.free(unsafe.Pointer(encryptionKey))
	defer universe.Del(encryptionKey._handle)
}
