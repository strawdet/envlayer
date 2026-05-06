// Package encryptor provides symmetric encryption and decryption of
// environment variable values using AES-256-GCM.
//
// Values are encrypted with a passphrase-derived key and stored as
// base64-encoded strings, making them safe to commit to version control
// when the passphrase is kept secret.
//
// Typical usage:
//
//	e := encryptor.New(os.Getenv("ENVLAYER_PASSPHRASE"))
//
//	// Encrypt a single value
//	enc, err := e.Encrypt("my-secret")
//
//	// Decrypt a single value
//	plain, err := e.Decrypt(enc)
//
//	// Encrypt all values in a merged env map
//	encryptedMap, err := e.EncryptMap(mergedEnv)
//
//	// Decrypt all values before exporting
//	decryptedMap, err := e.DecryptMap(encryptedMap)
package encryptor
