package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/streamingaead"
)

func main() {
	// A keyset created with "tinkey create-keyset --key-template=AES256_CTR_HMAC_SHA256_1MB". Note
	// that this keyset has the secret key information in cleartext.

	inPath := flag.String("file", "", "input file path")
	outPath := flag.String("output", "", "output file path")
	keysetPath := flag.String("keyset", "", "json key set path")

	flag.Parse()

	jsonKeysetFile, err := os.Open(*keysetPath)
	if err != nil {
		log.Fatal(err)
	}

	keysetHandle, err := insecurecleartextkeyset.Read(
		keyset.NewJSONReader(jsonKeysetFile))

	if err != nil {
		log.Fatal(err)
	}

	primitive, _ := streamingaead.New(keysetHandle)

	inputFile, _ := os.Open(*inPath)

	encryptedFile, err := os.Create(*outPath)

	if err != nil {
		log.Fatal(err)
	}

	w, err := primitive.NewEncryptingWriter(encryptedFile, []byte{})

	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(w, inputFile); err != nil {
		log.Fatal(err)
	}

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	if err := encryptedFile.Close(); err != nil {
		log.Fatal(err)
	}
	if err := inputFile.Close(); err != nil {
		log.Fatal(err)
	}
}
