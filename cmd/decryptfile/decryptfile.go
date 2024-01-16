package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/streamingaead"
)

func main() {

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

	// Retrieve the StreamingAEAD primitive we want to use from the keyset handle.
	primitive, err := streamingaead.New(keysetHandle)
	if err != nil {
		log.Fatal(err)
	}

	encryptedFile, err := os.Open(*inPath)
	if err != nil {
		log.Fatal(err)
	}

	decryptedPath := filepath.Join(*outPath)
	decryptedFile, err := os.Create(decryptedPath)
	if err != nil {
		log.Fatal(err)
	}

	r, err := primitive.NewDecryptingReader(encryptedFile, []byte{})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(decryptedFile, r); err != nil {
		log.Fatal(err)
	}
	if err := decryptedFile.Close(); err != nil {
		log.Fatal(err)
	}
	if err := encryptedFile.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished")
}
