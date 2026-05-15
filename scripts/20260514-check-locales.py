#!/usr/bin/env -S uv run

import os
from os import path, stat
from pathlib import Path
import json

def find_git_root() -> Path:
    dir = start = Path(os.getcwd())
    while dir != "/":
        if Path(dir).joinpath(".git").is_dir():
            return dir
        if dir == dir.parent:
            break
        dir = dir.parent

    raise FileNotFoundError(f"Git root not found above {start}")


def flatten_dict(dick: dict, prefix: str = ''):
    result = {}
    for (key, value) in dick.items():
        merged_key = key if prefix == '' else f"{prefix}.{key}"

        if isinstance(value, dict):
            result |= flatten_dict(value, merged_key)

        else:
            result[merged_key] = value

    return result

def key_present(key: str, other: dict) -> bool:
    if key in other:
        return True

    segments = key.split('.')
    is_variant = segments[-1] in ALL_VARIANTS
    if is_variant:
        base_key = '.'.join(segments[:-1]) if is_variant else key
        if base_key in other:
            return True
        for variant in ALL_VARIANTS:
            variant_key = f"{base_key}.{variant}"
            if variant_key in other:
                return True
    else:
        for variant in ALL_VARIANTS:
            variant_key = f"{key}.{variant}"
            if variant_key in other:
                return True
    return False

locales_dir = Path(find_git_root()).joinpath("i18n")
resolved_data = {}

ALL_VARIANTS = ['many', 'other', 'few', 'one']


def main():
    for file in Path(locales_dir).glob("*.json"):
        locale = file.stem

        with open(file) as f:
            data = json.load(f)

        resolved_data[locale] = flatten_dict(data)

    valid = True

    for locale in resolved_data.keys():
        other_locales = resolved_data.keys() - [locale]

        mine = resolved_data[locale].keys()

        for other in other_locales:
            other_dict = resolved_data[other]
            for key in mine:
                if key_present(key, other_dict):
                    continue

                valid = False
                print(f"Key {key} is present in locale {locale} but not in {other}")

    if not valid:
        exit(1)

if __name__ == '__main__':
    main()
