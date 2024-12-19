package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

type FontCache map[string]*truetype.Font

func (fc FontCache) Store(fd draw2d.FontData, f *truetype.Font) {
	fc[fd.Name] = f
}

func (fc FontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s not found", fd.Name)
	}
	return font, nil
}
