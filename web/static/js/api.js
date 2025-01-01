/**
 * @typedef {Object} Article
 * @property {string} slug
 * @property {string} title
 * @property {string} subtitle
 *
 * @typedef {Object} SearchResult
 * @property {Article[]} articles
 *
 * @param {string} q
 * @returns {Promise<SearchResult>}
 */
export async function search(q) {
  const query = new URLSearchParams({ q });
  const res = await fetch(`/api/search?${query.toString()}`);
  const data = await res.json();
  return data;
}
