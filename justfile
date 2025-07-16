build:
  goreleaser build --clean --snapshot

install version:
  curl -sSfL https://github.com/moonlight8978/cmd-runner/releases/download/{{version}}/cmd-runner_darwin_arm64.tar.gz | tar -xzf - -C tmp
  sudo mv tmp/c7r /usr/local/bin/c7r
