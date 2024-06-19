package shell

import (
	"github.com/c-bata/go-prompt"
	"github.com/mattn/go-shellwords"
)

// CursorState represents the state of the cursor in the prompt.
type CursorState struct {
	textBefore    string
	currentWord   string
	previousWords []string
}

// Update updates the state of the cursor based on the current prompt document.
func (c *CursorState) Update(in prompt.Document) {
	c.textBefore = in.TextBeforeCursor()
	c.currentWord = in.GetWordBeforeCursor() + in.GetWordAfterCursor()
	words, _ := shellwords.Parse(c.textBefore + "#")
	c.previousWords = words[:len(words)-1]
}

// CurrentWord returns the current word under the cursor.
func (c *CursorState) CurrentWord() string {
	return c.currentWord
}

// PreviousWords returns the words before the current word under the cursor.
func (c *CursorState) PreviousWords() []string {
	return c.previousWords
}

// PreviousWordsN returns the number of words before the current word under the cursor.
func (c *CursorState) PreviousWordsN() int {
	return len(c.previousWords)
}

// PreviousWordsAre returns true if the words before the current word under the cursor match the given arguments.
func (c *CursorState) PreviousWordsAre(words ...string) bool {
	if len(c.previousWords) != len(words) {
		return false
	}
	for i, v := range words {
		if c.previousWords[i] != v {
			return false
		}
	}
	return true
}
