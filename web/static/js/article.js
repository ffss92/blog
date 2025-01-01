/**
 * Builds the table of contents for the article.
 *
 * @typedef {Object} Header
 * @property {string} id
 * @property {string} title
 * @property {string} paddingClass
 */
/** @type {HTMLDivElement} */
const toc = document.getElementById("toc");
/** @type {Header[]} */
const headers = [];
document.querySelectorAll("article > h2, h3").forEach((headerEl) => {
  const header = {
    id: headerEl.id,
    title: headerEl.textContent,
    paddingClass: "pl-0",
  };
  switch (headerEl.tagName) {
    case "H2":
      headers.push(header);
      break;
    case "H3":
      header.paddingClass = "pl-6";
      headers.push(header);
      break;
  }
});

if (headers.length) {
  toc.className = "space-y-2";

  const title = document.createElement("h3");
  title.className = "text-xl font-bold";
  title.textContent = "In this article";
  toc.appendChild(title);

  const list = document.createElement("ol");
  list.className = "flex flex-col gap-1";
  toc.appendChild(list);

  for (const header of headers) {
    const item = document.createElement("li");
    list.appendChild(item);

    const link = document.createElement("a");
    link.className = `${header.paddingClass} text-sm`;
    link.textContent = header.title;
    link.href = `#${header.id}`;
    item.appendChild(link);
  }
}

/**
 * Adds `COPY` button to code blocks.
 */
document.querySelectorAll("pre > code").forEach((el) => {
  el.parentElement.classList.add("relative", "group");

  const copyBtn = document.createElement("button");
  copyBtn.className =
    "p-2 rounded-md absolute hidden md:group-hover:block text-sm top-2 right-2";
  copyBtn.innerText = "COPY";

  let id;
  copyBtn.addEventListener("click", async () => {
    await navigator.clipboard.writeText(el.textContent);
    copyBtn.innerText = "COPIED";
    copyBtn.disabled = true;

    clearTimeout(id);
    id = setTimeout(() => {
      copyBtn.disabled = false;
      copyBtn.innerText = "COPY";
    }, 1000);
  });

  el.parentElement.appendChild(copyBtn);
});
