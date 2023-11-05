package pnginfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPNGInfo(t *testing.T) {
	info, err := ReadPNGInfoFromFile("assets/00002-2416382472.png")
	assert.Nil(t, err)
	assert.Equal(t, info.Width, 1024)
	assert.Equal(t, info.Height, 1024)
	assert.Equal(t, info.Model, "sd_xl_base_1.0_0.9vae")
	assert.Equal(t, info.Version, "v1.6.0")
	assert.Equal(t, info.Steps, "40")
	assert.Equal(t, info.Sampler, "DPM++ 2M Karras")
	assert.Equal(t, info.CFGscale, "7")
	assert.Equal(t, info.Seed, "2416382472")
	assert.Equal(t, info.Size, "1024x1024")
	assert.Equal(t, info.Prompt, "draw a beautiful woman swimming in the water,")
	assert.Equal(t, info.ModelHash, "e6bb9ea85b")

	info2 := PNGInfo{}
	info2.DecodeParameters(`\u0000\u0000<lora:ShenHe_TongRen)-=杭州:0.65>,shenhe_genshin,1girl,(white hair, long hair),blue eyes,hair ornament,white short shirt,
	long sleeves,off shoulder,blue long dress,standing,looking at viewer,(shy),upper body,(cowboy shot,realistic, photorealistic),(masterpiece, best quality, high quality),(colorful),(delicate eyes and face),volumatic light,ray tracing,extremely detailed CG unity 8k wallpaper,indoors,church,paintings,vases,flowers,table,
	Negative prompt: ng_deepnegative_v1_75t,(badhandv4:1.5),(blurry background:1.3),(depth of field:1.7),(holding:2),(worst quality:2),(low quality:2),(normal quality:2),lowres,bad anatomy,bad hands,
	Steps: 40, Sampler: DPM++ 2M Karras, CFG scale: 7, Seed: 824974617, Size: 1024x1024, Model hash: e6bb9ea85b, Model: sd_xl_base_1.0_0.9vae, Version: v1.6.0`)

	assert.Equal(t, info2.Lora, "ShenHe_TongRen)-=杭州:0.65")
	assert.Contains(t, info2.NegativePrompt, "ng_deepnegative_v1_75t,(badhandv4:1.5)")
	assert.Contains(t, info2.Steps, "40")

	info3 := PNGInfo{}
	info3.DecodeParameters(`<lora:>`)

	assert.Equal(t, info3.Lora, "")
}
