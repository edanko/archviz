package adapters

import (
	"context"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/go-faster/errors"
	"github.com/goccy/go-yaml"
	"github.com/helshabini/fsbroker"
	"github.com/sourcegraph/conc/pool"
	"oss.terrastruct.com/d2/lib/log"

	"github.com/edanko/archviz/cmd/archviz/internal/domain"
)

const (
	ext = ".yaml"
)

// LocalDirectorySource is an in-memory data source for components and contexts.
type LocalDirectorySource struct {
	paths []string
	root  root
	graph *domain.ComponentGraph
	views domain.Views

	isLoaded bool
	mu       sync.Mutex

	watcher  *fsbroker.FSBroker
	onChange func(s *LocalDirectorySource)
}

// NewLocalDirectorySource creates a new LocalDirectorySource.
func NewLocalDirectorySource(
	ctx context.Context,
	paths []string,
) (*LocalDirectorySource, error) {
	ds := &LocalDirectorySource{
		paths: paths,
		mu:    sync.Mutex{},
	}

	ds.onChange = func(s *LocalDirectorySource) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.isLoaded = false
	}

	err := ds.setupWatcher(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup file watcher")
	}

	return ds, nil
}

// setupWatcher sets up the file watcher for the LocalDirectorySource
// and starts it if onChange is set.
func (s *LocalDirectorySource) setupWatcher(ctx context.Context) error {
	config := fsbroker.DefaultFSConfig()
	broker, err := fsbroker.NewFSBroker(config)
	if err != nil {
		return errors.Wrap(err, "failed to create FS Broker")
	}

	for _, path := range s.paths {
		err := broker.AddRecursiveWatch(path)
		if err != nil {
			log.Error(
				ctx,
				"error adding watch",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
		}
	}

	broker.Start()

	go s.runWatcher(ctx, broker)

	return nil
}

// runWatcher runs the file watcher for the LocalDirectorySource.
func (s *LocalDirectorySource) runWatcher(
	ctx context.Context,
	broker *fsbroker.FSBroker,
) {
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-broker.Next():
			// NOTE: we only care about yaml files.
			if filepath.Ext(event.Path) != ext {
				continue
			}

			log.Info(
				ctx,
				"fs event has occurred",
				slog.String("path", event.Path),
				slog.String("type", event.Type.String()),
				// slog.Time("timestamp", event.Timestamp),
				slog.Any("properties", event.Properties),
			)

			s.onChange(s)

		case err := <-broker.Error():
			log.Error(
				ctx,
				"fs event error has occurred",
				slog.String("error", err.Error()),
			)
		}
	}
}

// Close closes the LocalDirectorySource.
func (s *LocalDirectorySource) Close() error {
	if s.watcher != nil {
		s.watcher.Stop()
	}
	return nil
}

// GetGraph loads the data source from the given directories and returns the
// built graph. The loading is done lazily, so if the graph is already loaded,
// this function will return the existing graph.
func (s *LocalDirectorySource) GetGraph(ctx context.Context) (*domain.ComponentGraph, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isLoaded {
		var err error
		s.root, err = load(ctx, s.paths...)
		if err != nil {
			return nil, err
		}
		s.graph = buildGraph(&s.root)
		s.views = buildViews(&s.root)
		s.isLoaded = true
	}
	return s.graph, nil
}

// buildGraph builds a graph from the given root.
func buildGraph(source *root) *domain.ComponentGraph {
	if source == nil {
		return nil
	}

	g := domain.NewComponentGraph()

	for componentID, component := range source.Components {
		g.AddComponent(
			componentID,
			component.Title,
			domain.WithDescription(component.Description),
			domain.WithProperties(map[string]any{
				"class":        component.Class,
				"shape":        component.Shape,
				"technologies": component.Technologies,
			}),
		)
	}

	for _, link := range source.Links {
		if len(link.Via) == 0 {
			g.AddLink(
				link.From,
				link.To,
				link.Title,
			)
			continue
		}

		from := link.From
		to := link.Via[0]
		for _, via := range link.Via[0:] {
			g.AddLink(
				from,
				to,
				link.Title,
			)
			from = to
			to = via
		}
	}

	g.Check()

	return g
}

// ListViews returns a list of views.
func (s *LocalDirectorySource) ListViews(ctx context.Context) (domain.Views, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isLoaded {
		var err error
		s.root, err = load(ctx, s.paths...)
		if err != nil {
			return nil, err
		}
		s.graph = buildGraph(&s.root)
		s.views = buildViews(&s.root)
		s.isLoaded = true
	}
	return s.views, nil
}

// GetView returns the view with the given ID.
func (s *LocalDirectorySource) GetView(ctx context.Context, id string) (domain.View, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isLoaded {
		var err error
		s.root, err = load(ctx, s.paths...)
		if err != nil {
			return domain.View{}, err
		}
		s.graph = buildGraph(&s.root)
		s.views = buildViews(&s.root)
		s.isLoaded = true
	}

	view, exists := s.views[id]
	if !exists {
		return domain.View{}, errors.Errorf("view %s not found", id)
	}

	return view, nil
}

// buildViews builds a map of views from the given root.
func buildViews(source *root) map[string]domain.View {
	if source == nil {
		return nil
	}

	if source.Views == nil {
		return nil
	}

	views := make(map[string]domain.View, len(source.Views))
	for viewID, view := range source.Views {
		views[viewID] = domain.View{
			ID:         viewID,
			Title:      view.Title,
			Location:   view.Location,
			Components: view.Components,
			MaxDepth:   view.Depth,
		}
	}
	return views
}

type component struct {
	Title        string   `yaml:"title"`
	Shape        string   `yaml:"entity"`
	Technologies []string `yaml:"technologies,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Class        string   `yaml:"class,omitempty"`
}

type link struct {
	From  string   `yaml:"from"`
	To    string   `yaml:"to"`
	Title string   `yaml:"title"`
	Via   []string `yaml:"via,omitempty"`
}

type Context struct {
	Title      string   `yaml:"title"`
	Location   string   `yaml:"location"`
	Components []string `yaml:"components"`
	Depth      int      `yaml:"depth"`
}

type root struct {
	Components map[string]component `yaml:"components"`
	Views      map[string]Context   `yaml:"contexts"`
	Links      []link               `yaml:"links"`
}

var rootMu sync.Mutex

func (r *root) mergeWith(other *root) {
	rootMu.Lock()
	defer rootMu.Unlock()
	// should we validate something here?
	maps.Insert(r.Components, maps.All(other.Components))
	maps.Insert(r.Views, maps.All(other.Views))
	r.Links = slices.Concat(r.Links, other.Links)
}

func load(ctx context.Context, paths ...string) (root, error) {
	mainRoot := root{
		Components: make(map[string]component),
		Views:      make(map[string]Context),
	}

	sourcesPool := pool.NewWithResults[root]().WithContext(ctx)

	for _, path := range paths {
		source := os.DirFS(path)
		err := fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ext {
				return nil
			}

			sourcesPool.Go(func(_ context.Context) (root, error) {
				f, err := source.Open(path)
				if err != nil {
					return root{}, err
				}
				defer f.Close()

				var currentRoot root
				err = yaml.NewDecoder(f).Decode(&currentRoot)
				if err != nil {
					return root{}, err
				}
				return currentRoot, nil
			})

			return nil
		})
		if err != nil {
			return root{}, err
		}
	}

	results, err := sourcesPool.Wait()
	if err != nil {
		return root{}, errors.Wrap(err, "failed to load sources")
	}

	for idx := range results {
		mainRoot.mergeWith(&results[idx])
	}
	return mainRoot, nil
}
