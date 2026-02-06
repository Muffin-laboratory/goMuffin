package builders

import "github.com/bwmarrin/discordgo"

// A Thumbnail is builder for thumbnail
type Thumbnail struct {
	thumbnail *discordgo.Thumbnail
}

// ThumbnailBuilder creates a new Thumbnail
func ThumbnailBuilder(mediaURL string) *Thumbnail {
	return &Thumbnail{
		thumbnail: &discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: mediaURL,
			},
		},
	}
}

// SetMedia sets media
func (t *Thumbnail) SetMedia(mediaURL string) *Thumbnail {
	t.thumbnail.Media.URL = mediaURL
	return t
}

// SetDescription sets description
func (t *Thumbnail) SetDescription(description string) *Thumbnail {
	t.thumbnail.Description = &description
	return t
}

// SetSpoiler sets spoiler
func (t *Thumbnail) SetSpoiler(spoiler bool) *Thumbnail {
	t.thumbnail.Spoiler = spoiler
	return t
}

// Build returns discordgo.Thumbnail(discordgo.MessageComponent)
func (t *Thumbnail) Build() discordgo.MessageComponent {
	return t.thumbnail
}
