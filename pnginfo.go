package pnginfo

import (
	"io"
	"strings"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	riimage "github.com/dsoprea/go-utility/v2/image"
)

type PNGInfoItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PNGInfo struct {
	Width          int           `json:"width"`
	Height         int           `json:"height"`
	Model          string        `json:"model"`
	Lora           string        `json:"lora,omitempty"`
	ModelHash      string        `json:"modelHash,omitempty"`
	Version        string        `json:"version,omitempty"`
	Prompt         string        `json:"prompt,omitempty"`
	NegativePrompt string        `json:"negativePrompt,omitempty"`
	Seed           string        `json:"seed,omitempty"`
	Sampler        string        `json:"sampler,omitempty"`
	Steps          string        `json:"steps,omitempty"`
	CFGscale       string        `json:"cfgscale,omitempty"`
	Size           string        `json:"size,omitempty"`
	Values         []PNGInfoItem `json:"values,omitempty"`
	Parameters     string        `json:"parameters,omitempty"`
}

func (info *PNGInfo) DecodeParameters(text string) error {
	info.Parameters = text
	lines := strings.Split(text, "\n")

	var extra string
	if len(lines) >= 2 {
		extra = lines[len(lines)-1]
		lines = lines[0 : len(lines)-1]
	}
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimLeftFunc(line, func(r rune) bool {
			return r == '\u0000'
		}))
		if strings.HasPrefix(line, "Negative prompt:") {
			info.NegativePrompt = strings.TrimSpace(strings.TrimPrefix(line, "Negative prompt:"))
		} else {
			if info.Prompt != "" {
				info.Prompt += "\n"
			}
			info.Prompt += line
			loraPos := strings.Index(line, "<lora:")
			if loraPos >= 0 {
				loraEndPos := strings.Index(line[loraPos:], ">")
				if loraEndPos >= 0 {
					info.Lora = strings.TrimSpace(line[loraPos+6 : loraPos+loraEndPos])
				}
			}
		}
	}
	if extra != "" {
		items := strings.Split(extra, ",")
		for _, item := range items {
			kv := strings.SplitN(item, ":", 2)
			if len(kv) == 2 {
				name := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				info.Values = append(info.Values, PNGInfoItem{
					Name:  name,
					Value: value,
				})
				switch strings.ToLower(name) {
				case "model":
					info.Model = value
				case "seed":
					info.Seed = value
				case "sampler":
					info.Sampler = value
				case "size":
					info.Size = value
				case "cfg scale":
					info.CFGscale = value
				case "steps":
					info.Steps = value
				case "version":
					info.Version = value
				case "model hash":
					info.ModelHash = value
				}
			}
		}
	}
	return nil
}

// Decode PNG exif info from file
func ReadPNGInfoFromFile(f string) (*PNGInfo, error) {
	pmp := pngstructure.NewPngMediaParser()
	intfc, err := pmp.ParseFile(f)
	if err != nil {
		return nil, err
	}
	return readPNGInfo(intfc)
}

// Decode PNG exif info from reader
func ReadPNGInfo(rs io.ReadSeeker, size int) (*PNGInfo, error) {
	pmp := pngstructure.NewPngMediaParser()
	intfc, err := pmp.Parse(rs, size)
	if err != nil {
		return nil, err
	}
	return readPNGInfo(intfc)
}

func readPNGInfo(intfc riimage.MediaContext) (*PNGInfo, error) {
	cs := intfc.(*pngstructure.ChunkSlice)
	index := cs.Index()

	var width, height int
	if ihdrRawSlice, ok := index["IHDR"]; ok {
		cd := pngstructure.NewChunkDecoder()
		if ihdrRaw, err := cd.Decode(ihdrRawSlice[0]); err == nil {
			ihdr := ihdrRaw.(*pngstructure.ChunkIHDR)
			width = int(ihdr.Width)
			height = int(ihdr.Height)
		}
	}
	info := PNGInfo{
		Width:  width,
		Height: height,
	}

	var tEXT []*pngstructure.Chunk
	var ok bool
	if tEXT, ok = index["tEXt"]; !ok {
		tEXT = index["iTXt"]
	}

	for _, raw := range tEXT {
		if (raw.Type == "tEXt" || raw.Type == "iTXt") && raw.Data != nil {
			text := string(raw.Data)
			if !strings.HasPrefix(text, "parameters\x00") {
				continue
			}
			text = strings.TrimPrefix(text, "parameters\x00")
			if info.DecodeParameters(text) == nil {
				break
			}
		}
	}
	return &info, nil
}
