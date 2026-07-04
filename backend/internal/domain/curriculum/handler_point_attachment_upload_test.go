package curriculum

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"
)

func multipartFileHeaderForTest(t *testing.T, filename, headerContentType string, data []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	if headerContentType != "" {
		partHeader.Set("Content-Type", headerContentType)
	}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("CreatePart() error = %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("part.Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/upload", &body)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(body.Len()) + 1024); err != nil {
		t.Fatalf("ParseMultipartForm() error = %v", err)
	}
	files := req.MultipartForm.File["file"]
	if len(files) != 1 {
		t.Fatalf("files len = %d, want 1", len(files))
	}
	return files[0]
}

func TestDetectUploadedAttachmentContentTypePrefersSniffedContent(t *testing.T) {
	fileHeader := multipartFileHeaderForTest(t, "not-image.png", "image/png", []byte("plain text, not a png"))
	contentType, err := detectUploadedAttachmentContentType(fileHeader)
	if err != nil {
		t.Fatalf("detectUploadedAttachmentContentType() error = %v", err)
	}
	if contentType != "text/plain; charset=utf-8" {
		t.Fatalf("contentType = %q, want sniffed text/plain", contentType)
	}
}

func TestValidateResearchMaterialUploadContentRejectsSpoofedImageMIME(t *testing.T) {
	if err := validateResearchMaterialUploadContent("not-image.png", "image", "text/plain; charset=utf-8"); err == nil {
		t.Fatal("validateResearchMaterialUploadContent() error = nil, want error")
	}
}
