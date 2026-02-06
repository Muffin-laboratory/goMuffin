package builders

import "github.com/bwmarrin/discordgo"

// A MediaGallery is builder for media gallery
type MediaGallery struct {
	mediaGallery *discordgo.MediaGallery
}

// A MediaGalleryItem is item for media gallery
type MediaGalleryItem struct {
	Media       string
	Description string
	Spoiler     bool
}

// MediaGallery creates a new media gallery
func MediaGalleryBuilder() *MediaGallery {
	return &MediaGallery{
		mediaGallery: &discordgo.MediaGallery{},
	}
}

// AddItem adds media gallery items
func (m *MediaGallery) AddItem(items ...MediaGalleryItem) *MediaGallery {
	for _, item := range items {
		m.mediaGallery.Items = append(m.mediaGallery.Items, discordgo.MediaGalleryItem{
			Media: discordgo.UnfurledMediaItem{
				URL: item.Media,
			},
			Description: &item.Description,
			Spoiler:     item.Spoiler,
		})
	}

	return m
}

// Build returns a discordgo.MediaGallery(discordgo.MessageComponent)
func (m *MediaGallery) Build() discordgo.MessageComponent {
	return m.mediaGallery
}
