/**
 * @typedef {Object} Article
 * @property {string} slug
 * @property {string} title
 * @property {string} subtitle
 *
 * @typedef {Object} SearchResult
 * @property {Article[]} articles
 *
 * @param {string} value
 * @returns {Promise<SearchResult>}
 */
export async function search(value) {
  const query = new URLSearchParams({
    q: value,
  });
  const res = await fetch(`/api/search?${query.toString()}`);
  const data = await res.json();
  return data;
}
