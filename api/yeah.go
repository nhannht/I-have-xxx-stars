package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "GitHub token is not set", http.StatusUnauthorized)
		return
	}

	// Parse the webhook payload
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse webhook payload", http.StatusBadRequest)
		return
	}

	// Extract repository data
	repo, ok := payload["repository"].(map[string]interface{})
	if !ok {
		http.Error(w, "Failed to extract repository data", http.StatusBadRequest)
		return
	}

	// Get the star count
	stargazersCount, ok := repo["stargazers_count"].(float64)
	if !ok {
		http.Error(w, "Failed to get stargazers count", http.StatusBadRequest)
		return
	}

	// Update the repository description using a PATCH request
	repoFullName, ok := repo["full_name"].(string)
	if !ok {
		http.Error(w, "Failed to get repository full name", http.StatusBadRequest)
		return
	}

	newDescription := fmt.Sprintf("I have %.0f stars", stargazersCount)
	if err := updateRepoDescription(token, repoFullName, newDescription); err != nil {
		http.Error(w, "Failed to update repository description", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Repository description updated to '%s'"}`, newDescription)))
}

func updateRepoDescription(token, repoFullName, newDescription string) error {
	url := "https://api.github.com/repos/" + repoFullName
	reqBody := map[string]string{"description": newDescription}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API returned status: %s, body: %s", resp.Status, string(body))
	}

	return nil
}
