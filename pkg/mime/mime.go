package mime

var mimeByExt = []map[string]string{
	{"ext": ".bmp", "mimeType": "image/bmp"},
	{"ext": ".gif", "mimeType": "image/gif"},
	{"ext": ".ico", "mimeType": "image/vnd.microsoft.icon"},
	{"ext": ".jpeg", "mimeType": "image/jpeg"},
	{"ext": ".jpg", "mimeType": "image/jpeg"},
	{"ext": ".png", "mimeType": "image/png"},
	{"ext": ".svg", "mimeType": "image/svg+xml"},
	{"ext": ".tif", "mimeType": "image/tiff"},
	{"ext": ".tiff", "mimeType": "image/tiff"},
	{"ext": ".webp", "mimeType": "image/webp"},
}

// GetMimeByExtention definition
func GetMimeByExtension(ext string) string {
	for _, each := range mimeByExt {
		if each["ext"] == ext {
			return each["mimeType"]
		}
	}

	return ""
}
