. "$(dirname -- "$0")/functions.sh"
setup
install

${HUSKY_GO} install

# Test core.hooksPath
expect_hooksPath_to_be ".husky"

# Test pre-commit with 127 exit code
git add package.json
${HUSKY_GO} add .husky/pre-commit "exit 127"
expect 1 "git commit -m foo"
