import type { CreateBlogPostInput, UpdateBlogPostInput } from "~/hooks";

export interface BlogPostFormValues {
  title: string;
  slug: string;
  language: string;
  body: string;
}

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value;
}

export function toCreateBlogPostInput(values: BlogPostFormValues): CreateBlogPostInput {
  return {
    title: values.title,
    slug: values.slug,
    language: values.language,
    body: blankToNull(values.body),
  };
}

export function toUpdateBlogPostInput(values: BlogPostFormValues): UpdateBlogPostInput {
  return {
    title: values.title,
    slug: values.slug,
    language: values.language,
    body: blankToNull(values.body),
  };
}
