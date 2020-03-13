// Package sounds provides the basic types for Sounds within this system,
// plus multiple implementations of different sounds that can be used.
package sounds

import (
	"time"
)

const (
	// CyclesPerSecond is the sample rate of each sound stream.
	CyclesPerSecond = 44100.0

	// SecondsPerCycle is the inverse sample rate.
	SecondsPerCycle = 1.0 / CyclesPerSecond

	// DurationPerCycle = int64(SecondsPerCycle * 1e9) * time.Nanosecond
	DurationPerCycle = 22675 * time.Nanosecond // BIG HACK

	// MaxLength represents the number of samples in the maximum duration.
	MaxLength = uint64(406750706825295)

	// MaxDuration is used for unending sounds.
	MaxDuration = time.Duration(int64(float64(MaxLength)*SecondsPerCycle*1e9)) * time.Nanosecond
)

// A Sound is a model of a physical sound wave as a series of pressure changes over time.
//
// Each Sound contains a channel of samples in the range [-1, 1] of the intensity at each time step,
// as well as a count of samples, which then also defines how long the sound lasts.
//
// Sounds also provide a way to start and stop when the samples are written, and reset to an initial state.
type Sound interface {
	// Sound wave samples for the sound - only valid after Start() and before Stop()
	// NOTE: Only one sink should read from GetSamples(). Otherwise it will not receive every sample.
	GetSamples() <-chan float64

	// Number of samples in this sound, MaxLength if unlimited.
	Length() uint64

	// Length of time this goes for. Convenience method, should always be SamplesToDuration(Length())
	Duration() time.Duration

	// Start begins writing the sound wave to the samples channel.
	Start()

	// Running indicates whether a sound has Start()'d but not yet Stop()'d
	Running() bool

	// Stop ceases writing samples, and closes the channel.
	Stop()

	// Reset converts the sound back to the pre-Start() state. Can only be called on a Stop()'d Sound.
	Reset()
}

// SamplesToDuration converts a sample count to a duration of time.
func SamplesToDuration(sampleCount uint64) time.Duration {
	return time.Duration(int64(float64(sampleCount)*1e9*SecondsPerCycle)) * time.Nanosecond
}

// DurationToSamples converts a duration of time to a sample count
func DurationToSamples(duration time.Duration) uint64 {
	return uint64(float64(duration.Nanoseconds()) * 1e-9 * CyclesPerSecond)
}
