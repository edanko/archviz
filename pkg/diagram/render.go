package diagram

import (
	"context"
	"runtime"
	"sync"

	"github.com/go-faster/errors"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2plugin"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2target"
	"oss.terrastruct.com/d2/d2themes"
	"oss.terrastruct.com/d2/lib/imgbundler"
	"oss.terrastruct.com/d2/lib/png"
	"oss.terrastruct.com/d2/lib/simplelog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

// postProcessor is a function that post-processes an SVG.
type postProcessor func(context.Context, []byte) ([]byte, error)

// RenderOption is a function that applies an option to a render.
type RenderOption func(*Render) error

// Render is a diagram renderer.
type Render struct {
	compileOpts *d2lib.CompileOptions
	renderOpts  *d2svg.RenderOpts

	postProcessors []postProcessor
}

// NewRender creates a new render.
func NewRender(opts ...RenderOption) (*Render, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ruler")
	}

	r := &Render{
		compileOpts: &d2lib.CompileOptions{
			LayoutResolver: func(_ string) (d2graph.LayoutGraph, error) {
				return d2dagrelayout.DefaultLayout, nil
			},
			RouterResolver: func(_ string) (d2graph.RouteEdges, error) {
				return d2layouts.DefaultRouter, nil
			},
			Ruler: ruler,
		},
		renderOpts: &d2svg.RenderOpts{
			// Center:   go2.Pointer(true),
			NoXMLTag: go2.Pointer(true),
		},
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Render renders the diagram to SVG.
func (r *Render) Render(
	ctx context.Context,
	script string,
) ([]byte, error) {
	diagram, _, err := d2lib.Compile(ctx, script, r.compileOpts, r.renderOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile diagram")
	}

	svg, err := d2svg.Render(diagram, r.renderOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render diagram to svg")
	}

	for _, postProcessor := range r.postProcessors {
		svg, err = postProcessor(ctx, svg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to post process svg")
		}
	}

	return svg, nil
}

// PNG renders the diagram to PNG.
// Uses playwright under the hood to render.
func (r *Render) PNG(ctx context.Context, script string) ([]byte, error) {
	pw, err := png.InitPlaywright()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init playwright")
	}
	defer func() {
		err = pw.Cleanup()
		if err != nil {
			err = errors.Wrap(err, "failed to cleanup playwright")
		}
	}()

	svg, err := r.Render(ctx, script)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render diagram to svg")
	}

	pngBytes, err := png.ConvertSVG(pw.Page, svg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert svg to png")
	}

	return pngBytes, err
}

// WithSketch sets whether the diagram should be rendered as a sketch.
// Ref: https://d2lang.com/tour/sketch.
func WithSketch(sketch bool) RenderOption {
	return func(r *Render) error {
		r.renderOpts.Sketch = &sketch
		// ???
		// d.renderOpts.Font = string(d2fonts.HandDrawn)
		// d.compileOpts.FontFamily = go2.Pointer(d2fonts.HandDrawn)
		return nil
	}
}

// WithTheme sets the theme from the d2themes catalog.
// Ref: https://d2lang.com/tour/themes.
func WithTheme(theme *d2themes.Theme) RenderOption {
	return func(r *Render) error {
		r.renderOpts.ThemeID = &theme.ID
		return nil
	}
}

// WithThemeOverrides overrides colors from the theme.
func WithThemeOverrides(themeOverrides *d2target.ThemeOverrides) RenderOption {
	return func(r *Render) error {
		r.renderOpts.ThemeOverrides = themeOverrides
		return nil
	}
}

// WithPadding sets how much padding will be used around the diagram.
func WithPadding(padding int64) RenderOption {
	return func(r *Render) error {
		r.renderOpts.Pad = &padding
		return nil
	}
}

// WithELKLayoutEngine uses the ELK layout engine.
// Ref: https://d2lang.com/tour/elk.
func WithELKLayoutEngine() RenderOption {
	return func(r *Render) error {
		r.compileOpts.Layout = go2.Pointer("elk")
		r.compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
			return d2elklayout.DefaultLayout, nil
		}
		r.compileOpts.RouterResolver = func(_ string) (d2graph.RouteEdges, error) {
			return d2layouts.DefaultRouter, nil
		}
		return nil
	}
}

