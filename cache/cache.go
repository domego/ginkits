package cachekits

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
)

const (
	// GzipMinSize gzip min size
	GzipMinSize = 1024
	// CacheFormatRaw raw
	CacheFormatRaw = 0
	// CacheFormatRawGzip raw gzip
	CacheFormatRawGzip = 1
	// CacheFormatJSON json
	CacheFormatJSON = 10
	// CacheFormatJSONGzip json gzip
	CacheFormatJSONGzip = 11
)

type ModelCacheItem struct {
	Data []byte
	Flag uint32
}

// ToGzipJSON convert model to gziped json
func ToGzipJSON(obj interface{}) (gziped bool, data []byte, err error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return
	}
	if len(bs) <= GzipMinSize {
		return false, bs, nil
	}
	buf := &bytes.Buffer{}
	gzipWriter := gzip.NewWriter(buf)
	_, err = gzipWriter.Write(bs)
	gzipWriter.Close()
	if err != nil {
		return
	}
	return true, buf.Bytes(), nil
}

func FromGzipJSON(data []byte, obj interface{}) (err error) {
	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	defer gzipReader.Close()
	bs, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return
	}
	return json.Unmarshal(bs, obj)
}
