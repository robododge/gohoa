# Obtain owners directory

## Pull full directory
1. Log into the associateion voice website of your neighborhood
2. Open the dev tools of your browser, I used chrome
3. In the javascript console, Prepare a search settings object that will bring back the entire direcory, set the pageSize to be above the max neighbor count in your neighborhood
```
const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
  g_search: "",
  search_letter: "",
  search_all: "0",
  page: 1,
  pageSize: 750
})
};
```
4. Now, post search request in the browser's javascript console
```
fetch('https://SITED_DNS_NAME/Member/SearchData', requestOptions)
.then(response => response.json())
.then(data => console.log("my data",data))
```
5. Select the "data" object in the console and right click to "copy object"
6. Using a Mac, go to comman line and type `pbpaste > hoa_dir.json` to create a new file with the raw json.

## Cleaning and populating
Next, we need to clean the raw json in to entires that are represented by the Go Structs in this project
1.  run the cmd/clean/cleaner.go program, cd to the cmd/clean directory and run `go run cleaner.go`

