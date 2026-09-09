// chatinfo dumps server-side membership for a chat (diagnostic only).
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kongesque/line-cli/internal/session"
	"github.com/kongesque/line-cli/pkg/line"
)

func main() {
	unlock, err := session.Lock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lock:", err)
		os.Exit(1)
	}
	defer unlock()
	m := session.NewManager(session.KeychainStore{})
	err = m.Do(func(api session.API) error {
		resp, err := api.GetChats([]string{os.Args[1]}, true, true)
		if err != nil {
			return err
		}
		for _, c := range resp.Chats {
			out := map[string]any{"chatMid": c.ChatMid, "chatName": c.ChatName, "type": c.Type}
			if g := c.Extra.GroupExtra; g != nil {
				out["creator"] = g.CreatorMid
				out["memberMids"] = g.MemberMids
				out["inviteeMids"] = g.InviteeMids
			} else {
				out["groupExtra"] = nil
			}
			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "chatinfo:", err)
		os.Exit(1)
	}
	_ = line.Chat{}
}
