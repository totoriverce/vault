while IFS= read -r line; do
found=$(grep -r \"github.com/hashicorp/$line\" ./)
num_lines=$(echo -e "$found" | wc -l)

if [ "$num_lines" -eq 1 ]; then
echo "$line NOT in minimal"
else
    echo "$line in minimal"
fi
echo "\n$found\n"

done < "$1"
