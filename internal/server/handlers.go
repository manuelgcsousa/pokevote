package server

import (
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	pokemon := s.Database.GetTwoRandomPokemon()
	s.Tmpl.ExecuteTemplate(w, "index", pokemon)
}

func (s *Server) pokemonHandler(w http.ResponseWriter, r *http.Request) {
	pokemon := s.Database.GetTwoRandomPokemon()
	s.Tmpl.ExecuteTemplate(w, "pokemon", pokemon)
}

type VoteFormData struct {
	PokemonIds []string
	VoteId     string
}

func (s *Server) voteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		slog.Default().Error("Error while parsing form data: " + err.Error())
		return
	}

	voteData := VoteFormData{
		PokemonIds: r.Form["pokemonIds"],
		VoteId:     r.Form.Get("voteId"),
	}

	// Extract ID with vote won
	var pokemonIdWithVoteWon int
	pokemonIdWithVoteWon, err = strconv.Atoi(voteData.VoteId)
	if err != nil {
		slog.Default().Error("Error while extracting vote won: " + err.Error())
		return
	}

	// Extract ID with vote lost
	var pokemonIdWithVoteLost int
	for _, id := range voteData.PokemonIds {
		if id != voteData.VoteId {
			pokemonIdWithVoteLost, err = strconv.Atoi(id)
			if err != nil {
				slog.Default().Error("Error while extracting vote lost: " + err.Error())
				return
			}
			break
		}
	}

	// Update both vote counters
	// (background process, any errors are logged)
	go s.Database.UpdatePokemonVoteCounters(pokemonIdWithVoteWon, pokemonIdWithVoteLost)

	// Return two new pokemon
	pokemon := s.Database.GetTwoRandomPokemon()
	s.Tmpl.ExecuteTemplate(w, "pokemon", pokemon)
}

func (s *Server) resultsHandler(w http.ResponseWriter, r *http.Request) {
	pokemon := s.Database.GetAllPokemon()
	s.Tmpl.ExecuteTemplate(w, "results", pokemon)
}
