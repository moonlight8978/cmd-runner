build:
  goreleaser build --clean --snapshot

install:
  curl -sSfL https://github.com/moonlight8978/cmd-runner/releases/download/v1.0.2/cmd-runner_darwin_arm64.tar.gz | tar -xzf - -C tmp
  sudo mv tmp/c7r /usr/local/bin/c7r
