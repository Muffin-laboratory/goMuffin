package builders

import "github.com/bwmarrin/discordgo"

// A FileComponent is builder for file component
type FileComponent struct {
	file *discordgo.FileComponent
}

// FileComponentBuilder creates a new file component
func FileComponentBuilder() *FileComponent {
	return &FileComponent{
		file: &discordgo.FileComponent{},
	}
}

// SetFile sets a file.
// The argument should be `attachment://<filename>` syntax.
func (f *FileComponent) SetFile(file string) *FileComponent {
	f.file.File = discordgo.UnfurledMediaItem{
		URL: file,
	}

	return f
}

// SetSpoiler sets spoiler
func (f *FileComponent) SetSpoiler(spoiler bool) *FileComponent {
	f.file.Spoiler = spoiler
	return f
}

// Build returns a discordgo.FileComponent(discordgo.MessageComponent)
func (f *FileComponent) Build() discordgo.MessageComponent {
	return f.file
}
