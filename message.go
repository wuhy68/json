package mailer

type Message struct {
	MimeVersion         MimeVersion
	Headers             map[string]string
	FromName            string
	FromAddr            string
	ToAddrs             []string
	CcAddrs             []string
	BccAddrs            []string
	Date                string
	ReplyToAddr         string
	Subject             string
	BoundaryMixed       string
	BoundaryAlternative string
	ContentType         ContentType
	Charset             Charset
	Body                string
	Attachments         []*Attachment
	builded             []byte
}
