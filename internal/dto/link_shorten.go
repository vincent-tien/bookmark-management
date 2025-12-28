package dto

const (
	DefaultExpInSeconds = 3600
)

// LinkShortenRequestDto represents request payload for creating a shortened link
//
// swagger:model LinkShortenRequestDto
type LinkShortenRequestDto struct {
	// Expiration time in seconds for the shortened link
	// Must be greater than or equal to 1
	// description: Time-to-live of the shortened URL (TTL)
	// minimum: 1
	// example: 3600
	ExpInSeconds int `json:"exp" binding:"omitempty"`

	// Original URL that will be shortened
	// Must be a valid URL format (http or https)
	//
	// required: true
	// format: url
	// example: https://example.com
	Url string `json:"url" binding:"required,url"`
}

func (req *LinkShortenRequestDto) Prepare() {
	// If 0 (missing or explicitly 0), set to default
	if req.ExpInSeconds == 0 {
		req.ExpInSeconds = DefaultExpInSeconds
	}
}

// LinkShortenResponseDto represents shorten link response
//
// swagger:model LinkShortenResponseDto
type LinkShortenResponseDto struct {
	// Short code of the generated URL
	// example: abc123
	Code string `json:"code"`

	// Success message
	// example: Shorten URL generated successfully!
	Message string `json:"message"`
}
