package main 

import("strings")

func clean_string(msg string) (cleaned_string string){
	words := strings.Split(msg," ")
	banned_words := map[string]struct{}{"kerfuffle":{},"sharbert":{},"fornax":{}}
	const censor  = "****"
	for i,word := range words {
		if _, ok := banned_words[strings.ToLower(word)]; ok {
			words[i] = censor
		}
	}
	cleaned_string = strings.Join(words," ")
	return cleaned_string
}