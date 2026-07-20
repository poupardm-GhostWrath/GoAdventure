go build \
-ldflags=" \
-X github.com/poupardm-GhostWrath/GoAdventure/internal/build.GitCommit=$(git rev-parse --short HEAD) \
-X github.com/poupardm-GhostWrath/GoAdventure/internal/build.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
-o ./GoAdventure \
&& ./GoAdventure