package functions

import (
	"encoding/base64"
)
//just a simple Xor encrypt/decrypt function
func Xorify(input string, key string) (string, error) {
	if key == "" {
		return "Error: Key cannot be empty", nil
	}
	keyBytes := []byte(key)

	//attempting to decode input
	data, err := base64.StdEncoding.DecodeString(input)
	decrypting := (err == nil)

	if !decrypting {
		// treats as plaintext input for encryption
		data = []byte(input)
	}

	output := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		output[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}

	if decrypting {
		//decrypted plaintext
		return string(output), nil
	} else {
		// returns encrypted base64
		return base64.StdEncoding.EncodeToString(output), nil
	}
}
//ROT13 encrypt/decrypt function.
func Rot13ify(input string) string {
	output := make([]rune, len(input))
	for i, r := range input {
		switch {
			case r >= 'A' && r <= 'Z':
				output[i] = 'A' + (r-'A'+13)%26
			case r >= 'a' && r <= 'z':
				output[i] = 'a' + (r-'a'+13)%26
			default:
				output[i] = r
		}
	}

	return string(output)
}

