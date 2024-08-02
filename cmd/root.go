package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

type State struct {
	Prefix  string
	Threads int
	Tries   int
	ID      string
}

var state State
var writer = uilive.New()

var replacements = map[string]string{
	"0": "a",
	"1": "b",
	"2": "c",
	"3": "d",
	"4": "e",
	"5": "f",
	"6": "g",
	"7": "h",
	"8": "i",
	"9": "j",
	"a": "k",
	"b": "l",
	"c": "m",
	"d": "n",
	"e": "o",
	"f": "p",
}

var cmd = &cobra.Command{
	Use:   "ceid-generator",
	Short: "Generate Google Chrome Extension ID",
	Long:  `A Google Chrome extension ID generator that focuses on finding the first characters.`,
	Run: func(cmd *cobra.Command, args []string) {
		writer.Start()
		state.Tries = 0

		allowedCharacters := "abcdefghijklmnop"

		possibility := (1 / math.Pow(16, float64(len(state.Prefix)))) * 100
		fmt.Printf("Possibility: %f%%\n", possibility)

		estimatedTries := 100 / possibility
		fmt.Printf("Estimated tries: %d\n", int(estimatedTries))

		for _, char := range state.Prefix {
			if !strings.Contains(allowedCharacters, string(char)) {
				fmt.Printf("Invalid character: %s\n", string(char))
				os.Exit(1)
			}
		}

		for i := 0; i < state.Threads; i++ {
			go generate()
		}

		select {}
	},
}

func generate() error {
	if state.ID != "" {
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		return err
	}

	id := calculateChromeID(publicKeyBytes)

	state.Tries++

	fmt.Fprintf(writer, "Tries: %d\n", state.Tries)

	if strings.HasPrefix(id, state.Prefix) {
		writer.Stop()
		state.ID = id
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Public Key: %s\n", base64.RawStdEncoding.EncodeToString(publicKeyBytes))

		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})

		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		os.WriteFile("private.pem", privateKeyPEM, 0644)
		os.WriteFile("public.pem", publicKeyPEM, 0644)

		os.Exit(0)
	} else {
		generate()
	}

	return nil
}

func calculateChromeID(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	sum := hex.EncodeToString(hash[:])

	id := sum[:32]

	for k, v := range replacements {
		id = strings.ReplaceAll(id, k, v)
	}

	return id
}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().StringVarP(&state.Prefix, "prefix", "p", "", "Prefix to find")
	cmd.Flags().IntVarP(&state.Threads, "threads", "t", 1, "Number of threads to use")
}
