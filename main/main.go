package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Vote représente les résultats d'un vote pour un département.
type Vote struct {
	Candidates  map[string]int // Carte des candidats avec leurs votes.
	Departement string         // Le département concerné.
	NbVote      int            // Nombre total de votes dans ce département.
}

func main() {
	// Ouverture du fichier des résultats.
	readFile, err := os.Open("resultats-par-niveau-burvot-t1-france-entiere.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readFile.Close()

	// Initialisation des variables pour stocker les résultats.
	totalVotes := 0
	firstLine := true
	votesByCandidate := make(map[string]int)             // Carte pour stocker les votes par candidat.
	votesByCandidateByDepartment := make(map[string]int) // Carte pour stocker les votes par candidat et par département.
	rankingByDepartment := make(map[string]int)          // Carte pour stocker le classement par département.

	// Scanner pour lire le fichier ligne par ligne.
	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		if firstLine {
			firstLine = false
			continue // Ignorer la première ligne (entête).
		}

		// Créer un objet Vote à partir de la ligne lue.
		vote := createEntryFromString(fileScanner.Text())

		// Mettre à jour les résultats.
		for key, value := range vote.Candidates {
			candidateDepartmentKey := key + "_" + vote.Departement
			votesByCandidateByDepartment[candidateDepartmentKey] += value
			votesByCandidate[key] += value
		}
		rankingByDepartment[vote.Departement] += vote.NbVote
		totalVotes += vote.NbVote
	}

	// Afficher le total des votes.
	fmt.Println("Nombre de votes :", totalVotes)

	// Afficher les votes par candidat.
	printVotesByCandidate(votesByCandidate)

	// Afficher les votes par candidat et par département.
	printVotesByCandidateByDepartment(votesByCandidateByDepartment)

	// Afficher le classement par département.
	printRankingByDepartment(rankingByDepartment)
}

// Afficher les votes par candidat.
func printVotesByCandidate(votesByCandidate map[string]int) {
	for candidate, votes := range votesByCandidate {
		fmt.Println("Candidat :", candidate, "Nombre de votes :", votes)
	}
}

// Afficher les votes par candidat et par département.
func printVotesByCandidateByDepartment(votesByCandidateByDepartment map[string]int) {
	for key, votes := range votesByCandidateByDepartment {
		splitedKey := strings.Split(key, "_")
		candidate := splitedKey[0]
		department := splitedKey[1]
		fmt.Println("Département :", department, "Candidat :", candidate, "Nombre de votes :", votes)
	}
}

// Afficher le classement par département.
func printRankingByDepartment(rankingByDepartment map[string]int) {
	// Convertir la carte en une liste pour pouvoir trier.
	ranking := make([]string, 0, len(rankingByDepartment))
	for key := range rankingByDepartment {
		ranking = append(ranking, key)
	}

	// Trier la liste en fonction du nombre de votes.
	sort.SliceStable(ranking, func(i, j int) bool {
		return rankingByDepartment[ranking[i]] > rankingByDepartment[ranking[j]]
	})

	// Afficher le classement.
	for i, department := range ranking {
		fmt.Println("#", i+1, ":", department)
	}
}

// Créer un objet Vote à partir d'une chaîne de données.
func createEntryFromString(data string) Vote {
	splitedData := strings.Split(data, ";")
	vote := Vote{}
	vote.Candidates = parseCandidates(splitedData)
	vote.Departement = splitedData[1]
	if nbVote, err := strconv.Atoi(splitedData[10]); err == nil {
		vote.NbVote = nbVote
	} else {
		fmt.Println(err)
	}
	return vote
}

// Parser les données des candidats à partir d'une ligne de données.
func parseCandidates(splitedData []string) map[string]int {
	candidates := make(map[string]int)
	voteByCandidateStartIndex := 23
	numberOfRowsBetweenCandidates := 7
	for i := voteByCandidateStartIndex; i < len(splitedData); i += numberOfRowsBetweenCandidates {
		if i >= len(splitedData)-1 {
			break
		}
		candidateName := splitedData[i]
		candidateVoteStr := splitedData[i+2]
		if candidateVote, err := strconv.Atoi(candidateVoteStr); err == nil {
			candidates[candidateName] = candidateVote
		} else {
			fmt.Println(err)
		}
	}
	return candidates
}
