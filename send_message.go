package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
	"time"
)

type SendMessageService struct {
	mailer  *Mailer
	message *Message
}

type Attachment struct {
	FileName       string
	Content        []byte
	Inline         bool
	ContentType    string
	Encoding       Encoding
	EncodedContent string
}

func NewSendMessageService(e *Mailer) *SendMessageService {
	return &SendMessageService{
		mailer: e,
		message: &Message{
			MimeVersion:         MimeVersion1,
			Date:                time.Now().Format(time.RFC1123Z),
			ContentType:         ContentTypeTextPlain,
			Charset:             CharsetUTF8,
			BoundaryMixed:       RandomBoundary(),
			BoundaryAlternative: RandomBoundary(),
			Headers:             make(map[string]string),
		},
	}
}
func (e *SendMessageService) From(name, address string) *SendMessageService {
	e.message.FromName = name
	e.message.FromAddr = address
	return e
}

func (e *SendMessageService) To(to ...string) *SendMessageService {
	e.message.ToAddrs = append(e.message.ToAddrs, to...)
	return e
}

func (e *SendMessageService) Cc(cc ...string) *SendMessageService {
	e.message.CcAddrs = append(e.message.CcAddrs, cc...)
	return e
}

func (e *SendMessageService) Bcc(bcc ...string) *SendMessageService {
	e.message.BccAddrs = append(e.message.BccAddrs, bcc...)
	return e
}

func (e *SendMessageService) Subject(subject string) *SendMessageService {
	e.message.Subject = subject
	return e
}

func (e *SendMessageService) Body(contentType ContentType, body string) *SendMessageService {
	e.message.ContentType = contentType
	e.message.Body = body
	return e
}

func (e *SendMessageService) Date(date time.Time) *SendMessageService {
	e.message.Date = date.Format(time.RFC1123Z)
	return e
}

func (e *SendMessageService) Header(key string, value string) *SendMessageService {
	e.message.Headers[key] = value
	return e
}

func (e *SendMessageService) Attachment(content []byte, inline bool, fileName string) *SendMessageService {
	e.message.Attachments = append(e.message.Attachments, &Attachment{
		ContentType:    GetMimeType(fileName),
		Content:        content,
		Inline:         inline,
		FileName:       fileName,
		Encoding:       EncodingBase64,
		EncodedContent: base64.StdEncoding.EncodeToString(content),
	})

	return e
}

func (e *SendMessageService) Template(path, name string, reload bool) *SendMessageService {
	key := fmt.Sprintf("%s/%s", path, name)

	var result bytes.Buffer
	var err error

	if _, found := templates[key]; !found {
		e.mailer.mux.Lock()
		defer e.mailer.mux.Unlock()
		templates[key], err = ReadFile(key, nil)
		if err != nil {
			log.Error(err)
			return e
		}
	}

	t := template.New(name)
	t, err = t.Parse(string(templates[key]))
	if err == nil {
		if err := t.ExecuteTemplate(&result, name, e.message); err != nil {
			log.Error(err)
			return e
		}

	} else {
		log.Error(err)
		return e
	}

	e.message.builded = result.Bytes()

	return e
}

func (e *SendMessageService) Execute() ([]string, error) {
	if e.message.builded == nil {
		dir, err := os.Getwd()
		if err != nil {
			return []string{}, err
		}
		e.Template(dir+"/templates", "email.template", false)
	}

	return SendMail(
		fmt.Sprintf("%s:%s", e.mailer.config.Host, e.mailer.config.Port),
		e.mailer.auth,
		e.message.FromAddr,
		append(append(e.message.ToAddrs, e.message.CcAddrs...), e.message.BccAddrs...),
		e.message.builded,
	)
}
