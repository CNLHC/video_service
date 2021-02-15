package sdk

import "testing"

func TestSliceIdGen(t *testing.T) {
	sample := []string{
		"aaaaaaaaa",
		"aaaaaaaab",
		"aaaaaaaac",
		"aaaaaaaad",
		"aaaaaaaae",
		"aaaaaaaaf",
		"aaaaaaaag",
		"aaaaaaaah",
	}
	sdk := &XunfeiSDK{}
	sdk.ResetSliceId()
	for i := 0; i < 5; i += 1 {
		if sample[i] != sdk.cur_sliceid {
			t.Errorf("unexpected %s (%s)", sample[i], sdk.cur_sliceid)
		}
		_ = sdk.GetNextSliceId()
	}
}

func TestSignature(t *testing.T) {
	sk := "d9f4aa7ea6d94faca62cd88a28fd5234"
	appid := "595f23df"
	ts := "1512041814"
	signa := signature(appid, sk, ts)
	if signa != "IrrzsJeOFk1NGfJHW6SkHUoN9CU=" {
		t.Errorf("wrong signature")
	}

}
