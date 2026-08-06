export function formatPrice(amount: number | string, currency = "PLN"): string {
  return new Intl.NumberFormat("en", { style: "currency", currency }).format(Number(amount));
}
