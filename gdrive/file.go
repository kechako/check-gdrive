package gdrive

import (
	"fmt"

	drive "google.golang.org/api/drive/v3"
)

type File struct {
	g *GDrive
	drive.File
	Path string
}

func (f *File) IsFolder() bool {
	return f.MimeType == "application/vnd.google-apps.folder"
}

func (f *File) IsFile() bool {
	return !f.IsFolder()
}

func (f *File) GetFiles() ([]*File, error) {
	if !f.IsFolder() {
		return nil, fmt.Errorf("%s (%s) is not a folder.", f.Name, f.Id)
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", f.Id)

	files := make([]*File, 0, 20)
	pageToken := ""
	for {
		req := f.g.service.Files.List().Q(query).Fields("files(id,md5Checksum,mimeType,name),nextPageToken").OrderBy("name")
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			return nil, err
		}

		if len(r.Files) > 0 {
			for _, i := range r.Files {
				file := &File{
					File: *i,
					g:    f.g,
					Path: f.Join(i.Name),
				}
				files = append(files, file)
			}
		}
		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	return files, nil
}

func (f *File) Join(name string) string {
	if f.Path == "" {
		return name
	}

	return fmt.Sprintf("%s/%s", f.Path, name)
}
