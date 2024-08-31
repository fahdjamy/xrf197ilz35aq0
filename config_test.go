package xrf197ilz35aq0

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

var validYAML = []byte(`
field:
  innerField: field
  anotherField: another
`)

func TestNewConfig(t *testing.T) {
	t.Run("should return a valid config", func(t *testing.T) {
		config, err := NewConfig()
		AssertNoError(t, err)
		assertConfigIsNotNil(t, config)
	})

	t.Run("should read config once a valid config", func(t *testing.T) {
		routines := 10
		var wg sync.WaitGroup
		wg.Add(routines)

		// This is an atomic.Int32, which means it's a 32-bit integer that can be safely modified by multiple
		// goroutines concurrently without data races.
		// It's initialized to 0, signifying that the configuration hasn't been set yet.
		var configSetCounter atomic.Int32

		for i := 0; i < routines; i++ {
			go func() {
				defer wg.Done()

				config, err := NewConfig()
				AssertNoError(t, err)

				// configSetCount.CompareAndSwap(0, 1): This is an atomic operation:
				// checks if the current value of configSetCount is 0.
				// If it is 0, it atomically changes the value to 1.
				// returns true if the swap was successful (i.e., the value was 0 and is now 1), false otherwise

				assertConfigIsNotNil(t, config)
				if configSetCounter.CompareAndSwap(0, 1) {
					// This comment explains the intent clearly there's no code in here just to make sure
					// configSetCounter is set once to 1 indicating config was called once
					// check if config is nil before incrementing
					// Successfully incremented the counter, indicating the config was set for the first time
				}
			}()
		}

		wg.Wait()

		// After all calls, configSetCount should be exactly 1
		if configSetCounter.Load() != 1 {
			t.Errorf("Configurations should be set exactly once. Was set: %d", configSetCounter.Load())
		}
	})
}

func TestReadConfiguration(t *testing.T) {
	t.Run("should read configurations from file", func(t *testing.T) {
		mock := &mockFileDataCopier{content: validYAML}
		config, err := readConfiguration(mock)

		AssertNoError(t, err)
		assert.True(t, mock.closed)
		assertNonNilConfigPtr(t, config)
	})

	t.Run("should fail if file yaml content is invalid", func(t *testing.T) {
		invalidYAML := []byte(`invalid yaml`)
		mock := &mockFileDataCopier{content: invalidYAML}
		config, err := readConfiguration(mock)
		AssertError(t, err)
		assert.True(t, mock.closed)
		assertNilConfigPtr(t, config)
	})
}

func assertConfigIsNotNil(t testing.TB, config Config) {
	t.Helper()
	if reflect.DeepEqual(config, Config{}) {
		t.Error("Configuration should not be nil")
	}
}

func assertConfigIsNil(t testing.TB, config Config) {
	t.Helper()
	if !reflect.DeepEqual(config, Config{}) {
		t.Error("Configuration should be nil")
	}
}

func assertNilConfigPtr(t testing.TB, config *Config) {
	t.Helper()
	if config != nil {
		t.Error("Configuration should be nil")
	}
}

func assertNonNilConfigPtr(t testing.TB, config *Config) {
	t.Helper()
	if config == nil {
		t.Error("Configuration should not be nil")
	}
}
