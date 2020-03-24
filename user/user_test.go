package user

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

// TestMain clean up user.db after running unit test
func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func cleanDb(b *testing.B) {
	os.Remove(dbPath)

	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Jhon",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		b.Fatalf("Error saving a record: %s", err)
	}
	b.ResetTimer()
}

func BenchmarkCreate(b *testing.B) {
	cleanDb(b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create the user so it can be read
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		b.StartTimer()
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retriving a record: %s", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	cleanDb(b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create the user so it can be updated
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		b.StartTimer()

		u.Role = "developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cleanDb(b)

	for i := 0; i < b.N; i++ {
		// Create the user so it can be deleted
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		b.StartTimer()

		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error removing a recrod: %s", err)
		}
	}
}

func TestCRUD(t *testing.T) {
	// Create user and save to DB
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Jhon",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	// Read data from DB
	t.Log("READ")
	// Creates the u2 user using the u.ID by grabbing it from the storm DB.
	u2, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retriving a record: %s", err)
	}
	if !reflect.DeepEqual(u2, u) {
		t.Error("Records do not match")
	}

	// Update OneUser in DB
	t.Log("Update")
	u.Role = "developer"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	// Delete a user from DB
	t.Log("Delete")
	u3, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retriving a record: %s", err)
	}
	if !reflect.DeepEqual(u3, u) {
		t.Error("Records do not match")
	}

	err = Delete(u.ID)
	if err != nil {
		t.Fatalf("Error removing a recrod: %s", err)
	}

	_, err = One(u.ID)
	if err == nil {
		t.Fatalf("Record should not exist anymore")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("Error retriving non-existing record: %s", err)
	}

	// Get all users
	t.Log("Read All")
	// Update the u2 & u3 ID and save to DB
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()
	u2.Save()
	u3.Save()
	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	users, err := All()
	if err != nil {
		t.Fatalf("Error reading all records: %s", users)
	}
	// TODO make a better comparison then number of enteries
	if len(users) != 2 {
		t.Errorf("Different number of records retrived. Expected 2 Got: %d, %v", len(users), users)
	}

}
