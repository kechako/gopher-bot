package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeTestDir() (string, func(), error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot-store-test")
	if err != nil {
		return "", nil, err
	}

	return dir, func() {
		os.RemoveAll(dir)
	}, nil
}

type TestData struct {
	Name  string
	Value int
}

func Test_Store(t *testing.T) {
	dir, remove, err := makeTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer remove()

	store, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("New must not be return nil")
	}

	testData := &TestData{
		Name:  "Test",
		Value: 12345,
	}

	err = store.Update(func(tx *Tx) error {
		var err error

		var data TestData
		err = tx.Get("KEY_NOT_EXISTS", &data)
		if err != ErrKeyNotFound {
			t.Error("Tx.Get must be return ErrKeyNotFound")
		}

		err = tx.Set("TEST_KEY", testData)
		if err != nil {
			t.Error(err)
		}

		err = tx.Get("TEST_KEY", &data)
		if err != nil {
			t.Error(err)
		}

		diff := cmp.Diff(&data, testData)
		if diff != "" {
			t.Errorf("failed to read data from database: (-got +want)\n%s", diff)
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	err = store.View(func(tx *Tx) error {
		var err error

		var data TestData

		err = tx.Get("TEST_KEY", &data)
		if err != nil {
			t.Error(err)
		}

		diff := cmp.Diff(&data, testData)
		if diff != "" {
			t.Errorf("failed to read data from database: (-got +want)\n%s", diff)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	testData2 := &TestData{
		Name:  "Test2",
		Value: 67890,
	}

	err = store.Update(func(tx *Tx) error {
		var err error

		var data TestData
		err = tx.Set("TEST_KEY", testData2)
		if err != nil {
			t.Error(err)
		}

		err = tx.Get("TEST_KEY", &data)
		if err != nil {
			t.Error(err)
		}

		diff := cmp.Diff(&data, testData2)
		if diff != "" {
			t.Errorf("failed to read data from database: (-got +want)\n%s", diff)
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	err = store.Update(func(tx *Tx) error {
		var err error

		err = tx.Delete("KEY_NOT_EXISTS")
		if err != nil {
			t.Error(err)
		}

		err = tx.Delete("TEST_KEY")
		if err != nil {
			t.Error(err)
		}

		var data TestData
		err = tx.Get("TEST_KEY", &data)
		if err != ErrKeyNotFound {
			t.Error("Tx.Get must be return ErrKeyNotFound")
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}

}
