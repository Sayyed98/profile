const API_BASE = import.meta.env.VITE_API_URL || "";

export async function fetchPortfolio() {
  const res = await fetch(`${API_BASE}/api/v1/portfolio`);
  if (!res.ok) {
    throw new Error("portfolio request failed");
  }
  return res.json();
}

export async function sendContact(payload) {
  const res = await fetch(`${API_BASE}/api/v1/contact`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || "contact request failed");
  }
  return data;
}

export function formatRange(start, end) {
  const fmt = (value) => {
    if (!value) return "Present";
    const d = new Date(value);
    return d.toLocaleDateString("en-US", { month: "short", year: "numeric", timeZone: "UTC" });
  };
  return `${fmt(start)} — ${fmt(end)}`;
}
