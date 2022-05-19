package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
)

func Test_csvWriter_Write(t *testing.T) {
	f, err := ioutil.TempFile("", "gitextractor")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test for Write",
			args{item: &code.Commit{
				Sha:            "ffwefef3f34f",
				Additions:      3,
				Deletions:      4,
				DevEq:          7,
				Message:        "",
				AuthorName:     "",
				AuthorEmail:    "",
				AuthoredDate:   time.Now(),
				AuthorId:       "",
				CommitterName:  "",
				CommitterEmail: "",
				CommittedDate:  time.Now(),
				CommitterId:    "",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gitextractor")
			if err != nil {
				t.Error(err)
			}
			defer os.RemoveAll(dir)
			w, err := newCsvWriter(filepath.Join(dir, "commits.csv"), &code.Commit{})
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			if err := w.Write(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
