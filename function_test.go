package spider

import "testing"

func Test_spiderGps(t *testing.T) {
	context := &Context{
		Spider: NewSpiderMock(),
	}
	if err := spiderGps(context); err != nil {
		t.Fatalf("spiderGps fail %v", err)
	}
}

func Test_spiderMemory(t *testing.T) {
	context := &Context{
		Spider: NewSpiderMock(),
	}
	if err := spiderMemory(context); err != nil {
		t.Fatalf("spiderMemory fail %v", err)
	}
}

func Test_spiderStack(t *testing.T) {
	context := &Context{
		Spider: NewSpiderMock(),
	}
	if err := spiderStack(context); err != nil {
		t.Fatalf("spiderStack fail %v", err)
	}
}
