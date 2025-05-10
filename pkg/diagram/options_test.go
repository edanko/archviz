package diagram

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/minify/v2"
)

func TestWithMinifier(t *testing.T) {
	t.Parallel()

	t.Run("successfully minify SVG", func(t *testing.T) {
		t.Parallel()

		// arrange
		d, err := NewDiagram(WithMinifier())
		require.NoError(t, err)

		input := []byte(`<svg xmlns="http://www.w3.org/2000/svg">
			<rect width="100" height="100"/>
		</svg>`)

		ctx := context.Background()
		postProcessor := d.svgPostProcessors[0]

		// act
		result, err := postProcessor(ctx, input)
		require.NoError(t, err)

		// assert
		want := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100"/></svg>`)
		require.Equal(t, want, result)
	})

	t.Run("concurrency limit", func(t *testing.T) {
		t.Parallel()

		// arrange
		const concurrency = 4
		const totalRequests = 5

		block := make(chan struct{})
		proceed := make(chan struct{}, totalRequests)

		m := minify.New()
		m.AddFunc("image/svg+xml", func(_ *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
			proceed <- struct{}{}
			<-block
			_, err := io.Copy(w, r)
			return err
		})

		sem := make(chan struct{}, concurrency)
		d, err := NewDiagram(withMinifier(m, sem))
		require.NoError(t, err)

		postProcessor := d.svgPostProcessors[0]

		// act
		var wg sync.WaitGroup
		wg.Add(totalRequests)
		for range totalRequests {
			go func() {
				defer wg.Done()
				_, err := postProcessor(context.Background(), []byte("<svg></svg>"))
				require.NoError(t, err)
			}()
		}

		// assert
		for range concurrency {
			<-proceed
		}
		require.Equal(t, concurrency, len(sem), "Semaphore should be full")

		close(block)

		wg.Wait()
		require.Equal(t, 0, len(sem), "Semaphore should be empty")
	})

	t.Run("context cancellation before acquire", func(t *testing.T) {
		t.Parallel()

		// arrange
		d, err := NewDiagram(WithMinifier())
		require.NoError(t, err)

		input := []byte(`<svg xmlns="http://www.w3.org/2000/svg">
			<rect width="100" height="100"/>
		</svg>`)

		postProcessor := d.svgPostProcessors[0]
		ctx, cancel := context.WithCancel(context.Background())

		// act
		cancel()
		_, err = postProcessor(ctx, input)

		// assert
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("context cancellation while waiting", func(t *testing.T) {
		t.Parallel()

		// arrange
		sem := make(chan struct{}, 4)
		for range cap(sem) {
			sem <- struct{}{}
		}

		m := minify.New()
		m.AddFunc("image/svg+xml", func(_ *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
			_, err := io.Copy(w, r)
			return err
		})

		d, err := NewDiagram(withMinifier(m, sem))
		require.NoError(t, err)

		postProcessor := d.svgPostProcessors[0]

		// act
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err = postProcessor(ctx, []byte("<svg></svg>"))

		// assert
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("context cancellation during processing", func(t *testing.T) {
		t.Parallel()

		// arrange
		block := make(chan struct{})
		defer close(block)
		m := minify.New()
		m.AddFunc("image/svg+xml", func(_ *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
			<-block
			_, err := io.Copy(w, r)
			return err
		})

		sem := make(chan struct{}, 4)
		d, err := NewDiagram(withMinifier(m, sem))
		require.NoError(t, err)

		postProcessor := d.svgPostProcessors[0]

		ctx, cancel := context.WithCancel(context.Background())
		resultCh := make(chan error, 1)

		go func() {
			_, err := postProcessor(ctx, []byte("<svg></svg>"))
			resultCh <- err
		}()

		select {
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for semaphore acquisition")
		case <-sem:
			sem <- struct{}{}
		}

		// act
		cancel()

		// assert
		select {
		case err := <-resultCh:
			require.ErrorIs(t, err, context.Canceled)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for result")
		}
	})

	t.Run("return error if minifier fails", func(t *testing.T) {
		t.Parallel()

		// arrange
		m := minify.New()
		m.AddFunc(
			"image/svg+xml",
			func(_ *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
				return errors.New("some error")
			},
		)

		sem := make(chan struct{}, 4)

		d, err := NewDiagram(withMinifier(m, sem))
		require.NoError(t, err)

		postProcessor := d.svgPostProcessors[0]
		ctx := context.Background()
		input := []byte(`invalid-svg`)

		// act
		_, err = postProcessor(ctx, input)

		// assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to minify svg")
	})
}
