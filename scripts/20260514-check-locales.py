#!/usr/bin/env -S uv run
import os
from os import path
from pathlib import Path
import json

locales_dir = path.abspath('../i18n')

resolved_data = {}


def flatten_dict(dick: dict, prefix: str = ''):
    result = {}
    for (key, value) in dick.items():
        merged_key = key if prefix == '' else f"{prefix}.{key}"

        if isinstance(value, dict):
            result |= flatten_dict(value, merged_key)

        else:
            result[merged_key] = value

    return result


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
            missing_keys: list[str] = mine - resolved_data[other].keys()
            if len(missing_keys) > 0:
                valid = False
                for key in missing_keys:
                    is_variant = key.split('.')[-1] in ['many', 'other', 'few']
                    if not is_variant:
                        print(f"Key {key} is present in locale {locale} but not in {other}")

    if not valid:
        exit(1)


if __name__ == '__main__':
    main()
