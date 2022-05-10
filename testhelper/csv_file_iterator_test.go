package testhelper


func ExampleCsvFileIterator() {
   iter := NewCsvFileIterator("/path/to/foobar.csv")
   defer iter.Close()
   for iter.HasNext() {
	 row := iter.Fetch()
	 println(row["name"])   // foobar
	 println(row["json"])   // {"url": "https://example.com"}
   }
}
