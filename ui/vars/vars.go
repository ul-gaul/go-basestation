package vars

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/go-fonts/dejavu/dejavusans"
	"github.com/go-fonts/dejavu/dejavusansbold"
	"github.com/go-fonts/dejavu/dejavusansboldoblique"
	"github.com/go-fonts/dejavu/dejavusansmono"
	"github.com/go-fonts/dejavu/dejavusansmonobold"
	"github.com/go-fonts/dejavu/dejavusansmonoboldoblique"
	"github.com/go-fonts/dejavu/dejavusansmonooblique"
	"github.com/go-fonts/dejavu/dejavusansoblique"
	"github.com/go-fonts/liberation/liberationmonobold"
	"github.com/go-fonts/liberation/liberationmonobolditalic"
	"github.com/go-fonts/liberation/liberationmonoitalic"
	"github.com/go-fonts/liberation/liberationmonoregular"
	"github.com/go-fonts/liberation/liberationsansbold"
	"github.com/go-fonts/liberation/liberationsansbolditalic"
	"github.com/go-fonts/liberation/liberationsansitalic"
	"github.com/go-fonts/liberation/liberationsansregular"
	"github.com/ul-gaul/go-basestation/cfg"
	"github.com/ul-gaul/go-basestation/data"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/utils"
	"sync"
)

const (
	InitialWindowHeight = 600
	InitialWindowWidth  = 1000
	MinWindowHeight     = 420
	MinWindowWidth      = 840
	AppTitle            = "GAUL - Ground Station"
)

var (
	Theme  = material.NewTheme(append(dejavu(), append(liberation(), gofont.Collection()...)...))
	Window = app.NewWindow(
		app.Title(AppTitle),
		app.MinSize(unit.Dp(MinWindowWidth), unit.Dp(MinWindowHeight)),
		app.Size(unit.Dp(InitialWindowWidth), unit.Dp(InitialWindowHeight)),
	)

	DataRoller *data.Roller
	once sync.Once
)

func Initialize() {
	once.Do(func() {
		Theme.TextSize = Theme.TextSize.Scale(0.9)
		
		DataRoller = data.NewRoller(
			cfg.Frontend.DataRoller.MinTimeGap,
			cfg.Frontend.DataRoller.Limit)
		manager.AddDataHandler(DataRoller)
	})
}

func dejavu() []text.FontFace {
	var collection []text.FontFace

	register := func(fnt text.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		utils.CheckErr(err)
		fnt.Typeface = "Dejavu"
		collection = append(collection, text.FontFace{Font: fnt, Face: face})
	}

	register(text.Font{}, dejavusans.TTF)
	register(text.Font{Style: text.Italic}, dejavusansoblique.TTF)
	register(text.Font{Weight: text.Bold}, dejavusansbold.TTF)
	register(text.Font{Style: text.Italic, Weight: text.Bold}, dejavusansboldoblique.TTF)
	register(text.Font{Variant: "Mono"}, dejavusansmono.TTF)
	register(text.Font{Variant: "Mono", Style: text.Italic}, dejavusansmonooblique.TTF)
	register(text.Font{Variant: "Mono", Weight: text.Bold}, dejavusansmonobold.TTF)
	register(text.Font{Variant: "Mono", Style: text.Italic, Weight: text.Bold}, dejavusansmonoboldoblique.TTF)

	return collection
}

func liberation() []text.FontFace {
	var collection []text.FontFace

	register := func(fnt text.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		utils.CheckErr(err)
		fnt.Typeface = "Liberation"
		collection = append(collection, text.FontFace{Font: fnt, Face: face})
	}

	register(text.Font{}, liberationsansregular.TTF)
	register(text.Font{Style: text.Italic}, liberationsansitalic.TTF)
	register(text.Font{Weight: text.Bold}, liberationsansbold.TTF)
	register(text.Font{Style: text.Italic, Weight: text.Bold}, liberationsansbolditalic.TTF)
	register(text.Font{Variant: "Mono"}, liberationmonoregular.TTF)
	register(text.Font{Variant: "Mono", Style: text.Italic}, liberationmonoitalic.TTF)
	register(text.Font{Variant: "Mono", Weight: text.Bold}, liberationmonobold.TTF)
	register(text.Font{Variant: "Mono", Style: text.Italic, Weight: text.Bold}, liberationmonobolditalic.TTF)

	return collection
}
