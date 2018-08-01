package mailer

type Encoding string
type ContentType string
type MimeVersion string
type Charset string

const (
	Base64               Encoding    = "base64"
	ContentTypePlainText ContentType = "text/plain"
	MimeVersion1         MimeVersion = "1.0"
	UTF8                 Charset     = "UTF-8"
)
