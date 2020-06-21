package cmd

import "fmt"

type configValidator func(*Config) error

func validateConfig(config *Config) error {
	validators := [...]configValidator{
		validateUniqueTargetOutfiles,
	}

	for _, validator := range validators {
		if err := validator(config); err != nil {
			return err
		}
	}
	return nil
}

func validateUniqueTargetOutfiles(config *Config) error {
	outfiles := map[string]bool{}
	for _, target := range config.Targets {
		if outfiles[target.Outfile] {
			return fmt.Errorf("two targets have the same outfile, %s", target.Outfile)
		}
		outfiles[target.Outfile] = true
	}
	return nil
}

func validateOutfilePresent(config *Config) error {
	for i, target := range config.Targets {
		if target.Outfile == "" {
			return fmt.Errorf("target %d doesn't have an 'outfile' specified", i)
		}
	}
	return nil
}

func validateAtLeastOneSource(config *Config) error {
	for i, target := range config.Targets {
		if len(target.Sources) == 0 {
			return fmt.Errorf("target %d doesn't have any 'sources' specified", i)
		}
	}
	return nil
}
