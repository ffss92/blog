import { search } from "./api.js";

/** @type {HTMLButtonElement} */
const searchToggle = document.getElementById("search-toggle");
/** @type {HTMLDivElement} */
const searchModal = document.getElementById("search-modal");
/** @type {HTMLDivElement} */
const searchContent = document.getElementById("search-content");
/** @type {HTMLInputElement} */
const searchInput = document.getElementById("search-input");
/** @type {HTMLDivElement} */
const searchResults = document.getElementById("search-results");
/** @type {HTMLDivElement} */
const loadingIcon = document.getElementById("search-loading-icon");
/** @type {HTMLDivElement} */
const idleIcon = document.getElementById("search-idle-icon");

searchResults.appendChild(createEmptyResult());

let open = false;

document.body.appendChild(searchModal);
document.body.addEventListener("keydown", (e) => {
  if (e.key === "k" && (e.ctrlKey || e.metaKey)) {
    e.preventDefault();
    open ? closeModal() : openModal();
  }
});

/**
 * @param {KeyboardEvent} e
 */
function handleKeydown(e) {
  const focusable = searchModal.querySelectorAll(
    "a[href], button, input, textarea, [tabindex]"
  );
  const firstEl = focusable[0];
  const lastEl = focusable[focusable.length - 1];

  switch (e.key) {
    case "Tab":
      if (e.shiftKey) {
        if (document.activeElement === firstEl) {
          e.preventDefault();
          lastEl.focus();
        }
      } else {
        if (document.activeElement === lastEl) {
          e.preventDefault();
          firstEl.focus();
        }
      }
      break;
    case "Escape":
      closeModal();
      break;
  }
}

/**
 * @param {MouseEvent} e
 */
function handleClickOutside(e) {
  if (!searchContent.contains(e.target)) {
    closeModal();
  }
}

function openModal() {
  open = true;
  searchModal.addEventListener("click", handleClickOutside);
  searchModal.addEventListener("keydown", handleKeydown);
  searchModal.hidden = false;
  document.body.classList.add("overflow-hidden");
  searchInput.focus();
}

function closeModal() {
  open = false;
  searchInput.value = "";
  resetResults();
  searchModal.removeEventListener("click", handleClickOutside);
  searchModal.removeEventListener("keydown", handleKeydown);
  searchModal.hidden = true;
  document.body.classList.remove("overflow-hidden");
}

searchToggle.addEventListener("click", () => {
  openModal();
});

let debounceId;
searchInput.addEventListener("input", async () => {
  clearTimeout(debounceId);

  showLoadingIcon();
  debounceId = setTimeout(async () => {
    if (!searchInput.value) {
      resetResults();
      showIdleIcon();
      return;
    }

    try {
      const result = await search(searchInput.value);

      if (result.articles.length === 0) {
        resetResults();
        return;
      }

      const list = document.createElement("ul");
      list.className = "flex flex-col divide-y";
      for (const article of result.articles) {
        const entry = createResultEntry(article);
        list.appendChild(entry);
      }

      searchResults.innerHTML = "";
      searchResults.appendChild(list);
    } catch (err) {
      console.error(err);
    } finally {
      showIdleIcon();
    }
  }, 250);
});

/**
 * @param {Article} article
 */
function createResultEntry(article) {
  const item = document.createElement("li");
  item.className = "p-2";
  const link = document.createElement("a");
  link.href = `/articles/${article.slug}`;
  item.appendChild(link);

  const title = document.createElement("p");
  title.innerText = article.title;
  title.className = "text-sm font-semibold";
  link.appendChild(title);

  const subtitle = document.createElement("p");
  subtitle.innerHTML = article.subtitle;
  subtitle.className = "text-xs text-stone-700";
  link.appendChild(subtitle);

  return item;
}

function createEmptyResult() {
  const p = document.createElement("p");
  p.innerText = "No results...";
  p.classList = "text-center italic text-stone-700 text-sm p-2";
  return p;
}

function resetResults() {
  searchResults.innerHTML = "";
  const p = createEmptyResult();
  searchResults.appendChild(p);
}

function showIdleIcon() {
  idleIcon.hidden = false;
  loadingIcon.hidden = true;
}

function showLoadingIcon() {
  idleIcon.hidden = true;
  loadingIcon.hidden = false;
}
