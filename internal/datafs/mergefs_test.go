package datafs

import (
	"context"
	"net/url"
	"testing"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/stretchr/testify/assert"
)

// func TestReadMerge(t *testing.T) {
// 	ctx := context.Background()

// 	jsonContent := `{"hello": "world"}`
// 	yamlContent := "hello: earth\ngoodnight: moon\n"
// 	arrayContent := `["hello", "world"]`

// 	mergedContent := "goodnight: moon\nhello: world\n"

// 	fsys := fstest.MapFS{}
// 	fsys["tmp"] = &fstest.MapFile{Mode: fs.ModeDir | 0777}
// 	fsys["tmp/jsonfile.json"] = &fstest.MapFile{Data: []byte(jsonContent)}
// 	fsys["tmp/array.json"] = &fstest.MapFile{Data: []byte(arrayContent)}
// 	fsys["tmp/yamlfile.yaml"] = &fstest.MapFile{Data: []byte(yamlContent)}
// 	fsys["tmp/textfile.txt"] = &fstest.MapFile{Data: []byte(`plain text...`)}

// 	// workding dir with volume name trimmed
// 	wd, _ := os.Getwd()
// 	vol := filepath.VolumeName(wd)
// 	wd = wd[len(vol)+1:]

// 	fsys[path.Join(wd, "jsonfile.json")] = &fstest.MapFile{Data: []byte(jsonContent)}
// 	fsys[path.Join(wd, "array.json")] = &fstest.MapFile{Data: []byte(arrayContent)}
// 	fsys[path.Join(wd, "yamlfile.yaml")] = &fstest.MapFile{Data: []byte(yamlContent)}
// 	fsys[path.Join(wd, "textfile.txt")] = &fstest.MapFile{Data: []byte(`plain text...`)}

// 	fsmux := fsimpl.NewMux()
// 	fsmux.Add(fsimpl.WrappedFSProvider(&fsys, "file"))
// 	ctx = datafs.ContextWithFSProvider(ctx, fsmux)

// 	source := &Source{Alias: "foo", URL: mustParseURL("merge:file:///tmp/jsonfile.json|file:///tmp/yamlfile.yaml")}
// 	d := &Data{
// 		Sources: map[string]*Source{
// 			"foo":       source,
// 			"bar":       {Alias: "bar", URL: mustParseURL("file:///tmp/jsonfile.json")},
// 			"baz":       {Alias: "baz", URL: mustParseURL("file:///tmp/yamlfile.yaml")},
// 			"text":      {Alias: "text", URL: mustParseURL("file:///tmp/textfile.txt")},
// 			"badscheme": {Alias: "badscheme", URL: mustParseURL("bad:///scheme.json")},
// 			"badtype":   {Alias: "badtype", URL: mustParseURL("file:///tmp/textfile.txt?type=foo/bar")},
// 			"array":     {Alias: "array", URL: mustParseURL("file:///tmp/array.json?type=" + url.QueryEscape(jsonArrayMimetype))},
// 		},
// 		Ctx: ctx,
// 	}

// 	actual, err := d.readMerge(ctx, source)
// 	assert.NoError(t, err)
// 	assert.Equal(t, mergedContent, string(actual))

// 	source.URL = mustParseURL("merge:bar|baz")
// 	actual, err = d.readMerge(ctx, source)
// 	assert.NoError(t, err)
// 	assert.Equal(t, mergedContent, string(actual))

// 	source.URL = mustParseURL("merge:./jsonfile.json|baz")
// 	actual, err = d.readMerge(ctx, source)
// 	assert.NoError(t, err)
// 	assert.Equal(t, mergedContent, string(actual))

// 	source.URL = mustParseURL("merge:file:///tmp/jsonfile.json")
// 	_, err = d.readMerge(ctx, source)
// 	assert.Error(t, err)

// 	source.URL = mustParseURL("merge:bogusalias|file:///tmp/jsonfile.json")
// 	_, err = d.readMerge(ctx, source)
// 	assert.Error(t, err)

// 	source.URL = mustParseURL("merge:file:///tmp/jsonfile.json|badscheme")
// 	_, err = d.readMerge(ctx, source)
// 	assert.Error(t, err)

// 	source.URL = mustParseURL("merge:file:///tmp/jsonfile.json|badtype")
// 	_, err = d.readMerge(ctx, source)
// 	assert.Error(t, err)

// 	source.URL = mustParseURL("merge:file:///tmp/jsonfile.json|array")
// 	_, err = d.readMerge(ctx, source)
// 	assert.Error(t, err)
// }

func TestMergeData(t *testing.T) {
	// TODO: turn this back on!
	t.Skip()

	def := map[string]interface{}{
		"f": true,
		"t": false,
		"z": "def",
	}
	out, err := mergeData([]map[string]interface{}{def})
	assert.NoError(t, err)
	assert.Equal(t, "f: true\nt: false\nz: def\n", string(out))

	over := map[string]interface{}{
		"f": false,
		"t": true,
		"z": "over",
	}
	out, err = mergeData([]map[string]interface{}{over, def})
	assert.NoError(t, err)
	assert.Equal(t, "f: false\nt: true\nz: over\n", string(out))

	over = map[string]interface{}{
		"f": false,
		"t": true,
		"z": "over",
		"m": map[string]interface{}{
			"a": "aaa",
		},
	}
	out, err = mergeData([]map[string]interface{}{over, def})
	assert.NoError(t, err)
	assert.Equal(t, "f: false\nm:\n  a: aaa\nt: true\nz: over\n", string(out))

	uber := map[string]interface{}{
		"z": "über",
	}
	out, err = mergeData([]map[string]interface{}{uber, over, def})
	assert.NoError(t, err)
	assert.Equal(t, "f: false\nm:\n  a: aaa\nt: true\nz: über\n", string(out))

	uber = map[string]interface{}{
		"m": "notamap",
		"z": map[string]interface{}{
			"b": "bbb",
		},
	}
	out, err = mergeData([]map[string]interface{}{uber, over, def})
	assert.NoError(t, err)
	assert.Equal(t, "f: false\nm: notamap\nt: true\nz:\n  b: bbb\n", string(out))

	uber = map[string]interface{}{
		"m": map[string]interface{}{
			"b": "bbb",
		},
	}
	out, err = mergeData([]map[string]interface{}{uber, over, def})
	assert.NoError(t, err)
	assert.Equal(t, "f: false\nm:\n  a: aaa\n  b: bbb\nt: true\nz: over\n", string(out))
}

func TestMergeFS_Open(t *testing.T) {
	u, _ := url.Parse("merge:")
	fsys, err := NewMergeFS(u)
	assert.NoError(t, err)
	assert.IsType(t, &mergeFS{}, fsys)

	mux := fsimpl.NewMux()
	mux.Add(MergeFS)

	ctx := context.Background()
	ctx = ContextWithFSProvider(ctx, mux)

	fsys = fsimpl.WithContextFS(ctx, fsys)

	_, err = fsys.Open("/")
	assert.Error(t, err)

	_, err = fsys.Open("just/one/part")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "need at least 2 datasources to merge")

	// unknown aliases, fallback to relative files, but there's no FS registered
	// for the empty scheme
	_, err = fsys.Open("a|b")
	assert.ErrorContains(t, err, "no filesystem registered for scheme \"\"")
}
