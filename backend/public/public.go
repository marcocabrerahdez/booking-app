package public

import "embed"

/**
 * Embed the frontend app into the binary
 */

//go:embed *
var FrontendApp embed.FS
