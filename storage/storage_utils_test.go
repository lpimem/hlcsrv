package storage

import "github.com/lpimem/hlcsrv/util"

func InitTestDb() {
	InitStorage(util.GetAbsRunDirPath() + "/db/test.db")
	err := cleanDb()
	if err != nil {
		panic(err)
	}
	err = seedTestDb()
	if err != nil {
		panic(err)
	}
}

func cleanDb() error {
	_, err := storage.DB.Exec(`
	delete from hlc_range;
	delete from hlc_user;
	delete from hlc_page;
	delete from hlc_comments;
	`)
	return err
}

func seedTestDb() error {
	_, err := storage.DB.Exec(`
	insert into hlc_user(id, name, email, password, _slt) 
		values (1, "Bob", "bob@example.com", "unsafe", "unsafe");

	insert into hlc_page(id, title, url)
		values (1, "example", "http://example.com");

	insert into hlc_range(id, anchor, start, startOffset, end, endOffset, text, page, author)
		values (1, "#c", "#c/1", 0, "#c/12", 32, "This is the selected text", 1, 1);
	`)
	return err
}

func ResetTestDb() {
	cleanDb()
	seedTestDb()
	util.Log("DB@", storage.path, " is reseted")
}

func init() {
	InitTestDb()
}
