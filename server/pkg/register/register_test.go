package register

import "testing"

func TestRegister(t *testing.T) {
	rram := InitRRAM()

	val, err := rram.GetValue(0)
	if val != 0 || err != nil {
		t.Error("Can not read")
	}

	rram.PutValue(0, 100)
	val, err = rram.GetValue(0)
	if val != 100 || err != nil {
		t.Error("Can not write")
	}

	t.Log("Hey! I'm too lazy to write tests atm. Well, that's the problem of a future me")
}
