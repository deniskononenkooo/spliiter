package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deniskononenkooo/split-migrations/config"
)

const (
	migrationCommand  = "platform-itsm-migration-cmd -config %s -type %s -log %d-%d-%s-migration-%s.log -partners %s"
	validationCommand = "platform-itsm-migration-cmd -config %s -type %s -log %d-%d-%s-validation-%s.log -validate -partners %s"
)

func main() {
	config, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	ids, err := partnerIDs(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	resultFileName := fmt.Sprintf("%s-commands.txt", config.MigrationType)

	f, err := os.OpenFile(resultFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	result := splitPartners(ids, config)
	for _, r := range result {
		r += "\n\n"
		_, err = f.WriteString(r)
		if err != nil {
			fmt.Println("failed to write to file:", err)
			return
		}
	}

	fmt.Println("Finished")
}

func partnerIDs(c *config.Config) ([]string, error) {
	f, err := os.Open(c.PartnerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file with partner ids: %v", err)
	}
	defer f.Close()

	r, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file with partner ids: %v", err)
	}

	return strings.Split(string(r), ","), err
}

func splitPartners(ids []string, c *config.Config) []string {
	result := make([]string, 0)

	for i, j := 0, c.BatchSize; i < len(ids); i, j = i+c.BatchSize, j+c.BatchSize {
		if j > len(ids) {
			j = len(ids)
		}

		partners := strings.Join(ids[i:j], ",")

		var s string
		if c.Validation {
			s = fmt.Sprintf(validationCommand, c.MigrationConfig, c.MigrationType, i, j, c.MigrationType, c.MigrationEnv, partners)
		} else {
			s = fmt.Sprintf(migrationCommand, c.MigrationConfig, c.MigrationType, i, j, c.MigrationType, c.MigrationEnv, partners)
		}

		result = append(result, s)
	}

	return result
}
