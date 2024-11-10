#!/usr/bin/env bash

if [ -z "$POSTGRES_URL" ]; then
    echo "Error: POSTGRES_URL environment variable is not set"
    exit 1
fi

get_command_name() {
    local pkg_name=$1
    case $pkg_name in
        "nodejs_22"|"nodejs_20"|"nodejs_18") echo "node" ;;
        "go_1_23" | "go_1_22" | "go_1_21") echo "go" ;;
        "php83" | "php81") echo "php" ;;
        "ruby_3_2"  | "ruby_3_3") echo "ruby" ;;
        "gnutar") echo "tar" ;;
        "gnused") echo "sed" ;;
        "gnumake") echo "make" ;;
        "valgrind-light") echo "valgrind" ;;
        "protobuf_26") echo "protoc" ;;
        "db62") echo "db_verify" ;;
        "sqlite") echo "sqlite3" ;;
        "rocksdb") echo "ldb" ;;
        "jmespath") echo "jp" ;;
        "lua" | "lua5_4_compat" | "lua5_3_compat") echo "lua" ;;
        "python311" | "python312") echo "python" ;;
        "crystal" | "crystal_1_9") echo "crystal" ;;
        "julia" | "julia_1_9") echo "julia" ;;
        "perl" | "perl536") echo "perl" ;;
        *) echo "$pkg_name" ;;
    esac
}

missing_paths=()
db_update_failed=()

packages=$(psql ${POSTGRES_URL} -t -c "SELECT package_id, name, language FROM packages;")
if [ -z "$packages" ]; then
    echo "No packages found in the database."
    exit 0
fi

while IFS="|" read -r package_id name language; do
    name=$(echo $name | xargs)
    language=$(echo $language | xargs)
    if [ -z "$language" ]; then
        echo "Running nix-shell for $name (type: system)..."
        base_name=$(echo $name | cut -d'_' -f1)
        cmd_name=$(get_command_name $name)
        echo "Trying command name: $cmd_name"
        pkg_path=$(nix-shell -p $name --run "which $cmd_name 2>/dev/null || echo ''")

        if [ -z "$pkg_path" ] && echo $name | grep -q '[0-9]'; then
            stripped_name=$(echo $base_name | sed 's/[0-9]//g')
            if [ "$stripped_name" != "$cmd_name" ]; then
                echo "Trying alternate name: $stripped_name"
                pkg_path=$(nix-shell -p $name --run "which $stripped_name 2>/dev/null || echo ''")
            fi
        fi

        if [ ! -z "$pkg_path" ]; then
            echo "Found path: $pkg_path for $name"
            pkg_dir=$(dirname $pkg_path)
            if ! psql ${POSTGRES_URL} -c "UPDATE packages SET store_path='$pkg_dir' WHERE package_id='$package_id'" 2>/dev/null; then
                db_update_failed+=("$name")
            fi
        else
            echo "Warning: Could not find path for $name"
            missing_paths+=("$name")
        fi
    else
        echo "Running nix-shell for $language.$name (type: language)..."
        nix-shell -p $language.$name --run "exit"
    fi
done <<< "$packages"

echo -e "\nAll packages processed.\n"

if [ ${#missing_paths[@]} -ne 0 ]; then
    echo "Missing commands for these packages:"
    for pkg in "${missing_paths[@]}"; do
        echo "- $pkg"
    done
fi

if [ ${#db_update_failed[@]} -ne 0 ]; then
    echo -e "\nDatabase update failed for these packages:"
    for pkg in "${db_update_failed[@]}"; do
        echo "- $pkg"
    done
fi

if [ ${#missing_paths[@]} -eq 0 ] && [ ${#db_update_failed[@]} -eq 0 ]; then
    echo "All packages processed successfully with no errors."
fi