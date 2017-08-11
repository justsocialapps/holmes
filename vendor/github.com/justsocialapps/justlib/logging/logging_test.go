package logging

import (
	"bytes"
	"log"
	"testing"

	"github.com/justsocialapps/assert"
)

func TestLevelFromStringShouldBeCaseInsensitive(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("debug")
	a.Equal(*l, DEBUG, "Wrong level returned")

	l, _ = LevelFromString("DEBUG")
	a.Equal(*l, DEBUG, "Wrong level returned")
}

func TestLevelFromStringShouldReturnNilWithUnknownString(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("oink")
	a.Nil(l, "Level is not nil")
}

func TestLevelFromStringShouldReturnNilWithEmptyString(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("")
	a.Nil(l, "Level is not nil")
}

func TestDefaultPrefixedLoggerShouldWritePrefixCorrectly(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	logger := NewPrefixedLogger("pre")

	logger.Debug("did that")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[DEBUG\\] \\[pre\\] did that\n$", buf.String(), "Unexpected log output")
}

func TestDefaultDebugShouldLogWithDebugLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	logger := NewLogger()

	logger.Debug("did that")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[DEBUG\\] did that\n$", buf.String(), "Unexpected log output")
}

func TestDefaultInfoShouldLogWithInfoLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	logger := NewLogger()

	logger.Info("did it")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[INFO\\] did it\n$", buf.String(), "Unexpected log output")
}

func TestDefaultWarnShouldLogWithWarnLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	logger := NewLogger()

	logger.Warn("warning")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[WARN\\] warning\n$", buf.String(), "Unexpected log output")
}

func TestDefaultErrorShouldLogWithErrorLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	logger := NewLogger()

	logger.Error("erroaarrr!!")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[ERROR\\] erroaarrr!!\n$", buf.String(), "Unexpected log output")
}
