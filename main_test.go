package raml

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/acronis/go-raml/stacktrace"
)

func Test_main(t *testing.T) {
	start := time.Now()
	rml, err := ParseFromPath(`./tests/library.raml`, OptWithUnwrap(), OptWithValidate())
	if vErr, ok := stacktrace.Unwrap(err); ok {
		t.Logf("ParseFromPath error:\n%s", vErr.Sprint(stacktrace.WithEnsureDuplicates()))
		err = vErr
	}
	require.NoError(t, err)
	elapsed := time.Since(start)
	t.Logf("ParseFromPath took %d ms\n", elapsed.Milliseconds())
	fmt.Printf("Library location: %s\n", rml.entryPoint.GetLocation())

	shapesAll := rml.GetShapes()
	fmt.Printf("Total shapes: %d\n", len(shapesAll))

	conv := NewJSONSchemaConverter()
	for _, frag := range rml.fragmentsCache {
		switch f := frag.(type) {
		case *Library:
			for pair := f.AnnotationTypes.Oldest(); pair != nil; pair = pair.Next() {
				shape := pair.Value
				s := *shape
				conv.Convert(s)
				// b, err := json.MarshalIndent(schema, "", "  ")
				// if err != nil {
				// 	t.Errorf("StackTrace marshalling schema: %s", err)
				// }
				// os.WriteFile(fmt.Sprintf("./out/%s_%s.json", s.Base().Name, s.Base().Id), b, 0644)
				//fmt.Println(string(b))
			}
			for pair := f.Types.Oldest(); pair != nil; pair = pair.Next() {
				shape := pair.Value
				s := *shape
				conv.Convert(s)
				// b, err := json.MarshalIndent(schema, "", "  ")
				// if err != nil {
				// 	t.Errorf("StackTrace marshalling schema: %s", err)
				// }
				// os.WriteFile(fmt.Sprintf("./out/%s_%s.json", s.Base().Name, s.Base().Id), b, 0644)
				//fmt.Println(string(b))
			}
		case *DataType:
			s := *f.Shape
			conv.Convert(s)
			// b, err := json.MarshalIndent(schema, "", "  ")
			// if err != nil {
			// 	t.Errorf("StackTrace marshalling schema: %s", err)
			// }
			// os.WriteFile(fmt.Sprintf("./out/%s_%s.json", s.Base().Name, s.Base().Id), b, 0644)
			//fmt.Println(string(b))
		}
	}

	printMemUsage(t)
}

func printMemUsage(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("Alloc = %v MiB", m.Alloc/1024/1024)
	t.Logf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	t.Logf("\tSys = %v MiB", m.Sys/1024/1024)
	t.Logf("\tNumGC = %v\n", m.NumGC)
}
