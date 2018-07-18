package gojsonfilterbench

import (
	"testing"

	"github.com/couchbaselabs/gojsonsm"
)

func BenchmarkJsonSM(b *testing.B) {
	data, totalBytes, err := generateRandomData(1)
	if err != nil || len(data) == 0 {
		b.Fatalf("Data generation error: %s", err)
	}

	// name["first"]="Brett" OR (age<50 AND isActive=True)
	matchJson := []byte(`
	["or",
	  ["equals",
	    ["field", "name", "first"],
	    ["value", "Brett"]
	  ],
	  ["and",
	    ["lessthan",
	      ["field", "age"],
	      ["value", 50]
	    ],
	    ["equals",
	      ["field", "isActive"],
	      ["value", true]
	    ]
	  ]
    ]`)
	expr, err := gojsonsm.ParseJsonExpression(matchJson)
	if err != nil {
		b.Errorf("Failed to parse expression: %s", err)
		return
	}

	var trans gojsonsm.Transformer
	matchDef := trans.Transform([]gojsonsm.Expression{expr})
	m := gojsonsm.NewMatcher(matchDef)

	b.SetBytes(int64(totalBytes))
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		for i := 0; i < len(data); i++ {
			_, err := m.Match(data[i])

			if err != nil {
				b.Fatalf("Matcher error: %s", err)
			}
		}
	}
}

func BenchmarkJsonSMSlowMatcher(b *testing.B) {
	data, totalBytes, err := generateRandomData(1)
	if err != nil || len(data) == 0 {
		b.Fatalf("Data generation error: %s", err)
	}

	matchJson := []byte(`
	["or",
	  ["equals",
	    ["field", "name", "first"],
	    ["value", "Brett"]
	  ],
	  ["and",
	    ["lessthan",
	      ["field", "age"],
	      ["value", 50]
	    ],
	    ["equals",
	      ["field", "isActive"],
	      ["value", true]
	    ]
	  ]
    ]`)
	expr, err := gojsonsm.ParseJsonExpression(matchJson)
	if err != nil {
		b.Errorf("Failed to parse expression: %s", err)
		return
	}

	m := gojsonsm.NewSlowMatcher([]gojsonsm.Expression{expr})

	b.SetBytes(int64(totalBytes))
	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(data); i++ {
			_, err := m.Match(data[i])

			if err != nil {
				b.Fatalf("Slow matcher error: %s", err)
			}
		}
	}
}
