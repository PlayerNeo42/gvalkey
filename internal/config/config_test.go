package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite defines a test suite for config package
type ConfigTestSuite struct {
	suite.Suite
}

// SetupTest runs before each test
func (s *ConfigTestSuite) SetupTest() {
	s.cleanupEnv()
}

// TearDownTest runs after each test
func (s *ConfigTestSuite) TearDownTest() {
	s.cleanupEnv()
}

// TestDefaultValues tests loading configuration with default values
func (s *ConfigTestSuite) TestDefaultValues() {
	config, err := Load()

	s.Require().NoError(err, "Loading default config should not fail")
	s.Equal("0.0.0.0", config.Host, "Default host should be 0.0.0.0")
	s.Equal(6379, config.Port, "Default port should be 6379")
	s.Equal("INFO", config.LogLevel, "Default log level should be INFO")
}

// TestCustomValues tests loading configuration with custom environment variables
func (s *ConfigTestSuite) TestCustomValues() {
	// Set custom environment variables
	os.Setenv("GVK_HOST", "127.0.0.1")
	os.Setenv("GVK_PORT", "8080")
	os.Setenv("GVK_LOG_LEVEL", "debug")

	config, err := Load()

	s.Require().NoError(err, "Loading custom config should not fail")
	s.Equal("127.0.0.1", config.Host, "Host should match custom value")
	s.Equal(8080, config.Port, "Port should match custom value")
	s.Equal("DEBUG", config.LogLevel, "Log level should be converted to uppercase")
}

// TestLogLevelUpperCaseConversion tests that log levels are converted to uppercase
func (s *ConfigTestSuite) TestLogLevelUpperCaseConversion() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase debug", "debug", "DEBUG"},
		{"lowercase info", "info", "INFO"},
		{"lowercase warn", "warn", "WARN"},
		{"lowercase error", "error", "ERROR"},
		{"uppercase debug", "DEBUG", "DEBUG"},
		{"mixed case info", "Info", "INFO"},
		{"mixed case warn", "WaRn", "WARN"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.cleanupEnv()
			os.Setenv("GVK_LOG_LEVEL", tc.input)

			config, err := Load()

			s.Require().NoError(err, "Loading config should not fail")
			s.Equal(tc.expected, config.LogLevel,
				"Log level should be converted to uppercase")
		})
	}
}

// TestInvalidPorts tests various invalid port configurations
func (s *ConfigTestSuite) TestInvalidPorts() {
	testCases := []struct {
		name string
		port string
	}{
		{"zero port", "0"},
		{"negative port", "-1"},
		{"port exceeds maximum", "65536"},
		{"non-numeric port", "abc"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.cleanupEnv()
			os.Setenv("GVK_PORT", tc.port)

			_, err := Load()
			s.Error(err, "Invalid port should return error")
		})
	}
}

// TestInvalidHosts tests various invalid host configurations
func (s *ConfigTestSuite) TestInvalidHosts() {
	testCases := []struct {
		name string
		host string
	}{
		{"invalid IP", "999.999.999.999"},
		{"invalid hostname", "invalid..hostname"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.cleanupEnv()
			os.Setenv("GVK_HOST", tc.host)

			_, err := Load()
			s.Error(err, "Invalid host should return error")
		})
	}
}

// TestEmptyHost tests that empty host string uses default value
func (s *ConfigTestSuite) TestEmptyHost() {
	os.Setenv("GVK_HOST", "")

	config, err := Load()

	s.Require().NoError(err, "Empty host should use default value")
	s.Equal("0.0.0.0", config.Host,
		"Empty host should fallback to default value")
}

// TestInvalidLogLevels tests various invalid log level configurations
func (s *ConfigTestSuite) TestInvalidLogLevels() {
	testCases := []string{
		"INVALID",
		"TRACE",
		"FATAL",
		"123",
	}

	for _, logLevel := range testCases {
		s.Run(logLevel, func() {
			s.cleanupEnv()
			os.Setenv("GVK_LOG_LEVEL", logLevel)

			_, err := Load()
			s.Error(err, "Invalid log level should return error")
		})
	}
}

