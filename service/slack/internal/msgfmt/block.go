// Package msgfmt parses message block from message text.
package msgfmt

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type BlockType int

var SpaceBlock = &Block{
	Type:    TextBlock,
	Content: " ",
}

const (
	TextBlock BlockType = iota
	ChannelBlock
	UserBlock
	CommandBlock
	LinkBlock
)

type Block struct {
	Type    BlockType
	Label   string
	Content string
}

func (b *Block) String() string {
	if b.Label != "" {
		return b.Label
	}

	if b.Type == CommandBlock {
		return b.commandString()
	}

	return b.Content
}

func (b *Block) commandString() string {
	// TODO : support date format, subteam format
	switch b.Content {
	case "!everyone", "!channel", "!hear":
		return b.Content
	}

	return b.Content
}

func (b *Block) URL() (*url.URL, error) {
	if b.Type != LinkBlock {
		return nil, errors.New("the block is not a link block")
	}

	u, err := url.Parse(b.Content)
	if err != nil {
		return nil, fmt.Errorf("the block content is not valid URL: %w", err)
	}

	return u, nil
}

var blockRegexp = regexp.MustCompile("<(.*?)>")

func Parse(text string) []*Block {
	matches := blockRegexp.FindAllStringSubmatchIndex(text, -1)

	var blocks []*Block
	tail := 0
	for _, match := range matches {
		t := text[tail:match[0]]
		if t != "" {
			t := Unescape(t)
			blocks = append(blocks, &Block{
				Type:    TextBlock,
				Content: t,
			})
		}

		tail = match[1]

		fields := strings.SplitN(text[match[2]:match[3]], "|", 2)
		content := fields[0]

		block := &Block{}
		if len(fields) == 2 {
			block.Label = Unescape(fields[1])
		}

		switch {
		case strings.HasPrefix(content, "#C"):
			block.Type = ChannelBlock
			block.Content = content[1:]
		case strings.HasPrefix(content, "@U"):
			block.Type = UserBlock
			block.Content = content[1:]
		case strings.HasPrefix(content, "!"):
			block.Type = CommandBlock
			block.Content = content[1:]
		default:
			block.Type = LinkBlock
		}

		blocks = append(blocks, block)
	}

	if tail < len(text) {
		t := Unescape(text[tail:])
		blocks = append(blocks, &Block{
			Type:    TextBlock,
			Content: t,
		})
	}

	return blocks
}

func Format(blocks ...*Block) string {
	var s strings.Builder
	for _, b := range blocks {
		switch b.Type {
		case TextBlock:
			text := b.Content
			if text == "" {
				text = b.Label
			}
			if text == "" {
				continue
			}
			s.WriteString(Escape(text))
		case ChannelBlock, UserBlock, LinkBlock:
			if b.Content == "" {
				continue
			}
			s.WriteRune('<')
			s.WriteString(b.Content)
			if b.Label != "" {
				s.WriteRune('|')
				s.WriteString(Escape(b.Label))
			}
			s.WriteRune('>')
		case CommandBlock:
			// TODO
		}
	}

	return s.String()
}
