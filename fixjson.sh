#!/bin/bash

function modify_const {
    local json="$1"
    
    # If the current item is an object and has the key 'const', modify its structure.
    json=$(echo "$json" | jq 'if type == "object" and has("const") then
        .const |= { "@type": (type), "value": . }
    else
        .
    end')
    
    # If it's an object, iterate through its properties and process recursively.
    keys=$(echo "$json" | jq -r 'if type == "object" then keys[] else empty end')
    for key in $keys; do
        value=$(echo "$json" | jq ".\"$key\"")
        modified_value=$(modify_const "$value")
        json=$(echo "$json" | jq ".\"$key\" = $modified_value")
    done
    
    # If it's an array, iterate through its items and process recursively.
    length=$(echo "$json" | jq 'if type == "array" then length else -1 end')
    if [ "$length" -ge 0 ]; then
        for i in $(seq 0 $(($length-1))); do
            item=$(echo "$json" | jq ".[$i]")
            modified_item=$(modify_const "$item")
            json=$(echo "$json" | jq ".[$i] = $modified_item")
        done
    fi
    
    echo "$json"
}

for file in ./commands/*.json; do
    echo "Processing $file..."
    json=$(cat "$file")
    modified_json=$(modify_const "$json")
    echo "$modified_json" > "$file"
done