// WithTALALayoutEngine uses the TALA layout engine.
// TALA is a Proprietary layout engine developed by Terrastruct,
// designed specifically for software architecture diagrams.
// Ref: https://d2lang.com/tour/tala.
// See https://github.com/terrastruct/TALA for install instructions.
func WithTALALayoutEngine() RenderOption {
	return func(r *Render) error {
		plugins, err := d2plugin.ListPlugins(context.Background())
		if err != nil {
			return errors.Wrap(err, "failed to list d2 plugins")
		}
		talaPlugin, err := d2plugin.FindPlugin(context.Background(), plugins, "tala")
		if err != nil {
			return errors.Wrap(err, "failed to find TALA plugin")
		}

		r.compileOpts.Layout = go2.Pointer("tala")
		r.compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
			return talaPlugin.Layout, nil
		}
		r.compileOpts.RouterResolver = func(_ string) (d2graph.RouteEdges, error) {
			return talaPlugin.(d2plugin.RoutingPlugin).RouteEdges, nil
		}

		r.postProcessors = append(r.postProcessors, talaPlugin.PostProcess)

		return nil
	}
}

// WithImageBundler enables image bundling.
func WithImageBundler(cacheImages bool) RenderOption {
	return func(r *Render) error {
		postProcessor := func(ctx context.Context, svg []byte) ([]byte, error) {
			logger := simplelog.FromLibLog(ctx)
			// NOTE: input path is ignored for a local bundle.
			svg, err := imgbundler.BundleLocal(ctx, logger, "asdf", svg, cacheImages)
			if err != nil {
				return nil, err
			}
			svg, err = imgbundler.BundleRemote(ctx, logger, svg, cacheImages)
			if err != nil {
				return nil, err
			}

			return svg, nil
		}
		r.postProcessors = append(r.postProcessors, postProcessor)
		return nil
	}
}

// minifierConcurrency is the number of goroutines can be used to minify SVGs.
var minifierConcurrency = runtime.GOMAXPROCS(0)

// minifierSemaphore is used to limit the number of minifier goroutines.
var minifierSemaphore chan struct{}

var initMinifier = sync.OnceValue(func() *minify.M {
	minifierSemaphore = make(chan struct{}, minifierConcurrency)

	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	return m
})

// WithMinifier enables SVG minification.
func WithMinifier() RenderOption {
	m := initMinifier()
	return withMinifier(m, minifierSemaphore)
}

// minifiedResult is the result of minifying SVG.
type minifiedResult struct {
	data []byte
	err  error
}

// withMinifier enables SVG minification with custom minifier.
// Useful for testing.
func withMinifier(m *minify.M, sem chan struct{}) RenderOption {
	return func(r *Render) error {
		postProcessor := func(ctx context.Context, input []byte) ([]byte, error) {
			resultCh := make(chan minifiedResult, 1)

			select {
			case sem <- struct{}{}:
				// proceed
			case <-ctx.Done():
				return nil, ctx.Err()
			}

			go func() {
				defer func() { <-sem }()

				minified, err := m.Bytes("image/svg+xml", input)
				resultCh <- minifiedResult{
					data: minified,
					err:  err,
				}
			}()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case res := <-resultCh:
				if res.err != nil {
					return nil, errors.Wrap(res.err, "failed to minify svg")
				}
				return res.data, nil
			}
		}

		r.postProcessors = append(r.postProcessors, postProcessor)
		return nil
	}
}

func WithAppendix() RenderOption {
	return func(r *Render) error {
		postProcessor := func(_ context.Context, _ []byte) ([]byte, error) {
			// TODO: implement when PDF export is implemented.
			// return appendix.Append()
			panic("not implemented")
		}

		r.postProcessors = append(r.postProcessors, postProcessor)
		return nil
	}
}
