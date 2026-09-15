package kidwords

import (
	"errors"
)

type SeparatorFunc func() []byte

type writerOptions struct {
	separator  SeparatorFunc
	dictionary *Dictionary
}

type WriterOption interface {
	applyWriterOption(*writerOptions) error
}

type readerOptions struct {
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
	dictionary *Dictionary
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

func WithDictionary(d *Dictionary) Option {
	return &dictionaryOption{dictionary: d}
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
