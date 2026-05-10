package services

import (
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input := "**Teks Tebal**\n- Item 1\n- Item 2"

	output, err := RenderMarkdown(input)
	if err != nil {
		t.Fatalf("Diharapkan tidak ada error, tapi mendapat: %v", err)
	}

	if !strings.Contains(output, "<strong>Teks Tebal</strong>") {
		t.Errorf("Hasil HTML tidak mengandung tag strong. Output: %s", output)
	}
	if !strings.Contains(output, "<li>Item 1</li>") {
		t.Errorf("Hasil HTML tidak mengandung tag li. Output: %s", output)
	}
}
