// mkgroup creates a LINE group. pkg/line has CreateChat but no CLI verb wires it.
package main

import (
	"fmt"
	"os"

	"github.com/kongesque/line-cli/internal/session"
	"github.com/kongesque/line-cli/pkg/line"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: mkgroup NAME   (creates a group with no other members)")
		os.Exit(2)
	}
	name := os.Args[1]

	unlock, err := session.Lock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lock:", err)
		os.Exit(1)
	}
	defer unlock()

	m := session.NewManager(session.KeychainStore{})
	// Mutate, not Do: Do retries, and a replayed createChat could make two groups.
	err = m.Mutate(func(api session.API) error {
		client, ok := api.(*line.Client)
		if !ok {
			return fmt.Errorf("unexpected client type %T", api)
		}
		chat, err := client.CreateChat([]string{}, name, 0) // Chat.Type 0 = GROUP
		if err != nil {
			return err
		}
		fmt.Printf("created %q  id=%s  type=%d\n", chat.ChatName, chat.ChatMid, chat.Type)
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkgroup:", err)
		os.Exit(1)
	}
}
