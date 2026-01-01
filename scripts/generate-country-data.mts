#!/usr/bin/env -S deno run --allow-write

import path from "node:path";

function isCountry(isoCode: string): boolean {
  const intl = new Intl.DisplayNames("en", { type: "region" });
  const formatted = intl.of(isoCode);
  return Boolean(formatted && formatted !== isoCode);
}

function isException(isoCode: string): boolean {
  const exceptions = [
    "BU",
    "DD",
    "DY",
    "EU",
    "EZ",
    "FX",
    "RS",
    "SU",
    "UK",
    "NH",
    "HV",
    "RH",
    "AN",
    "VN",
    "TP",
    "CQ",
    "UN",
    "YD",
    "YU",
    "ZR",
    "ZZ",
  ];

  return (
    (isoCode >= "QM" && isoCode <= "QZ") ||
    (isoCode >= "XA" && isoCode <= "XZ") ||
    exceptions.includes(isoCode)
  );
}

const codes: string[] = [];

for (let x = 0; x < 26; x++) {
  for (let y = 0; y < 26; y++) {
    const code = String.fromCharCode(65 + x) + String.fromCharCode(65 + y);
    if (isCountry(code) && !isException(code)) {
      codes.push(code);
    }
  }
}

function formatCountry(locale: string, isoCode: string): string {
  return new Intl.DisplayNames(locale, { type: "region" }).of(isoCode)!;
}

const options = codes.map((iso) => ({
  value: iso,
  labelPl: formatCountry("pl", iso),
  labelEn: formatCountry("en-GB", iso),
}));

const json = JSON.stringify(options, null, 2);
const resolved = path.join(
  path.resolve(process.cwd()),
  "internal/countries/countries.json",
);

const duplicates = {};

for (const entry of options) {
  if (duplicates[entry.labelPl] && duplicates[entry.labelPl] !== entry.value) {
    console.log(duplicates[entry.labelPl], entry.value, entry.labelPl);
    continue;
  }

  duplicates[entry.labelPl] = entry.value;
}

await Deno.writeTextFile(resolved, json);
