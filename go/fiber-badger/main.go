package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type list struct {
	Name      string `json:"name"`
	TodoCount int    `json:"todoCount"`
}

type Store struct {
	db *badger.DB
}

func (s *Store) createListForUser(username string, listname string) {
	err := s.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("user/%s/list/%s", username, listname)
		listMeta := list{
			Name:      listname,
			TodoCount: 0,
		}
		val, _ := json.Marshal(listMeta)
		return txn.Set([]byte(key), val)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) getListById(username string, listname string) (*list, error) {
	var result list
	err := s.db.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("user/%s/list/%s", username, listname)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &result)
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Store) getAllListForUser(username string) ([]*list, error) {
	result := make([]*list, 0)
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(fmt.Sprintf("user/%s/list/", username))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				list := list{}
				err := json.Unmarshal(v, &list)
				result = append(result, &list)
				return err
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	s := &Store{db}

	// list1, err := s.getListById("ashwin", "list1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%v\n", list1)

	// s.createListForUser("ashwin", "list1")
	// s.createListForUser("ashwin", "list2")
	// s.createListForUser("ashwin", "list3")

	allLists, _ := s.getAllListForUser("ashwin")

	for i, v := range allLists {
		fmt.Printf("===========\nItem %d\n=========\n%v\n===========\n", i+1, v)
	}

}
