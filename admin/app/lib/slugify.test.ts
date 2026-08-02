import { describe, expect, test } from "vitest";
import { slugify } from "./slugify";

describe(slugify, () => {
  test("handles slug generation", () => {
    const EXAMPLES = [["Zażółć gęślą jaźń", "zazolc-gesla-jazn"]];

    for (const [input, expected] of EXAMPLES) {
      expect(slugify(input)).toEqual(expected);
    }
  });
});
