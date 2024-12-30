document.querySelectorAll("pre > code").forEach((el) => {
  el.parentElement.classList.add("relative", "group");

  const copyBtn = document.createElement("button");
  copyBtn.className =
    "p-2 rounded-md absolute hidden group-hover:block text-sm top-2 right-2";
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
