package integrations

import (
	"reflect"
	"testing"
)

func TestUniqueCIDRs(t *testing.T) {
	got := uniqueCIDRs([]string{"192.168.1.0/24", "", " 10.0.0.0/24 ", "192.168.1.0/24"})
	want := []string{"192.168.1.0/24", "10.0.0.0/24"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
