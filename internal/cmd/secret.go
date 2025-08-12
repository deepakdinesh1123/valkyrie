package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

var GenerateKeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "valkyrie encryption key",
	Long:  `generate valkyrie encryption key`,
	RunE:  genKeyExec,
}

func genKeyExec(cmd *cobra.Command, args []string) error {
	key := make([]byte, 32)

	if _, err := rand.Reader.Read(key); err != nil {
		fmt.Println("error generating random encryption key ", err)
		return err
	}

	encodedHex := hex.EncodeToString(key)
	fmt.Printf("Valkyrie encryption key: %s\n", encodedHex)

	return nil
}
