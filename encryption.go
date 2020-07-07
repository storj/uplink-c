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

//export derive_encryption_key
// derive_encryption_key derives a salted encryption key for passphrase using the
// salt.
//
// This function is useful for deriving a salted encryption key for users when
// implementing multitenancy in a single app bucket.
func derive_encryption_key(passphrase *C.char, salt unsafe.Pointer, length C.size_t) C.EncryptionKeyResult {
	if passphrase == nil {
		return C.EncryptionKeyResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.EncryptionKeyResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var goSalt []byte
	*(*reflect.SliceHeader)(unsafe.Pointer(&goSalt)) = reflect.SliceHeader{
		Data: uintptr(salt),
		Len:  ilength,
		Cap:  ilength,
	}

	encKey, err := uplink.DeriveEncryptionKey(C.GoString(passphrase), goSalt)
	if err != nil {
		return C.EncryptionKeyResult{
			error: mallocError(err),
		}
	}

	return C.EncryptionKeyResult{
		encryption_key: (*C.EncryptionKey)(mallocHandle(universe.Add(&EncryptionKey{encKey}))),
	}
}

//export free_encryption_key_result
// free_encryption_key_result frees the resources associated with encryption key.
func free_encryption_key_result(result C.EncryptionKeyResult) {
	free_error(result.error)
	freeEncryptionKey(result.encryption_key)
}

func freeEncryptionKey(encryptionKey *C.EncryptionKey) {
	if encryptionKey == nil {
		return
	}
	defer C.free(unsafe.Pointer(encryptionKey))
	defer universe.Del(encryptionKey._handle)
}
