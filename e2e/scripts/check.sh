#!/bin/bash
set -euo pipefail

CCACHE_CONFIG="$HOME/.bitrise/cache/ccache/config.json"
XCELERATE_CONFIG="$HOME/.bitrise-xcelerate/config.json"
GRADLE_INIT="$HOME/.gradle/init.d/bitrise-build-cache.init.gradle.kts"

assert_file_exists() {
  if [ ! -f "$1" ]; then
    echo "❌ Expected file not found: $1"
    exit 1
  fi
  echo "✅ Found: $1"
}

assert_file_absent() {
  if [ -f "$1" ]; then
    echo "❌ File should not exist when feature is disabled: $1"
    exit 1
  fi
  echo "✅ Correctly absent: $1"
}

assert_json_field_nonempty() {
  local file="$1" field="$2"
  local value
  value=$(jq -r "$field" "$file")
  if [ -z "$value" ] || [ "$value" = "null" ]; then
    echo "❌ $file: $field is empty or null"
    exit 1
  fi
  echo "✅ $file: $field = $value"
}

assert_json_field_equals() {
  local file="$1" field="$2" expected="$3"
  local value
  value=$(jq -r "$field" "$file")
  if [ "$value" != "$expected" ]; then
    echo "❌ $file: $field should be '$expected', got '$value'"
    exit 1
  fi
  echo "✅ $file: $field = $value"
}

assert_env_equals() {
  local var="$1" expected="$2"
  local value="${!var:-}"
  if [ "$value" != "$expected" ]; then
    echo "❌ \$$var should be '$expected', got '$value'"
    exit 1
  fi
  echo "✅ \$$var=$value"
}

assert_env_nonempty() {
  local var="$1"
  local value="${!var:-}"
  if [ -z "$value" ]; then
    echo "❌ \$$var is not set or empty"
    exit 1
  fi
  echo "✅ \$$var=$value"
}

# C++ cache checks (cpp is enabled whenever gradle is enabled)
if [ "$GRADLE_CACHE_ENABLED" = "true" ]; then
  echo "--- C++ cache ---"
  assert_file_exists "$CCACHE_CONFIG"
  assert_json_field_nonempty "$CCACHE_CONFIG" ".ipcEndpoint"
  assert_json_field_equals   "$CCACHE_CONFIG" ".enabled"     "true"
  assert_json_field_equals   "$CCACHE_CONFIG" ".pushEnabled" "true"

  assert_env_equals   "CCACHE_NOHASHDIR"            "true"
  assert_env_equals   "CCACHE_REMOTE_ONLY"          "true"
  assert_env_nonempty "CCACHE_REMOTE_STORAGE"
  assert_env_equals   "CMAKE_CXX_COMPILER_LAUNCHER" "ccache"
  assert_env_equals   "CMAKE_C_COMPILER_LAUNCHER"   "ccache"
else
  assert_file_absent "$CCACHE_CONFIG"
fi

# Xcode cache checks
if [ "$XCODE_CACHE_ENABLED" = "true" ]; then
  echo "--- Xcode cache ---"
  assert_file_exists "$XCELERATE_CONFIG"
  assert_json_field_nonempty "$XCELERATE_CONFIG" ".proxySocketPath"
  assert_json_field_equals   "$XCELERATE_CONFIG" ".buildCacheEnabled" "true"
else
  assert_file_absent "$XCELERATE_CONFIG"
fi

# Gradle cache checks
if [ "$GRADLE_CACHE_ENABLED" = "true" ]; then
  echo "--- Gradle cache ---"
  assert_file_exists "$GRADLE_INIT"
  grep -q "BitriseBuildCache" "$GRADLE_INIT" || {
    echo "❌ $GRADLE_INIT does not contain expected cache plugin reference"
    exit 1
  }
  echo "✅ $GRADLE_INIT contains BitriseBuildCache plugin reference"
else
  assert_file_absent "$GRADLE_INIT"
fi

echo ""
echo "✅ All checks passed"
