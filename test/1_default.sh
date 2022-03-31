. "$(dirname "$0")/functions.sh"
setup

${HUSKY_GO} install

# Test core.hooksPath
expect_hooksPath_to_be ".husky"

# Test pre-commit
git add package.json
${HUSKY_GO} add .husky/pre-commit "echo \"pre-commit\" && exit 1"
expect 1 "git commit -m foo"

# Uninstall
${HUSKY_GO} uninstall
expect 1 "git config core.hooksPath"
