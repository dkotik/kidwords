package kidwords

import (
	"errors"
	"fmt"

	"github.com/dkotik/kidwords/dictionary"
)

// type SplitFunc func()
type SeparatorFunc func() []byte

type writerOptions struct {
	separator  SeparatorFunc
	dictionary *dictionary.Dictionary
}

type WriterOption interface {
	applyWriterOption(*writerOptions) error
}

type readerOptions struct {
	// split SplitFunc
	dictionary map[string]byte
}

type ReaderOption interface {
	applyReaderOption(*readerOptions) error
}

type Option interface {
	ReaderOption
	WriterOption
}

type dictionaryOption struct {
	dictionary *dictionary.Dictionary
}

func (d *dictionaryOption) validate() error {
	if d == nil || d.dictionary == nil {
		return errors.New("cannot use a <nil> dictionary")
	}
	return d.dictionary.Validate()
}

func (d *dictionaryOption) applyWriterOption(o *writerOptions) error {
	if err := d.validate(); err != nil {
		return err
	}
	if o.dictionary != nil {
		return errors.New("dictionary is already set")
	}
	o.dictionary = d.dictionary
	return nil
}

func (d *dictionaryOption) applyReaderOption(o *readerOptions) error {
	if err := d.validate(); err != nil {
		return err
	}
	if o.dictionary != nil {
		return errors.New("dictionary is already set")
	}
	o.dictionary = d.dictionary.Reverse()
	return nil
}

func WithDictionary(d *dictionary.Dictionary) Option {
	return &dictionaryOption{dictionary: d}
}

type dictionaryFileOption string

func (d dictionaryFileOption) applyWriterOption(o *writerOptions) error {
	if d == "" {
		return errors.New("cannot use an empty dictionary file path")
	}
	dictionary, err := dictionary.LoadFile(string(d))
	if err != nil {
		return fmt.Errorf("cannot load dictionary %q: %w", d, err)
	}
	return (&dictionaryOption{dictionary: &dictionary}).applyWriterOption(o)
}

func (d dictionaryFileOption) applyReaderOption(o *readerOptions) error {
	if d == "" {
		return errors.New("cannot use an empty dictionary file path")
	}
	dictionary, err := dictionary.LoadFile(string(d))
	if err != nil {
		return fmt.Errorf("cannot load dictionary %q: %w", d, err)
	}
	return (&dictionaryOption{dictionary: &dictionary}).applyReaderOption(o)
}

type separatorOption SeparatorFunc

func (s separatorOption) applyWriterOption(o *writerOptions) error {
	if s == nil {
		return errors.New("cannot us a <nil> separator function")
	}
	if o.separator != nil {
		return errors.New("separator function is already set")
	}
	o.separator = SeparatorFunc(s)
	return nil
}

func WithSeparator(f SeparatorFunc) WriterOption {
	return separatorOption(f)
}
