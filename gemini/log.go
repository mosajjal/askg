package gemini

import "github.com/rs/zerolog"

// Log is a wrapper around zerolog.Logger to make it compatible with resty
type Log struct {
	*zerolog.Logger
}

// Debugf implements resty.Logger interface
func (l Log) Debugf(msg string, args ...interface{}) {
	l.Debug().Msgf(msg, args...)
}

// Infof implements resty.Logger interface
func (l Log) Infof(msg string, args ...interface{}) {
	l.Info().Msgf(msg, args...)
}

// Warnf implements resty.Logger interface
func (l Log) Warnf(msg string, args ...interface{}) {
	l.Warn().Msgf(msg, args...)
}

// Errorf implements resty.Logger interface
func (l Log) Errorf(msg string, args ...interface{}) {
	l.Error().Msgf(msg, args...)
}

// Fatalf implements resty.Logger interface
func (l Log) Fatalf(msg string, args ...interface{}) {
	l.Fatal().Msgf(msg, args...)
}
