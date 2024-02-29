package rules

import (
	"testing"
)

func add(i1, i2 int) (result int) {
	result = i1 + i2
	return 
}

func TestAdd(t *testing.T) {
	result := add(3, 4)
	expected := 7
	if result != expected {
		t.Errorf("Add function returned incorrect result, got: %d, want: %d", result, expected)
	}
}


func TestRecovery(t *testing.T) {
	srv := NewRecoverySrv()
	srv.init()
	size := 0
	ban_list.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	if size != 1 {
		t.Errorf("TestRecovery, got: %d, want: %d", size, 1)
	}
	//srv.Add("172.168.1.1", 8889, "iptables", "test...")
}