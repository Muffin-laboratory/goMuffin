package builders

import "github.com/bwmarrin/discordgo"

type Thumbnail struct {
	*discordgo.Thumbnail
}

func ThumbnailBuilder(mediaURL string) *Thumbnail {
	return &Thumbnail{
		Thumbnail: &discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: mediaURL,
			},
		},
	}
}

func (t *Thumbnail) SetMedia(mediaURL string) *Thumbnail {
	t.Thumbnail.Media.URL = mediaURL
	return t
}

func (t *Thumbnail) SetDescription(description string) *Thumbnail {
	t.Thumbnail.Description = &description
	return t
}

func (t *Thumbnail) SetSpoiler(spoiler bool) *Thumbnail {
	t.Thumbnail.Spoiler = spoiler
	return t
}

func (t *Thumbnail) Build() discordgo.MessageComponent {
	return t.Thumbnail
}
