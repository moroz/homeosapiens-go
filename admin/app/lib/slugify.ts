export function slugify(title: string): string {
  return title
    .normalize("NFD")
    .toLocaleLowerCase()
    .replaceAll("ł", "l")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-");
}
