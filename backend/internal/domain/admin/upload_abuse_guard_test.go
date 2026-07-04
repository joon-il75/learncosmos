package admin

import (
	"encoding/binary"
	"strings"
	"testing"
)

func fakeVP8XWebP(width, height int) []byte {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	payload := make([]byte, 10)
	w := width - 1
	h := height - 1
	payload[4] = byte(w)
	payload[5] = byte(w >> 8)
	payload[6] = byte(w >> 16)
	payload[7] = byte(h)
	payload[8] = byte(h >> 8)
	payload[9] = byte(h >> 16)

	out := make([]byte, 0, 30)
	out = append(out, []byte("RIFF")...)
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(4+8+len(payload)))
	out = append(out, size...)
	out = append(out, []byte("WEBP")...)
	out = append(out, []byte("VP8X")...)
	binary.LittleEndian.PutUint32(size, uint32(len(payload)))
	out = append(out, size...)
	out = append(out, payload...)
	return out
}

func TestValidateLumiAssetUploadPayloadRejectsNonWebP(t *testing.T) {
	if err := validateLumiAssetUploadPayload([]byte("plain text")); err == nil {
		t.Fatal("validateLumiAssetUploadPayload() error = nil, want invalid webp error")
	}
}

func TestValidateLumiAssetUploadPayloadRejectsHugeDimensions(t *testing.T) {
	err := validateLumiAssetUploadPayload(fakeVP8XWebP(lumiAssetMaxWidth+1, 128))
	if err == nil || !strings.Contains(err.Error(), "dimensions") {
		t.Fatalf("validateLumiAssetUploadPayload() error = %v, want dimensions error", err)
	}
}

func TestValidateLumiAssetUploadPayloadAllowsBoundedWebP(t *testing.T) {
	if err := validateLumiAssetUploadPayload(fakeVP8XWebP(1024, 1024)); err != nil {
		t.Fatalf("validateLumiAssetUploadPayload() error = %v", err)
	}
}

func TestValidatePlanetTextureMapDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
	}{
		{name: "valid atlas", width: 2048, height: 768},
		{name: "not divisible by columns", width: 2049, height: 768, wantErr: true},
		{name: "not divisible by rows", width: 2048, height: 769, wantErr: true},
		{name: "too wide", width: planetTextureMapMaxWidth + 4, height: 768, wantErr: true},
		{name: "too tall", width: 2048, height: planetTextureMapMaxHeight + 3, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlanetTextureMapDimensions(tt.width, tt.height)
			if tt.wantErr && err == nil {
				t.Fatal("validatePlanetTextureMapDimensions() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("validatePlanetTextureMapDimensions() error = %v", err)
			}
		})
	}
}

func TestReadWebPDimensionsRejectsMalformedChunk(t *testing.T) {
	payload := fakeVP8XWebP(64, 64)
	payload[16] = 100
	if _, _, err := readWebPDimensions(payload); err == nil {
		t.Fatal("readWebPDimensions() error = nil, want malformed chunk error")
	}
}
