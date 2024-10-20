# FAT

## Info
Use API to get nutrition information from USDA.
  For the API lets try fat secret. I don't want to use the usda. I want macros / portion sizes primarily. 
- [x] Construct search function
- [x] Make json struct for search requests
- [x] Display top 1 calories by default - Generic food_types only.
- [ ] Pretty print information
- [ ] add flag parameters
  - [ ] change total # of results (default: 2)
  - [ ] allow entering desired number of serving size (2 would be 2oz/2g/2servings depending on the item)
  - [ ] Search with ID option (Create function and search request)
  - [ ] Can i configure a menu that allows a user to pick which option they'll do the serving size calc on?

## Links

### References
[Go - On building URL strings](https://www.jacoelho.com/blog/2021/04/go-on-building-url-strings/)

[DigitalOcean - How to Make HTTP Requests in GO](https://www.digitalocean.com/community/tutorials/how-to-make-http-requests-in-go)
[Optional function parameter pattern](https://engineering.01cloud.com/2023/04/13/optional-function-parameter-pattern/)
[Parsing JSON files With Golang](https://tutorialedge.net/golang/parsing-json-with-golang/)

[Better way to read and write JSON file in Golang](https://medium.com/kanoteknologi/better-way-to-read-and-write-json-file-in-golang-9d575b7254f2)

### Imports
[go-querystring](https://github.com/google/go-querystring)
