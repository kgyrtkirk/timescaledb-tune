package tstune

import "testing"

func TestParseLineForSharedLibResult(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  *sharedLibResult
	}{
		{
			desc: "initial config value",
			input: "#shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "extra commented out",
			input: "###shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "commented with space after",
			input: "# shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "extra commented with space after",
			input: "## shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "initial config value, uncommented",
			input: "shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    false,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "initial config value, uncommented with leading space",
			input: "  shared_preload_libraries = ''		# (change requires restart)",
			want: &sharedLibResult{
				commented:    false,
				hasTimescale: false,
				libs:         "",
			},
		},
		{
			desc: "timescaledb already there but commented",
			input: "#shared_preload_libraries = 'timescaledb'		# (change requires restart)",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: true,
				libs:         "timescaledb",
			},
		},
		{
			desc:  "other libraries besides timescaledb, commented",
			input: "#shared_preload_libraries = 'pg_stats' # (change requires restart)   ",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: false,
				libs:         "pg_stats",
			},
		},
		{
			desc:  "no string after the quotes",
			input: "shared_preload_libraries = 'pg_stats,timescaledb'",
			want: &sharedLibResult{
				commented:    false,
				hasTimescale: true,
				libs:         "pg_stats,timescaledb",
			},
		},
		{
			desc: "don't be greedy with things between single quotes",
			input: "#shared_preload_libraries = 'timescaledb'		# comment with single quote ' test",
			want: &sharedLibResult{
				commented:    true,
				hasTimescale: true,
				libs:         "timescaledb",
			},
		},
		{
			desc:  "not shared preload line",
			input: "data_dir = '/path/to/data'",
			want:  nil,
		},
	}
	for _, c := range cases {
		res := parseLineForSharedLibResult(c.input)
		if res == nil && c.want != nil {
			t.Errorf("%s: result was unexpectedly nil: want %v", c.desc, c.want)
		} else if res != nil && c.want == nil {
			t.Errorf("%s: result was unexpectedly non-nil: got %v", c.desc, res)
		} else if c.want != nil {
			if got := res.commented; got != c.want.commented {
				t.Errorf("%s: incorrect commented: got %v want %v", c.desc, got, c.want.commented)
			}
			if got := res.hasTimescale; got != c.want.hasTimescale {
				t.Errorf("%s: incorrect hasTimescale: got %v want %v", c.desc, got, c.want.hasTimescale)
			}
			if got := res.libs; got != c.want.libs {
				t.Errorf("%s: incorrect libs: got %s want %s", c.desc, got, c.want.libs)
			}
		}
	}
}
