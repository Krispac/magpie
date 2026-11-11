package codexcat

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/yetone/magpie/internal/catalog"
)

// A Claude model's entry says multi-agent V2 and lists the clock and
// request_user_input_async, as Codex's own GPT-6 entries do, so Codex hands
// it the subagent tools, the clock and request_user_input_async. It sets no
// tool_mode: Claude gets the tools directly. Other models made an entry
// from the same template are left as they were.
func TestClaudeEntriesCarryAgentTools(t *testing.T) {
	v1Home(t, `{"etag":"W/\"a\"","models":[]}`)
	ms := []catalog.Model{
		{ID: "claude/claude-opus-5-5", Name: "Claude Opus 5.5"},
		{ID: "ollama/glm-5.3:cloud", Name: "GLM 5.3"},
		{ID: "copilot/gpt-6-luna", Name: "GPT-6 Luna"},
	}
	var got struct {
		Models []map[string]any `json:"models"`
	}
	if err := json.Unmarshal(Catalog(ms), &got); err != nil {
		t.Fatal(err)
	}
	for _, m := range got.Models {
		slug := m["slug"].(string)
		tools, _ := m["experimental_supported_tools"].([]any)
		if slug == "claude/claude-opus-5-5" {
			if m["multi_agent_version"] != "v2" || !reflect.DeepEqual(tools, []any{"send_user_message_async", "clock"}) {
				t.Errorf("%s: multi_agent_version %v, tools %v", slug, m["multi_agent_version"], tools)
			}
			if _, ok := m["tool_mode"]; ok {
				t.Errorf("%s: tool_mode %v, want none", slug, m["tool_mode"])
			}
			continue
		}
		if m["multi_agent_version"] != nil || len(tools) != 0 {
			t.Errorf("%s: multi_agent_version %v, tools %v, want neither", slug, m["multi_agent_version"], tools)
		}
	}
}
