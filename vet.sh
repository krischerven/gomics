#!/bin/bash -e
echo "Running go vet:"
cd src
x=$(go vet 2>&1 | perl -wln -M'Term::ANSIColor' -e 'print "\e[1;91m", "$_", "\e[0m"')
if [ -z "$x" ]; then
	# if there were no bad things, tell us that
	echo "  No issues found ✔️"
else
	echo "$x"
fi
