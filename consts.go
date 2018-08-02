package mailer

type Encoding string
type ContentType string
type MimeVersion string
type Charset string

const (
	EncodingBase64 Encoding = "base64"

	ContentTypeJSON          = "application/json"
	ContentTypeJavaScript    = "application/javascript"
	ContentTypeXML           = "application/xml"
	ContentTypeTextXML       = "text/xml"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeProtobuf      = "application/protobuf"
	ContentTypeMsgpack       = "application/msgpack"
	ContentTypeTextHTML      = "text/html"
	ContentTypeTextPlain     = "text/plain"
	ContentTypeMultipartForm = "multipart/form-data"
	ContentTypeOctetStream   = "application/octet-stream"

	MimeVersion1 MimeVersion = "1.0"
	CharsetUTF8  Charset     = "UTF-8"
)
