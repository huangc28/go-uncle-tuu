package main

// This worker fetch unimported excel file from GCS and imported to DB.
// The worker will first read column `import_status` (pending, imported, import_failed) from database.
// Read all file names that are with `pending` status. Fetch GCS object from google cloud storage into a io.Reader.
//
// We will use `excelize` to read column data from this io.Reader.

func main() {}
