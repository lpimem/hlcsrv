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
	delete from hlc_comment;
	delete from hlc_google_auth;
	delete from hlc_session;
	`)
	return err
}

func seedTestDb() error {
	_, err := storage.DB.Exec(`
	insert into hlc_user(id, name, email, password, _slt) 
		values (1, "Bob", "bob@example.com", "unsafe", "unsafe");
		values (10, "ExAm", "example@gmail.com", "unsafe", "unsafe");

	insert into hlc_page(id, title, url)
		values (1, "example", "http://example.com");

	insert into hlc_range(id, anchor, start, startOffset, end, endOffset, text, page, author)
		values (1, "#c", "#c/1", 0, "#c/12", 32, "This is the selected text", 1, 1);

	insert into hlc_google_auth(google_id, uid) values ("100000", 1);
	insert into hlc_google_auth(google_id, uid) values ("example@gmail.com", 10);

	insert into hlc_session(id, uid, last_access) values ("fake_session_id", 10, CURRENT_TIMESTAMP);
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