// TestEmptyLogLevel tests that empty log level string uses default value
func (s *ConfigTestSuite) TestEmptyLogLevel() {
	os.Setenv("GVK_LOG_LEVEL", "")

	config, err := Load()

	s.Require().NoError(err, "Empty log level should use default value")
	s.Equal("INFO", config.LogLevel,
		"Empty log level should fallback to default value")
}

// TestValidHosts tests various valid host configurations
func (s *ConfigTestSuite) TestValidHosts() {
	testCases := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"192.168.1.1",
		"example.com",
		"sub.example.com",
	}

	for _, host := range testCases {
		s.Run(host, func() {
			s.cleanupEnv()
			os.Setenv("GVK_HOST", host)

			config, err := Load()

			s.Require().NoError(err, "Valid host should not return error")
			s.Equal(host, config.Host, "Host should match set value")
		})
	}
}

// TestValidPorts tests various valid port configurations
func (s *ConfigTestSuite) TestValidPorts() {
	testCases := []struct {
		port     string
		expected int
	}{
		{"1", 1},
		{"80", 80},
		{"443", 443},
		{"6379", 6379},
		{"8080", 8080},
		{"65535", 65535},
	}

	for _, tc := range testCases {
		s.Run(tc.port, func() {
			s.cleanupEnv()
			os.Setenv("GVK_PORT", tc.port)

			config, err := Load()

			s.Require().NoError(err, "Valid port should not return error")
			s.Equal(tc.expected, config.Port, "Port should match expected value")
		})
	}
}

// TestValidateConfigWithValidInput tests validateConfig function with valid configuration
func (s *ConfigTestSuite) TestValidateConfigWithValidInput() {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		LogLevel: "INFO",
	}

	err := validateConfig(config)
	s.NoError(err, "Valid config should pass validation")
}

// TestValidateConfigWithInvalidInput tests validateConfig function with invalid configurations
func (s *ConfigTestSuite) TestValidateConfigWithInvalidInput() {
	testCases := []struct {
		name   string
		config *Config
	}{
		{
			name: "invalid host",
			config: &Config{
				Host:     "invalid..host",
				Port:     8080,
				LogLevel: "INFO",
			},
		},
		{
			name: "invalid port",
			config: &Config{
				Host:     "localhost",
				Port:     0,
				LogLevel: "INFO",
			},
		},
		{
			name: "invalid log level",
			config: &Config{
				Host:     "localhost",
				Port:     8080,
				LogLevel: "INVALID",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := validateConfig(tc.config)
			s.Error(err, "Invalid config should fail validation")
		})
	}
}

// cleanupEnv cleans up environment variables used in tests
func (s *ConfigTestSuite) cleanupEnv() {
	envVars := []string{"GVK_HOST", "GVK_PORT", "GVK_LOG_LEVEL"}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

// TestConfigSuite runs the config test suite
func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Standalone tests for backward compatibility and additional coverage

// TestLoad_Standalone tests the Load function independently
func TestLoad_Standalone(t *testing.T) {
	// Clean environment
	envVars := []string{"GVK_HOST", "GVK_PORT", "GVK_LOG_LEVEL"}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
	defer func() {
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}
	}()

	config, err := Load()

	require.NoError(t, err)
	require.Equal(t, "0.0.0.0", config.Host)
	require.Equal(t, 6379, config.Port)
	require.Equal(t, "INFO", config.LogLevel)
}

// TestValidateConfig_Standalone tests the validateConfig function independently
func TestValidateConfig_Standalone(t *testing.T) {
	// Test valid config
	validConfig := &Config{
		Host:     "localhost",
		Port:     8080,
		LogLevel: "DEBUG",
	}

	err := validateConfig(validConfig)
	require.NoError(t, err)

	// Test invalid config
	invalidConfig := &Config{
		Host:     "",
		Port:     0,
		LogLevel: "INVALID",
	}

	err = validateConfig(invalidConfig)
	require.Error(t, err)
}
