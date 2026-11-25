package passes

import "testing"

func TestRemoveUnreachable(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "return after goto",
			src: `func f(x int) {
				for {
					x++
					goto end
				}
				goto end
				return
			end:
				x--
			}`,
			want: `func f(x int) {
				for {
					x++
					goto end
				}
				goto end
			end:
				x--
			}`,
		},
		{
			name: "statements after return up to a label",
			src: `func f(x int) int {
				if x > 0 {
					return x
					x++
				}
				goto out
				x++
				x--
			out:
				return x
			}`,
			want: `func f(x int) int {
				if x > 0 {
					return x
				}
				goto out
			out:
				return x
			}`,
		},
		{
			name: "after panic",
			src: `func f(x int) {
				panic("unreachable")
				x++
			}`,
			want: `func f(x int) {
				panic("unreachable")
			}`,
		},
		{
			name: "reachable code stays",
			src: `func f(x int) {
				if x > 0 {
					return
				}
				x++
			}`,
			want: `func f(x int) {
				if x > 0 {
					return
				}
				x++
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := parseFunc(t, tt.src)
			RemoveUnreachable(fn)
			got := formatFunc(t, fn)
			want := normalizeFunc(t, tt.want)
			if got != want {
				t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
			}
		})
	}
}
