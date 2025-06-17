package env

import (
	"bufio"
	"os"
	"strings"
)

// Manager handles operations related to the .env file
type Manager struct {
	filePath string
}

// NewManager creates a new instance of the environment variables manager
func NewManager(filePath string) *Manager {
	return &Manager{
		filePath: filePath,
	}
}

// Load loads environment variables from the file
func (m *Manager) Load() (map[string]string, error) {
	file, err := os.Open(m.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envVars[key] = value
		}
	}

	return envVars, scanner.Err()
}

// Save saves environment variables to the file
func (m *Manager) Save(vars map[string]string) error {
	file, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range vars {
		if _, err := writer.WriteString(key + "=" + value + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
