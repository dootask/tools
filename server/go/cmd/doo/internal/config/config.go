// Package config 读写 doo 的本地配置文件（含 server 与 token）。
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 是落盘的 CLI 配置。
type Config struct {
	Server string `json:"server,omitempty"`
	Token  string `json:"token,omitempty"`
}

// Dir 返回配置目录：$XDG_CONFIG_HOME/doo 或 ~/.config/doo。
func Dir() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "doo")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".doo"
	}
	return filepath.Join(home, ".config", "doo")
}

// Path 返回配置文件路径。
func Path() string {
	return filepath.Join(Dir(), "config.json")
}

// Load 读取配置；文件不存在时返回空配置而非错误。
func Load() (Config, error) {
	var c Config
	b, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return c, err
	}
	if len(b) == 0 {
		return c, nil
	}
	err = json.Unmarshal(b, &c)
	return c, err
}

// Save 以 0600 权限写入配置（目录 0700）。
func Save(c Config) error {
	if err := os.MkdirAll(Dir(), 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), b, 0o600)
}
