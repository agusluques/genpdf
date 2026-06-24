/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package model

import (
	"bytes"
	"testing"
)

// utf16BEDate encodes an ASCII date string as UTF-16BE with a leading BOM,
// matching what producers like Aspose.PDF emit.
func utf16BEDate(ascii string) string {
	var b bytes.Buffer
	b.Write([]byte{0xFE, 0xFF})
	for _, r := range ascii {
		b.WriteByte(byte(r >> 8))
		b.WriteByte(byte(r & 0xFF))
	}
	return b.String()
}

func TestNewPdfDate(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
		year    int64
	}{
		{name: "ascii utc", in: "D:20240809160228Z", year: 2024},
		{name: "ascii offset", in: "D:20221231235959+02'00'", year: 2022},
		{name: "utf16be bom", in: utf16BEDate("D:20240809160228Z"), year: 2024},
		{name: "garbage", in: "not a date", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := NewPdfDate(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got date %+v", d)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.year != tc.year {
				t.Fatalf("year = %d, want %d", d.year, tc.year)
			}
		})
	}
}
