package builders

import "github.com/bwmarrin/discordgo"

// A FileUpload is builder of file upload
type FileUpload struct {
	fileUpload *discordgo.FileUpload
}

// FileUploadBuilder creates a new file upload
func FileUploadBuilder() *FileUpload {
	return &FileUpload{
		fileUpload: &discordgo.FileUpload{},
	}
}

// SetCustomID sets a custom id
func (f *FileUpload) SetCustomID(customID string) *FileUpload {
	f.fileUpload.CustomID = customID
	return f
}

// SetMinValues sets a min values
func (f *FileUpload) SetMinValues(minValues int) *FileUpload {
	f.fileUpload.MinValues = &minValues
	return f
}

// SetMaxValues sets a max values
func (f *FileUpload) SetMaxValues(maxValues int) *FileUpload {
	f.fileUpload.MaxValues = maxValues
	return f
}

// SetRequired sets required
func (f *FileUpload) SetRequired(required bool) *FileUpload {
	f.fileUpload.Required = &required
	return f
}

// Build returns a discordgo.FileUpload(discordgo.MessageComponent)
func (f *FileUpload) Build() discordgo.MessageComponent {
	return f.fileUpload
}
