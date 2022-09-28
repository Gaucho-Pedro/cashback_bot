package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

/*
// Retrieve a token, saves the token, then returns the generated client.

	func getClient(config *oauth2.Config) *http.Client {
		// The file token.json stores the user's access and refresh tokens, and is
		// created automatically when the authorization flow completes for the first
		// time.
		tokFile := "token.json"
		tok, err := tokenFromFile(tokFile)
		if err != nil {
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		}
		return config.Client(context.Background(), tok)
	}

// Request a token from the web, then returns the retrieved token.

	func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)

		var authCode string
		if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
		}

		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
		}
		return tok
	}

// Retrieves a token from a local file.

	func tokenFromFile(file string) (*oauth2.Token, error) {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		tok := &oauth2.Token{}
		err = json.NewDecoder(f).Decode(tok)
		return tok, err
	}

// Saves a token to a file path.

	func saveToken(path string, token *oauth2.Token) {
		fmt.Printf("Saving credential file to: %s\n", path)
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("Unable to cache oauth token: %v", err)
		}
		defer f.Close()
		json.NewEncoder(f).Encode(token)
	}
*/
func main() {
	ctx := context.Background()
	file := option.WithCredentialsFile("configs/credentials.json")
	driveService, err := drive.NewService(ctx, file)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	_, err = driveService.Drives.Get("1f__dmBg5wnvmIAYSVett6l_OxffmNvrT").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	/*	fmt.Println("Files:")
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				fmt.Printf("%s (%s)\n", i.Name, i.Id)
			}
		}*/

}
func sheetsTest(ctx context.Context, clientOption option.ClientOption, spreadsheetId string) {
	sheetsService, err := sheets.NewService(ctx, clientOption)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

	//spreadsheetId := "1DKQC4gZi3qsq8imK5x4TDXXKzJx-pT6NiiOVlJ-rMbU"
	values := [][]interface{}{
		{"Ячейка 1", "Ячейка 2"},
		{" ", 456}}

	valueRange := &sheets.ValueRange{
		MajorDimension:  "ROWS",
		Range:           "Лист1!A1:B",
		Values:          values,
		ServerResponse:  googleapi.ServerResponse{},
		ForceSendFields: nil,
		NullFields:      nil,
	}

	request := &sheets.BatchUpdateValuesRequest{
		Data:                         []*sheets.ValueRange{valueRange},
		IncludeValuesInResponse:      false,
		ResponseDateTimeRenderOption: "",
		ResponseValueRenderOption:    "",
		ValueInputOption:             "USER_ENTERED",
		ForceSendFields:              nil,
		NullFields:                   nil,
	}

	resp, err := sheetsService.Spreadsheets.Values.BatchUpdate(spreadsheetId, request).Do()

	//sheetsService.Spreadsheets.BatchUpdate()
	//sheetsService.Spreadsheets.Values.Update()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(resp)

	/*	spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
		readRange := "Class Data!A2:E"
		resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
		}

		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			fmt.Println("Name, Major:")
			for _, row := range resp.Values {
				// Print columns A and E, which correspond to indices 0 and 4.
				fmt.Printf("%s, %s\n", row[0], row[4])
			}
		}*/
}
