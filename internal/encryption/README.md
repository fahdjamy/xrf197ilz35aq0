## Encryption

#### Terms used

- Cipher Block: In cryptography, a cipher block refers to a fundamental unit of data that an encryption algorithm 
  operates on. Most symmetric encryption algorithms (like AES) process data in fixed-size blocks,
  typically 128 bits (16 bytes).
  - The line block, err := aes.NewCipher(key) creates a cipher block using the AES algorithm and the provided key.
  - This block object represents the core encryption/decryption functionality for processing individual 128-bit blocks of data.
- GCM Cipher (Galois/Counter Mode): is an advanced mode of operation for block ciphers that combines confidentiality (encryption) with authentication (integrity check).
  It ensures that data is both encrypted and protected against tampering.
  - aesgcm, err := cipher.NewGCM(block) creates a GCM cipher using the previously created AES cipher block.
  - This aesgcm object handles the entire GCM encryption/decryption process, including generating a unique counter for 
  each block, performing authentication, and combining the encrypted data with the authentication tag.
- Nonce (Number used once): A nonce is a unique value that's used only once for a specific encryption operation.
  It ensures that even if the same plaintext is encrypted multiple times with the same key, the resulting ciphertexts will be different, enhancing security.
  - nonce := make([]byte, aesgcm.NonceSize()) creates a nonce with the appropriate size for the GCM cipher.
  - io.ReadFull(rand.Reader, nonce) fills the nonce with random bytes, ensuring its uniqueness.
  - The nonce is then included in both the encryption (aesgcm.Seal) and decryption processes.


#### Code

- Separate Nonces: We now generate gcmNonce and aadNonce separately. 
- AAD Usage: The aadNonce is used as the additional authenticated data (AAD) in the aesgcm.Seal call. 
- Prepending Nonces: Both nonces are prepended to the ciphertext before returning it from the Encrypt function. 
- Decrypt Function: A Decrypt function is added to demonstrate how to retrieve the nonces and decrypt the ciphertext. 
