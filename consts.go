package mailer

type Encoding string
type ContentType string
type MimeVersion string
type Charset string

const (
	EncodingBase64 Encoding = "base64"

	ContentTypeJSON          ContentType = "application/json"
	ContentTypeJavaScript    ContentType = "application/javascript"
	ContentTypeXML           ContentType = "application/xml"
	ContentTypeTextXML       ContentType = "text/xml"
	ContentTypeForm          ContentType = "application/x-www-form-urlencoded"
	ContentTypeProtobuf      ContentType = "application/protobuf"
	ContentTypeMsgpack       ContentType = "application/msgpack"
	ContentTypeTextHTML      ContentType = "text/html"
	ContentTypeTextPlain     ContentType = "text/plain"
	ContentTypeMultipartForm ContentType = "multipart/form-data"
	ContentTypeOctetStream   ContentType = "application/octet-stream"

	MimeVersion1 MimeVersion = "1.0"
	CharsetUTF8  Charset     = "UTF-8"
)
