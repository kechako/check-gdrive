package md5sum

import "testing"

const testdata = "testdata/test.dat"
const testdataMd5 = "2907575916a7052813fa6ea9aa9e601f"

func TestSumFile(t *testing.T) {
	hash, err := SumFile(testdata)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if hash != testdataMd5 {
		t.Errorf("got %v\nwant %v", hash, testdataMd5)
	}
}
