package gdrive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kechako/check-gdrive/utils"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v2"
)

type GDrive struct {
	client  *http.Client
	service *drive.Service
}

func New(ctx context.Context) (*GDrive, error) {
	b, err := Asset("data/client_secret.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	client, err := getClient(ctx, config)
	if err != nil {
		return nil, err
	}

	service, err := drive.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve drive client %v", err)
	}

	gDrive := &GDrive{}
	gDrive.client = client
	gDrive.service = service

	return gDrive, nil
}

func (g *GDrive) GetFile(id string) (*File, error) {
	f, err := g.service.Files.Get(id).Do()
	if err != nil {
		return nil, err
	}

	file := &File{
		g:    g,
		File: *f,
		Path: "",
	}

	return file, nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		return nil, fmt.Errorf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(cacheFile, tok)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("State-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the"+
		"authorization code. \n%v\n", authURL)
	utils.OpenBrowser(authURL)

	fmt.Print("Authorization code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web %v", err)
	}
	return tok, nil
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape("check-gdrive.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any nead error encoutered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
