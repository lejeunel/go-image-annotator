#!/usr/bin/env bash

usage() {
    echo "Usage: $0 <phase>"
    echo "Where <phase> is 'before' or 'after'"
    exit 1
}

validate_phase() {
    case "$1" in
    "before" | "after")
        return 1
        ;;
    *)
        return 0
        ;;
    esac
}

ensure_env_var() {
    if [ -z "$1" ]; then
        echo "❌ $2 env var must be set prior to running this script"
        exit 1
    fi
    echo "✅ Found $2"
}

parse_nix_store() {
    find /nix/store -mindepth 1 -maxdepth 1 ! -name \*.drv ! -name ".links" |
        sort
}

make_nix_copy_query_arg() {
    if [ -z "$NIX_COMPRESS_CACHE" ]; then
        echo ""
    elif [ "$NIX_COMPRESS_CACHE" = true ]; then
        echo ""
    elif [ "$NIX_COMPRESS_CACHE" = false ]; then
        echo "?compression=none"
    else
        echo ""
    fi
}

after() {

    # Assemble list of new builds
    parse_nix_store >/nix/.after
    comm -13 /nix/.before /nix/.after >/nix/.new

    COPY_ARGS=$(make_nix_copy_query_arg)
    echo "$COPY_ARGS"

    [[ -s /nix/.new ]] ||
        exit 0

    echo -e "copying new items from /nix/store to $NIX_CACHE_DIR with args: \"$COPY_ARGS\""
    xargs -a /nix/.new nix copy --to "file://$NIX_CACHE_DIR$COPY_ARGS"
    echo -e "done"

}

before() {
    parse_nix_store >/nix/.before

    echo "extra-substituters = file://$NIX_CACHE_DIR?priority=10&trusted=true" \
        >>/etc/nix/nix.conf

}

main() {
    # Validate arguments
    if [[ $# -ne 1 ]]; then
        usage
    fi

    if validate_phase "$1"; then
        echo "Error: Argument must be 'before' or 'after'" >&2
    fi

    ensure_env_var "$NIX_CACHE_DIR" "NIX_CACHE_DIR"

    if [[ "$1" == "before" ]]; then
        before
    else
        after
    fi

}

main "$@"
