. "$(dirname "$0")/functions.sh"
setup
install

f=".husky/pre-commit"

${HUSKY_GO} install

${HUSKY_GO} add $f "foo"
grep -m 1 _ $f && grep foo $f && ok

${HUSKY_GO} add .husky/pre-commit "bar"
grep -m 1 _ $f && grep foo $f && grep bar $f && ok

${HUSKY_GO} set .husky/pre-commit "baz"
grep -m 1 _ $f && grep foo $f || grep bar $f || grep baz $f && ok
