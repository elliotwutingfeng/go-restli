package stdstructs

import "github.com/PapaCharlie/go-restli/protocol/restlicodec"

type EmptyRecord struct{}

func (e EmptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(restlicodec.Reader, string) error { return reader.Skip() })
}

func (e EmptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}
