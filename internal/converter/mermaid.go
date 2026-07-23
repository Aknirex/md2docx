package converter

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// --- MermaidInkRenderer: uses the public mermaid.ink API ---

const defaultMermaidInkURL = "https://mermaid.ink"

// MermaidInkRenderer renders mermaid diagrams via the mermaid.ink HTTP API.
// Zero external dependencies — only requires network access.
type MermaidInkRenderer struct {
	BaseURL string // defaults to https://mermaid.ink
	Theme   string // mermaid theme: "default", "neutral", "dark", "forest"; defaults to "default"
}

// Render implements MermaidRenderer by calling the mermaid.ink API.
func (r *MermaidInkRenderer) Render(diagram string) ([]byte, int, int, error) {
	baseURL := r.BaseURL
	if baseURL == "" {
		baseURL = defaultMermaidInkURL
	}
	theme := r.Theme
	if theme == "" {
		theme = "default"
	}

	encoded := encodeMermaidInk(diagram, theme)
	url := fmt.Sprintf("%s/img/%s?type=png", baseURL, encoded)

	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("mermaid.ink request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, 0, 0, fmt.Errorf("mermaid.ink returned %d: %s", resp.StatusCode, string(body))
	}

	pngBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("reading mermaid.ink response: %w", err)
	}

	w, h, err := readPNGDimensions(pngBytes)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("decoding mermaid PNG dimensions: %w", err)
	}

	return pngBytes, w, h, nil
}

// encodeMermaidInk encodes a mermaid diagram for the mermaid.ink API.
// Format: JSON → raw deflate → base64url (no padding).
func encodeMermaidInk(diagram string, theme string) string {
	payload := map[string]interface{}{
		"code": diagram,
		"mermaid": map[string]interface{}{
			"theme": theme,
		},
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.BestCompression)
	w.Write(jsonBytes)
	w.Close()

	return base64.RawURLEncoding.EncodeToString(buf.Bytes())
}

// readPNGDimensions reads width and height from PNG IHDR chunk without full decoding.
// PNG format: 8-byte signature, then chunks of [4-byte length][4-byte type][data][4-byte CRC].
// The IHDR chunk is always first and contains width (bytes 0-3) and height (bytes 4-7) as big-endian uint32.
func readPNGDimensions(data []byte) (int, int, error) {
	if len(data) < 24 {
		return 0, 0, fmt.Errorf("png data too short: %d bytes", len(data))
	}
	// Verify PNG signature
	if !bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return 0, 0, fmt.Errorf("invalid png signature")
	}
	// IHDR chunk starts at offset 8: [4-byte length][4-byte "IHDR"]
	chunkType := string(data[12:16])
	if chunkType != "IHDR" {
		return 0, 0, fmt.Errorf("expected IHDR chunk, got %q", chunkType)
	}
	width := binary.BigEndian.Uint32(data[16:20])
	height := binary.BigEndian.Uint32(data[20:24])
	if width == 0 || height == 0 {
		return 0, 0, fmt.Errorf("invalid png dimensions: %dx%d", width, height)
	}
	return int(width), int(height), nil
}

// --- MermaidCLIRenderer: uses a local mermaid-cli (mmdc) installation ---

// MermaidCLIRenderer renders mermaid diagrams via the local mmdc command.
// Requires @mermaid-js/mermaid-cli to be installed (npm i -g @mermaid-js/mermaid-cli).
type MermaidCLIRenderer struct {
	MMDCPath string // path to mmdc; defaults to "mmdc" (looked up from PATH)
	Theme    string // mermaid theme; defaults to "default"
}

// Render implements MermaidRenderer by calling mmdc with --outputFormat png and piping to stdout.
func (r *MermaidCLIRenderer) Render(diagram string) ([]byte, int, int, error) {
	mmdc := r.MMDCPath
	if mmdc == "" {
		mmdc = "mmdc"
	}
	theme := r.Theme
	if theme == "" {
		theme = "default"
	}

	cmd := exec.Command(mmdc,
		"--theme", theme,
		"--outputFormat", "png",
		"--backgroundColor", "white",
		"-", // read from stdin
	)
	cmd.Stdin = strings.NewReader(diagram)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, 0, 0, fmt.Errorf("mmdc failed: %w\nstderr: %s", err, stderr.String())
	}

	pngBytes := stdout.Bytes()
	if len(pngBytes) == 0 {
		return nil, 0, 0, fmt.Errorf("mmdc produced empty output")
	}

	w, h, err := readPNGDimensions(pngBytes)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("decoding mmdc PNG dimensions: %w", err)
	}

	return pngBytes, w, h, nil
}
