package pkg_test


import (
  "io"
  "os"
  "path/filepath"
  "runtime"
  "testing"

  "github.com/tronicum/punchbag-cube-testsuite/multitool/pkg"
)


func copyConfig(t *testing.T, profile string) string {
  t.Helper()
  tmp := t.TempDir()
  mtconfig := filepath.Join(tmp, "multitool/.mtconfig", profile)
  err := os.MkdirAll(mtconfig, 0o755)
  if err != nil {
    t.Fatalf("failed to create temp config dir: %v", err)
  }
  // Find repo root using this test file's location
  _, thisFile, _, _ := runtime.Caller(0)
  // Find the punchbag-cube-testsuite root
  dir := filepath.Dir(thisFile)
  for i := 0; i < 5; i++ {
    if filepath.Base(dir) == "punchbag-cube-testsuite" {
      break
    }
    dir = filepath.Dir(dir)
  }
  repoRoot := dir
  src := filepath.Join(repoRoot, "examples", "mtconfig", profile, "config.yaml")
  dst := filepath.Join(mtconfig, "config.yaml")
  in, err := os.Open(src)
  if err != nil {
    t.Fatalf("failed to open example config: %v", err)
  }
  defer in.Close()
  out, err := os.Create(dst)
  if err != nil {
    t.Fatalf("failed to create temp config: %v", err)
  }
  defer out.Close()
  if _, err := io.Copy(out, in); err != nil {
    t.Fatalf("failed to copy config: %v", err)
  }
  return tmp
}

func TestLoadMTConfig_DefaultProfile(t *testing.T) {
  os.Setenv("AWS_ACCESS_KEY_ID", "test-access")
  os.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
  tmp := copyConfig(t, "default")
  oldwd, _ := os.Getwd()
  defer os.Chdir(oldwd)
  os.Chdir(tmp)
  cfg, err := pkg.LoadMTConfig("default")
  if err != nil {
  t.Fatalf("failed to load default profile: %v", err)
  }
  if cfg.Provider != "aws" {
  t.Errorf("expected provider aws, got %s", cfg.Provider)
  }
  if cfg.Region != "us-east-1" {
  t.Errorf("expected region us-east-1, got %s", cfg.Region)
  }
  if cfg.Credentials["access_key"] != "test-access" {
  t.Errorf("expected access_key from env, got %s", cfg.Credentials["access_key"])
  }
}


func TestLoadMTConfig_AwsDevProfile(t *testing.T) {
  os.Setenv("AWS_ACCESS_KEY_ID", "dev-access")
  os.Setenv("AWS_SECRET_ACCESS_KEY", "dev-secret")
  tmp := copyConfig(t, "aws-dev")
  oldwd, _ := os.Getwd()
  defer os.Chdir(oldwd)
  os.Chdir(tmp)
  cfg, err := pkg.LoadMTConfig("aws-dev")
  if err != nil {
  t.Fatalf("failed to load aws-dev profile: %v", err)
  }
  if cfg.Region != "eu-central-1" {
  t.Errorf("expected region eu-central-1, got %s", cfg.Region)
  }
  if v, ok := cfg.Settings["versioning"].(bool); !ok || !v {
  t.Errorf("expected versioning true, got %v", cfg.Settings["versioning"])
  }
  if cfg.Credentials["access_key"] != "dev-access" {
  t.Errorf("expected access_key from env, got %s", cfg.Credentials["access_key"])
  }
}
